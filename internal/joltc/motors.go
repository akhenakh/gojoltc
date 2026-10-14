//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

func springToC(s SpringSettings) C.JPH_SpringSettings {
	return C.JPH_SpringSettings{
		mode:                 C.JPH_SpringMode(s.Mode),
		frequencyOrStiffness: C.float(s.FrequencyOrStiffness),
		damping:              C.float(s.Damping),
	}
}

func springFromC(c C.JPH_SpringSettings) SpringSettings {
	return SpringSettings{
		Mode:                 SpringMode(c.mode),
		FrequencyOrStiffness: float32(c.frequencyOrStiffness),
		Damping:              float32(c.damping),
	}
}

func motorToC(m MotorSettings) C.JPH_MotorSettings {
	return C.JPH_MotorSettings{
		springSettings: springToC(m.SpringSettings),
		minForceLimit:  C.float(m.MinForceLimit),
		maxForceLimit:  C.float(m.MaxForceLimit),
		minTorqueLimit: C.float(m.MinTorqueLimit),
		maxTorqueLimit: C.float(m.MaxTorqueLimit),
	}
}

func motorFromC(c C.JPH_MotorSettings) MotorSettings {
	return MotorSettings{
		SpringSettings: springFromC(c.springSettings),
		MinForceLimit:  float32(c.minForceLimit),
		MaxForceLimit:  float32(c.maxForceLimit),
		MinTorqueLimit: float32(c.minTorqueLimit),
		MaxTorqueLimit: float32(c.maxTorqueLimit),
	}
}

// --- Hinge ------------------------------------------------------------------

// HingeSetTargetAngle sets the angle the position motor drives to.
func (c *Constraint) HingeSetTargetAngle(angle float32) {
	C.JPH_HingeConstraint_SetTargetAngle((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), C.float(angle))
}

