package jolt

import (
	"math"
	"testing"
)

// Ported from third_party/joltc/tests/test_constraints.cpp.

func newConstraintFixture(t *testing.T) (f *physicsFixture, body1, body2 *Body) {
	t.Helper()
	f = newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shape := NewSphereShape(0.5)
	bs1 := NewBodyCreationSettings(shape, Vec3{0, 0, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	bs2 := NewBodyCreationSettings(shape, Vec3{2, 0, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)

	body1 = bi.CreateBody(bs1)
	body2 = bi.CreateBody(bs2)
	bi.AddBody(body1.ID(), ActivationActivate)
	bi.AddBody(body2.ID(), ActivationActivate)

	bs1.Destroy()
	bs2.Destroy()
	shape.Destroy()

	t.Cleanup(func() {
		bi.RemoveBody(body1.ID())
		bi.RemoveBody(body2.ID())
		bi.DestroyBody(body1.ID())
		bi.DestroyBody(body2.ID())
	})
	return f, body1, body2
}

func TestConstraintsFixedCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewFixedConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	if c == nil {
		t.Fatal("constraint is nil")
	}
	defer c.Destroy()

	if c.SubType() != ConstraintSubTypeFixed {
		t.Errorf("subtype = %d, want fixed", c.SubType())
	}
	equalFloat(t, "lambda.x", c.FixedTotalLambdaPosition().X, 0)
}

func TestConstraintsDistanceCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewDistanceConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeDistance {
		t.Errorf("subtype = %d, want distance", c.SubType())
	}
}

func TestConstraintsDistanceSetDistance(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewDistanceConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	c.DistanceSetDistance(1, 3)
	equalFloat(t, "min", c.DistanceMinDistance(), 1)
	equalFloat(t, "max", c.DistanceMaxDistance(), 3)
}

func TestConstraintsPointCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewPointConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypePoint {
		t.Errorf("subtype = %d, want point", c.SubType())
	}
	_ = c.PointLocalSpacePoint1()
	_ = c.PointLocalSpacePoint2()
}

func TestConstraintsHingeCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewHingeConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeHinge {
		t.Errorf("subtype = %d, want hinge", c.SubType())
	}
	if a := c.HingeCurrentAngle(); math.IsNaN(float64(a)) {
		t.Error("current angle is NaN")
	}
}

func TestConstraintsHingeLimits(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewHingeConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	c.HingeSetLimits(-Pi/4, Pi/4)
	if !c.HingeHasLimits() {
		t.Error("expected limits")
	}
	nearFloat(t, "min", c.HingeLimitsMin(), -Pi/4)
	nearFloat(t, "max", c.HingeLimitsMax(), Pi/4)
}

func TestConstraintsHingeMotor(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewHingeConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	c.HingeSetMotorState(MotorStateVelocity)
	if c.HingeMotorState() != MotorStateVelocity {
		t.Error("motor state not velocity")
	}
	c.HingeSetTargetAngularVelocity(1)
	equalFloat(t, "target angular velocity", c.HingeTargetAngularVelocity(), 1)
}

func TestConstraintsSliderCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewSliderConstraintSettings()
	defer settings.Destroy()
	settings.SetSliderAxis(Vec3{1, 0, 0})

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeSlider {
		t.Errorf("subtype = %d, want slider", c.SubType())
	}
}

func TestConstraintsSliderLimits(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewSliderConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	c.SliderSetLimits(-1, 1)
	if !c.SliderHasLimits() {
		t.Error("expected limits")
	}
	equalFloat(t, "min", c.SliderLimitsMin(), -1)
	equalFloat(t, "max", c.SliderLimitsMax(), 1)
}

func TestConstraintsConeCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewConeConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeCone {
		t.Errorf("subtype = %d, want cone", c.SubType())
	}

	c.ConeSetHalfConeAngle(Pi / 6)
	nearFloat(t, "cos half cone angle", c.ConeCosHalfConeAngle(), float32(math.Cos(float64(Pi/6))))
}

func TestConstraintsSwingTwistCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewSwingTwistConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeSwingTwist {
		t.Errorf("subtype = %d, want swing twist", c.SubType())
	}
	if a := c.SwingTwistNormalHalfConeAngle(); math.IsNaN(float64(a)) {
		t.Error("half cone angle is NaN")
	}
}

func TestConstraintsSixDOFCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewSixDOFConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeSixDOF {
		t.Errorf("subtype = %d, want six DOF", c.SubType())
	}
}

func TestConstraintsSixDOFSettingsAxes(t *testing.T) {
	settings := NewSixDOFConstraintSettings()
	defer settings.Destroy()

	settings.MakeFreeAxis(SixDOFConstraintAxisTranslationX)
	if !settings.IsFreeAxis(SixDOFConstraintAxisTranslationX) {
		t.Error("expected free axis")
	}

	settings.MakeFixedAxis(SixDOFConstraintAxisTranslationY)
	if !settings.IsFixedAxis(SixDOFConstraintAxisTranslationY) {
		t.Error("expected fixed axis")
	}

	settings.SetLimitedAxis(SixDOFConstraintAxisTranslationZ, -1, 1)
}

func TestConstraintsGearCreate(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewGearConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()
	if c.SubType() != ConstraintSubTypeGear {
		t.Errorf("subtype = %d, want gear", c.SubType())
	}
	equalFloat(t, "total lambda", c.GearTotalLambda(), 0)
}

func TestConstraintsEnableDisable(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewFixedConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	if !c.Enabled() {
		t.Error("expected enabled")
	}
	c.SetEnabled(false)
	if c.Enabled() {
		t.Error("expected disabled")
	}
	c.SetEnabled(true)
	if !c.Enabled() {
		t.Error("expected enabled")
	}
}

func TestConstraintsUserData(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewFixedConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	c.SetUserData(12345)
	if got := c.UserData(); got != 12345 {
		t.Errorf("user data = %d, want 12345", got)
	}
}

func TestConstraintsPriority(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewFixedConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	c.SetPriority(100)
	if got := c.Priority(); got != 100 {
		t.Errorf("priority = %d, want 100", got)
	}
}

func TestConstraintsTwoBodyGetBodies(t *testing.T) {
	_, b1, b2 := newConstraintFixture(t)
	settings := NewFixedConstraintSettings()
	defer settings.Destroy()

	c := settings.Create(b1, b2)
	defer c.Destroy()

	if c.Body1().ID() != b1.ID() {
		t.Error("body1 mismatch")
	}
	if c.Body2().ID() != b2.ID() {
		t.Error("body2 mismatch")
	}
}
