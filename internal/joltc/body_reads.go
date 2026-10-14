//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// Compile-time checks that the Go math types are byte-compatible with the C
// structs that Batched reads write into caller-owned slices. If either layout
// drifts, the build fails rather than silently corrupting memory.
var (
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Vec3{}) - unsafe.Sizeof(C.JPH_Vec3{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Vec4{}) - unsafe.Sizeof(C.JPH_Vec4{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Quat{}) - unsafe.Sizeof(C.JPH_Quat{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.Mat4{}) - unsafe.Sizeof(C.JPH_Mat4{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.RVec3{}) - unsafe.Sizeof(C.JPH_RVec3{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(vecmath.RMat4{}) - unsafe.Sizeof(C.JPH_RMat4{}))]struct{}{}
)

// The *Into methods read many bodies in one sweep without allocating: results
// are written into caller-owned slices, whose element addresses are passed to
// C directly. Use these in per-frame loops instead of the per-call accessors,
// which heap-allocate their C out-parameter.

// PositionsInto writes the positions of ids into dst and returns the number
// written (min(len(ids), len(dst))).
func (bi *BodyInterface) PositionsInto(ids []BodyID, dst []vecmath.RVec3) int {
	n := min(len(ids), len(dst))
	for i := range n {
		C.JPH_BodyInterface_GetPosition(bi.ptr, C.JPH_BodyID(ids[i]), (*C.JPH_RVec3)(unsafe.Pointer(&dst[i])))
	}
	return n
}

// RotationsInto writes the rotations of ids into dst and returns the number
// written (min(len(ids), len(dst))).
func (bi *BodyInterface) RotationsInto(ids []BodyID, dst []vecmath.Quat) int {
	n := min(len(ids), len(dst))
	for i := range n {
		C.JPH_BodyInterface_GetRotation(bi.ptr, C.JPH_BodyID(ids[i]), (*C.JPH_Quat)(unsafe.Pointer(&dst[i])))
	}
	return n
}

// LinearVelocitiesInto writes the linear velocities of ids into dst and returns
// the number written (min(len(ids), len(dst))).
func (bi *BodyInterface) LinearVelocitiesInto(ids []BodyID, dst []vecmath.Vec3) int {
	n := min(len(ids), len(dst))
	for i := range n {
		C.JPH_BodyInterface_GetLinearVelocity(bi.ptr, C.JPH_BodyID(ids[i]), (*C.JPH_Vec3)(unsafe.Pointer(&dst[i])))
	}
	return n
}

// PositionsAndRotationsInto writes the pose of ids into dstP and dstR and
// returns the number written (min(len(ids), len(dstP), len(dstR))).
func (bi *BodyInterface) PositionsAndRotationsInto(ids []BodyID, dstP []vecmath.RVec3, dstR []vecmath.Quat) int {
	n := min(len(ids), len(dstP), len(dstR))
	for i := range n {
		C.JPH_BodyInterface_GetPositionAndRotation(bi.ptr, C.JPH_BodyID(ids[i]), (*C.JPH_RVec3)(unsafe.Pointer(&dstP[i])), (*C.JPH_Quat)(unsafe.Pointer(&dstR[i])))
	}
	return n
}

// SetPositionsInto sets the positions of ids from src and returns the number
// written (min(len(ids), len(src))).
func (bi *BodyInterface) SetPositionsInto(ids []BodyID, src []vecmath.RVec3, activation Activation) int {
	n := min(len(ids), len(src))
	for i := range n {
		C.JPH_BodyInterface_SetPosition(bi.ptr, C.JPH_BodyID(ids[i]), (*C.JPH_RVec3)(unsafe.Pointer(&src[i])), C.JPH_Activation(activation))
	}
	return n
}
