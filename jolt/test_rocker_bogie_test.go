package jolt

import (
	"math"
	"testing"
)

// A rocker-bogie chassis driven over gentle uneven ground: the mechanism the
// wrapper's constraint work was built for. It exercises hinge frames with travel
// limits, mass overrides with shape-derived inertia, and six driven wheels in
// ground contact — the combination a real rover depends on.
//
// Layout, in rover-local coordinates (x forward, y up, z out to the side): a
// rocker beam pivots on the chassis at its centre and is a symmetric seesaw —
// the two rear wheels hang off its -x arms, the two bogies off its +x arms, all
// the same distance below the pivot. Each bogie likewise hinges at its centre,
// with the middle and front wheels on its arms. That is the classic
// rocker-bogie arrangement: rear wheels on the rocker, the other four on the
// bogies.
//
// The beam spans both sides, i.e. the differential is locked (a spool). A rover
// that needs the two rockers to swing independently needs a real differential,
// and Jolt's gear constraint cannot express one hung on a moving chassis — see
// TestGearCouplesTwoHingesAsDifferentialPivot for where a gear *is* correct.
const (
	rbWheelRadius = 0.12
	rbPivotY      = 0.18 // height of the chassis/rocker pivot line above the axles
	rbHalfZ       = 0.20 // lateral offset of the link row
	rbWheelZ      = 0.26 // outboard of the links, so the wheels clear their boxes
	rbAxleX       = 0.00
	rbRearX       = -0.35
	rbBogieX      = 0.35
	rbMiddleX     = 0.10
	rbFrontX      = 0.60
)

// rbTerrainAt is a gentle undulation: enough to make the bogies articulate,
// shallow enough (±0.11 m, ~6° slopes) that a rigid one-metre linkage always has
// ground under every wheel.
func rbTerrainAt(x, z float64) float32 {
	return float32(0.06*math.Sin(x*1.7+0.4) + 0.05*math.Cos(z*1.3))
}

