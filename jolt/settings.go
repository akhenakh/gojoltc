package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Collision group and solver-settings types, re-exported from the bindings.
type (
	// PhysicsSettings holds the solver tuning parameters of a physics system.
	PhysicsSettings = joltc.PhysicsSettings
	// CollisionGroup is a collision group id plus an optional filter.
	CollisionGroup = joltc.CollisionGroup
	// GroupFilter decides whether two collision groups may collide.
	GroupFilter = joltc.GroupFilter
	// SpringSettings configures a constraint spring.
	SpringSettings = joltc.SpringSettings
	// MotorSettings configures a constraint motor.
	MotorSettings = joltc.MotorSettings
	// LinearCurve is a piecewise-linear function (wheel friction curves).
	LinearCurve = joltc.LinearCurve
)

// NewGroupFilterTable creates a group filter table with numSubGroups sub-groups
// (0 sizes it automatically).
func NewGroupFilterTable(numSubGroups uint32) *GroupFilter {
	return joltc.NewGroupFilterTable(numSubGroups)
}

// NewPlaneShapeWithMaterial creates a plane shape with a physics material.
func NewPlaneShapeWithMaterial(plane Plane, material *PhysicsMaterial, halfExtent float32) *Shape {
	return joltc.NewPlaneShapeWithMaterial(plane, material, halfExtent)
}

// DefaultSpringSettings returns Jolt's spring defaults.
func DefaultSpringSettings() SpringSettings { return joltc.DefaultSpringSettings() }

// DefaultMotorSettings returns Jolt's motor defaults (a 2 Hz spring with
// unlimited force and torque).
func DefaultMotorSettings() MotorSettings { return joltc.DefaultMotorSettings() }

// NewLinearCurve creates an empty piecewise-linear curve.
func NewLinearCurve() *LinearCurve { return joltc.NewLinearCurve() }

// NewTaperedCylinderShapeSettingsWithMaterial creates tapered cylinder settings
// with a physics material.
func NewTaperedCylinderShapeSettingsWithMaterial(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius float32, material *PhysicsMaterial) *TaperedCylinderShapeSettings {
	return joltc.NewTaperedCylinderShapeSettingsWithMaterial(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius, material)
}
