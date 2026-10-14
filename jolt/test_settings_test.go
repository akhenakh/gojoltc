package jolt

import "testing"

func TestPhysicsSettings(t *testing.T) {
	f := newPhysicsFixture(t)

	def := f.system.PhysicsSettings()
	if def.NumVelocitySteps == 0 || def.NumPositionSteps == 0 {
		t.Fatalf("unexpected defaults: %+v", def)
	}

	def.NumVelocitySteps = 8
	def.NumPositionSteps = 3
	def.AllowSleeping = false
	f.system.SetPhysicsSettings(def)

	got := f.system.PhysicsSettings()
	if got.NumVelocitySteps != 8 || got.NumPositionSteps != 3 {
		t.Errorf("steps not applied: velocity=%d position=%d", got.NumVelocitySteps, got.NumPositionSteps)
	}
	if got.AllowSleeping {
		t.Error("allow sleeping not applied")
	}
}

func TestPhysicsMaterialOnShape(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	material := NewPhysicsMaterial("ice", Color(0xff0000ff))
	defer material.Destroy()

	shape := NewPlaneShapeWithMaterial(Plane{Normal: Vec3{0, 1, 0}}, material, 10)
	defer shape.Destroy()

	if m := shape.Material(0); m == nil || m.DebugName() != "ice" {
		t.Fatalf("shape material = %v", m)
	}

	settings := NewBodyCreationSettings(shape, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	id := bi.CreateAndAddBody(settings, ActivationDontActivate)
	settings.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })

	if m := bi.GetMaterial(id, 0); m == nil || m.DebugName() != "ice" {
		t.Errorf("body material = %v", m)
	}
}

// dropSphereOverFloor builds a floor (sub-group 0) and a sphere (sub-group 1)
// with the given collision group filter and returns the sphere's height after
// two seconds. Note Jolt's table filter never lets two bodies in the *same*
// sub-group collide, so the pair under test must differ.
func dropSphereOverFloor(t *testing.T, group CollisionGroup) Real {
	t.Helper()
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	floorGroup := group
	floorGroup.SubGroupID = 0
	sphereGroup := group
	sphereGroup.SubGroupID = 1

	floorShape := NewBoxShape(Vec3{50, 1, 50}, DefaultConvexRadius)
	floorSettings := NewBodyCreationSettings(floorShape, RVec3{0, -1, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	floorSettings.SetCollisionGroup(floorGroup)
	floorID := bi.CreateAndAddBody(floorSettings, ActivationDontActivate)
	floorSettings.Destroy()
	floorShape.Destroy()

	sphereShape := NewSphereShape(0.5)
	sphereSettings := NewBodyCreationSettings(sphereShape, RVec3{0, 2, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	sphereSettings.SetCollisionGroup(sphereGroup)
	sphereID := bi.CreateAndAddBody(sphereSettings, ActivationActivate)
	sphereSettings.Destroy()
	sphereShape.Destroy()
	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(sphereID)
		bi.RemoveAndDestroyBody(floorID)
	})

	f.system.OptimizeBroadPhase()
	for range 180 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	return bi.Position(sphereID).Y
}

func TestCollisionGroupDisablesPair(t *testing.T) {
	filter := NewGroupFilterTable(2)
	defer filter.Destroy()
	filter.DisableCollision(0, 1)

	y := dropSphereOverFloor(t, CollisionGroup{Filter: filter, GroupID: 0})
	if y > 0 {
		t.Errorf("sphere should fall through with the pair disabled: y = %v", y)
	}
}

func TestCollisionGroupEnablesPair(t *testing.T) {
	filter := NewGroupFilterTable(2)
	defer filter.Destroy()
	filter.EnableCollision(0, 1)

	if !filter.IsCollisionEnabled(0, 1) {
		t.Fatal("pair should be enabled")
	}
	y := dropSphereOverFloor(t, CollisionGroup{Filter: filter, GroupID: 0})
	if y < 0 {
		t.Errorf("sphere should rest on the floor with the pair enabled: y = %v", y)
	}
}

func TestBodyInterfaceBatchActivation(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	defer shape.Destroy()

	ids := make([]BodyID, 3)
	for i := range ids {
		settings := NewBodyCreationSettings(shape, RVec3{Real(i) * 2, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		ids[i] = bi.CreateAndAddBody(settings, ActivationActivate)
		settings.Destroy()
	}
	t.Cleanup(func() {
		for _, id := range ids {
			bi.RemoveAndDestroyBody(id)
		}
	})

	bi.DeactivateBodies(ids)
	for _, id := range ids {
		if bi.IsActive(id) {
			t.Errorf("body %d should be inactive", id)
		}
	}

	bi.ActivateBodies(ids)
	for _, id := range ids {
		if !bi.IsActive(id) {
			t.Errorf("body %d should be active", id)
		}
	}
}

func TestBodyInterfaceVelocityAndBuoyancy(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	settings := NewBodyCreationSettings(shape, RVec3{0, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	id := bi.CreateAndAddBody(settings, ActivationActivate)
	settings.Destroy()
	shape.Destroy()
	t.Cleanup(func() {
		if bi.IsAdded(id) {
			bi.RemoveAndDestroyBody(id)
		}
	})

	bi.SetLinearAndAngularVelocity(id, Vec3{}, Vec3{})
	bi.AddLinearAndAngularVelocity(id, Vec3{1, 0, 0}, Vec3{0, 0, 1})
	lv, av := bi.GetLinearAndAngularVelocity(id)
	equalFloat(t, "linear.x", lv.X, 1)
	equalFloat(t, "angular.z", av.Z, 1)

	// Buoyancy at a surface should push the body up.
	bi.SetLinearAndAngularVelocity(id, Vec3{}, Vec3{})
	if !bi.ApplyBuoyancyImpulse(id, RVec3{0, 5, 0}, Vec3{0, 1, 0}, 1.0, 0.1, 0.1, Vec3{}, Vec3{0, -10, 0}, 1.0/60.0) {
		t.Error("ApplyBuoyancyImpulse reported no surface hit")
	}
	if v := bi.LinearVelocity(id); v.Y <= 0 {
		t.Errorf("buoyancy did not push up: %v", v)
	}
}

func TestBodyWithoutIDAssign(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	defer shape.Destroy()
	settings := NewBodyCreationSettings(shape, RVec3{0, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBodyWithoutID(settings)
	if body == nil {
		t.Fatal("body is nil")
	}
	if !bi.AssignBodyID(body) {
		t.Fatal("AssignBodyID failed")
	}
	id := body.ID()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })

	bi.AddBody(id, ActivationActivate)
	if !bi.IsAdded(id) {
		t.Error("body should be added")
	}
}
