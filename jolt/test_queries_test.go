package jolt

import "testing"

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
