// Package vecmath holds the pure-Go math types shared by the public jolt
// package and the internal cgo bindings.
//
// It never imports "C". World-space types (RVec3, RMat4) and the Real scalar
// type depend on the precision build tag: single precision (default) aliases
// them to the float32 types, while the jolt_double tag makes them float64.
package vecmath

// Vec3 is a 3-component single-precision vector (Jolt's Vec3).
type Vec3 struct {
	X, Y, Z float32
}

// Vec4 is a 4-component single-precision vector (Jolt's Vec4).
type Vec4 struct {
	X, Y, Z, W float32
}

// Quat is a quaternion (Jolt's Quat), components in x, y, z, w order.
type Quat struct {
	X, Y, Z, W float32
}

// Mat4 is a 4x4 matrix stored as four column vectors (Jolt's Mat44).
type Mat4 struct {
	Col [4]Vec4
}

// Point is a 2-component vector (Jolt's Point).
type Point struct {
	X, Y float32
}

// AABox is an axis-aligned bounding box (Jolt's AABox).
type AABox struct {
	Min, Max Vec3
}

// Plane is a plane defined by a normal and a distance (Jolt's Plane).
type Plane struct {
	Normal   Vec3
	Distance float32
}

// Color is a packed RGBA color (Jolt's Color, a uint32).
type Color uint32