func TestRockerBogieDrivesOverUnevenGround(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	// A rigid rocker-bogie is a closed kinematic chain resting on six contacts,
	// which the default 10 velocity / 2 position solver passes cannot satisfy
	// (the contacts win and drag the linkage over).
	ps := f.system.PhysicsSettings()
	ps.NumVelocitySteps, ps.NumPositionSteps = 24, 8
	f.system.SetPhysicsSettings(ps)

	// --- terrain -------------------------------------------------------------
	const samples, size = 65, 16.0
	step := float32(size) / float32(samples-1)
	heights := make([]float32, samples*samples)
	for j := range samples {
		for i := range samples {
			x := (float64(i) - samples/2) * float64(step)
			z := (float64(j) - samples/2) * float64(step)
			heights[j*samples+i] = rbTerrainAt(x, z)
		}
	}
	terrain := NewHeightFieldShape(heights, Vec3{-size / 2, 0, -size / 2}, Vec3{step, 1, step})
	if terrain == nil {
		t.Fatal("terrain shape is nil")
	}
	defer terrain.Destroy()

	terrainSettings := NewBodyCreationSettings(terrain, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	terrainSettings.SetFriction(1)
	terrainID := bi.CreateAndAddBody(terrainSettings, ActivationDontActivate)
	terrainSettings.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(terrainID) })

	// --- the mechanism: 10 rigid parts ---------------------------------------
	// A rocker-bogie is one mechanism and Jolt's constraints have no "collide
	// connected" switch, so the parts share a group filter table with every pair
	// disabled (same GroupID, distinct SubGroupIDs).
	const numParts = 10
	parts := NewGroupFilterTable(numParts)
	defer parts.Destroy()
	for i := range numParts {
		for j := i + 1; j < numParts; j++ {
			parts.DisableCollision(CollisionSubGroupID(i), CollisionSubGroupID(j))
		}
	}
	group := 0
	nextGroup := func() CollisionSubGroupID {
		g := CollisionSubGroupID(group)
		group++
		return g
	}

	// Spawn just above the highest ground under the rover's outline.
	spawnY := float64(rbTerrainAt(0, 0))
	for _, x := range []float64{rbRearX - 0.1, rbFrontX + 0.1} {
		for _, z := range []float64{-rbWheelZ - 0.1, rbWheelZ + 0.1} {
			spawnY = math.Max(spawnY, float64(rbTerrainAt(x, z)))
		}
	}
	spawnY += rbWheelRadius + 0.015
	spawnAt := func(local Vec3) RVec3 {
		return RVec3{Real(local.X), Real(spawnY + float64(local.Y)), Real(local.Z)}
	}

	type link struct {
		pivot Vec3
		body  *Body
		id    BodyID
	}
	links := map[string]*link{}
	addLink := func(name string, half, pivot Vec3, mass float32) {
		shape := NewBoxShape(half, DefaultConvexRadius)
		defer shape.Destroy()
		s := NewBodyCreationSettings(shape, spawnAt(pivot), QuatIdentity(), MotionTypeDynamic, layerMoving)
		s.SetFriction(0.6)
		s.SetAllowSleeping(false)
		s.SetCollisionGroup(CollisionGroup{Filter: parts, GroupID: 1, SubGroupID: nextGroup()})
		s.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
		s.SetMassPropertiesOverride(NewMassProperties(mass))
		id := bi.CreateAndAddBody(s, ActivationActivate)
		s.Destroy()
		t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })
		links[name] = &link{pivot: pivot, body: f.system.BodyPtr(id), id: id}
	}

	// The rear wheels are heavy on purpose: they are the load on the beam's short
	// arms while the bogie ends carry three bodies each, and a rocker balances its
	// arms through its pivot.
	addLink("chassis", Vec3{0.22, 0.05, 0.16}, Vec3{0, 0.24, 0}, 4.0)
	addLink("rocker", Vec3{0.35, 0.03, rbHalfZ}, Vec3{rbAxleX, rbPivotY, 0}, 2.4)
	addLink("bogieL", Vec3{0.25, 0.03, 0.03}, Vec3{rbBogieX, 0, -rbHalfZ}, 0.4)
	addLink("bogieR", Vec3{0.25, 0.03, 0.03}, Vec3{rbBogieX, 0, rbHalfZ}, 0.4)

	// addJoint attaches a hinge whose point is shared, expressed in each body's
	// local space. Axes are per body because a wheel body is rotated.
	addJoint := func(a, b *link, aLocal, bLocal Vec3, aAxis, bAxis Vec3, limited bool, min, max float32) *Constraint {
		s := NewHingeConstraintSettings()
		defer s.Destroy()
		s.SetSpace(ConstraintSpaceLocalToBodyCOM)
		s.SetLocalSpacePoint1(RVec3{Real(aLocal.X), Real(aLocal.Y), Real(aLocal.Z)})
		s.SetLocalSpacePoint2(RVec3{Real(bLocal.X), Real(bLocal.Y), Real(bLocal.Z)})
		s.SetHingeAxis1(aAxis)
		s.SetHingeAxis2(bAxis)
		s.SetNormalAxis1(Vec3{1, 0, 0})
		s.SetNormalAxis2(Vec3{1, 0, 0})
		if limited {
			// Leaving the limits unset keeps the settings' default +/-pi range,
			// which Jolt reads as "no limits"; any narrower range — including a
			// single point like 0,0 — activates them and locks the joint.
			s.SetInitialLimits(min, max)
		}
		c := s.Create(a.body, b.body)
		f.system.AddConstraint(c)
		t.Cleanup(func() {
			f.system.RemoveConstraint(c)
			c.Destroy()
		})
		return c
	}

	lateral := Vec3{0, 0, 1}
	chassis, rocker := links["chassis"], links["rocker"]
	rockerHinge := addJoint(chassis, rocker,
		Vec3{rocker.pivot.X - chassis.pivot.X, rocker.pivot.Y - chassis.pivot.Y, rocker.pivot.Z - chassis.pivot.Z},
		Vec3{0, 0, 0}, lateral, lateral, true, -0.25, 0.25)

	// The bogies hang off the beam's far arms, which drop to the axle line just
	// like the rear wheels on the other side, so both hinge points land on the
	// bogie's own body line. (Frames that disagree at creation make the solver
	// fight a permanent error, which flings the whole mechanism about.)
	bogieL := addJoint(rocker, links["bogieL"], Vec3{rbBogieX, -rbPivotY, -rbHalfZ}, Vec3{0, 0, 0}, lateral, lateral, true, -0.35, 0.35)
	bogieR := addJoint(rocker, links["bogieR"], Vec3{rbBogieX, -rbPivotY, rbHalfZ}, Vec3{0, 0, 0}, lateral, lateral, true, -0.35, 0.35)

	// --- six cylinder wheels, one drive hinge each ---------------------------
	// Jolt's cylinder is Y-aligned, so each wheel body is rotated a quarter turn
	// about X to put its axis on Z; its own local hinge axis is therefore Y.
	wheelShape := NewCylinderShape(0.03, rbWheelRadius)
	defer wheelShape.Destroy()
	wheelRot := QuatRotation(Vec3{1, 0, 0}, Pi/2)

	wheels := make([]*Constraint, 0, 6)
	addWheel := func(carrier string, local Vec3, mass float32) {
		c := links[carrier]
		s := NewBodyCreationSettings(wheelShape, spawnAt(local), wheelRot, MotionTypeDynamic, layerMoving)
		s.SetFriction(1)
		s.SetAllowSleeping(false)
		s.SetCollisionGroup(CollisionGroup{Filter: parts, GroupID: 1, SubGroupID: nextGroup()})
		s.SetOverrideMassProperties(OverrideMassPropertiesCalculateInertia)
		s.SetMassPropertiesOverride(NewMassProperties(mass))
		id := bi.CreateAndAddBody(s, ActivationActivate)
		s.Destroy()
		t.Cleanup(func() { bi.RemoveAndDestroyBody(id) })

		w := &link{pivot: local, body: f.system.BodyPtr(id), id: id}
		aLocal := Vec3{local.X - c.pivot.X, local.Y - c.pivot.Y, local.Z - c.pivot.Z}
		wheels = append(wheels, addJoint(c, w, aLocal, Vec3{0, 0, 0}, lateral, Vec3{0, 1, 0}, false, 0, 0))
	}
	addWheel("rocker", Vec3{rbRearX, 0, -rbWheelZ}, 1.0)
	addWheel("rocker", Vec3{rbRearX, 0, rbWheelZ}, 1.0)
	addWheel("bogieL", Vec3{rbMiddleX, 0, -rbWheelZ}, 0.25)
	addWheel("bogieR", Vec3{rbMiddleX, 0, rbWheelZ}, 0.25)
	addWheel("bogieL", Vec3{rbFrontX, 0, -rbWheelZ}, 0.25)
	addWheel("bogieR", Vec3{rbFrontX, 0, rbWheelZ}, 0.25)

	f.system.OptimizeBroadPhase()

	// Let it settle onto the terrain.
	for range 90 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	settled := bi.Position(chassis.id)
	if settled.Y < 0 {
		t.Fatalf("rover fell through the terrain: y = %v", settled.Y)
	}

	// --- drive: every wheel motored, the linkage passive ---------------------
	motor := DefaultMotorSettings()
	motor.MinTorqueLimit, motor.MaxTorqueLimit = -0.6, 0.6
	for _, h := range wheels {
		h.HingeSetMotorSettings(motor)
		h.HingeSetMotorState(MotorStateVelocity)
		h.HingeSetTargetAngularVelocity(-2) // spins the wheels forward (+x)
	}

	maxRockerAngle, maxBogieAngle, maxWheelAngle := 0.0, 0.0, 0.0
	for range 300 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
		maxRockerAngle = math.Max(maxRockerAngle, math.Abs(float64(rockerHinge.HingeCurrentAngle())))
		for _, h := range []*Constraint{bogieL, bogieR} {
			maxBogieAngle = math.Max(maxBogieAngle, math.Abs(float64(h.HingeCurrentAngle())))
		}
		for _, h := range wheels {
			maxWheelAngle = math.Max(maxWheelAngle, math.Abs(float64(h.HingeCurrentAngle())))
		}
	}

	end := bi.Position(chassis.id)
	travelled := math.Hypot(float64(end.X-settled.X), float64(end.Z-settled.Z))
	up := bi.Rotation(chassis.id).RotateAxisY()

	t.Logf("travelled %.2f m; up.y %.2f; max wheel %.1f rad; max rocker %.3f rad; max bogie %.3f rad",
		travelled, up.Y, maxWheelAngle, maxRockerAngle, maxBogieAngle)

	if travelled < 0.6 {
		t.Errorf("rover did not drive: moved %.3f m", travelled)
	}
	if up.Y < 0.8 {
		t.Errorf("rover tipped over: up.y = %v", up.Y)
	}
	if float64(end.Y) < float64(rbTerrainAt(float64(end.X), float64(end.Z)))-0.5 {
		t.Errorf("rover sank through the terrain: y = %v", end.Y)
	}
	if maxWheelAngle < 2 {
		t.Errorf("wheels barely turned: max %.2f rad", maxWheelAngle)
	}
	if maxBogieAngle < 0.01 {
		t.Errorf("terrain never articulated the linkage: max bogie angle %.4f rad", maxBogieAngle)
	}
	// The travel stops hold under load: the beam is limited to +/-0.25 rad, and it
	// is the load-bearing constraint here (mass 4 kg chassis + 6 driven wheels).
	if maxRockerAngle > 0.30 {
		t.Errorf("rocker travel stop did not hold: reached %.4f rad (limit 0.25)", maxRockerAngle)
	}
}
