package jolt

import (
	"math"
	"testing"
)

// sensorScene builds a static floor (NON_MOVING) and a dynamic box (MOVING)
// directly above it, returning the system, the interface and both body ids.
func sensorScene(t *testing.T) (*physicsFixture, *BodyInterface, BodyID, BodyID) {
	t.Helper()
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	floorShape := NewBoxShape(Vec3{20, 1, 20}, DefaultConvexRadius)
	floorSettings := NewBodyCreationSettings(floorShape, RVec3{0, -1, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	floorID := bi.CreateAndAddBody(floorSettings, ActivationDontActivate)
	floorSettings.Destroy()
	floorShape.Destroy()

	boxShape := NewBoxShape(Vec3{0.5, 0.5, 0.5}, DefaultConvexRadius)
	boxSettings := NewBodyCreationSettings(boxShape, RVec3{0, 3, 0}, QuatIdentity(), MotionTypeStatic, layerMoving)
	boxID := bi.CreateAndAddBody(boxSettings, ActivationActivate)
	boxSettings.Destroy()
	boxShape.Destroy()

	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(boxID)
		bi.RemoveAndDestroyBody(floorID)
	})
	f.system.OptimizeBroadPhase()
	return f, bi, floorID, boxID
}

func TestRayCastFilters(t *testing.T) {
	f, _, floorID, boxID := sensorScene(t)
	query := f.system.NarrowPhaseQuery()
	origin := RVec3{0, 10, 0}
	down := Vec3{0, -20, 0}

	// Unfiltered: the box (higher up) is hit first.
	hit, ok := query.CastRayClosest(origin, down, DefaultRayCastSettings(), nil)
	if !ok || hit.BodyID != boxID {
		t.Fatalf("unfiltered closest = %v (ok=%v), want the box %d", hit.BodyID, ok, boxID)
	}

	// Body filter excluding the box: the floor is hit instead.
	bodyFilter := NewBodyFilter(&BodyFilterCallbacks{
		ShouldCollide:       func(id BodyID) bool { return id != boxID },
		ShouldCollideLocked: func(body *Body) bool { return body.ID() != boxID },
	})
	defer bodyFilter.Destroy()
	hit, ok = query.CastRayClosest(origin, down, DefaultRayCastSettings(), &QueryOptions{Body: bodyFilter})
	if !ok || hit.BodyID != floorID {
		t.Fatalf("body-filtered closest = %v (ok=%v), want the floor %d", hit.BodyID, ok, floorID)
	}

	// Object layer filter that only accepts NON_MOVING: the floor again.
	layerFilter := NewObjectLayerFilter(&ObjectLayerFilterCallbacks{
		ShouldCollide: func(layer ObjectLayer) bool { return layer == layerNonMoving },
	})
	defer layerFilter.Destroy()
	hit, ok = query.CastRayClosest(origin, down, DefaultRayCastSettings(), &QueryOptions{ObjectLayer: layerFilter})
	if !ok || hit.BodyID != floorID {
		t.Fatalf("layer-filtered closest = %v (ok=%v), want the floor %d", hit.BodyID, ok, floorID)
	}

	// Object layer filter that only accepts MOVING: the box.
	movingOnly := NewObjectLayerFilter(&ObjectLayerFilterCallbacks{
		ShouldCollide: func(layer ObjectLayer) bool { return layer == layerMoving },
	})
	defer movingOnly.Destroy()
	hit, ok = query.CastRayClosest(origin, down, DefaultRayCastSettings(), &QueryOptions{ObjectLayer: movingOnly})
	if !ok || hit.BodyID != boxID {
		t.Fatalf("moving-only closest = %v (ok=%v), want the box %d", hit.BodyID, ok, boxID)
	}
}

func TestRayCastAllHits(t *testing.T) {
	f, _, floorID, boxID := sensorScene(t)
	query := f.system.NarrowPhaseQuery()

	hits := query.CastRayAll(RVec3{0, 10, 0}, Vec3{0, -20, 0}, DefaultRayCastSettings(), nil)
	if len(hits) < 2 {
		t.Fatalf("got %d hits, want at least 2 (box + floor)", len(hits))
	}
	seen := map[BodyID]bool{}
	for _, h := range hits {
		seen[h.BodyID] = true
	}
	if !seen[boxID] || !seen[floorID] {
		t.Errorf("hits %v missing box(%d) or floor(%d)", seen, boxID, floorID)
	}
}

func TestCollidePointQuery(t *testing.T) {
	f, _, _, boxID := sensorScene(t)
	query := f.system.NarrowPhaseQuery()

	hits := query.CollidePoint(RVec3{0, 3, 0}, nil) // inside the box
	if len(hits) == 0 {
		t.Fatal("point inside the box returned no hits")
	}
	found := false
	for _, h := range hits {
		if h.BodyID == boxID {
			found = true
		}
	}
	if !found {
		t.Errorf("collide point did not report the box %d", boxID)
	}

	if n := len(query.CollidePoint(RVec3{50, 50, 50}, nil)); n != 0 {
		t.Errorf("point in empty space reported %d hits", n)
	}
}

