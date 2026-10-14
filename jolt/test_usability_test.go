package jolt

import "testing"

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
