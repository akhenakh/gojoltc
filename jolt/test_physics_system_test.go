package jolt

import "testing"

// Ported from third_party/joltc/tests/test_physics_system.cpp.

const (
	layerNonMoving      ObjectLayer = 0
	layerMoving         ObjectLayer = 1
	numObjectLayers                 = 2
	numBroadPhaseLayers             = 2
	invalidBodyID       BodyID      = 0xFFFFFFFF
)

type physicsFixture struct {
	jobSystem *JobSystem
	bpLayer   *BroadPhaseLayerInterface
	pair      *ObjectLayerPairFilter
	objBp     *ObjectVsBroadPhaseLayerFilter
	system    *PhysicsSystem
}

func newPhysicsFixture(t *testing.T) *physicsFixture {
	t.Helper()

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	f := &physicsFixture{}
	f.jobSystem = NewJobSystemThreadPool(nil)
	if f.jobSystem == nil {
		t.Fatal("job system is nil")
	}

	f.bpLayer = NewBroadPhaseLayerInterfaceTable(numObjectLayers, numBroadPhaseLayers)
	f.bpLayer.MapObjectToBroadPhaseLayer(layerNonMoving, 0)
	f.bpLayer.MapObjectToBroadPhaseLayer(layerMoving, 1)

	f.pair = NewObjectLayerPairFilterTable(numObjectLayers)
	f.pair.EnableCollision(layerNonMoving, layerMoving)
	f.pair.EnableCollision(layerMoving, layerMoving)

	f.objBp = NewObjectVsBroadPhaseLayerFilterTable(f.bpLayer, numBroadPhaseLayers, f.pair, numObjectLayers)

	f.system = NewPhysicsSystem(&PhysicsSystemSettings{
		MaxBodies:                     1024,
		NumBodyMutexes:                0,
		MaxBodyPairs:                  1024,
		MaxContactConstraints:         1024,
		BroadPhaseLayerInterface:      f.bpLayer,
		ObjectLayerPairFilter:         f.pair,
		ObjectVsBroadPhaseLayerFilter: f.objBp,
	})
	if f.system == nil {
		t.Fatal("physics system is nil")
	}

	t.Cleanup(func() {
		f.system.Destroy()
		f.jobSystem.Destroy()
		Shutdown()
	})
	return f
}

func TestPhysicsSystemCreate(t *testing.T) {
	f := newPhysicsFixture(t)
	if f.system == nil {
		t.Fatal("physics system is nil")
	}
}

func TestPhysicsSystemGetBodyInterface(t *testing.T) {
	f := newPhysicsFixture(t)
	if f.system.BodyInterface() == nil {
		t.Fatal("body interface is nil")
	}
}

func TestPhysicsSystemGetBodyInterfaceNoLock(t *testing.T) {
	f := newPhysicsFixture(t)
	if f.system.BodyInterfaceNoLock() == nil {
		t.Fatal("no-lock body interface is nil")
	}
}

func TestPhysicsSystemSetGravity(t *testing.T) {
	f := newPhysicsFixture(t)
	f.system.SetGravity(Vec3{0, -10, 0})
	g := f.system.Gravity()
	equalFloat(t, "gravity.x", g.X, 0)
	equalFloat(t, "gravity.y", g.Y, -10)
	equalFloat(t, "gravity.z", g.Z, 0)
}

func TestPhysicsSystemNumBodiesInitial(t *testing.T) {
	f := newPhysicsFixture(t)
	if n := f.system.NumBodies(); n != 0 {
		t.Fatalf("num bodies = %d, want 0", n)
	}
}

func TestPhysicsSystemNumActiveBodiesInitial(t *testing.T) {
	f := newPhysicsFixture(t)
	if n := f.system.NumActiveBodies(BodyTypeRigid); n != 0 {
		t.Fatalf("num active bodies = %d, want 0", n)
	}
}

func TestPhysicsSystemMaxBodies(t *testing.T) {
	f := newPhysicsFixture(t)
	if n := f.system.MaxBodies(); n != 1024 {
		t.Fatalf("max bodies = %d, want 1024", n)
	}
}

