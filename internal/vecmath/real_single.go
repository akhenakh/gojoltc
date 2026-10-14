//go:build !jolt_double

package vecmath

// Real is the scalar type used for world-space positions. Single precision by
// default; the jolt_double build tag selects float64.
type Real = float32

// RVec3 is a world-space position (Jolt's RVec3). In single precision it is
// identical to Vec3.
type RVec3 = Vec3

// RMat4 is a world-space matrix (Jolt's RMat44). In single precision it is
// identical to Mat4.
type RMat4 = Mat4

// ToRVec3 converts a local-space Vec3 to a world-space RVec3.
func (v Vec3) ToRVec3() RVec3 { return v }

// ToVec3 converts a world-space RVec3 to a local-space Vec3.
func (r RVec3) ToVec3() Vec3 { return r }
