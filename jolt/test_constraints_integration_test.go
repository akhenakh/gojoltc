package jolt

import (
	"math"
	"testing"
)

// Behavioural test: a constraint only does anything once it is added to the
// system. Creating it (as the ported upstream tests do) proves nothing.

func TestConstraintActuallyConstrains(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	// Static anchor and a dynamic body offset to the side.
	anchorShape := NewSphereShape(0.2)
	anchorSettings := NewBodyCreationSettings(anchorShape, RVec3{0, 5, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	anchorID := bi.CreateAndAddBody(anchorSettings, ActivationDontActivate)
	anchorSettings.Destroy()
	anchorShape.Destroy()

	bodyShape := NewSphereShape(0.2)
	bodySettings := NewBodyCreationSettings(bodyShape, RVec3{2, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	bodyID := bi.CreateAndAddBody(bodySettings, ActivationActivate)
	bodySettings.Destroy()
	bodyShape.Destroy()

	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(bodyID)
		bi.RemoveAndDestroyBody(anchorID)
	})

	anchor := f.system.BodyPtr(anchorID)
	body := f.system.BodyPtr(bodyID)

	settings := NewFixedConstraintSettings()
	defer settings.Destroy()
	constraint := settings.Create(anchor, body)
	defer constraint.Destroy()

	// Without adding it, gravity would drop the body; with it, it hangs.
	f.system.AddConstraint(constraint)
	if n := f.system.NumConstraints(); n != 1 {
		t.Fatalf("num constraints = %d, want 1", n)
	}

	f.system.OptimizeBroadPhase()
	for range 90 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	pos := body.Position()
	if math.Abs(float64(pos.Y)-5) > 0.5 {
		t.Errorf("constrained body fell: y = %v, want ~5", pos.Y)
	}
	if math.Abs(float64(pos.X)-2) > 0.5 {
		t.Errorf("constrained body drifted: x = %v, want ~2", pos.X)
	}

	f.system.RemoveConstraint(constraint)
	if n := f.system.NumConstraints(); n != 0 {
		t.Errorf("num constraints after removal = %d, want 0", n)
	}
}

func TestConstraintsPluralAddRemove(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.2)
	defer shape.Destroy()

	// A chain: anchor -> b1 -> b2, held by two fixed constraints added at once.
	ids := make([]BodyID, 3)
	for i := range 3 {
		pos := RVec3{float32(i) * 1.5, 5, 0}
		motion := MotionTypeDynamic
		if i == 0 {
			motion = MotionTypeStatic
		}
		settings := NewBodyCreationSettings(shape, pos, QuatIdentity(), motion, layerMoving)
		ids[i] = bi.CreateAndAddBody(settings, ActivationActivate)
		settings.Destroy()
	}
	t.Cleanup(func() {
		for _, id := range ids {
			bi.RemoveAndDestroyBody(id)
		}
	})

	s1 := NewFixedConstraintSettings()
	defer s1.Destroy()
	s2 := NewFixedConstraintSettings()
	defer s2.Destroy()

	c1 := s1.Create(f.system.BodyPtr(ids[0]), f.system.BodyPtr(ids[1]))
	defer c1.Destroy()
	c2 := s2.Create(f.system.BodyPtr(ids[1]), f.system.BodyPtr(ids[2]))
	defer c2.Destroy()

	f.system.AddConstraints([]*Constraint{c1, c2})
	if n := f.system.NumConstraints(); n != 2 {
		t.Fatalf("num constraints = %d, want 2", n)
	}

	f.system.OptimizeBroadPhase()
	for range 60 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	// All three bodies should still hang together near y=5.
	for i, id := range ids {
		if y := bi.Position(id).Y; math.Abs(float64(y)-5) > 1.0 {
			t.Errorf("body %d fell: y = %v", i, y)
		}
	}

	f.system.RemoveConstraints([]*Constraint{c1, c2})
	if n := f.system.NumConstraints(); n != 0 {
		t.Errorf("num constraints after removal = %d, want 0", n)
	}
}
