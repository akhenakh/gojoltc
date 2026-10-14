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

// NewWheelSettingsWV creates wheel settings for a wheeled vehicle. The returned
// settings use the base WheelSettings accessors plus the WV-specific ones.
func NewWheelSettingsWV() *WheelSettings {
	return &WheelSettings{ptr: (*C.JPH_WheelSettings)(unsafe.Pointer(C.JPH_WheelSettingsWV_Create()))}
}

func (s *WheelSettings) wv() *C.JPH_WheelSettingsWV {
	return (*C.JPH_WheelSettingsWV)(unsafe.Pointer(s.ptr))
}

// SetMaxSteerAngle sets the steering limit (wheeled vehicles).
func (s *WheelSettings) SetMaxSteerAngle(v float32) {
	C.JPH_WheelSettingsWV_SetMaxSteerAngle(s.wv(), C.float(v))
}

// MaxSteerAngle returns the steering limit.
func (s *WheelSettings) MaxSteerAngle() float32 {
	return float32(C.JPH_WheelSettingsWV_GetMaxSteerAngle(s.wv()))
}

// SetMaxBrakeTorque sets the brake torque.
func (s *WheelSettings) SetMaxBrakeTorque(v float32) {
	C.JPH_WheelSettingsWV_SetMaxBrakeTorque(s.wv(), C.float(v))
}

// SetMaxHandBrakeTorque sets the hand-brake torque.
func (s *WheelSettings) SetMaxHandBrakeTorque(v float32) {
	C.JPH_WheelSettingsWV_SetMaxHandBrakeTorque(s.wv(), C.float(v))
}

// SetInertia sets the wheel inertia.
func (s *WheelSettings) SetInertia(v float32) { C.JPH_WheelSettingsWV_SetInertia(s.wv(), C.float(v)) }

// SetAngularDamping sets the wheel angular damping.
func (s *WheelSettings) SetAngularDamping(v float32) {
	C.JPH_WheelSettingsWV_SetAngularDamping(s.wv(), C.float(v))
}

// VehicleConstraint is a vehicle: a chassis body with wheels and a controller.
type VehicleConstraint struct{ ptr *C.JPH_VehicleConstraint }

// NewVehicleConstraint creates a vehicle constraint for a chassis body.
//
// The wheel settings and the controller must stay alive for as long as the
// constraint does.
func NewVehicleConstraint(body *Body, settings *VehicleConstraintSettings, wheels []*WheelSettings, controller *WheeledVehicleControllerSettings) *VehicleConstraint {
	var c C.JPH_VehicleConstraintSettings
	C.JPH_VehicleConstraintSettings_Init(&c)
	c.up = vec3ToC(settings.Up)
	c.forward = vec3ToC(settings.Forward)
	c.maxPitchRollAngle = C.float(settings.MaxPitchRollAngle)

	if n := len(wheels); n > 0 {
		buf := (**C.JPH_WheelSettings)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(uintptr(0)))))
		arr := unsafe.Slice(buf, n)
		for i, w := range wheels {
			arr[i] = w.ptr
		}
		c.wheels = buf
		c.wheelsCount = C.uint32_t(n)
		defer C.free(unsafe.Pointer(buf))
	}
	if controller != nil {
		c.controller = (*C.JPH_VehicleControllerSettings)(unsafe.Pointer(controller.ptr))
	}

	return &VehicleConstraint{ptr: C.JPH_VehicleConstraint_Create(body.ptr, &c)}
}

// AsConstraint views the vehicle as a constraint so it can be added to a system.
func (v *VehicleConstraint) AsConstraint() *Constraint {
	return &Constraint{ptr: (*C.JPH_Constraint)(unsafe.Pointer(v.ptr))}
}

// Destroy releases the vehicle constraint. Remove it from the system first.
func (v *VehicleConstraint) Destroy() {
	if v.ptr != nil {
		C.JPH_Constraint_Destroy((*C.JPH_Constraint)(unsafe.Pointer(v.ptr)))
		v.ptr = nil
	}
}

