//go:build cgo && !jolt_double

package joltc

/*
#include <joltc.h>
*/
import "C"

import "github.com/akhenakh/gojoltc/internal/vecmath"

// In single precision, JPH_RVec3/JPH_RMat4 are aliases of JPH_Vec3/JPH_Mat4.

func rvec3ToC(v vecmath.RVec3) C.JPH_RVec3   { return vec3ToC(v) }
func rvec3FromC(c C.JPH_RVec3) vecmath.RVec3 { return vec3FromC(c) }
func rmat4ToC(m vecmath.RMat4) C.JPH_RMat4   { return mat4ToC(m) }
func rmat4FromC(c C.JPH_RMat4) vecmath.RMat4 { return mat4FromC(c) }
