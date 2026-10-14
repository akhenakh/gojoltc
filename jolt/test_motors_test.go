package jolt

import (
	"math"
	"testing"
)

// addConstraint attaches c and registers cleanup in the only correct order:
// remove from the system first, then destroy.
func addConstraint(t *testing.T, f *physicsFixture, c *Constraint) {
	t.Helper()
	f.system.AddConstraint(c)
	t.Cleanup(func() {
		f.system.RemoveConstraint(c)
		c.Destroy()
	})
}

// twoBodyRig (anchor static, arm dynamic) with both bodies added to the world.
func twoBodyRig(t *testing.T, armPos RVec3) (*physicsFixture, *BodyInterface, BodyID, *Body, *Body) {
	t.Helper()
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	anchorShape := NewBoxShape(Vec3{0.25, 0.25, 0.25}, DefaultConvexRadius)
	anchorSettings := NewBodyCreationSettings(anchorShape, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	anchorID := bi.CreateAndAddBody(anchorSettings, ActivationDontActivate)
	anchorSettings.Destroy()
	anchorShape.Destroy()

	armShape := NewBoxShape(Vec3{0.75, 0.1, 0.1}, DefaultConvexRadius)
	armSettings := NewBodyCreationSettings(armShape, armPos, QuatIdentity(), MotionTypeDynamic, layerMoving)
	armID := bi.CreateAndAddBody(armSettings, ActivationActivate)
	armSettings.Destroy()
	armShape.Destroy()

	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(armID)
		bi.RemoveAndDestroyBody(anchorID)
	})
	f.system.OptimizeBroadPhase()
	return f, bi, armID, f.system.BodyPtr(anchorID), f.system.BodyPtr(armID)
}

func TestHingeMotorDrivesToTargetAngle(t *testing.T) {
	f, _, _, anchor, arm := twoBodyRig(t, RVec3{2, 0, 0})

	settings := NewHingeConstraintSettings()
	defer settings.Destroy()
	c := settings.Create(anchor, arm)
	addConstraint(t, f, c)

	motor := DefaultMotorSettings()
	motor.MinTorqueLimit, motor.MaxTorqueLimit = -2000, 2000
	c.HingeSetMotorSettings(motor)
	c.HingeSetMotorState(MotorStatePosition)
	c.HingeSetTargetAngle(0.6)

	for range 180 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if got := c.HingeCurrentAngle(); math.Abs(float64(got)-0.6) > 0.2 {
		t.Errorf("hinge angle = %v, want ~0.6", got)
	}
	if got := c.HingeTargetAngle(); math.Abs(float64(got)-0.6) > 1e-5 {
		t.Errorf("target angle = %v, want 0.6", got)
	}
}

func TestSliderMotorDrivesToTargetPosition(t *testing.T) {
	f, _, _, anchor, arm := twoBodyRig(t, RVec3{0, 0, 0})

	settings := NewSliderConstraintSettings()
	defer settings.Destroy()
	settings.SetSliderAxis(Vec3{1, 0, 0})
	c := settings.Create(anchor, arm)
	addConstraint(t, f, c)

	motor := DefaultMotorSettings()
	motor.MinForceLimit, motor.MaxForceLimit = -5000, 5000
	c.SliderSetMotorSettings(motor)
	c.SliderSetMotorState(MotorStatePosition)
	c.SliderSetTargetPosition(1.5)

	for range 180 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if got := c.SliderCurrentPosition(); math.Abs(float64(got)-1.5) > 0.25 {
		t.Errorf("slider position = %v, want ~1.5", got)
	}
	if got := c.SliderMotorState(); got != MotorStatePosition {
		t.Errorf("motor state = %v, want position", got)
	}
}

func TestSixDOFMotorTranslates(t *testing.T) {
	f, bi, armID, anchor, arm := twoBodyRig(t, RVec3{0, 0, 0})

	settings := NewSixDOFConstraintSettings()
	defer settings.Destroy()
	settings.MakeFreeAxis(SixDOFConstraintAxisTranslationX)
	settings.MakeFixedAxis(SixDOFConstraintAxisTranslationY)
	settings.MakeFixedAxis(SixDOFConstraintAxisTranslationZ)
	settings.MakeFixedAxis(SixDOFConstraintAxisRotationX)
	settings.MakeFixedAxis(SixDOFConstraintAxisRotationY)
	settings.MakeFixedAxis(SixDOFConstraintAxisRotationZ)

	c := settings.Create(anchor, arm)
	addConstraint(t, f, c)

	if !c.SixDOFIsFreeAxis(SixDOFConstraintAxisTranslationX) {
		t.Fatal("translation X should be free")
	}
	if !c.SixDOFIsFixedAxis(SixDOFConstraintAxisTranslationY) {
		t.Fatal("translation Y should be fixed")
	}

	// A velocity motor with Jolt's default (unlimited) force drives the axis.
	c.SixDOFSetMotorState(SixDOFConstraintAxisTranslationX, MotorStateVelocity)
	c.SixDOFSetTargetVelocityCS(Vec3{1, 0, 0})

	for range 120 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}

	if x := bi.Position(armID).X; x < 0.2 {
		t.Errorf("6DOF velocity motor did not translate the arm: x = %v", x)
	}
}

func TestLinearCurve(t *testing.T) {
	initShapes(t)

	c := NewLinearCurve()
	defer c.Destroy()
	c.Reserve(2)
	c.AddPoint(1, 2)
	c.AddPoint(0, 0)
	c.Sort()

	if n := c.PointCount(); n != 2 {
		t.Fatalf("point count = %d, want 2", n)
	}
	equalFloat(t, "min x", c.MinX(), 0)
	equalFloat(t, "max x", c.MaxX(), 1)
	nearFloat(t, "value at 0.5", c.Value(0.5), 1)
	nearFloat(t, "value below range", c.Value(-1), 0)
	nearFloat(t, "value above range", c.Value(2), 2)

	if p := c.Point(1); p.X != 1 || p.Y != 2 {
		t.Errorf("point 1 = %v, want (1,2)", p)
	}
	pts := make([]Point, 4)
	if n := c.Points(pts); n != 2 {
		t.Errorf("Points wrote %d, want 2", n)
	}
}

func TestWheelFrictionCurves(t *testing.T) {
	initShapes(t)

	wheel := NewWheelSettingsWV()
	defer wheel.Destroy()

	curve := NewLinearCurve()
	curve.AddPoint(0, 1)
	curve.AddPoint(1, 0.4)
	curve.Sort()

	wheel.SetLongitudinalFrictionCurve(curve)
	wheel.SetLateralFrictionCurve(curve)
	curve.Destroy() // the setter copies, so this must be safe

	got := wheel.LongitudinalFrictionCurve()
	if n := got.PointCount(); n != 2 {
		t.Fatalf("longitudinal curve point count = %d, want 2", n)
	}
	nearFloat(t, "longitudinal value at 0", got.Value(0), 1)

	// The returned curve is owned by the wheel: Destroy must not free it.
	got.Destroy()
	if n := wheel.LongitudinalFrictionCurve().PointCount(); n != 2 {
		t.Errorf("wheel curve was freed by the borrowed handle: count = %d", n)
	}
	nearFloat(t, "lateral value at 0", wheel.LateralFrictionCurve().Value(0), 1)
}
