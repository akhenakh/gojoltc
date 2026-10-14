package jolt

import "github.com/akhenakh/gojoltc/internal/vecmath"

// Math types are re-exported from the pure-Go internal/vecmath package so that
// callers use them through the public jolt package.
type (
	// Real is the world-space scalar type: float32 by default, float64 with
	// the jolt_double build tag.
	Real = vecmath.Real
	// Vec3 is a 3-component single-precision vector.
	Vec3 = vecmath.Vec3
	// Vec4 is a 4-component single-precision vector.
	Vec4 = vecmath.Vec4
	// Quat is a quaternion.
	Quat = vecmath.Quat
	// Mat4 is a 4x4 matrix.
	Mat4 = vecmath.Mat4
	// RVec3 is a world-space position.
	RVec3 = vecmath.RVec3
	// RMat4 is a world-space matrix.
	RMat4 = vecmath.RMat4
	// Point is a 2-component vector.
	Point = vecmath.Point
	// AABox is an axis-aligned bounding box.
	AABox = vecmath.AABox
	// Plane is a plane defined by a normal and a distance.
	Plane = vecmath.Plane
	// Color is a packed RGBA color.
	Color = vecmath.Color
)
