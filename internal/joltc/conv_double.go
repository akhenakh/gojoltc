//go:build cgo && jolt_double

package joltc

/*
#include <joltc.h>
*/
import "C"

import "github.com/akhenakh/gojoltc/internal/vecmath"

// In double precision, JPH_RVec3 holds three doubles and JPH_RMat4 is three
// float32 columns plus a double translation column.

func rvec3ToC(v vecmath.RVec3) C.JPH_RVec3 {
	return C.JPH_RVec3{x: C.double(v.X), y: C.double(v.Y), z: C.double(v.Z)}
}

func rvec3FromC(c C.JPH_RVec3) vecmath.RVec3 {
	return vecmath.RVec3{X: float64(c.x), Y: float64(c.y), Z: float64(c.z)}
}

func rmat4ToC(m vecmath.RMat4) C.JPH_RMat4 {
	var out C.JPH_RMat4
	for i := 0; i < 3; i++ {
		out.column[i] = vec4ToC(m.Col[i])
	}
	out.column3 = rvec3ToC(m.Col3)
	return out
}

func rmat4FromC(c C.JPH_RMat4) vecmath.RMat4 {
	var out vecmath.RMat4
	for i := 0; i < 3; i++ {
		out.Col[i] = vec4FromC(c.column[i])
	}
	out.Col3 = rvec3FromC(c.column3)
	return out
}
