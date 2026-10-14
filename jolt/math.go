package jolt

import "github.com/akhenakh/gojoltc/internal/vecmath"

// Pi is the value of pi as float32.
const Pi = vecmath.Pi

// RayCast is a ray defined by an origin point and a direction vector.
type RayCast = vecmath.RayCast

// Vec3AxisX returns the unit X axis.
func Vec3AxisX() Vec3 { return vecmath.Vec3AxisX() }

// Vec3AxisY returns the unit Y axis.
func Vec3AxisY() Vec3 { return vecmath.Vec3AxisY() }

// Vec3AxisZ returns the unit Z axis.
func Vec3AxisZ() Vec3 { return vecmath.Vec3AxisZ() }

// QuatIdentity returns the identity quaternion.
func QuatIdentity() Quat { return vecmath.QuatIdentity() }

// QuatRotation returns a rotation of angle radians around a normalized axis.
func QuatRotation(axis Vec3, angle float32) Quat { return vecmath.QuatRotation(axis, angle) }

// QuatFromTo returns the rotation that takes from to to.
func QuatFromTo(from, to Vec3) Quat { return vecmath.QuatFromTo(from, to) }

// QuatEulerAngles builds a quaternion from Euler angles in radians.
func QuatEulerAngles(angles Vec3) Quat { return vecmath.QuatEulerAngles(angles) }

// Mat4Identity returns the identity matrix.
func Mat4Identity() Mat4 { return vecmath.Mat4Identity() }

// Mat4Zero returns the zero matrix.
func Mat4Zero() Mat4 { return vecmath.Mat4Zero() }

// Mat4Rotation returns a rotation matrix from a quaternion.
func Mat4Rotation(q Quat) Mat4 { return vecmath.Mat4Rotation(q) }

// Mat4RotationAxisAngle returns a rotation matrix of angle radians around axis.
func Mat4RotationAxisAngle(axis Vec3, angle float32) Mat4 {
	return vecmath.Mat4RotationAxisAngle(axis, angle)
}

// Mat4Translation returns a translation matrix.
func Mat4Translation(t Vec3) Mat4 { return vecmath.Mat4Translation(t) }

// Mat4RotationTranslation returns a rotation followed by a translation.
func Mat4RotationTranslation(q Quat, t Vec3) Mat4 { return vecmath.Mat4RotationTranslation(q, t) }

// Mat4InverseRotationTranslation returns the inverse of Mat4RotationTranslation.
func Mat4InverseRotationTranslation(q Quat, t Vec3) Mat4 {
	return vecmath.Mat4InverseRotationTranslation(q, t)
}

// Mat4Scale returns a scale matrix.
func Mat4Scale(s Vec3) Mat4 { return vecmath.Mat4Scale(s) }

// RayCastGetPointOnRay returns the point at fraction along the ray defined by
// origin and direction.
func RayCastGetPointOnRay(origin, direction Vec3, fraction float32) Vec3 {
	return vecmath.RayCastGetPointOnRay(origin, direction, fraction)
}
