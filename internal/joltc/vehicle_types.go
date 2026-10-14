package joltc

import "github.com/akhenakh/gojoltc/internal/vecmath"

// VehicleEngineSettings mirrors JPH_VehicleEngineSettings (the normalizedTorque
// curve is omitted).
type VehicleEngineSettings struct {
	MaxTorque      float32
	MinRPM         float32
	MaxRPM         float32
	Inertia        float32
	AngularDamping float32
}

// VehicleDifferentialSettings mirrors JPH_VehicleDifferentialSettings.
type VehicleDifferentialSettings struct {
	LeftWheel         int32
	RightWheel        int32
	DifferentialRatio float32
	LeftRightSplit    float32
	LimitedSlipRatio  float32
	EngineTorqueRatio float32
}

// VehicleAntiRollBar mirrors JPH_VehicleAntiRollBar.
type VehicleAntiRollBar struct {
	LeftWheel  int32
	RightWheel int32
	Stiffness  float32
}

// VehicleConstraintSettings mirrors the value part of JPH_VehicleConstraintSettings.
type VehicleConstraintSettings struct {
	Up                vecmath.Vec3
	Forward           vecmath.Vec3
	MaxPitchRollAngle float32
}
