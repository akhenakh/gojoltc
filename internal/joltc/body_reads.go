//go:build cgo

package joltc

/*
#include <joltc.h>

extern void joltc_body_positions(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_RVec3* out);
extern void joltc_body_rotations(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_Quat* out);
extern void joltc_body_linear_velocities(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_Vec3* out);
extern void joltc_body_position_and_rotation(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_RVec3* outPos, JPH_Quat* outRot);
extern void joltc_body_set_positions(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_RVec3* in, JPH_Activation activation);
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// Compile-time checks that the Go math types are byte-compatible with the C
// structs that the batched reads write into caller-owned slices. If either
// layout drifts, the build fails rather than silently corrupting memory.
var (
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Vec3{}) - unsafe.Sizeof(C.JPH_Vec3{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Vec4{}) - unsafe.Sizeof(C.JPH_Vec4{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Quat{}) - unsafe.Sizeof(C.JPH_Quat{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Mat4{}) - unsafe.Sizeof(C.JPH_Mat4{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.RVec3{}) - unsafe.Sizeof(C.JPH_RVec3{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.RMat4{}) - unsafe.Sizeof(C.JPH_RMat4{}))]struct{}{}
)

// The *Into methods read or write many bodies in a single cgo call (the C loop
// lives in shim.c) without allocating: results go into caller-owned slices.
// Use these in per-frame loops instead of the per-call accessors, which
// heap-allocate their C out-parameter and cross cgo once per body.

// PositionsInto writes the positions of ids into dst and returns the number
// written (min(len(ids), len(dst))).
func (bi *BodyInterface) PositionsInto(ids []BodyID, dst []vecmath.RVec3) int {
	n := min(len(ids), len(dst))
	if n == 0 {
		return 0
	}
	C.joltc_body_positions(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n), (*C.JPH_RVec3)(unsafe.Pointer(&dst[0])))
	return n
}

// RotationsInto writes the rotations of ids into dst and returns the number
// written (min(len(ids), len(dst))).
func (bi *BodyInterface) RotationsInto(ids []BodyID, dst []vecmath.Quat) int {
	n := min(len(ids), len(dst))
	if n == 0 {
		return 0
	}
	C.joltc_body_rotations(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n), (*C.JPH_Quat)(unsafe.Pointer(&dst[0])))
	return n
}

// LinearVelocitiesInto writes the linear velocities of ids into dst and returns
// the number written (min(len(ids), len(dst))).
func (bi *BodyInterface) LinearVelocitiesInto(ids []BodyID, dst []vecmath.Vec3) int {
	n := min(len(ids), len(dst))
	if n == 0 {
		return 0
	}
	C.joltc_body_linear_velocities(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n), (*C.JPH_Vec3)(unsafe.Pointer(&dst[0])))
	return n
}

// PositionsAndRotationsInto writes the pose of ids into dstP and dstR and
// returns the number written (min(len(ids), len(dstP), len(dstR))).
func (bi *BodyInterface) PositionsAndRotationsInto(ids []BodyID, dstP []vecmath.RVec3, dstR []vecmath.Quat) int {
	n := min(len(ids), len(dstP), len(dstR))
	if n == 0 {
		return 0
	}
	C.joltc_body_position_and_rotation(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n), (*C.JPH_RVec3)(unsafe.Pointer(&dstP[0])), (*C.JPH_Quat)(unsafe.Pointer(&dstR[0])))
	return n
}

// SetPositionsInto sets the positions of ids from src and returns the number
// written (min(len(ids), len(src))).
func (bi *BodyInterface) SetPositionsInto(ids []BodyID, src []vecmath.RVec3, activation Activation) int {
	n := min(len(ids), len(src))
	if n == 0 {
		return 0
	}
	C.joltc_body_set_positions(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n), (*C.JPH_RVec3)(unsafe.Pointer(&src[0])), C.JPH_Activation(activation))
	return n
}
