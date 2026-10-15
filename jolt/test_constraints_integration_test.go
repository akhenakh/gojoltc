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
		pos := RVec3{Real(i) * 1.5, 5, 0}
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

// A hinge is a plain C struct with public fields, so until these setters existed
// a hinge could not be aimed at all (the slider had an axis setter, the hinge had
// none). This pins down an explicitly constructed one: the angles are measured
// about the axis it was given, and the limits it was created with hold.
func TestHingeSettingsExplicitAxesAndLimits(t *testing.T) {
	f, _, _, anchor, arm := twoBodyRig(t, RVec3{1.25, 0, 0})

	settings := NewHingeConstraintSettings()
	defer settings.Destroy()
	settings.SetSpace(ConstraintSpaceLocalToBodyCOM)
	settings.SetLocalSpacePoint1(RVec3{0.5, 0, 0})   // anchor local: world (0.5, 0, 0)
	settings.SetLocalSpacePoint2(RVec3{-0.75, 0, 0}) // arm's inner end: the same world point
	settings.SetHingeAxis1(Vec3{0, 1, 0})
	settings.SetHingeAxis2(Vec3{0, 1, 0})
	settings.SetNormalAxis1(Vec3{1, 0, 0})
	settings.SetNormalAxis2(Vec3{1, 0, 0})
	settings.SetInitialLimits(-0.4, 0.4)

	c := settings.Create(anchor, arm)
	addConstraint(t, f, c)

	// The settings reached the constraint.
	nearFloat(t, "limits min", c.HingeLimitsMin(), -0.4)
	nearFloat(t, "limits max", c.HingeLimitsMax(), 0.4)
	if !c.HingeHasLimits() {
		t.Error("hinge was created without limits")
	}

	// Ask for more than the limit allows: the joint must stop at it.
	motor := DefaultMotorSettings()
	c.HingeSetMotorSettings(motor)
	c.HingeSetMotorState(MotorStatePosition)
	c.HingeSetTargetAngle(1.0)

	for range 240 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if got := c.HingeCurrentAngle(); math.Abs(math.Abs(float64(got))-0.4) > 0.06 {
		t.Errorf("hinge angle = %v, want it held at the limit (0.4)", got)
	}

	// It turned about the axis it was given (world Y through the pivot), so the arm
	// kept its height and its distance from the pivot, and swung sideways in XZ.
	pos := arm.Position()
	if math.Abs(float64(pos.Y)) > 0.02 {
		t.Errorf("arm y = %v, want it level: rotation must be about Y", pos.Y)
	}
	if r := math.Hypot(float64(pos.X-0.5), float64(pos.Z)); math.Abs(r-0.75) > 0.02 {
		t.Errorf("arm is %v from the pivot, want 0.75", r)
	}
	if math.Abs(float64(pos.Z)) < 0.1 {
		t.Errorf("arm z = %v: it did not swing about the configured axis", pos.Z)
	}
}

// A rocker-bogie differential pivot: two bodies hinged about the same axis and
// geared together, so driving one drives the other the opposite way. Jolt defines
// the ratio as rotation1 = -ratio * rotation2, so ratio 1 is equal and opposite;
// SetRatio's doc records that (-1 would drive them the *same* way).
func TestGearCouplesTwoHingesAsDifferentialPivot(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	rotorShape := NewBoxShape(Vec3{0.3, 0.1, 0.1}, DefaultConvexRadius)
	defer rotorShape.Destroy()
	anchorShape := NewBoxShape(Vec3{0.2, 0.2, 0.2}, DefaultConvexRadius)
	defer anchorShape.Destroy()

	// A static anchor well below the pivot so nothing touches it.
	anchorSettings := NewBodyCreationSettings(anchorShape, RVec3{0, -1, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	anchorID := bi.CreateAndAddBody(anchorSettings, ActivationDontActivate)
	anchorSettings.Destroy()

	rotorIDs := make([]BodyID, 2)
	for i, pos := range []RVec3{{0.5, 0, 0}, {-0.5, 0, 0}} {
		s := NewBodyCreationSettings(rotorShape, pos, QuatIdentity(), MotionTypeDynamic, layerMoving)
		rotorIDs[i] = bi.CreateAndAddBody(s, ActivationActivate)
		s.Destroy()
	}
	t.Cleanup(func() {
		for _, id := range rotorIDs {
			bi.RemoveAndDestroyBody(id)
		}
		bi.RemoveAndDestroyBody(anchorID)
	})

	anchor := f.system.BodyPtr(anchorID)
	rotors := []*Body{f.system.BodyPtr(rotorIDs[0]), f.system.BodyPtr(rotorIDs[1])}
	f.system.OptimizeBroadPhase()

	// One hinge per rotor, both about world Y through the origin.
	hinges := make([]*Constraint, 2)
	for i, localPivot := range []RVec3{{-0.5, 0, 0}, {0.5, 0, 0}} {
		s := NewHingeConstraintSettings()
		s.SetSpace(ConstraintSpaceLocalToBodyCOM)
		s.SetLocalSpacePoint1(RVec3{0, 1, 0}) // anchor COM is (0,-1,0); the pivot is the origin
		s.SetLocalSpacePoint2(localPivot)
		s.SetHingeAxis1(Vec3{0, 1, 0})
		s.SetHingeAxis2(Vec3{0, 1, 0})
		s.SetNormalAxis1(Vec3{1, 0, 0})
		s.SetNormalAxis2(Vec3{1, 0, 0})
		hinges[i] = s.Create(anchor, rotors[i])
		s.Destroy()
		addConstraint(t, f, hinges[i])
	}

	// The differential itself: the two rotors geared one to one (world-space axes,
	// as Jolt's GearConstraintTest does).
	gearSettings := NewGearConstraintSettings()
	defer gearSettings.Destroy()
	gearSettings.SetHingeAxis1(Vec3{0, 1, 0})
	gearSettings.SetHingeAxis2(Vec3{0, 1, 0})
	gearSettings.SetRatio(1)
	gear := gearSettings.Create(rotors[0], rotors[1])
	gear.GearSetConstraints(hinges[0], hinges[1])
	addConstraint(t, f, gear)

	// Drive one side; the gear must turn the other the other way.
	motor := DefaultMotorSettings()
	hinges[0].HingeSetMotorSettings(motor)
	hinges[0].HingeSetMotorState(MotorStatePosition)
	hinges[0].HingeSetTargetAngle(0.6)

	for range 240 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	driven, followed := hinges[0].HingeCurrentAngle(), hinges[1].HingeCurrentAngle()
	t.Logf("driven hinge %v rad, geared hinge %v rad", driven, followed)
	if math.Abs(float64(driven)-0.6) > 0.25 {
		t.Errorf("driven hinge angle = %v, want ~0.6", driven)
	}
	if math.Abs(float64(followed)+0.6) > 0.25 {
		t.Errorf("geared hinge angle = %v, want ~-0.6 (ratio 1 is equal and opposite)", followed)
	}
	if sum := math.Abs(float64(driven + followed)); sum > 0.1 {
		t.Errorf("the two hinges are not locked to ratio 1: |a1 + a2| = %v", sum)
	}
}
