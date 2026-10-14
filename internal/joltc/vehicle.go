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

// WheelSettings holds settings for a wheel.
type WheelSettings struct{ ptr *C.JPH_WheelSettings }

// NewWheelSettings creates wheel settings with defaults.
func NewWheelSettings() *WheelSettings { return &WheelSettings{ptr: C.JPH_WheelSettings_Create()} }

// Destroy releases the settings.
func (s *WheelSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_WheelSettings_Destroy(s.ptr)
		s.ptr = nil
	}
}

// Position returns the wheel position.
func (s *WheelSettings) Position() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_WheelSettings_GetPosition(s.ptr, &o)
	return vec3FromC(o)
}

// SetPosition sets the wheel position.
func (s *WheelSettings) SetPosition(v vecmath.Vec3) {
	c := vec3ToC(v)
	C.JPH_WheelSettings_SetPosition(s.ptr, &c)
}

// Radius returns the wheel radius.
func (s *WheelSettings) Radius() float32 { return float32(C.JPH_WheelSettings_GetRadius(s.ptr)) }

// SetRadius sets the wheel radius.
func (s *WheelSettings) SetRadius(v float32) { C.JPH_WheelSettings_SetRadius(s.ptr, C.float(v)) }

// Width returns the wheel width.
func (s *WheelSettings) Width() float32 { return float32(C.JPH_WheelSettings_GetWidth(s.ptr)) }

// SetWidth sets the wheel width.
func (s *WheelSettings) SetWidth(v float32) { C.JPH_WheelSettings_SetWidth(s.ptr, C.float(v)) }

// SuspensionMinLength returns the minimum suspension length.
func (s *WheelSettings) SuspensionMinLength() float32 {
	return float32(C.JPH_WheelSettings_GetSuspensionMinLength(s.ptr))
}

// SetSuspensionMinLength sets the minimum suspension length.
func (s *WheelSettings) SetSuspensionMinLength(v float32) {
	C.JPH_WheelSettings_SetSuspensionMinLength(s.ptr, C.float(v))
}

// SuspensionMaxLength returns the maximum suspension length.
func (s *WheelSettings) SuspensionMaxLength() float32 {
	return float32(C.JPH_WheelSettings_GetSuspensionMaxLength(s.ptr))
}

// SetSuspensionMaxLength sets the maximum suspension length.
func (s *WheelSettings) SetSuspensionMaxLength(v float32) {
	C.JPH_WheelSettings_SetSuspensionMaxLength(s.ptr, C.float(v))
}

// Wheel is a runtime wheel instance.
type Wheel struct{ ptr *C.JPH_Wheel }

// NewWheel creates a wheel from settings.
func NewWheel(settings *WheelSettings) *Wheel { return &Wheel{ptr: C.JPH_Wheel_Create(settings.ptr)} }

// Destroy releases the wheel.
func (w *Wheel) Destroy() {
	if w.ptr != nil {
		C.JPH_Wheel_Destroy(w.ptr)
		w.ptr = nil
	}
}

// AngularVelocity returns the wheel angular velocity.
func (w *Wheel) AngularVelocity() float32 { return float32(C.JPH_Wheel_GetAngularVelocity(w.ptr)) }

// SetAngularVelocity sets the wheel angular velocity.
func (w *Wheel) SetAngularVelocity(v float32) { C.JPH_Wheel_SetAngularVelocity(w.ptr, C.float(v)) }

// RotationAngle returns the wheel rotation angle.
func (w *Wheel) RotationAngle() float32 { return float32(C.JPH_Wheel_GetRotationAngle(w.ptr)) }

// SetRotationAngle sets the wheel rotation angle.
func (w *Wheel) SetRotationAngle(v float32) { C.JPH_Wheel_SetRotationAngle(w.ptr, C.float(v)) }

// SteerAngle returns the wheel steer angle.
func (w *Wheel) SteerAngle() float32 { return float32(C.JPH_Wheel_GetSteerAngle(w.ptr)) }

// SetSteerAngle sets the wheel steer angle.
func (w *Wheel) SetSteerAngle(v float32) { C.JPH_Wheel_SetSteerAngle(w.ptr, C.float(v)) }

// VehicleCollisionTester is a vehicle ground collision tester.
type VehicleCollisionTester struct{ ptr *C.JPH_VehicleCollisionTester }

// Destroy releases the tester.
func (t *VehicleCollisionTester) Destroy() {
	if t.ptr != nil {
		C.JPH_VehicleCollisionTester_Destroy(t.ptr)
		t.ptr = nil
	}
}

// ObjectLayer returns the tester's object layer.
func (t *VehicleCollisionTester) ObjectLayer() ObjectLayer {
	return ObjectLayer(C.JPH_VehicleCollisionTester_GetObjectLayer(t.ptr))
}

// NewVehicleCollisionTesterRay creates a ray-based collision tester.
func NewVehicleCollisionTesterRay(layer ObjectLayer, up vecmath.Vec3, maxSlopeAngle float32) *VehicleCollisionTester {
	u := vec3ToC(up)
	return &VehicleCollisionTester{ptr: (*C.JPH_VehicleCollisionTester)(unsafe.Pointer(C.JPH_VehicleCollisionTesterRay_Create(C.JPH_ObjectLayer(layer), &u, C.float(maxSlopeAngle))))}
}

