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

// worldAnchor returns the world position of a body-local point (relative to the
// center of mass), in float64 so the assertion reads the same in both precisions.
func worldAnchor(t *testing.T, bi *BodyInterface, id BodyID, local Vec3) (float64, float64, float64) {
	t.Helper()
	p := bi.Position(id)
	off := bi.Rotation(id).Rotate(local)
	return float64(p.X) + float64(off.X), float64(p.Y) + float64(off.Y), float64(p.Z) + float64(off.Z)
}

// The distance settings' anchors default to the world origin, so a rod only
// holds two bodies apart once its ends are set. This is the difference between a
// usable tie-rod (or differential-bar linkage) and one that welds both bodies to
// (0,0,0).
func TestDistanceConstraintAnchors(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewBoxShape(Vec3{0.5, 0.5, 0.5}, DefaultConvexRadius)
	defer shape.Destroy()
	ids := make([]BodyID, 2)
	for i, start := range []RVec3{{-1, 5, 0}, {1, 5, 0}} {
		s := NewBodyCreationSettings(shape, start, QuatIdentity(), MotionTypeDynamic, layerMoving)
		s.SetAllowSleeping(false)
		ids[i] = bi.CreateAndAddBody(s, ActivationActivate)
		s.Destroy()
	}
	t.Cleanup(func() {
		for _, id := range ids {
			bi.RemoveAndDestroyBody(id)
		}
	})
	f.system.OptimizeBroadPhase()

	// A rigid rod between the two inner faces: they are 1 m apart to start.
	settings := NewDistanceConstraintSettings()
	defer settings.Destroy()
	settings.SetSpace(ConstraintSpaceLocalToBodyCOM)
	settings.SetLocalSpacePoint1(RVec3{0.5, 0, 0})
	settings.SetLocalSpacePoint2(RVec3{-0.5, 0, 0})
	settings.SetInitialDistance(1, 1)
	c := settings.Create(f.system.BodyPtr(ids[0]), f.system.BodyPtr(ids[1]))
	addConstraint(t, f, c)

	for range 120 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	x0, y0, z0 := worldAnchor(t, bi, ids[0], Vec3{0.5, 0, 0})
	x1, y1, z1 := worldAnchor(t, bi, ids[1], Vec3{-0.5, 0, 0})
	if d := math.Sqrt((x1-x0)*(x1-x0) + (y1-y0)*(y1-y0) + (z1-z0)*(z1-z0)); math.Abs(d-1) > 0.05 {
		t.Errorf("anchor separation = %v, want 1", d)
	}
	// And they fell together rather than being dragged to the anchors' default
	// world origin.
	for i, id := range ids {
		if p := bi.Position(id); math.Abs(math.Hypot(float64(p.X), float64(p.Z))) < 0.5 {
			t.Errorf("body %d was pulled towards the origin: x = %v, z = %v", i, p.X, p.Z)
		}
	}
	if y0 > 4 {
		t.Errorf("the pair did not fall under gravity: y = %v", y0)
	}
}