func TestCollideShapeProximity(t *testing.T) {
	f, _, _, boxID := sensorScene(t)
	query := f.system.NarrowPhaseQuery()

	probe := NewSphereShape(0.2)
	defer probe.Destroy()

	// 0.3 m above the box top (box top at y = 3.5), so no overlap but within 0.5.
	transform := RMat4Translation(RVec3{0, 3.8, 0})

	if hits := query.CollideShape(probe, Vec3{1, 1, 1}, transform, DefaultCollideShapeSettings(), nil); len(hits) != 0 {
		t.Errorf("expected no overlap without separation distance, got %d", len(hits))
	}

	settings := CollideShapeSettings{MaxSeparationDistance: 0.5}
	hits := query.CollideShape(probe, Vec3{1, 1, 1}, transform, settings, nil)
	if len(hits) == 0 {
		t.Fatal("proximity query with max separation 0.5 found nothing")
	}
	found := false
	for _, h := range hits {
		if h.BodyID2 == boxID {
			found = true
		}
	}
	if !found {
		t.Errorf("proximity query did not report the box %d", boxID)
	}
}

func TestCastShapeSweep(t *testing.T) {
	f, _, _, boxID := sensorScene(t)
	query := f.system.NarrowPhaseQuery()

	probe := NewSphereShape(0.2)
	defer probe.Destroy()

	hits := query.CastShape(probe, RMat4Translation(RVec3{0, 10, 0}), Vec3{0, -20, 0}, DefaultShapeCastSettings(), nil)
	if len(hits) == 0 {
		t.Fatal("shape sweep found nothing")
	}
	found := false
	for _, h := range hits {
		if h.BodyID2 == boxID && h.Fraction > 0 && h.Fraction < 1 {
			found = true
		}
	}
	if !found {
		t.Errorf("sweep did not report the box %d with a fraction in (0,1): %+v", boxID, hits)
	}
}

// The native fast path must agree with the collector path ray for ray: it is the
// same query with the closest-hit early-out left to Jolt (no per-hit callback),
// and it is what a LIDAR sweep uses.
func TestCastRaysFiltered(t *testing.T) {
	f, _, floorID, boxID := sensorScene(t)
	query := f.system.NarrowPhaseQuery()

	const n = 32
	origins := make([]RVec3, n)
	dirs := make([]Vec3, n)
	for i := range n {
		// A fan of rays from above: some hit the floor, some the box, some miss.
		x := Real(i)*0.35 - 5
		origins[i] = RVec3{x, 10, 0}
		dirs[i] = Vec3{0, -20, 0} // the direction's length is the ray's
	}
	out := make([]RayCastResult, n)

	hits := query.CastRaysFiltered(origins, dirs, out, nil)
	if hits == 0 {
		t.Fatal("no ray hit anything")
	}
	agreed := 0
	for i := range n {
		want, ok := query.CastRayClosest(origins[i], dirs[i], DefaultRayCastSettings(), nil)
		if !ok {
			if out[i].BodyID != InvalidBodyID {
				t.Errorf("ray %d: batch hit %v, collector missed", i, out[i].BodyID)
			}
			continue
		}
		if !ok || out[i].BodyID != want.BodyID || math.Abs(float64(out[i].Fraction-want.Fraction)) > 1e-4 {
			t.Errorf("ray %d: batch (%v, %v) != collector (%v, %v)", i, out[i].BodyID, out[i].Fraction, want.BodyID, want.Fraction)
		}
		agreed++
	}
	if hits != agreed {
		t.Errorf("batch reported %d hits, the collector agrees on %d", hits, agreed)
	}

	// The single-ray fast path agrees, and honours a body filter: with the floor
	// excluded, the box is hit instead, and away from the box nothing is left.
	one, ok := query.CastRayFiltered(RVec3{0, 10, 0}, Vec3{0, -20, 0}, nil)
	if !ok || one.BodyID != boxID {
		t.Errorf("CastRayFiltered = (%v, %v), want the box %d", one.BodyID, ok, boxID)
	}
	floorOut := NewBodyFilter(&BodyFilterCallbacks{
		ShouldCollide:       func(id BodyID) bool { return id != floorID },
		ShouldCollideLocked: func(body *Body) bool { return body.ID() != floorID },
	})
	defer floorOut.Destroy()
	opts := &QueryOptions{Body: floorOut}
	if hit, ok := query.CastRayFiltered(RVec3{0, 10, 0}, Vec3{0, -20, 0}, opts); !ok || hit.BodyID != boxID {
		t.Errorf("with the floor excluded the fast path = (%v, %v), want the box %d", hit.BodyID, ok, boxID)
	}
	if _, ok := query.CastRayFiltered(RVec3{-5, 10, 0}, Vec3{0, -20, 0}, opts); ok {
		t.Error("the fast path still hit the excluded floor")
	}

	// A shape filter is not part of the native cast, so it falls back to the
	// collector path rather than silently ignoring the filter.
	shapeFilter := NewShapeFilter(&ShapeFilterCallbacks{
		ShouldCollide: func(shape2 *Shape, subShapeID2 SubShapeID) bool { return true },
	})
	defer shapeFilter.Destroy()
	fallback, ok := query.CastRayFiltered(RVec3{0, 10, 0}, Vec3{0, -20, 0}, &QueryOptions{Shape: shapeFilter})
	if !ok || fallback.BodyID != boxID {
		t.Errorf("shape-filtered fallback = (%v, %v), want the box %d", fallback.BodyID, ok, boxID)
	}
}
