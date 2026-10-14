package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Vehicle types and value settings, re-exported from the bindings.
type (
	// WheelSettings holds settings for a wheel.
	WheelSettings = joltc.WheelSettings
	// Wheel is a runtime wheel instance.
	Wheel = joltc.Wheel
	// VehicleCollisionTester is a vehicle ground collision tester.
	VehicleCollisionTester = joltc.VehicleCollisionTester
	// VehicleTransmissionSettings holds transmission settings.
	VehicleTransmissionSettings = joltc.VehicleTransmissionSettings
	// WheeledVehicleControllerSettings holds settings for a wheeled vehicle controller.
	WheeledVehicleControllerSettings = joltc.WheeledVehicleControllerSettings
	// VehicleEngineSettings mirrors JPH_VehicleEngineSettings.
	VehicleEngineSettings = joltc.VehicleEngineSettings
	// VehicleDifferentialSettings mirrors JPH_VehicleDifferentialSettings.
	VehicleDifferentialSettings = joltc.VehicleDifferentialSettings
	// VehicleAntiRollBar mirrors JPH_VehicleAntiRollBar.
	VehicleAntiRollBar = joltc.VehicleAntiRollBar
	// VehicleConstraintSettings mirrors JPH_VehicleConstraintSettings.
	VehicleConstraintSettings = joltc.VehicleConstraintSettings
)

// NewWheelSettings creates wheel settings with defaults.
func NewWheelSettings() *WheelSettings { return joltc.NewWheelSettings() }

// NewWheel creates a wheel from settings.
func NewWheel(settings *WheelSettings) *Wheel { return joltc.NewWheel(settings) }

// NewVehicleCollisionTesterRay creates a ray-based collision tester.
func NewVehicleCollisionTesterRay(layer ObjectLayer, up Vec3, maxSlopeAngle float32) *VehicleCollisionTester {
	return joltc.NewVehicleCollisionTesterRay(layer, up, maxSlopeAngle)
}

// NewVehicleCollisionTesterCastSphere creates a sphere-cast collision tester.
func NewVehicleCollisionTesterCastSphere(layer ObjectLayer, radius float32, up Vec3, maxSlopeAngle float32) *VehicleCollisionTester {
	return joltc.NewVehicleCollisionTesterCastSphere(layer, radius, up, maxSlopeAngle)
}

// NewVehicleCollisionTesterCastCylinder creates a cylinder-cast collision tester.
func NewVehicleCollisionTesterCastCylinder(layer ObjectLayer, convexRadiusFraction float32) *VehicleCollisionTester {
	return joltc.NewVehicleCollisionTesterCastCylinder(layer, convexRadiusFraction)
}

// NewVehicleEngineSettings returns engine settings with Jolt defaults.
func NewVehicleEngineSettings() VehicleEngineSettings { return joltc.NewVehicleEngineSettings() }

// NewVehicleDifferentialSettings returns differential settings with Jolt defaults.
func NewVehicleDifferentialSettings() VehicleDifferentialSettings {
	return joltc.NewVehicleDifferentialSettings()
}

// NewVehicleAntiRollBar returns an anti-roll bar with Jolt defaults.
func NewVehicleAntiRollBar() VehicleAntiRollBar { return joltc.NewVehicleAntiRollBar() }

// NewVehicleConstraintSettings returns vehicle constraint settings with defaults.
func NewVehicleConstraintSettings() VehicleConstraintSettings {
	return joltc.NewVehicleConstraintSettings()
}

// NewVehicleTransmissionSettings creates transmission settings with defaults.
func NewVehicleTransmissionSettings() *VehicleTransmissionSettings {
	return joltc.NewVehicleTransmissionSettings()
}

// NewWheeledVehicleControllerSettings creates wheeled controller settings with defaults.
func NewWheeledVehicleControllerSettings() *WheeledVehicleControllerSettings {
	return joltc.NewWheeledVehicleControllerSettings()
}
