//go:build jolt_double

package vecmath

// Real is the scalar type used for world-space positions. Single precision by
// default; the jolt_double build tag selects float64.
type Real = float64

// RVec3 is a world-space position (Jolt's RVec3). In double precision it is a
// 3-component float64 vector.
type RVec3 struct {
	X, Y, Z Real
}

// RMat4 is a world-space matrix (Jolt's RMat44): three float32 columns plus a
// float64 translation column.
type RMat4 struct {
	Col  [3]Vec4
	Col3 RVec3
}

// ToRVec3 converts a local-space Vec3 to a world-space RVec3.
func (v Vec3) ToRVec3() RVec3 { return RVec3{float64(v.X), float64(v.Y), float64(v.Z)} }

// ToVec3 converts a world-space RVec3 to a local-space Vec3.
func (r RVec3) ToVec3() Vec3 { return Vec3{float32(r.X), float32(r.Y), float32(r.Z)} }