// HingeTargetAngle returns the position motor's target angle.
func (c *Constraint) HingeTargetAngle() float32 {
	return float32(C.JPH_HingeConstraint_GetTargetAngle((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeSetMaxFrictionTorque sets the friction torque applied when the motor is off.
func (c *Constraint) HingeSetMaxFrictionTorque(v float32) {
	C.JPH_HingeConstraint_SetMaxFrictionTorque((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), C.float(v))
}

// HingeMaxFrictionTorque returns the friction torque.
func (c *Constraint) HingeMaxFrictionTorque() float32 {
	return float32(C.JPH_HingeConstraint_GetMaxFrictionTorque((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeSetMotorSettings sets the motor spring and force limits.
func (c *Constraint) HingeSetMotorSettings(m MotorSettings) {
	cm := motorToC(m)
	C.JPH_HingeConstraint_SetMotorSettings((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), &cm)
}

// HingeMotorSettings returns the motor settings.
func (c *Constraint) HingeMotorSettings() MotorSettings {
	var cm C.JPH_MotorSettings
	C.JPH_HingeConstraint_GetMotorSettings((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), &cm)
	return motorFromC(cm)
}

// HingeSetLimitsSpringSettings sets the spring used to enforce the limits.
func (c *Constraint) HingeSetLimitsSpringSettings(s SpringSettings) {
	cs := springToC(s)
	C.JPH_HingeConstraint_SetLimitsSpringSettings((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), &cs)
}

// HingeTotalLambdaMotor returns the motor impulse applied this step (force feedback).
func (c *Constraint) HingeTotalLambdaMotor() float32 {
	return float32(C.JPH_HingeConstraint_GetTotalLambdaMotor((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Slider -----------------------------------------------------------------

// SliderSetMotorState sets the slider motor mode.
func (c *Constraint) SliderSetMotorState(s MotorState) {
	C.JPH_SliderConstraint_SetMotorState((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), C.JPH_MotorState(s))
}

// SliderMotorState returns the slider motor mode.
func (c *Constraint) SliderMotorState() MotorState {
	return MotorState(C.JPH_SliderConstraint_GetMotorState((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderSetTargetVelocity sets the velocity motor target.
func (c *Constraint) SliderSetTargetVelocity(v float32) {
	C.JPH_SliderConstraint_SetTargetVelocity((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), C.float(v))
}

// SliderTargetVelocity returns the velocity motor target.
func (c *Constraint) SliderTargetVelocity() float32 {
	return float32(C.JPH_SliderConstraint_GetTargetVelocity((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderSetTargetPosition sets the position motor target.
func (c *Constraint) SliderSetTargetPosition(v float32) {
	C.JPH_SliderConstraint_SetTargetPosition((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), C.float(v))
}

// SliderTargetPosition returns the position motor target.
func (c *Constraint) SliderTargetPosition() float32 {
	return float32(C.JPH_SliderConstraint_GetTargetPosition((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderSetMaxFrictionForce sets the friction force applied when the motor is off.
func (c *Constraint) SliderSetMaxFrictionForce(v float32) {
	C.JPH_SliderConstraint_SetMaxFrictionForce((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), C.float(v))
}

// SliderMaxFrictionForce returns the friction force.
func (c *Constraint) SliderMaxFrictionForce() float32 {
	return float32(C.JPH_SliderConstraint_GetMaxFrictionForce((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderSetMotorSettings sets the motor spring and force limits.
func (c *Constraint) SliderSetMotorSettings(m MotorSettings) {
	cm := motorToC(m)
	C.JPH_SliderConstraint_SetMotorSettings((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), &cm)
}

// SliderSetLimitsSpringSettings sets the spring used to enforce the limits.
func (c *Constraint) SliderSetLimitsSpringSettings(s SpringSettings) {
	cs := springToC(s)
	C.JPH_SliderConstraint_SetLimitsSpringSettings((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), &cs)
}

// SliderCurrentPosition returns the current position along the slider axis.
func (c *Constraint) SliderCurrentPosition() float32 {
	return float32(C.JPH_SliderConstraint_GetCurrentPosition((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderTotalLambdaMotor returns the motor impulse applied this step.
func (c *Constraint) SliderTotalLambdaMotor() float32 {
	return float32(C.JPH_SliderConstraint_GetTotalLambdaMotor((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Six DOF ----------------------------------------------------------------

// SixDOFSetMotorState sets the motor mode of one axis.
func (c *Constraint) SixDOFSetMotorState(axis SixDOFConstraintAxis, state MotorState) {
	C.JPH_SixDOFConstraint_SetMotorState((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), C.JPH_SixDOFConstraintAxis(axis), C.JPH_MotorState(state))
}

// SixDOFMotorState returns the motor mode of one axis.
func (c *Constraint) SixDOFMotorState(axis SixDOFConstraintAxis) MotorState {
	return MotorState(C.JPH_SixDOFConstraint_GetMotorState((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), C.JPH_SixDOFConstraintAxis(axis)))
}

// SixDOFSetTargetVelocityCS sets the target linear velocity in constraint space.
func (c *Constraint) SixDOFSetTargetVelocityCS(v vecmath.Vec3) {
	cv := vec3ToC(v)
	C.JPH_SixDOFConstraint_SetTargetVelocityCS((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &cv)
}

// SixDOFSetTargetAngularVelocityCS sets the target angular velocity in constraint space.
func (c *Constraint) SixDOFSetTargetAngularVelocityCS(v vecmath.Vec3) {
	cv := vec3ToC(v)
	C.JPH_SixDOFConstraint_SetTargetAngularVelocityCS((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &cv)
}

// SixDOFSetTargetPositionCS sets the target position in constraint space.
func (c *Constraint) SixDOFSetTargetPositionCS(v vecmath.Vec3) {
	cv := vec3ToC(v)
	C.JPH_SixDOFConstraint_SetTargetPositionCS((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &cv)
}

// SixDOFSetTargetOrientationCS sets the target orientation in constraint space.
func (c *Constraint) SixDOFSetTargetOrientationCS(q vecmath.Quat) {
	cq := quatToC(q)
	C.JPH_SixDOFConstraint_SetTargetOrientationCS((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &cq)
}

// SixDOFSetTargetOrientationBS sets the target orientation in body space.
func (c *Constraint) SixDOFSetTargetOrientationBS(q vecmath.Quat) {
	cq := quatToC(q)
	C.JPH_SixDOFConstraint_SetTargetOrientationBS((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &cq)
}

// SixDOFSetMaxFriction sets the friction of one axis.
func (c *Constraint) SixDOFSetMaxFriction(axis SixDOFConstraintAxis, friction float32) {
	C.JPH_SixDOFConstraint_SetMaxFriction((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), C.JPH_SixDOFConstraintAxis(axis), C.float(friction))
}

// SixDOFSetLimitsSpringSettings sets the spring used to enforce the limits.
func (c *Constraint) SixDOFSetLimitsSpringSettings(s SpringSettings, axis SixDOFConstraintAxis) {
	cs := springToC(s)
	C.JPH_SixDOFConstraint_SetLimitsSpringSettings((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &cs, C.JPH_SixDOFConstraintAxis(axis))
}

// SixDOFIsFreeAxis reports whether an axis is free.
func (c *Constraint) SixDOFIsFreeAxis(axis SixDOFConstraintAxis) bool {
	return bool(C.JPH_SixDOFConstraint_IsFreeAxis((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), C.JPH_SixDOFConstraintAxis(axis)))
}

// SixDOFIsFixedAxis reports whether an axis is fixed.
func (c *Constraint) SixDOFIsFixedAxis(axis SixDOFConstraintAxis) bool {
	return bool(C.JPH_SixDOFConstraint_IsFixedAxis((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), C.JPH_SixDOFConstraintAxis(axis)))
}

// SixDOFTotalLambdaMotorTranslation returns the motor impulses on the translation axes.
func (c *Constraint) SixDOFTotalLambdaMotorTranslation() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_SixDOFConstraint_GetTotalLambdaMotorTranslation((*C.JPH_SixDOFConstraint)(unsafe.Pointer(c.ptr)), &o)
	return vec3FromC(o)
}

// --- Distance ---------------------------------------------------------------

// DistanceSetLimitsSpringSettings sets the spring used to enforce the distance limits.
func (c *Constraint) DistanceSetLimitsSpringSettings(s SpringSettings) {
	cs := springToC(s)
	C.JPH_DistanceConstraint_SetLimitsSpringSettings((*C.JPH_DistanceConstraint)(unsafe.Pointer(c.ptr)), &cs)
}
