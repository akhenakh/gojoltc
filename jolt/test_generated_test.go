package jolt

import "testing"

// Exercises accessors emitted by the generator (internal/joltc/gen_accessors.go).

func TestGeneratedBodyCreationSettingsAccessors(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()

	s := NewBodyCreationSettings(shape, RVec3{0, 1, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer s.Destroy()

	s.SetLinearDamping(0.25)
	equalFloat(t, "linear damping", s.LinearDamping(), 0.25)
	s.SetAngularDamping(0.75)
	equalFloat(t, "angular damping", s.AngularDamping(), 0.75)
	s.SetGravityFactor(0.5)
	equalFloat(t, "gravity factor", s.GravityFactor(), 0.5)
	s.SetIsSensor(true)
	if !s.IsSensor() {
		t.Error("is sensor")
	}
	s.SetMotionQuality(MotionQualityLinearCast)
	if s.MotionQuality() != MotionQualityLinearCast {
		t.Error("motion quality")
	}
	s.SetAllowedDOFs(AllowedDOFsPlane2D)
	if s.AllowedDOFs() != AllowedDOFsPlane2D {
		t.Error("allowed DOFs")
	}
	s.SetMaxLinearVelocity(42)
	equalFloat(t, "max linear velocity", s.MaxLinearVelocity(), 42)
	s.SetNumVelocityStepsOverride(3)
	if s.NumVelocityStepsOverride() != 3 {
		t.Error("num velocity steps override")
	}
	if s.MotionType() != MotionTypeDynamic {
		t.Error("motion type")
	}
	if s.ObjectLayer() != layerMoving {
		t.Error("object layer")
	}
}

func TestGeneratedBodyAccessors(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	defer shape.Destroy()
	settings := NewBodyCreationSettings(shape, RVec3{1, 2, 3}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	defer settings.Destroy()

	body := bi.CreateBody(settings)
	id := body.ID()
	bi.AddBody(id, ActivationActivate)
	t.Cleanup(func() {
		bi.RemoveBody(id)
		bi.DestroyBody(id)
	})

	if body.BodyType() != BodyTypeRigid {
		t.Error("body type")
	}
	if body.MotionType() != MotionTypeDynamic {
		t.Error("motion type")
	}
	if !body.IsDynamic() {
		t.Error("is dynamic")
	}
	if body.ObjectLayer() != layerMoving {
		t.Error("object layer")
	}
	if !body.IsActive() {
		t.Error("is active")
	}

	p := body.Position()
	equalReal(t, "position.x", p.X, 1)
	equalReal(t, "position.y", p.Y, 2)
	equalReal(t, "position.z", p.Z, 3)

	r := body.Rotation()
	nearFloat(t, "rotation.w", r.W, 1)

	body.SetFriction(0.3)
	equalFloat(t, "friction", body.Friction(), 0.3)
	body.SetUserData(42)
	if body.UserData() != 42 {
		t.Error("user data")
	}
}

func TestGeneratedPhysicsSystemAccessors(t *testing.T) {
	f := newPhysicsFixture(t)
	if n := f.system.NumConstraints(); n != 0 {
		t.Errorf("num constraints = %d, want 0", n)
	}
	if n := f.system.MaxBodies(); n != 1024 {
		t.Errorf("max bodies = %d, want 1024", n)
	}
}
