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

// A configured swing-twist joint: the settings are plain struct fields with no
// C setters, so this is the proof they reach the created constraint.
func TestSwingTwistSettingsConfigureJoint(t *testing.T) {
	f, body1, body2 := newConstraintFixture(t)
	_ = f

	settings := NewSwingTwistConstraintSettings()
	defer settings.Destroy()
	settings.SetSpace(ConstraintSpaceLocalToBodyCOM)
	// The fixture's spheres sit at (0,0,0) and (2,0,0), so these two points (near
	// each body's +x face) land on the same world point. Points are measured from
	// each body's center of mass, not its origin.
	settings.SetPosition1(RVec3{0.5, 0, 0})
	settings.SetPosition2(RVec3{-1.5, 0, 0})
	settings.SetTwistAxis1(Vec3{1, 0, 0})
	settings.SetTwistAxis2(Vec3{1, 0, 0})
	settings.SetPlaneAxis1(Vec3{0, 1, 0})
	settings.SetPlaneAxis2(Vec3{0, 1, 0})
	settings.SetSwingType(SwingTypeCone)
	settings.SetNormalHalfConeAngle(0.8)
	settings.SetPlaneHalfConeAngle(0.8)
	settings.SetTwistMinAngle(-1.5)
	settings.SetTwistMaxAngle(1.5)

	c := settings.Create(body1, body2)
	defer c.Destroy()

	if c.SubType() != ConstraintSubTypeSwingTwist {
		t.Fatalf("subtype = %d, want swing twist", c.SubType())
	}
	if a := c.SwingTwistNormalHalfConeAngle(); math.Abs(float64(a)-0.8) > 1e-5 {
		t.Errorf("half cone angle = %v, want 0.8 (the setter did not reach the constraint)", a)
	}
}