// StepListener returns the vehicle's internal step listener. Do not Destroy it;
// remove it from the system with PhysicsSystem.RemoveStepListener instead.
func (v *VehicleConstraint) StepListener() *StepListener {
	return &StepListener{ptr: C.JPH_VehicleConstraint_AsPhysicsStepListener(v.ptr)}
}

// SetCollisionTester sets the tester used to find ground contacts.
func (v *VehicleConstraint) SetCollisionTester(t *VehicleCollisionTester) {
	C.JPH_VehicleConstraint_SetVehicleCollisionTester(v.ptr, t.ptr)
}

// OverrideGravity replaces the gravity the vehicle sees (empty gravity disables falling).
func (v *VehicleConstraint) OverrideGravity(g vecmath.Vec3) {
	c := vec3ToC(g)
	C.JPH_VehicleConstraint_OverrideGravity(v.ptr, &c)
}

// IsGravityOverridden reports whether the vehicle's gravity is overridden.
func (v *VehicleConstraint) IsGravityOverridden() bool {
	return bool(C.JPH_VehicleConstraint_IsGravityOverridden(v.ptr))
}

// SetMaxPitchRollAngle sets the maximum pitch/roll before the vehicle is forced upright.
func (v *VehicleConstraint) SetMaxPitchRollAngle(v2 float32) {
	C.JPH_VehicleConstraint_SetMaxPitchRollAngle(v.ptr, C.float(v2))
}

// WheelsCount returns the number of wheels.
func (v *VehicleConstraint) WheelsCount() uint32 {
	return uint32(C.JPH_VehicleConstraint_GetWheelsCount(v.ptr))
}

// Wheel returns the runtime wheel at index.
func (v *VehicleConstraint) Wheel(index uint32) *Wheel {
	return &Wheel{ptr: C.JPH_VehicleConstraint_GetWheel(v.ptr, C.uint32_t(index))}
}

// VehicleBody returns the chassis body.
func (v *VehicleConstraint) VehicleBody() *Body {
	return &Body{ptr: (*C.JPH_Body)(unsafe.Pointer(C.JPH_VehicleConstraint_GetVehicleBody(v.ptr)))}
}

// Controller returns the vehicle controller as a wheeled vehicle controller.
func (v *VehicleConstraint) Controller() *WheeledVehicleController {
	return &WheeledVehicleController{ptr: (*C.JPH_WheeledVehicleController)(unsafe.Pointer(C.JPH_VehicleConstraint_GetController(v.ptr)))}
}

// WheeledVehicleController drives a wheeled vehicle.
type WheeledVehicleController struct {
	ptr *C.JPH_WheeledVehicleController
}

// SetDriverInput sets the driver input: forward/right in [-1,1], brake and
// handBrake in [0,1].
func (c *WheeledVehicleController) SetDriverInput(forward, right, brake, handBrake float32) {
	C.JPH_WheeledVehicleController_SetDriverInput(c.ptr, C.float(forward), C.float(right), C.float(brake), C.float(handBrake))
}

// SetForwardInput sets the forward input in [-1,1].
func (c *WheeledVehicleController) SetForwardInput(v float32) {
	C.JPH_WheeledVehicleController_SetForwardInput(c.ptr, C.float(v))
}

// SetRightInput sets the steering input in [-1,1].
func (c *WheeledVehicleController) SetRightInput(v float32) {
	C.JPH_WheeledVehicleController_SetRightInput(c.ptr, C.float(v))
}

// ForwardInput returns the current forward input.
func (c *WheeledVehicleController) ForwardInput() float32 {
	return float32(C.JPH_WheeledVehicleController_GetForwardInput(c.ptr))
}

// RightInput returns the current steering input.
func (c *WheeledVehicleController) RightInput() float32 {
	return float32(C.JPH_WheeledVehicleController_GetRightInput(c.ptr))
}
