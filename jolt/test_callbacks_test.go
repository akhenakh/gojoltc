package jolt

import (
	"sync/atomic"
	"testing"
)

func TestContactListener(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	var added, removed atomic.Int64
	listener := NewContactListener(&ContactListenerCallbacks{
		OnAdded: func(body1, body2 *Body, manifold *ContactManifold, settings *ContactSettings) {
			added.Add(1)
		},
		OnRemoved: func(pair SubShapeIDPair) {
			removed.Add(1)
		},
	})
	f.system.SetContactListener(listener)

	// Ground.
	groundShape := NewBoxShape(Vec3{20, 1, 20}, DefaultConvexRadius)
	groundSettings := NewBodyCreationSettings(groundShape, RVec3{0, -1, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	groundID := bi.CreateAndAddBody(groundSettings, ActivationDontActivate)
	groundSettings.Destroy()
	groundShape.Destroy()

	// Falling sphere.
	sphereShape := NewSphereShape(0.5)
	sphereSettings := NewBodyCreationSettings(sphereShape, RVec3{0, 2, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	sphereID := bi.CreateAndAddBody(sphereSettings, ActivationActivate)
	sphereSettings.Destroy()
	sphereShape.Destroy()

	t.Cleanup(func() {
		if bi.IsAdded(sphereID) {
			bi.RemoveAndDestroyBody(sphereID)
		}
		bi.RemoveAndDestroyBody(groundID)
	})
	// The system borrows the listener: clear it before destroying the listener,
	// and do that before the system itself is destroyed.
	t.Cleanup(func() {
		f.system.SetContactListener(nil)
		listener.Destroy()
	})

	f.system.OptimizeBroadPhase()
	for i := 0; i < 120 && added.Load() == 0; i++ {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if added.Load() == 0 {
		t.Fatal("no contact added callback fired")
	}

	// Remove the sphere and step once: the contact should be reported removed.
	bi.RemoveAndDestroyBody(sphereID)
	f.system.Update(1.0/60.0, 1, f.jobSystem)
	if removed.Load() == 0 {
		t.Error("no contact removed callback fired")
	}
}

func TestBodyActivationListener(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	var activated atomic.Int64
	listener := NewBodyActivationListener(&BodyActivationCallbacks{
		OnActivated: func(bodyID BodyID, bodyUserData uint64) { activated.Add(1) },
	})
	f.system.SetBodyActivationListener(listener)

	shape := NewSphereShape(0.5)
	settings := NewBodyCreationSettings(shape, RVec3{0, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	id := bi.CreateAndAddBody(settings, ActivationActivate)
	settings.Destroy()
	shape.Destroy()

	t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })
	t.Cleanup(func() {
		f.system.SetBodyActivationListener(nil)
		listener.Destroy()
	})

	if activated.Load() == 0 {
		t.Error("activation callback did not fire")
	}
}

func TestStepListener(t *testing.T) {
	f := newPhysicsFixture(t)

	var steps atomic.Int64
	listener := NewStepListener(&StepListenerCallbacks{
		OnStep: func(ctx StepContext) {
			if ctx.DeltaTime <= 0 {
				t.Errorf("delta time = %v", ctx.DeltaTime)
			}
			steps.Add(1)
		},
	})
	defer listener.Destroy()
	f.system.AddStepListener(listener)

	for range 5 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	if got := steps.Load(); got != 5 {
		t.Errorf("steps = %d, want 5", got)
	}

	f.system.RemoveStepListener(listener)
	f.system.Update(1.0/60.0, 1, f.jobSystem)
	if got := steps.Load(); got != 5 {
		t.Errorf("steps after removal = %d, want 5", got)
	}
}

func TestDebugRendererDrawBodies(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewBoxShape(Vec3{1, 1, 1}, DefaultConvexRadius)
	settings := NewBodyCreationSettings(shape, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	id := bi.CreateAndAddBody(settings, ActivationDontActivate)
	settings.Destroy()
	shape.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })

	var lines, triangles atomic.Int64
	renderer := NewDebugRenderer(&DebugRendererCallbacks{
		DrawLine: func(from, to RVec3, color Color) { lines.Add(1) },
		DrawTriangle: func(v1, v2, v3 RVec3, color Color, castShadow bool) {
			triangles.Add(1)
		},
	})
	defer renderer.Destroy()

	renderer.NextFrame()
	f.system.DrawBodies(DefaultDrawSettings(), renderer)

	if lines.Load() == 0 && triangles.Load() == 0 {
		t.Fatal("debug renderer received no geometry")
	}
}
