package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Constraint types, re-exported from the bindings.
type (
	// Constraint is a physics joint between bodies.
	Constraint = joltc.Constraint
	// FixedConstraintSettings holds settings for a fixed constraint.
	FixedConstraintSettings = joltc.FixedConstraintSettings
	// DistanceConstraintSettings holds settings for a distance constraint.
	DistanceConstraintSettings = joltc.DistanceConstraintSettings
	// PointConstraintSettings holds settings for a point constraint.
	PointConstraintSettings = joltc.PointConstraintSettings
	// HingeConstraintSettings holds settings for a hinge constraint.
	HingeConstraintSettings = joltc.HingeConstraintSettings
	// SliderConstraintSettings holds settings for a slider constraint.
	SliderConstraintSettings = joltc.SliderConstraintSettings
	// ConeConstraintSettings holds settings for a cone constraint.
	ConeConstraintSettings = joltc.ConeConstraintSettings
	// SwingTwistConstraintSettings holds settings for a swing-twist constraint.
	SwingTwistConstraintSettings = joltc.SwingTwistConstraintSettings
	// SixDOFConstraintSettings holds settings for a six DOF constraint.
	SixDOFConstraintSettings = joltc.SixDOFConstraintSettings
	// GearConstraintSettings holds settings for a gear constraint.
	GearConstraintSettings = joltc.GearConstraintSettings
)

// NewFixedConstraintSettings creates fixed constraint settings.
func NewFixedConstraintSettings() *FixedConstraintSettings { return joltc.NewFixedConstraintSettings() }

// NewDistanceConstraintSettings creates distance constraint settings.
func NewDistanceConstraintSettings() *DistanceConstraintSettings {
	return joltc.NewDistanceConstraintSettings()
}

// NewPointConstraintSettings creates point constraint settings.
func NewPointConstraintSettings() *PointConstraintSettings { return joltc.NewPointConstraintSettings() }

// NewHingeConstraintSettings creates hinge constraint settings.
func NewHingeConstraintSettings() *HingeConstraintSettings { return joltc.NewHingeConstraintSettings() }

// NewSliderConstraintSettings creates slider constraint settings.
func NewSliderConstraintSettings() *SliderConstraintSettings {
	return joltc.NewSliderConstraintSettings()
}

// NewConeConstraintSettings creates cone constraint settings.
func NewConeConstraintSettings() *ConeConstraintSettings { return joltc.NewConeConstraintSettings() }

// NewSwingTwistConstraintSettings creates swing-twist constraint settings.
func NewSwingTwistConstraintSettings() *SwingTwistConstraintSettings {
	return joltc.NewSwingTwistConstraintSettings()
}

// NewSixDOFConstraintSettings creates six DOF constraint settings.
func NewSixDOFConstraintSettings() *SixDOFConstraintSettings {
	return joltc.NewSixDOFConstraintSettings()
}

// NewGearConstraintSettings creates gear constraint settings.
func NewGearConstraintSettings() *GearConstraintSettings { return joltc.NewGearConstraintSettings() }
