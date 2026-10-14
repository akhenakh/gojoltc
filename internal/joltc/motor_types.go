package joltc

import "math"

// Constraint motor/spring value types. They live here (not in the cgo file) so
// the cgo-less stub and the public package can reference them.

// SpringSettings configures a spring (JPH_SpringSettings).
type SpringSettings struct {
	Mode                 SpringMode
	FrequencyOrStiffness float32
	Damping              float32
}

// DefaultSpringSettings returns Jolt's spring defaults: a 2 Hz spring with 1.0 damping.
func DefaultSpringSettings() SpringSettings {
	return SpringSettings{Mode: SpringModeFrequencyAndDamping, FrequencyOrStiffness: 2, Damping: 1}
}

// MotorSettings configures a constraint motor (JPH_MotorSettings).
type MotorSettings struct {
	SpringSettings SpringSettings
	MinForceLimit  float32
	MaxForceLimit  float32
	MinTorqueLimit float32
	MaxTorqueLimit float32
}

// DefaultMotorSettings returns Jolt's motor defaults: a 2 Hz spring with
// unlimited force and torque. Note the zero value has *zero* limits (no force),
// so start from this when writing motor settings.
func DefaultMotorSettings() MotorSettings {
	return MotorSettings{
		SpringSettings: DefaultSpringSettings(),
		MinForceLimit:  -math.MaxFloat32,
		MaxForceLimit:  math.MaxFloat32,
		MinTorqueLimit: -math.MaxFloat32,
		MaxTorqueLimit: math.MaxFloat32,
	}
}