func TestPhysicsSystemCreateBodySphere(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	if shape == nil {
		t.Fatal("shape is nil")
	}
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	if body == nil {
		t.Fatal("body is nil")
	}
	bodyID := body.ID()
	if bodyID == invalidBodyID {
		t.Fatalf("body id is invalid (%#x)", uint32(bodyID))
	}

	if n := f.system.NumBodies(); n != 1 {
		t.Fatalf("num bodies = %d, want 1", n)
	}

	bi.AddBody(bodyID, ActivationActivate)
	if n := f.system.NumActiveBodies(BodyTypeRigid); n != 1 {
		t.Fatalf("num active bodies = %d, want 1", n)
	}

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemCreateBodyBox(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewBoxShape(Vec3{1, 1, 1}, DefaultConvexRadius)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationDontActivate)

	if n := f.system.NumActiveBodies(BodyTypeRigid); n != 0 {
		t.Fatalf("num active bodies = %d, want 0", n)
	}

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodyGetPosition(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{1, 2, 3}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	p := bi.Position(bodyID)
	equalFloat(t, "pos.x", p.X, 1)
	equalFloat(t, "pos.y", p.Y, 2)
	equalFloat(t, "pos.z", p.Z, 3)

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodySetPosition(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 0, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	bi.SetPosition(bodyID, Vec3{5, 10, 15}, ActivationActivate)
	p := bi.Position(bodyID)
	equalFloat(t, "pos.x", p.X, 5)
	equalFloat(t, "pos.y", p.Y, 10)
	equalFloat(t, "pos.z", p.Z, 15)

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodySetLinearVelocity(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	bi.SetLinearVelocity(bodyID, Vec3{1, 2, 3})
	v := bi.LinearVelocity(bodyID)
	equalFloat(t, "vel.x", v.X, 1)
	equalFloat(t, "vel.y", v.Y, 2)
	equalFloat(t, "vel.z", v.Z, 3)

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodyAddImpulse(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	initial := bi.LinearVelocity(bodyID)
	equalFloat(t, "initial vel.x", initial.X, 0)

	bi.AddImpulse(bodyID, Vec3{10, 0, 0})
	if v := bi.LinearVelocity(bodyID); v.X <= 0 {
		t.Fatalf("vel.x = %v, want > 0", v.X)
	}

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodyGetMotionType(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	if m := bi.MotionType(bodyID); m != MotionTypeDynamic {
		t.Fatalf("motion type = %d, want %d", m, MotionTypeDynamic)
	}

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodyActivateDeactivate(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()

	bi.AddBody(bodyID, ActivationDontActivate)
	if bi.IsActive(bodyID) {
		t.Fatal("body should be inactive")
	}

	bi.ActivateBody(bodyID)
	if !bi.IsActive(bodyID) {
		t.Fatal("body should be active")
	}

	bi.DeactivateBody(bodyID)
	if bi.IsActive(bodyID) {
		t.Fatal("body should be inactive")
	}

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodyRestitution(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	bi.SetRestitution(bodyID, 0.8)
	equalFloat(t, "restitution", bi.Restitution(bodyID), 0.8)

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemBodyFriction(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	bi.SetFriction(bodyID, 0.5)
	equalFloat(t, "friction", bi.Friction(bodyID), 0.5)

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemUpdateFallingBody(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(1)
	defer shape.Destroy()

	settings := NewBodyCreationSettings(shape, Vec3{0, 100, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	bodyID := body.ID()
	bi.AddBody(bodyID, ActivationActivate)

	f.system.SetGravity(Vec3{0, -10, 0})

	initial := bi.Position(bodyID)

	f.system.OptimizeBroadPhase()
	if err := f.system.Update(1.0/60.0, 1, f.jobSystem); err != PhysicsUpdateErrorNone {
		t.Fatalf("update error = %d", err)
	}

	if p := bi.Position(bodyID); p.Y >= initial.Y {
		t.Fatalf("body did not fall: y = %v, initial %v", p.Y, initial.Y)
	}

	bi.RemoveBody(bodyID)
	bi.DestroyBody(bodyID)
}

func TestPhysicsSystemMultipleBodies(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	const numBodies = 10
	bodyIDs := make([]BodyID, numBodies)

	for i := range numBodies {
		shape := NewSphereShape(1)
		settings := NewBodyCreationSettings(shape, Vec3{float32(i) * 3, 10, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		body := bi.CreateBody(settings)
		bodyIDs[i] = body.ID()
		bi.AddBody(bodyIDs[i], ActivationActivate)
		settings.Destroy()
		shape.Destroy()
	}

	if n := f.system.NumBodies(); n != numBodies {
		t.Fatalf("num bodies = %d, want %d", n, numBodies)
	}
	if n := f.system.NumActiveBodies(BodyTypeRigid); n != numBodies {
		t.Fatalf("num active bodies = %d, want %d", n, numBodies)
	}

	for _, id := range bodyIDs {
		bi.RemoveBody(id)
		bi.DestroyBody(id)
	}
}
