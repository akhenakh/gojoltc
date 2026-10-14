package jolt

import "testing"

// The batched *Into accessors write into caller-owned slices; verify they agree
// with the per-call accessors.

func TestBatchedBodyReads(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	defer shape.Destroy()

	for i := range 3 {
		settings := NewBodyCreationSettings(shape, RVec3{Real(i), 1, 2}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		id := bi.CreateAndAddBody(settings, ActivationActivate)
		settings.Destroy()
		t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })
	}

	ids := f.system.Bodies()
	if len(ids) != 3 {
		t.Fatalf("bodies = %d, want 3", len(ids))
	}

	positions := make([]RVec3, len(ids))
	if n := bi.PositionsInto(ids, positions); n != len(ids) {
		t.Fatalf("PositionsInto wrote %d, want %d", n, len(ids))
	}
	rotations := make([]Quat, len(ids))
	bi.RotationsInto(ids, rotations)
	velocities := make([]Vec3, len(ids))
	bi.LinearVelocitiesInto(ids, velocities)

	for i, id := range ids {
		want := bi.Position(id)
		if positions[i] != want {
			t.Errorf("body %d position = %v, want %v", i, positions[i], want)
		}
	}

	// Into variants never write past the destination length.
	short := make([]BodyID, 2)
	if n := f.system.BodyIDsInto(short); n != 2 {
		t.Errorf("BodyIDsInto wrote %d, want 2", n)
	}
	if n := f.system.ActiveBodyIDsInto(BodyTypeRigid, short); n != 2 {
		t.Errorf("ActiveBodyIDsInto wrote %d, want 2", n)
	}
}
