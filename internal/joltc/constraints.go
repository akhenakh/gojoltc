//go:build cgo

package joltc

/*
#include <joltc.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// Constraint is a physics joint between bodies.
type Constraint struct {
	ptr *C.JPH_Constraint
}

// Destroy releases the constraint.
func (c *Constraint) Destroy() {
	if c.ptr != nil {
		C.JPH_Constraint_Destroy(c.ptr)
		c.ptr = nil
	}
}

// SubType returns the constraint's specific type.
func (c *Constraint) SubType() ConstraintSubType {
	return ConstraintSubType(C.JPH_Constraint_GetSubType(c.ptr))
}

// Enabled reports whether the constraint is enabled.
func (c *Constraint) Enabled() bool { return bool(C.JPH_Constraint_GetEnabled(c.ptr)) }

// SetEnabled enables or disables the constraint.
func (c *Constraint) SetEnabled(v bool) { C.JPH_Constraint_SetEnabled(c.ptr, C.bool(v)) }

// UserData returns the constraint's user data.
func (c *Constraint) UserData() uint64 { return uint64(C.JPH_Constraint_GetUserData(c.ptr)) }

// SetUserData sets the constraint's user data.
func (c *Constraint) SetUserData(v uint64) { C.JPH_Constraint_SetUserData(c.ptr, C.uint64_t(v)) }

// Priority returns the constraint priority.
func (c *Constraint) Priority() uint32 { return uint32(C.JPH_Constraint_GetConstraintPriority(c.ptr)) }

// SetPriority sets the constraint priority.
func (c *Constraint) SetPriority(v uint32) {
	C.JPH_Constraint_SetConstraintPriority(c.ptr, C.uint32_t(v))
}

// Body1 returns the first body of a two-body constraint.
func (c *Constraint) Body1() *Body {
	return &Body{ptr: C.JPH_TwoBodyConstraint_GetBody1((*C.JPH_TwoBodyConstraint)(unsafe.Pointer(c.ptr)))}
}

// Body2 returns the second body of a two-body constraint.
func (c *Constraint) Body2() *Body {
	return &Body{ptr: C.JPH_TwoBodyConstraint_GetBody2((*C.JPH_TwoBodyConstraint)(unsafe.Pointer(c.ptr)))}
}

// --- Fixed ------------------------------------------------------------------

// FixedConstraintSettings holds settings for a fixed constraint.
type FixedConstraintSettings struct {
	ptr *C.JPH_FixedConstraintSettings
}

// NewFixedConstraintSettings creates fixed constraint settings.
func NewFixedConstraintSettings() *FixedConstraintSettings {
	p := (*C.JPH_FixedConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_FixedConstraintSettings{}))))
	C.JPH_FixedConstraintSettings_Init(p)
	return &FixedConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *FixedConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the fixed constraint between two bodies.
func (s *FixedConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_FixedConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// FixedTotalLambdaPosition returns the total position impulse of a fixed constraint.
func (c *Constraint) FixedTotalLambdaPosition() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_FixedConstraint_GetTotalLambdaPosition((*C.JPH_FixedConstraint)(unsafe.Pointer(c.ptr)), &o)
	return vec3FromC(o)
}

// --- Distance ---------------------------------------------------------------

// DistanceConstraintSettings holds settings for a distance constraint.
type DistanceConstraintSettings struct {
	ptr *C.JPH_DistanceConstraintSettings
}

// NewDistanceConstraintSettings creates distance constraint settings.
func NewDistanceConstraintSettings() *DistanceConstraintSettings {
	p := (*C.JPH_DistanceConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_DistanceConstraintSettings{}))))
	C.JPH_DistanceConstraintSettings_Init(p)
	return &DistanceConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *DistanceConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the distance constraint between two bodies.
func (s *DistanceConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_DistanceConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// DistanceSetDistance sets the distance limits.
func (c *Constraint) DistanceSetDistance(minDistance, maxDistance float32) {
	C.JPH_DistanceConstraint_SetDistance((*C.JPH_DistanceConstraint)(unsafe.Pointer(c.ptr)), C.float(minDistance), C.float(maxDistance))
}

// DistanceMinDistance returns the minimum distance.
func (c *Constraint) DistanceMinDistance() float32 {
	return float32(C.JPH_DistanceConstraint_GetMinDistance((*C.JPH_DistanceConstraint)(unsafe.Pointer(c.ptr))))
}

// DistanceMaxDistance returns the maximum distance.
func (c *Constraint) DistanceMaxDistance() float32 {
	return float32(C.JPH_DistanceConstraint_GetMaxDistance((*C.JPH_DistanceConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Point ------------------------------------------------------------------

// PointConstraintSettings holds settings for a point constraint.
type PointConstraintSettings struct {
	ptr *C.JPH_PointConstraintSettings
}

// NewPointConstraintSettings creates point constraint settings.
func NewPointConstraintSettings() *PointConstraintSettings {
	p := (*C.JPH_PointConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_PointConstraintSettings{}))))
	C.JPH_PointConstraintSettings_Init(p)
	return &PointConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *PointConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the point constraint between two bodies.
func (s *PointConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_PointConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// PointLocalSpacePoint1 returns the local space point on body 1.
func (c *Constraint) PointLocalSpacePoint1() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_PointConstraint_GetLocalSpacePoint1((*C.JPH_PointConstraint)(unsafe.Pointer(c.ptr)), &o)
	return vec3FromC(o)
}

// PointLocalSpacePoint2 returns the local space point on body 2.
func (c *Constraint) PointLocalSpacePoint2() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_PointConstraint_GetLocalSpacePoint2((*C.JPH_PointConstraint)(unsafe.Pointer(c.ptr)), &o)
	return vec3FromC(o)
}

// --- Hinge ------------------------------------------------------------------

// HingeConstraintSettings holds settings for a hinge constraint.
type HingeConstraintSettings struct {
	ptr *C.JPH_HingeConstraintSettings
}

// NewHingeConstraintSettings creates hinge constraint settings.
func NewHingeConstraintSettings() *HingeConstraintSettings {
	p := (*C.JPH_HingeConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_HingeConstraintSettings{}))))
	C.JPH_HingeConstraintSettings_Init(p)
	return &HingeConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *HingeConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the hinge constraint between two bodies.
func (s *HingeConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_HingeConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// HingeCurrentAngle returns the current hinge angle.
func (c *Constraint) HingeCurrentAngle() float32 {
	return float32(C.JPH_HingeConstraint_GetCurrentAngle((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeSetLimits sets the hinge angle limits.
func (c *Constraint) HingeSetLimits(min, max float32) {
	C.JPH_HingeConstraint_SetLimits((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), C.float(min), C.float(max))
}

// HingeHasLimits reports whether the hinge has limits.
func (c *Constraint) HingeHasLimits() bool {
	return bool(C.JPH_HingeConstraint_HasLimits((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeLimitsMin returns the minimum hinge angle.
func (c *Constraint) HingeLimitsMin() float32 {
	return float32(C.JPH_HingeConstraint_GetLimitsMin((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeLimitsMax returns the maximum hinge angle.
func (c *Constraint) HingeLimitsMax() float32 {
	return float32(C.JPH_HingeConstraint_GetLimitsMax((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeSetMotorState sets the hinge motor state.
func (c *Constraint) HingeSetMotorState(state MotorState) {
	C.JPH_HingeConstraint_SetMotorState((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), C.JPH_MotorState(state))
}

// HingeMotorState returns the hinge motor state.
func (c *Constraint) HingeMotorState() MotorState {
	return MotorState(C.JPH_HingeConstraint_GetMotorState((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// HingeSetTargetAngularVelocity sets the hinge target angular velocity.
func (c *Constraint) HingeSetTargetAngularVelocity(v float32) {
	C.JPH_HingeConstraint_SetTargetAngularVelocity((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr)), C.float(v))
}

// HingeTargetAngularVelocity returns the hinge target angular velocity.
func (c *Constraint) HingeTargetAngularVelocity() float32 {
	return float32(C.JPH_HingeConstraint_GetTargetAngularVelocity((*C.JPH_HingeConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Slider -----------------------------------------------------------------

// SliderConstraintSettings holds settings for a slider constraint.
type SliderConstraintSettings struct {
	ptr *C.JPH_SliderConstraintSettings
}

// NewSliderConstraintSettings creates slider constraint settings.
func NewSliderConstraintSettings() *SliderConstraintSettings {
	p := (*C.JPH_SliderConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_SliderConstraintSettings{}))))
	C.JPH_SliderConstraintSettings_Init(p)
	return &SliderConstraintSettings{ptr: p}
}

// SetSliderAxis sets the slider axis.
func (s *SliderConstraintSettings) SetSliderAxis(axis vecmath.Vec3) {
	c := vec3ToC(axis)
	C.JPH_SliderConstraintSettings_SetSliderAxis(s.ptr, &c)
}

// Destroy releases the settings.
func (s *SliderConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the slider constraint between two bodies.
func (s *SliderConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_SliderConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// SliderSetLimits sets the slider limits.
func (c *Constraint) SliderSetLimits(min, max float32) {
	C.JPH_SliderConstraint_SetLimits((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr)), C.float(min), C.float(max))
}

// SliderHasLimits reports whether the slider has limits.
func (c *Constraint) SliderHasLimits() bool {
	return bool(C.JPH_SliderConstraint_HasLimits((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderLimitsMin returns the minimum slider limit.
func (c *Constraint) SliderLimitsMin() float32 {
	return float32(C.JPH_SliderConstraint_GetLimitsMin((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// SliderLimitsMax returns the maximum slider limit.
func (c *Constraint) SliderLimitsMax() float32 {
	return float32(C.JPH_SliderConstraint_GetLimitsMax((*C.JPH_SliderConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Cone -------------------------------------------------------------------

// ConeConstraintSettings holds settings for a cone constraint.
type ConeConstraintSettings struct{ ptr *C.JPH_ConeConstraintSettings }

// NewConeConstraintSettings creates cone constraint settings.
func NewConeConstraintSettings() *ConeConstraintSettings {
	p := (*C.JPH_ConeConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_ConeConstraintSettings{}))))
	C.JPH_ConeConstraintSettings_Init(p)
	return &ConeConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *ConeConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the cone constraint between two bodies.
func (s *ConeConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_ConeConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// ConeSetHalfConeAngle sets the cone half angle.
func (c *Constraint) ConeSetHalfConeAngle(angle float32) {
	C.JPH_ConeConstraint_SetHalfConeAngle((*C.JPH_ConeConstraint)(unsafe.Pointer(c.ptr)), C.float(angle))
}

// ConeCosHalfConeAngle returns the cosine of the cone half angle.
func (c *Constraint) ConeCosHalfConeAngle() float32 {
	return float32(C.JPH_ConeConstraint_GetCosHalfConeAngle((*C.JPH_ConeConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Swing twist ------------------------------------------------------------

// SwingTwistConstraintSettings holds settings for a swing-twist constraint.
type SwingTwistConstraintSettings struct {
	ptr *C.JPH_SwingTwistConstraintSettings
}

// NewSwingTwistConstraintSettings creates swing-twist constraint settings.
func NewSwingTwistConstraintSettings() *SwingTwistConstraintSettings {
	p := (*C.JPH_SwingTwistConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_SwingTwistConstraintSettings{}))))
	C.JPH_SwingTwistConstraintSettings_Init(p)
	return &SwingTwistConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *SwingTwistConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the swing-twist constraint between two bodies.
func (s *SwingTwistConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_SwingTwistConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// SwingTwistNormalHalfConeAngle returns the normal half cone angle.
func (c *Constraint) SwingTwistNormalHalfConeAngle() float32 {
	return float32(C.JPH_SwingTwistConstraint_GetNormalHalfConeAngle((*C.JPH_SwingTwistConstraint)(unsafe.Pointer(c.ptr))))
}

// --- Six DOF ----------------------------------------------------------------

// SixDOFConstraintSettings holds settings for a six DOF constraint.
type SixDOFConstraintSettings struct {
	ptr *C.JPH_SixDOFConstraintSettings
}

// NewSixDOFConstraintSettings creates six DOF constraint settings.
func NewSixDOFConstraintSettings() *SixDOFConstraintSettings {
	p := (*C.JPH_SixDOFConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_SixDOFConstraintSettings{}))))
	C.JPH_SixDOFConstraintSettings_Init(p)
	return &SixDOFConstraintSettings{ptr: p}
}

// MakeFreeAxis makes an axis free.
func (s *SixDOFConstraintSettings) MakeFreeAxis(axis SixDOFConstraintAxis) {
	C.JPH_SixDOFConstraintSettings_MakeFreeAxis(s.ptr, C.JPH_SixDOFConstraintAxis(axis))
}

// IsFreeAxis reports whether an axis is free.
func (s *SixDOFConstraintSettings) IsFreeAxis(axis SixDOFConstraintAxis) bool {
	return bool(C.JPH_SixDOFConstraintSettings_IsFreeAxis(s.ptr, C.JPH_SixDOFConstraintAxis(axis)))
}

// MakeFixedAxis makes an axis fixed.
func (s *SixDOFConstraintSettings) MakeFixedAxis(axis SixDOFConstraintAxis) {
	C.JPH_SixDOFConstraintSettings_MakeFixedAxis(s.ptr, C.JPH_SixDOFConstraintAxis(axis))
}

// IsFixedAxis reports whether an axis is fixed.
func (s *SixDOFConstraintSettings) IsFixedAxis(axis SixDOFConstraintAxis) bool {
	return bool(C.JPH_SixDOFConstraintSettings_IsFixedAxis(s.ptr, C.JPH_SixDOFConstraintAxis(axis)))
}

// SetLimitedAxis limits an axis between min and max.
func (s *SixDOFConstraintSettings) SetLimitedAxis(axis SixDOFConstraintAxis, min, max float32) {
	C.JPH_SixDOFConstraintSettings_SetLimitedAxis(s.ptr, C.JPH_SixDOFConstraintAxis(axis), C.float(min), C.float(max))
}

// Destroy releases the settings.
func (s *SixDOFConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the six DOF constraint between two bodies.
func (s *SixDOFConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_SixDOFConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// --- Gear -------------------------------------------------------------------

// GearConstraintSettings holds settings for a gear constraint.
type GearConstraintSettings struct{ ptr *C.JPH_GearConstraintSettings }

// NewGearConstraintSettings creates gear constraint settings.
func NewGearConstraintSettings() *GearConstraintSettings {
	p := (*C.JPH_GearConstraintSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_GearConstraintSettings{}))))
	C.JPH_GearConstraintSettings_Init(p)
	return &GearConstraintSettings{ptr: p}
}

// Destroy releases the settings.
func (s *GearConstraintSettings) Destroy() {
	if s.ptr != nil {
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Create builds the gear constraint between two bodies.
func (s *GearConstraintSettings) Create(b1, b2 *Body) *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(C.JPH_GearConstraint_Create(s.ptr, b1.ptr, b2.ptr)))}
}

// GearTotalLambda returns the total gear impulse.
func (c *Constraint) GearTotalLambda() float32 {
	return float32(C.JPH_GearConstraint_GetTotalLambda((*C.JPH_GearConstraint)(unsafe.Pointer(c.ptr))))
}
