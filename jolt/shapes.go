package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Shape value types and settings, re-exported from the bindings.
type (
	// MassProperties holds a shape's mass and inertia tensor.
	MassProperties = joltc.MassProperties
	// RayCastResult describes a ray hit against a shape.
	RayCastResult = joltc.RayCastResult
	// SphereShapeSettings holds settings for a sphere shape.
	SphereShapeSettings = joltc.SphereShapeSettings
	// BoxShapeSettings holds settings for a box shape.
	BoxShapeSettings = joltc.BoxShapeSettings
	// CapsuleShapeSettings holds settings for a capsule shape.
	CapsuleShapeSettings = joltc.CapsuleShapeSettings
	// CylinderShapeSettings holds settings for a cylinder shape.
	CylinderShapeSettings = joltc.CylinderShapeSettings
	// TaperedCapsuleShapeSettings holds settings for a tapered capsule shape.
	TaperedCapsuleShapeSettings = joltc.TaperedCapsuleShapeSettings
	// TaperedCylinderShapeSettings holds settings for a tapered cylinder shape.
	TaperedCylinderShapeSettings = joltc.TaperedCylinderShapeSettings
	// ConvexHullShapeSettings holds settings for a convex hull shape.
	ConvexHullShapeSettings = joltc.ConvexHullShapeSettings
	// CompoundShapeSettings is the base for compound shape settings.
	CompoundShapeSettings = joltc.CompoundShapeSettings
	// StaticCompoundShapeSettings holds settings for a static compound shape.
	StaticCompoundShapeSettings = joltc.StaticCompoundShapeSettings
	// MutableCompoundShapeSettings holds settings for a mutable compound shape.
	MutableCompoundShapeSettings = joltc.MutableCompoundShapeSettings
	// EmptyShapeSettings holds settings for an empty shape.
	EmptyShapeSettings = joltc.EmptyShapeSettings
)

// NewMassProperties returns mass properties carrying a mass and an identity
// inertia. Pair it with OverrideMassPropertiesCalculateInertia to impose a real
// mass while Jolt still derives the inertia from the shape (and scales it to
// that mass); the inertia here is only used by
// OverrideMassPropertiesMassAndInertiaProvided.
func NewMassProperties(mass float32) MassProperties {
	return MassProperties{Mass: mass, Inertia: Mat4Identity()}
}

// NewTriangleShape creates a triangle shape.
func NewTriangleShape(v1, v2, v3 Vec3, convexRadius float32) *Shape {
	return joltc.NewTriangleShape(v1, v2, v3, convexRadius)
}

// NewPlaneShape creates a plane shape.
func NewPlaneShape(plane Plane, halfExtent float32) *Shape {
	return joltc.NewPlaneShape(plane, halfExtent)
}

// NewRotatedTranslatedShape creates a shape decorated with a translation and rotation.
func NewRotatedTranslatedShape(position Vec3, rotation Quat, inner *Shape) *Shape {
	return joltc.NewRotatedTranslatedShape(position, rotation, inner)
}

// NewScaledShape creates a shape decorated with a scale.
func NewScaledShape(inner *Shape, scale Vec3) *Shape { return joltc.NewScaledShape(inner, scale) }

// NewOffsetCenterOfMassShape creates a shape decorated with a center of mass offset.
func NewOffsetCenterOfMassShape(offset Vec3, inner *Shape) *Shape {
	return joltc.NewOffsetCenterOfMassShape(offset, inner)
}

// NewSphereShapeSettings creates sphere shape settings.
func NewSphereShapeSettings(radius float32) *SphereShapeSettings {
	return joltc.NewSphereShapeSettings(radius)
}

// NewBoxShapeSettings creates box shape settings.
func NewBoxShapeSettings(halfExtent Vec3, convexRadius float32) *BoxShapeSettings {
	return joltc.NewBoxShapeSettings(halfExtent, convexRadius)
}

// NewCapsuleShapeSettings creates capsule shape settings.
func NewCapsuleShapeSettings(halfHeightOfCylinder, radius float32) *CapsuleShapeSettings {
	return joltc.NewCapsuleShapeSettings(halfHeightOfCylinder, radius)
}

// NewCylinderShapeSettings creates cylinder shape settings.
func NewCylinderShapeSettings(halfHeight, radius, convexRadius float32) *CylinderShapeSettings {
	return joltc.NewCylinderShapeSettings(halfHeight, radius, convexRadius)
}

// NewTaperedCapsuleShapeSettings creates tapered capsule shape settings.
func NewTaperedCapsuleShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius float32) *TaperedCapsuleShapeSettings {
	return joltc.NewTaperedCapsuleShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius)
}

// NewTaperedCylinderShapeSettings creates tapered cylinder shape settings.
func NewTaperedCylinderShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius float32) *TaperedCylinderShapeSettings {
	return joltc.NewTaperedCylinderShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius)
}

// NewConvexHullShapeSettings creates convex hull shape settings from points.
func NewConvexHullShapeSettings(points []Vec3, maxConvexRadius float32) *ConvexHullShapeSettings {
	return joltc.NewConvexHullShapeSettings(points, maxConvexRadius)
}

// NewStaticCompoundShapeSettings creates empty static compound shape settings.
func NewStaticCompoundShapeSettings() *StaticCompoundShapeSettings {
	return joltc.NewStaticCompoundShapeSettings()
}

// NewMutableCompoundShapeSettings creates empty mutable compound shape settings.
func NewMutableCompoundShapeSettings() *MutableCompoundShapeSettings {
	return joltc.NewMutableCompoundShapeSettings()
}

// NewEmptyShapeSettings creates empty shape settings with the given center of mass.
func NewEmptyShapeSettings(centerOfMass Vec3) *EmptyShapeSettings {
	return joltc.NewEmptyShapeSettings(centerOfMass)
}