// The rover's rod end: a ball joint that has to articulate, not jam (two hinges
// per end over-constrain the linkage). Each case pins down one setting
// behaviourally, which a round-trip cannot: the rod hangs from a pivot at its own
// tip, so a wide cone lets it drop 0.5*sin(cone) while a narrow one holds it
// level.
func TestSwingTwistRodEndArticulates(t *testing.T) {
	type motion struct {
		drop, angle, twist, yaw float64
		maxTwist, maxYaw        float64
	}

	// rig hangs a horizontal rod (center at world (0.6,5,0)) from a pivot at its
	// near tip, world (0.1,5,0), and reports how it moved.
	rig := func(t *testing.T, swingType SwingType, normalCone, planeCone, twistLimit float32, spin Vec3) motion {
		t.Helper()
		f := newPhysicsFixture(t)
		bi := f.system.BodyInterface()

		// The anchor sits above the pivot so the rod swings clear of it: a pivot
		// inside the anchor's box blocks the swing at whatever angle the contact
		// allows.
		anchorShape := NewBoxShape(Vec3{0.1, 0.1, 0.1}, DefaultConvexRadius)
		defer anchorShape.Destroy()
		as := NewBodyCreationSettings(anchorShape, RVec3{0.1, 5.6, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
		anchorID := bi.CreateAndAddBody(as, ActivationDontActivate)
		as.Destroy()

		rodShape := NewBoxShape(Vec3{0.5, 0.04, 0.04}, DefaultConvexRadius)
		defer rodShape.Destroy()
		rs := NewBodyCreationSettings(rodShape, RVec3{0.6, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		rs.SetAllowSleeping(false)
		rs.SetAngularVelocity(spin)
		rodID := bi.CreateAndAddBody(rs, ActivationActivate)
		rs.Destroy()
		t.Cleanup(func() {
			bi.RemoveAndDestroyBody(rodID)
			bi.RemoveAndDestroyBody(anchorID)
		})

		s := NewSwingTwistConstraintSettings()
		defer s.Destroy()
		s.SetSpace(ConstraintSpaceLocalToBodyCOM)
		s.SetPosition1(RVec3{0, -0.6, 0}) // the anchor's pivot, world (0.1, 5, 0)
		s.SetPosition2(RVec3{-0.5, 0, 0}) // the rod's near tip, the same point
		s.SetTwistAxis1(Vec3{1, 0, 0})    // twist about the rod's length
		s.SetTwistAxis2(Vec3{1, 0, 0})
		s.SetPlaneAxis1(Vec3{0, 1, 0})
		s.SetPlaneAxis2(Vec3{0, 1, 0})
		s.SetSwingType(swingType)
		s.SetNormalHalfConeAngle(normalCone)
		s.SetPlaneHalfConeAngle(planeCone)
		s.SetTwistMinAngle(-twistLimit)
		s.SetTwistMaxAngle(twistLimit)

		c := s.Create(f.system.BodyPtr(anchorID), f.system.BodyPtr(rodID))
		addConstraint(t, f, c)

		// A joint with no motor and no friction swings between its limits, so track
		// the peaks rather than one instant.
		var maxTwist, maxYaw float64
		f.system.OptimizeBroadPhase()
		for range 150 {
			f.system.Update(1.0/60.0, 1, f.jobSystem)
			q := bi.Rotation(rodID)
			maxTwist = math.Max(maxTwist, math.Abs(2*math.Atan2(float64(q.X), float64(q.W))))
			maxYaw = math.Max(maxYaw, math.Abs(2*math.Atan2(float64(q.Y), float64(q.W))))
		}

		p, q := bi.Position(rodID), bi.Rotation(rodID)
		if math.IsNaN(float64(p.Y)) {
			t.Fatal("the joint diverged")
		}
		// The pivot must hold: a wrongly anchored joint yanks the bodies.
		px, py, pz := worldAnchor(t, bi, rodID, Vec3{-0.5, 0, 0})
		if d := math.Sqrt((px-0.1)*(px-0.1) + (py-5)*(py-5) + pz*pz); d > 0.03 {
			t.Errorf("the joint's pivot moved %v m from (0.1,5,0)", d)
		}
		return motion{
			drop:     5 - float64(p.Y),
			angle:    2 * math.Atan2(math.Hypot(float64(q.Y), float64(q.Z)), float64(q.W)),
			twist:    2 * math.Atan2(float64(q.X), float64(q.W)),
			yaw:      2 * math.Atan2(float64(q.Y), float64(q.W)),
			maxTwist: maxTwist,
			maxYaw:   maxYaw,
		}
	}

	// The plane cone bounds the pitch: wide, the rod drops 0.5*sin(0.8) = 0.359 m;
	// narrow, it hangs level even with a wide normal cone — so the two angles are
	// not interchangeable.
	wide := rig(t, SwingTypeCone, 0.05, 0.8, 1.5, Vec3{})
	if wide.drop < 0.28 || wide.drop > 0.45 {
		t.Errorf("wide plane cone: center dropped %v m, want ~0.36 (swings, then stops at the cone)", wide.drop)
	}
	if math.Abs(wide.angle-0.8) > 0.06 {
		t.Errorf("wide plane cone: rod rotated %v rad, want ~0.8", wide.angle)
	}
	if locked := rig(t, SwingTypeCone, 0.8, 0.05, 1.5, Vec3{}); locked.drop > 0.08 {
		t.Errorf("narrow plane cone: center dropped %v m, want ~0.02 — the plane angle did not bound the swing", locked.drop)
	}

	// The normal cone bounds the yaw: a spin about world Y turns the rod until that
	// angle stops it. The plane cone is nearly locked, so this measures the yaw.
	turned := rig(t, SwingTypePyramid, 0.3, 0.02, 1.5, Vec3{0, 6, 0})
	if math.Abs(turned.maxYaw-0.3) > 0.06 {
		t.Errorf("normal cone 0.3: rod yawed up to %v rad, want ~0.3", turned.maxYaw)
	}
	if held := rig(t, SwingTypePyramid, 0.05, 0.02, 1.5, Vec3{0, 6, 0}); held.maxYaw > 0.1 {
		t.Errorf("normal cone 0.05: rod yawed up to %v rad — the normal angle did not bound it", held.maxYaw)
	}

	// The twist range: with the swing nearly locked, a spin about the rod's own
	// axis stops at the 0.2 limit. The spin is modest on purpose — a limit is
	// enforced per step, so a fast joint overshoots it by up to one step's
	// rotation (6 rad/s at 60 Hz would pass it by 0.1 rad before the impulse lands).
	spun := rig(t, SwingTypeCone, 0.02, 0.02, 0.2, Vec3{2, 0, 0})
	if math.Abs(spun.maxTwist-0.2) > 0.05 {
		t.Errorf("twist peaked at %v rad, want it stopped at the 0.2 limit", spun.maxTwist)
	}
}

// The six-DOF route to a rod end: three free rotations, three fixed
// translations, anchored at an offset point. Two things are worth pinning here —
// Jolt's defaults are free on *every* axis, so the translations have to be fixed
// explicitly or the joint slides, and the anchor is only read correctly once the
// space is body-local.
func TestSixDOFBallJointHoldsAndSwings(t *testing.T) {
	type motion struct {
		drop, twist, yaw float64
	}

	// rig hangs a horizontal rod (center at world (0.6,5,0)) from a pivot at its
	// near tip, world (0.1,5,0), with all rotations free and the given rotation
	// limit applied about world Z when limit > 0 (the direction gravity swings it).
	rig := func(t *testing.T, limit float32, spin Vec3) motion {
		t.Helper()
		f := newPhysicsFixture(t)
		bi := f.system.BodyInterface()

		// The anchor sits above the pivot so the rod swings clear of it.
		anchorShape := NewBoxShape(Vec3{0.1, 0.1, 0.1}, DefaultConvexRadius)
		defer anchorShape.Destroy()
		as := NewBodyCreationSettings(anchorShape, RVec3{0.1, 5.6, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
		anchorID := bi.CreateAndAddBody(as, ActivationDontActivate)
		as.Destroy()

		rodShape := NewBoxShape(Vec3{0.5, 0.04, 0.04}, DefaultConvexRadius)
		defer rodShape.Destroy()
		rs := NewBodyCreationSettings(rodShape, RVec3{0.6, 5, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		rs.SetAllowSleeping(false)
		rs.SetAngularVelocity(spin)
		rodID := bi.CreateAndAddBody(rs, ActivationActivate)
		rs.Destroy()
		t.Cleanup(func() {
			bi.RemoveAndDestroyBody(rodID)
			bi.RemoveAndDestroyBody(anchorID)
		})

		s := NewSixDOFConstraintSettings()
		defer s.Destroy()
		s.SetSpace(ConstraintSpaceLocalToBodyCOM)
		s.SetPosition1(RVec3{0, -0.6, 0}) // the anchor's pivot, world (0.1, 5, 0)
		s.SetPosition2(RVec3{-0.5, 0, 0}) // the rod's near tip, the same point
		// Everything is free by default, so a ball joint is made by fixing the
		// three translations.
		for _, axis := range []SixDOFConstraintAxis{
			SixDOFConstraintAxisTranslationX,
			SixDOFConstraintAxisTranslationY,
			SixDOFConstraintAxisTranslationZ,
		} {
			s.MakeFixedAxis(axis)
		}
		if limit > 0 {
			s.SetLimitedAxis(SixDOFConstraintAxisRotationZ, -limit, limit)
		}
		c := s.Create(f.system.BodyPtr(anchorID), f.system.BodyPtr(rodID))
		addConstraint(t, f, c)

		// A free joint swings through its arc, and with a spin it keeps oscillating,
		// so track the peaks rather than one instant.
		var maxTwist, maxYaw, maxDrop float64
		f.system.OptimizeBroadPhase()
		for range 150 {
			f.system.Update(1.0/60.0, 1, f.jobSystem)
			q := bi.Rotation(rodID)
			maxTwist = math.Max(maxTwist, math.Abs(2*math.Atan2(float64(q.X), float64(q.W))))
			maxYaw = math.Max(maxYaw, math.Abs(2*math.Atan2(float64(q.Y), float64(q.W))))
			maxDrop = math.Max(maxDrop, 5-float64(bi.Position(rodID).Y))
		}

		p := bi.Position(rodID)
		if math.IsNaN(float64(p.Y)) {
			t.Fatal("the joint diverged")
		}
		// The translations are held: the pivot has not moved, which also fails
		// loudly if the positions were read in world space instead of body-local.
		px, py, pz := worldAnchor(t, bi, rodID, Vec3{-0.5, 0, 0})
		if d := math.Sqrt((px-0.1)*(px-0.1) + (py-5)*(py-5) + pz*pz); d > 0.03 {
			t.Errorf("the joint's pivot moved %v m from (0.1,5,0)", d)
		}
		return motion{drop: maxDrop, twist: maxTwist, yaw: maxYaw}
	}

	// A free ball joint: the rod hangs the whole way down (its center falls its
	// full 0.5 m offset), and the other two rotations are free too, so a spin about
	// both the rod's axis and the vertical keeps turning.
	free := rig(t, 0, Vec3{2, 2, 0})
	if free.drop < 0.45 || free.drop > 0.55 {
		t.Errorf("free joint: center dropped %v m, want ~0.5 (it hangs straight down)", free.drop)
	}
	if free.twist < 0.5 {
		t.Errorf("free joint: the rod barely twisted (%v rad); RotationX is not free", free.twist)
	}
	if free.yaw < 0.5 {
		t.Errorf("free joint: the rod barely yawed (%v rad); RotationY is not free", free.yaw)
	}

	// A limited rod end: the same rig with RotationZ bounded stops at the limit,
	// dropping 0.5*sin(0.4) = 0.195 m — the swing stop a real rod end needs.
	limited := rig(t, 0.4, Vec3{})
	if limited.drop < 0.14 || limited.drop > 0.26 {
		t.Errorf("limited joint: center dropped %v m, want ~0.195 (stopped at the 0.4 rad swing limit)", limited.drop)
	}
}