// NewVehicleCollisionTesterCastSphere creates a sphere-cast collision tester.
func NewVehicleCollisionTesterCastSphere(layer ObjectLayer, radius float32, up vecmath.Vec3, maxSlopeAngle float32) *VehicleCollisionTester {
	u := vec3ToC(up)
	return &VehicleCollisionTester{ptr: (*C.JPH_VehicleCollisionTester)(unsafe.Pointer(C.JPH_VehicleCollisionTesterCastSphere_Create(C.JPH_ObjectLayer(layer), C.float(radius), &u, C.float(maxSlopeAngle))))}
}

// NewVehicleCollisionTesterCastCylinder creates a cylinder-cast collision tester.
func NewVehicleCollisionTesterCastCylinder(layer ObjectLayer, convexRadiusFraction float32) *VehicleCollisionTester {
	return &VehicleCollisionTester{ptr: (*C.JPH_VehicleCollisionTester)(unsafe.Pointer(C.JPH_VehicleCollisionTesterCastCylinder_Create(C.JPH_ObjectLayer(layer), C.float(convexRadiusFraction))))}
}

// NewVehicleEngineSettings returns engine settings with Jolt defaults.
func NewVehicleEngineSettings() VehicleEngineSettings {
	var c C.JPH_VehicleEngineSettings
	C.JPH_VehicleEngineSettings_Init(&c)
	return engineFromC(c)
}

// NewVehicleDifferentialSettings returns differential settings with Jolt defaults.
func NewVehicleDifferentialSettings() VehicleDifferentialSettings {
	var c C.JPH_VehicleDifferentialSettings
	C.JPH_VehicleDifferentialSettings_Init(&c)
	return VehicleDifferentialSettings{
		LeftWheel:         int32(c.leftWheel),
		RightWheel:        int32(c.rightWheel),
		DifferentialRatio: float32(c.differentialRatio),
		LeftRightSplit:    float32(c.leftRightSplit),
		LimitedSlipRatio:  float32(c.limitedSlipRatio),
		EngineTorqueRatio: float32(c.engineTorqueRatio),
	}
}

// NewVehicleAntiRollBar returns an anti-roll bar with Jolt defaults.
func NewVehicleAntiRollBar() VehicleAntiRollBar {
	var c C.JPH_VehicleAntiRollBar
	C.JPH_VehicleAntiRollBar_Init(&c)
	return VehicleAntiRollBar{LeftWheel: int32(c.leftWheel), RightWheel: int32(c.rightWheel), Stiffness: float32(c.stiffness)}
}

// NewVehicleConstraintSettings returns vehicle constraint settings with defaults.
func NewVehicleConstraintSettings() VehicleConstraintSettings {
	var c C.JPH_VehicleConstraintSettings
	C.JPH_VehicleConstraintSettings_Init(&c)
	return VehicleConstraintSettings{
		Up:                vec3FromC(c.up),
		Forward:           vec3FromC(c.forward),
		MaxPitchRollAngle: float32(c.maxPitchRollAngle),
	}
}

func engineFromC(c C.JPH_VehicleEngineSettings) VehicleEngineSettings {
	return VehicleEngineSettings{
		MaxTorque:      float32(c.maxTorque),
		MinRPM:         float32(c.minRPM),
		MaxRPM:         float32(c.maxRPM),
		Inertia:        float32(c.inertia),
		AngularDamping: float32(c.angularDamping),
	}
}

func engineToC(s VehicleEngineSettings) C.JPH_VehicleEngineSettings {
	// Initialize first: Init fills the normalizedTorque curve pointer, which
	// SetEngine dereferences. Leaving it nil crashes.
	var c C.JPH_VehicleEngineSettings
	C.JPH_VehicleEngineSettings_Init(&c)
	c.maxTorque = C.float(s.MaxTorque)
	c.minRPM = C.float(s.MinRPM)
	c.maxRPM = C.float(s.MaxRPM)
	c.inertia = C.float(s.Inertia)
	c.angularDamping = C.float(s.AngularDamping)
	return c
}

// VehicleTransmissionSettings holds transmission settings.
type VehicleTransmissionSettings struct {
	ptr *C.JPH_VehicleTransmissionSettings
}

// NewVehicleTransmissionSettings creates transmission settings with defaults.
func NewVehicleTransmissionSettings() *VehicleTransmissionSettings {
	return &VehicleTransmissionSettings{ptr: C.JPH_VehicleTransmissionSettings_Create()}
}

// Destroy releases the settings.
func (s *VehicleTransmissionSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_VehicleTransmissionSettings_Destroy(s.ptr)
		s.ptr = nil
	}
}

// Mode returns the transmission mode.
func (s *VehicleTransmissionSettings) Mode() TransmissionMode {
	return TransmissionMode(C.JPH_VehicleTransmissionSettings_GetMode(s.ptr))
}

