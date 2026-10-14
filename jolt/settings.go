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

// NewTaperedCylinderShapeSettingsWithMaterial creates tapered cylinder settings
// with a physics material.
func NewTaperedCylinderShapeSettingsWithMaterial(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius float32, material *PhysicsMaterial) *TaperedCylinderShapeSettings {
	return joltc.NewTaperedCylinderShapeSettingsWithMaterial(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius, material)
}
