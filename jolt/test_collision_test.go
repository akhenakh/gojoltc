package jolt

import "testing"

// Ported from third_party/joltc/tests/test_collision.cpp.

func createSphereBody(t *testing.T, f *physicsFixture, x, y, z Real, radius float32, motionType MotionType, layer ObjectLayer) BodyID {
	t.Helper()
	bi := f.system.BodyInterface()
	shape := NewSphereShape(radius)
	settings := NewBodyCreationSettings(shape, RVec3{x, y, z}, QuatIdentity(), motionType, layer)
	body := bi.CreateBody(settings)
	id := body.ID()
	bi.AddBody(id, ActivationActivate)
	settings.Destroy()
	shape.Destroy()
	return id
}

func TestCollisionGetNarrowPhaseQuery(t *testing.T) {
	f := newPhysicsFixture(t)
	if f.system.NarrowPhaseQuery() == nil {
		t.Fatal("narrow phase query is nil")
	}
}

func TestCollisionGetNarrowPhaseQueryNoLock(t *testing.T) {
	f := newPhysicsFixture(t)
	if f.system.NarrowPhaseQueryNoLock() == nil {
		t.Fatal("no-lock narrow phase query is nil")
	}
}

func TestCollisionRayCastHit(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	sphereID := createSphereBody(t, f, 0, 0, 0, 1, MotionTypeStatic, layerNonMoving)
	f.system.OptimizeBroadPhase()

	query := f.system.NarrowPhaseQuery()
	hit, ok := query.CastRay(RVec3{-5, 0, 0}, Vec3{10, 0, 0})
	if !ok {
		t.Fatal("expected hit")
	}
	nearFloat(t, "fraction", hit.Fraction, 0.4)
	if hit.BodyID != sphereID {
		t.Errorf("body id = %d, want %d", hit.BodyID, sphereID)
	}

	bi.RemoveAndDestroyBody(sphereID)
}

func TestCollisionRayCastMiss(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	sphereID := createSphereBody(t, f, 0, 0, 0, 1, MotionTypeStatic, layerNonMoving)
	f.system.OptimizeBroadPhase()

	if _, ok := f.system.NarrowPhaseQuery().CastRay(RVec3{0, 10, 0}, Vec3{10, 0, 0}); ok {
		t.Error("expected miss")
	}

	bi.RemoveAndDestroyBody(sphereID)
}

func TestCollisionRayCastClosestHit(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	sphere1ID := createSphereBody(t, f, -3, 0, 0, 1, MotionTypeStatic, layerNonMoving)
	sphere2ID := createSphereBody(t, f, 3, 0, 0, 1, MotionTypeDynamic, layerMoving)
	f.system.OptimizeBroadPhase()

	hit, ok := f.system.NarrowPhaseQuery().CastRay(RVec3{-10, 0, 0}, Vec3{20, 0, 0})
	if !ok {
		t.Fatal("expected hit")
	}
	if hit.BodyID != sphere1ID {
		t.Errorf("body id = %d, want %d", hit.BodyID, sphere1ID)
	}

	bi.RemoveAndDestroyBody(sphere1ID)
	bi.RemoveAndDestroyBody(sphere2ID)
}

func TestCollisionGetBroadPhaseQuery(t *testing.T) {
	f := newPhysicsFixture(t)
	if f.system.BroadPhaseQuery() == nil {
		t.Fatal("broad phase query is nil")
	}
}