// SetMode sets the transmission mode.
func (s *VehicleTransmissionSettings) SetMode(m TransmissionMode) {
	C.JPH_VehicleTransmissionSettings_SetMode(s.ptr, C.JPH_TransmissionMode(m))
}

// GearRatioCount returns the number of gear ratios.
func (s *VehicleTransmissionSettings) GearRatioCount() uint32 {
	return uint32(C.JPH_VehicleTransmissionSettings_GetGearRatioCount(s.ptr))
}

// GearRatio returns the gear ratio at index.
func (s *VehicleTransmissionSettings) GearRatio(index uint32) float32 {
	return float32(C.JPH_VehicleTransmissionSettings_GetGearRatio(s.ptr, C.uint32_t(index)))
}

// SetGearRatios sets the gear ratios.
func (s *VehicleTransmissionSettings) SetGearRatios(ratios []float32) {
	if len(ratios) == 0 {
		C.JPH_VehicleTransmissionSettings_SetGearRatios(s.ptr, nil, 0)
		return
	}
	buf := (*C.float)(C.malloc(C.size_t(len(ratios)) * C.size_t(unsafe.Sizeof(C.float(0)))))
	defer C.free(unsafe.Pointer(buf))
	arr := unsafe.Slice(buf, len(ratios))
	for i, r := range ratios {
		arr[i] = C.float(r)
	}
	C.JPH_VehicleTransmissionSettings_SetGearRatios(s.ptr, buf, C.uint32_t(len(ratios)))
}

// SwitchTime returns the gear switch time.
func (s *VehicleTransmissionSettings) SwitchTime() float32 {
	return float32(C.JPH_VehicleTransmissionSettings_GetSwitchTime(s.ptr))
}

// SetSwitchTime sets the gear switch time.
func (s *VehicleTransmissionSettings) SetSwitchTime(v float32) {
	C.JPH_VehicleTransmissionSettings_SetSwitchTime(s.ptr, C.float(v))
}

// ClutchReleaseTime returns the clutch release time.
func (s *VehicleTransmissionSettings) ClutchReleaseTime() float32 {
	return float32(C.JPH_VehicleTransmissionSettings_GetClutchReleaseTime(s.ptr))
}

// SetClutchReleaseTime sets the clutch release time.
func (s *VehicleTransmissionSettings) SetClutchReleaseTime(v float32) {
	C.JPH_VehicleTransmissionSettings_SetClutchReleaseTime(s.ptr, C.float(v))
}

// ShiftUpRPM returns the upshift RPM.
func (s *VehicleTransmissionSettings) ShiftUpRPM() float32 {
	return float32(C.JPH_VehicleTransmissionSettings_GetShiftUpRPM(s.ptr))
}

// SetShiftUpRPM sets the upshift RPM.
func (s *VehicleTransmissionSettings) SetShiftUpRPM(v float32) {
	C.JPH_VehicleTransmissionSettings_SetShiftUpRPM(s.ptr, C.float(v))
}

// ShiftDownRPM returns the downshift RPM.
func (s *VehicleTransmissionSettings) ShiftDownRPM() float32 {
	return float32(C.JPH_VehicleTransmissionSettings_GetShiftDownRPM(s.ptr))
}

// SetShiftDownRPM sets the downshift RPM.
func (s *VehicleTransmissionSettings) SetShiftDownRPM(v float32) {
	C.JPH_VehicleTransmissionSettings_SetShiftDownRPM(s.ptr, C.float(v))
}

// WheeledVehicleControllerSettings holds settings for a wheeled vehicle controller.
type WheeledVehicleControllerSettings struct {
	ptr *C.JPH_WheeledVehicleControllerSettings
}

// NewWheeledVehicleControllerSettings creates wheeled controller settings with defaults.
func NewWheeledVehicleControllerSettings() *WheeledVehicleControllerSettings {
	return &WheeledVehicleControllerSettings{ptr: C.JPH_WheeledVehicleControllerSettings_Create()}
}

// Destroy releases the settings.
func (s *WheeledVehicleControllerSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_VehicleControllerSettings_Destroy((*C.JPH_VehicleControllerSettings)(unsafe.Pointer(s.ptr)))
		s.ptr = nil
	}
}

// Engine returns the engine settings.
func (s *WheeledVehicleControllerSettings) Engine() VehicleEngineSettings {
	var c C.JPH_VehicleEngineSettings
	C.JPH_WheeledVehicleControllerSettings_GetEngine(s.ptr, &c)
	return engineFromC(c)
}

// SetEngine sets the engine settings.
func (s *WheeledVehicleControllerSettings) SetEngine(e VehicleEngineSettings) {
	c := engineToC(e)
	C.JPH_WheeledVehicleControllerSettings_SetEngine(s.ptr, &c)
}

// AddDifferential connects a pair of wheels (indexes into the vehicle's wheel
// list) through a differential driven by the engine. Without at least one
// differential no engine torque reaches any wheel.
func (s *WheeledVehicleControllerSettings) AddDifferential(leftWheel, rightWheel int) {
	C.JPH_WheeledVehicleControllerSettings_AddDifferential(s.ptr, C.int(leftWheel), C.int(rightWheel))
}
