package jolt

import (
	"math"
	"testing"
)

// Behavioural tests for the usability pass: bodies driven through the extended
// BodyInterface, plus mesh and height field level geometry.

func TestBodyInterfaceExtras(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	settings := NewBodyCreationSettings(shape, RVec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	id := bi.CreateAndAddBody(settings, ActivationActivate)
	settings.Destroy()
	t.Cleanup(func() {
		if bi.IsAdded(id) {
			bi.RemoveAndDestroyBody(id)
		}
	})

	// GetShape returns the shape we created.
	if s := bi.GetShape(id); s == nil || s.SubType() != ShapeSubTypeSphere {
		t.Fatalf("GetShape mismatch: %v", s)
	}

	// Batch transform read.
	m := bi.WorldTransform(id)
	_ = m

	// Both velocities in one call.
	bi.SetLinearAndAngularVelocity(id, Vec3{1, 2, 3}, Vec3{0, 0, 1})
	lv, av := bi.GetLinearAndAngularVelocity(id)
	equalFloat(t, "linear.x", lv.X, 1)
	equalFloat(t, "angular.z", av.Z, 1)

	// Object layer round-trip.
	bi.SetObjectLayer(id, layerNonMoving)
	if got := bi.ObjectLayer(id); got != layerNonMoving {
		t.Errorf("object layer = %d, want %d", got, layerNonMoving)
	}

	// User data round-trip (generated accessor).
	bi.SetUserData(id, 7)
	if got := bi.UserData(id); got != 7 {
		t.Errorf("user data = %d, want 7", got)
	}

	// A force at a point produces torque: the body should start rotating.
	bi.SetAngularVelocity(id, Vec3{})
	bi.SetLinearVelocity(id, Vec3{})
	bi.AddForceAtPoint(id, Vec3{0, 0, 100}, joltRVec3(1, 0, 0))
	f.system.Update(1.0/60.0, 1, f.jobSystem)
	if _, av := bi.GetLinearAndAngularVelocity(id); av.Y == 0 {
		t.Error("force at a point produced no torque")
	}
}

func joltRVec3(x, y, z Real) RVec3 { return RVec3{x, y, z} }

func TestPhysicsMaterial(t *testing.T) {
	initShapes(t)
	m := NewPhysicsMaterial("ground", Color(0xff00ff00))
	defer m.Destroy()
	if got := m.DebugName(); got != "ground" {
		t.Errorf("name = %q, want ground", got)
	}
	if got := m.DebugColor(); got != Color(0xff00ff00) {
		t.Errorf("color = %#x", uint32(got))
	}
}

// A static mesh floor must hold a falling sphere.
func TestMeshShapeHoldsBody(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	// A 20x20 quad made of two triangles at y = 0.
	half := float32(20)
	vertices := []Vec3{
		{-half, 0, -half}, {half, 0, -half}, {half, 0, half}, {-half, 0, half},
	}
	// Wind counter-clockwise seen from above so the face normals point up (+Y).
	triangles := []IndexedTriangle{{I1: 0, I2: 2, I3: 1}, {I1: 0, I2: 3, I3: 2}}
	mesh := NewMeshShapeFromVertices(vertices, triangles)
	if mesh == nil {
		t.Fatal("mesh shape is nil")
	}
	defer mesh.Destroy()

	floorSettings := NewBodyCreationSettings(mesh, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	floorID := bi.CreateAndAddBody(floorSettings, ActivationDontActivate)
	floorSettings.Destroy()

	sphereShape := NewSphereShape(0.5)
	sphereSettings := NewBodyCreationSettings(sphereShape, RVec3{0, 3, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	sphereID := bi.CreateAndAddBody(sphereSettings, ActivationActivate)
	sphereSettings.Destroy()
	sphereShape.Destroy()
	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(sphereID)
		bi.RemoveAndDestroyBody(floorID)
	})

	f.system.OptimizeBroadPhase()
	for range 120 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if y := bi.Position(sphereID).Y; y < 0 {
		t.Errorf("sphere fell through the mesh floor: y = %v", y)
	}
}

// A height field must hold a falling sphere too.
func TestHeightFieldShapeHoldsBody(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	const n = 16
	samples := make([]float32, n*n) // flat, y = 0
	heightfield := NewHeightFieldShape(samples, Vec3{-n / 2, 0, -n / 2}, Vec3{1, 1, 1})
	if heightfield == nil {
		t.Fatal("height field shape is nil")
	}
	defer heightfield.Destroy()

	floorSettings := NewBodyCreationSettings(heightfield, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	floorID := bi.CreateAndAddBody(floorSettings, ActivationDontActivate)
	floorSettings.Destroy()

	sphereShape := NewSphereShape(0.5)
	sphereSettings := NewBodyCreationSettings(sphereShape, RVec3{0, 3, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	sphereID := bi.CreateAndAddBody(sphereSettings, ActivationActivate)
	sphereSettings.Destroy()
	sphereShape.Destroy()
	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(sphereID)
		bi.RemoveAndDestroyBody(floorID)
	})

	f.system.OptimizeBroadPhase()
	for range 120 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if y := bi.Position(sphereID).Y; y < -0.5 {
		t.Errorf("sphere fell through the height field: y = %v", y)
	}
}

// Behavioural test for the mass override. Jolt's default density would give this
// 1 m cube a mass of ~1000 kg, so imposing 11.3 kg is directly visible in the
// response to a known impulse; OverrideMassPropertiesCalculateInertia means the
// inertia must then come from the shape (scaled), not from the identity that
// NewMassProperties hands in.
func TestBodyMassOverride(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	const mass = 11.3
	shape := NewBoxShape(Vec3{0.5, 0.5, 0.5}, 0) // no convex radius: a plain 1 m cube
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, RVec3{0, 20, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()
	if got := settings.OverrideMassProperties(); got != OverrideMassPropertiesCalculateMassAndInertia {
		t.Fatalf("default override mode = %v, want CalculateMassAndInertia", got)
	}
	settings.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
	settings.SetMassPropertiesOverride(NewMassProperties(mass))

	id := bi.CreateAndAddBody(settings, ActivationActivate)
	t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })
	body := f.system.BodyPtr(id)

	// Linear: an impulse of `mass` N.s gives exactly 1 m/s to a mass-kg body.
	bi.AddImpulse(id, Vec3{mass, 0, 0})
	lin, _ := bi.GetLinearAndAngularVelocity(id)
	if math.Abs(float64(lin.X)-1) > 1e-3 {
		t.Errorf("vx after a %v N.s impulse = %v, want 1 (default density gives ~0.011)", float32(mass), lin.X)
	}

	// Angular: a solid 1 m cube of mass m has inertia m/6 about every axis, so a
	// torque impulse of m/6 gives 1 rad/s. The identity inertia that was passed in
	// would give ~1.88 instead.
	body.AddAngularImpulse(Vec3{0, mass / 6, 0})
	_, ang := bi.GetLinearAndAngularVelocity(id)
	if math.Abs(float64(ang.Y)-1) > 1e-3 {
		t.Errorf("wy after a %v N.m.s impulse = %v, want 1 (inertia must come from the shape)", float32(mass/6), ang.Y)
	}
}

// The code generator used to mistake the non-const pointer *inputs* of these
// mutators for out-parameters: it dropped the argument, handed Jolt an
// uninitialised value and returned that same value to the caller.
func TestBodyInterfaceVectorMutators(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	const mass = 11.3 // inertia m/6 about each axis for this 1 m cube
	shape := NewBoxShape(Vec3{0.5, 0.5, 0.5}, 0)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, RVec3{0, 20, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	settings.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
	settings.SetMassPropertiesOverride(NewMassProperties(mass))
	id := bi.CreateAndAddBody(settings, ActivationActivate)
	settings.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })

	// AddLinearVelocity adds its argument instead of discarding it.
	bi.AddLinearVelocity(id, Vec3{2, 0, 0})
	if lin, _ := bi.GetLinearAndAngularVelocity(id); math.Abs(float64(lin.X)-2) > 1e-4 {
		t.Errorf("vx after AddLinearVelocity = %v, want 2", lin.X)
	}

	// AddAngularImpulse: an impulse of I gives 1 rad/s.
	bi.AddAngularImpulse(id, Vec3{0, mass / 6, 0})
	if _, ang := bi.GetLinearAndAngularVelocity(id); math.Abs(float64(ang.Y)-1) > 1e-3 {
		t.Errorf("wy after AddAngularImpulse = %v, want 1", ang.Y)
	}

	// AddTorque accumulates a torque that the next step applies: dW = torque*dt/I.
	bi.AddTorque(id, Vec3{0, mass / 6 * 60, 0})
	f.system.Update(1.0/60.0, 1, f.jobSystem)
	if _, ang := bi.GetLinearAndAngularVelocity(id); math.Abs(float64(ang.Y)-2) > 0.1 {
		t.Errorf("wy after AddTorque = %v, want ~2 (1 from the impulse plus one step of torque)", ang.Y)
	}
}

// Pairs with TestBodyMassOverride: a body reports the mass and the shape-derived
// inertia Jolt is actually using, which is what makes the override verifiable
// (and what telemetry needs).
func TestBodyMassPropertiesReadBack(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewBoxShape(Vec3{0.5, 0.5, 0.5}, 0) // a plain 1 m cube
	defer shape.Destroy()

	dyn := NewBodyCreationSettings(shape, RVec3{0, 20, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	dyn.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
	dyn.SetMassPropertiesOverride(NewMassProperties(6))
	dynID := bi.CreateAndAddBody(dyn, ActivationActivate)
	dyn.Destroy()

	kin := NewBodyCreationSettings(shape, RVec3{5, 20, 0}, QuatIdentity(), MotionTypeKinematic, layerMoving)
	kinID := bi.CreateAndAddBody(kin, ActivationActivate)
	kin.Destroy()

	stat := NewBodyCreationSettings(shape, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	statID := bi.CreateAndAddBody(stat, ActivationDontActivate)
	stat.Destroy()

	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(dynID)
		bi.RemoveAndDestroyBody(kinID)
		bi.RemoveAndDestroyBody(statID)
	})

	body := f.system.BodyPtr(dynID)
	if got := body.InverseMass(); math.Abs(float64(got)-1.0/6.0) > 1e-4 {
		t.Errorf("inverse mass = %v, want 1/6", got)
	}

	mp := body.MassProperties()
	if math.Abs(float64(mp.Mass)-6) > 1e-3 {
		t.Errorf("mass = %v, want 6", mp.Mass)
	}
	// A solid 1 m cube of mass m has inertia m/6 about every axis, so with the
	// override at 6 the inertia must read 1 — the shape-derived value, not the
	// identity NewMassProperties hands in.
	for _, ax := range []struct {
		name string
		got  float32
	}{
		{"xx", mp.Inertia.Col[0].X},
		{"yy", mp.Inertia.Col[1].Y},
		{"zz", mp.Inertia.Col[2].Z},
	} {
		if math.Abs(float64(ax.got)-1) > 1e-3 {
			t.Errorf("inertia %s = %v, want 1 (m/6 for this cube)", ax.name, ax.got)
		}
	}

	// A kinematic body still carries mass (Jolt derives it from the shape's
	// density) even though forces never move it, and a static body carries none.
	// Neither may crash here.
	if got := f.system.BodyPtr(kinID).InverseMass(); got <= 0 {
		t.Errorf("kinematic inverse mass = %v, want a positive shape-derived value", got)
	}
	staticBody := f.system.BodyPtr(statID)
	if got := staticBody.InverseMass(); got != 0 {
		t.Errorf("static inverse mass = %v, want 0", got)
	}
	if got := staticBody.MassProperties().Mass; got != math.MaxFloat32 {
		t.Errorf("static mass = %v, want Jolt's FLT_MAX", got)
	}
}