// A point constraint is a ball joint: with the anchor at the arm's tip, the tip
// stays put while the arm swings about it.
func TestPointConstraintAnchors(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	anchorShape := NewBoxShape(Vec3{0.1, 0.1, 0.1}, DefaultConvexRadius)
	defer anchorShape.Destroy()
	as := NewBodyCreationSettings(anchorShape, RVec3{0, 5, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	anchorID := bi.CreateAndAddBody(as, ActivationDontActivate)
	as.Destroy()

	armShape := NewBoxShape(Vec3{0.75, 0.05, 0.05}, DefaultConvexRadius)
	defer armShape.Destroy()
	// Centered at x=1, so its -x tip is at x=0.25 and the anchor's face is at 0.1:
	// the joint pivot is chosen on the arm's tip and the anchor's center.
	arms := NewBodyCreationSettings(armShape, RVec3{1, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	arms.SetAllowSleeping(false)
	armID := bi.CreateAndAddBody(arms, ActivationActivate)
	arms.Destroy()
	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(armID)
		bi.RemoveAndDestroyBody(anchorID)
	})
	f.system.OptimizeBroadPhase()

	settings := NewPointConstraintSettings()
	defer settings.Destroy()
	settings.SetSpace(ConstraintSpaceLocalToBodyCOM)
	settings.SetLocalSpacePoint1(RVec3{0, 0, 0})  // the anchor's center, world (0, 5, 0)
	settings.SetLocalSpacePoint2(RVec3{-1, 0, 0}) // the arm's near tip, the same point
	c := settings.Create(f.system.BodyPtr(anchorID), f.system.BodyPtr(armID))
	addConstraint(t, f, c)

	for range 120 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	tipX, tipY, tipZ := worldAnchor(t, bi, armID, Vec3{-1, 0, 0})
	if d := math.Sqrt(tipX*tipX + (tipY-5)*(tipY-5) + tipZ*tipZ); d > 0.05 {
		t.Errorf("the joint's pivot moved %v m from (0,5,0)", d)
	}
	// The arm hung off the pivot instead of staying welded at the origin.
	centre := bi.Position(armID)
	if math.Hypot(float64(centre.X), float64(centre.Y-5)) < 0.5 {
		t.Errorf("arm did not swing about the anchor: centre = (%v,%v)", centre.X, centre.Y)
	}
}

// The chassis-mounted differential a rover needs — two rockers whose rotations
// stay equal and opposite — cannot be a gear (Jolt's gear couples absolute
// angular velocities, so it only works with its hinges anchored to the world; see
// TestGearCouplesTwoHingesAsDifferentialPivot). It can be a differential bar: a
// bar on a *longitudinal* hinge, tied to each rocker by a rigid rod, which turns
// the pair's average into the bar's rotation and forces the anti-symmetric
// coupling. This is what the new distance anchors are for.
func TestDifferentialBarCouplesRockers(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	// A static chassis frame: the two rockers hang off it at z = +/-0.25 and the
	// bar pivots on it above them.
	anchorShape := NewBoxShape(Vec3{0.1, 0.1, 0.1}, DefaultConvexRadius)
	defer anchorShape.Destroy()
	as := NewBodyCreationSettings(anchorShape, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	anchorID := bi.CreateAndAddBody(as, ActivationDontActivate)
	as.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(anchorID) })
	anchor := f.system.BodyPtr(anchorID)

	// Rockers: boxes centred on their own pivots (so gravity exerts no torque and
	// the only load is the rods), each hinged to the frame about the lateral axis.
	rockerShape := NewBoxShape(Vec3{0.35, 0.03, 0.03}, DefaultConvexRadius)
	defer rockerShape.Destroy()
	var rockers [2]*Body
	var rockerHinges [2]*Constraint
	for i, z := range []Real{-0.25, 0.25} {
		s := NewBodyCreationSettings(rockerShape, RVec3{0, 0, z}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		s.SetFriction(0.6)
		s.SetAllowSleeping(false)
		s.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
		s.SetMassPropertiesOverride(NewMassProperties(1.2))
		id := bi.CreateAndAddBody(s, ActivationActivate)
		s.Destroy()
		t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })
		rockers[i] = f.system.BodyPtr(id)

		hs := NewHingeConstraintSettings()
		hs.SetSpace(ConstraintSpaceLocalToBodyCOM)
		hs.SetLocalSpacePoint1(RVec3{0, 0, z}) // the frame's pivot, in the frame's frame
		hs.SetLocalSpacePoint2(RVec3{0, 0, 0})
		hs.SetHingeAxis1(Vec3{0, 0, 1}) // lateral
		hs.SetHingeAxis2(Vec3{0, 0, 1})
		hs.SetNormalAxis1(Vec3{1, 0, 0})
		hs.SetNormalAxis2(Vec3{1, 0, 0})
		rockerHinges[i] = hs.Create(anchor, rockers[i])
		hs.Destroy()
		addConstraint(t, f, rockerHinges[i])
	}

	// The bar: hinged to the frame about the *longitudinal* axis so its two arms
	// move in opposite vertical directions, with its arms reaching down over the
	// rockers' front ends.
	barShape := NewBoxShape(Vec3{0.03, 0.03, 0.25}, DefaultConvexRadius)
	defer barShape.Destroy()
	bs := NewBodyCreationSettings(barShape, RVec3{0.35, 0.5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	bs.SetAllowSleeping(false)
	bs.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
	bs.SetMassPropertiesOverride(NewMassProperties(0.4))
	barID := bi.CreateAndAddBody(bs, ActivationActivate)
	bs.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(barID) })
	bar := f.system.BodyPtr(barID)

	bhs := NewHingeConstraintSettings()
	bhs.SetSpace(ConstraintSpaceLocalToBodyCOM)
	bhs.SetLocalSpacePoint1(RVec3{0.35, 0.5, 0}) // the bar's own pivot, from the frame
	bhs.SetLocalSpacePoint2(RVec3{0, 0, 0})
	bhs.SetHingeAxis1(Vec3{1, 0, 0}) // longitudinal
	bhs.SetHingeAxis2(Vec3{1, 0, 0})
	bhs.SetNormalAxis1(Vec3{0, 1, 0}) // perpendicular to the hinge axis
	bhs.SetNormalAxis2(Vec3{0, 1, 0})
	barHinge := bhs.Create(anchor, bar)
	bhs.Destroy()
	addConstraint(t, f, barHinge)

	// The rods: rigid (min == max), vertical at rest, from each bar arm down to the
	// rocker's front end, so a rocker's rotation reaches the bar and back out.
	for i, z := range []Real{-0.25, 0.25} {
		ds := NewDistanceConstraintSettings()
		ds.SetSpace(ConstraintSpaceLocalToBodyCOM)
		ds.SetLocalSpacePoint1(RVec3{0, 0, z})    // bar arm
		ds.SetLocalSpacePoint2(RVec3{0.35, 0, 0}) // rocker's front end
		ds.SetInitialDistance(0.5, 0.5)           // the rods are 0.5 m at rest
		c := ds.Create(bar, rockers[i])
		ds.Destroy()
		addConstraint(t, f, c)
	}

	f.system.OptimizeBroadPhase()

	// Drive one rocker; the bar must pass it to the other, reversed.
	motor := DefaultMotorSettings()
	rockerHinges[0].HingeSetMotorSettings(motor)
	rockerHinges[0].HingeSetMotorState(MotorStatePosition)
	rockerHinges[0].HingeSetTargetAngle(0.35)

	for range 180 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	driven := rockerHinges[0].HingeCurrentAngle()
	followed := rockerHinges[1].HingeCurrentAngle()
	t.Logf("driven rocker %v rad, coupled rocker %v rad, bar %v rad", driven, followed, barHinge.HingeCurrentAngle())

	if math.Abs(float64(driven)-0.35) > 0.1 {
		t.Errorf("driven rocker = %v rad, want ~0.35", driven)
	}
	if math.Abs(float64(followed)+0.35) > 0.1 {
		t.Errorf("coupled rocker = %v rad, want ~-0.35 (the bar inverts it)", followed)
	}
	if sum := math.Abs(float64(driven + followed)); sum > 0.05 {
		t.Errorf("the differential bar did not hold the rockers mirrored: |a1 + a2| = %v", sum)
	}
}
