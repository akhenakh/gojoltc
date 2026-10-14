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

// BoxHalfExtent returns the half extent of a box shape.
func (s *Shape) BoxHalfExtent() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_BoxShape_GetHalfExtent((*C.JPH_BoxShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}

// BoxConvexRadius returns the convex radius of a box shape.
func (s *Shape) BoxConvexRadius() float32 {
	return float32(C.JPH_BoxShape_GetConvexRadius((*C.JPH_BoxShape)(unsafe.Pointer(s.ptr))))
}

// SphereRadius returns the radius of a sphere shape.
func (s *Shape) SphereRadius() float32 {
	return float32(C.JPH_SphereShape_GetRadius((*C.JPH_SphereShape)(unsafe.Pointer(s.ptr))))
}

// CapsuleHalfHeightOfCylinder returns the capsule cylinder half height.
func (s *Shape) CapsuleHalfHeightOfCylinder() float32 {
	return float32(C.JPH_CapsuleShape_GetHalfHeightOfCylinder((*C.JPH_CapsuleShape)(unsafe.Pointer(s.ptr))))
}

// CapsuleRadius returns the capsule radius.
func (s *Shape) CapsuleRadius() float32 {
	return float32(C.JPH_CapsuleShape_GetRadius((*C.JPH_CapsuleShape)(unsafe.Pointer(s.ptr))))
}

// CylinderHalfHeight returns the cylinder half height.
func (s *Shape) CylinderHalfHeight() float32 {
	return float32(C.JPH_CylinderShape_GetHalfHeight((*C.JPH_CylinderShape)(unsafe.Pointer(s.ptr))))
}

// CylinderRadius returns the cylinder radius.
func (s *Shape) CylinderRadius() float32 {
	return float32(C.JPH_CylinderShape_GetRadius((*C.JPH_CylinderShape)(unsafe.Pointer(s.ptr))))
}

// TriangleVertex1 returns the first vertex of a triangle shape.
func (s *Shape) TriangleVertex1() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_TriangleShape_GetVertex1((*C.JPH_TriangleShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}

// TriangleVertex2 returns the second vertex of a triangle shape.
func (s *Shape) TriangleVertex2() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_TriangleShape_GetVertex2((*C.JPH_TriangleShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}

// TriangleVertex3 returns the third vertex of a triangle shape.
func (s *Shape) TriangleVertex3() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_TriangleShape_GetVertex3((*C.JPH_TriangleShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}

// Plane returns the plane of a plane shape.
func (s *Shape) Plane() vecmath.Plane {
	var o C.JPH_Plane
	C.JPH_PlaneShape_GetPlane((*C.JPH_PlaneShape)(unsafe.Pointer(s.ptr)), &o)
	return planeFromC(o)
}

// PlaneHalfExtent returns the half extent of a plane shape.
func (s *Shape) PlaneHalfExtent() float32 {
	return float32(C.JPH_PlaneShape_GetHalfExtent((*C.JPH_PlaneShape)(unsafe.Pointer(s.ptr))))
}

// TaperedCapsuleTopRadius returns the tapered capsule top radius.
func (s *Shape) TaperedCapsuleTopRadius() float32 {
	return float32(C.JPH_TaperedCapsuleShape_GetTopRadius((*C.JPH_TaperedCapsuleShape)(unsafe.Pointer(s.ptr))))
}

// TaperedCapsuleBottomRadius returns the tapered capsule bottom radius.
func (s *Shape) TaperedCapsuleBottomRadius() float32 {
	return float32(C.JPH_TaperedCapsuleShape_GetBottomRadius((*C.JPH_TaperedCapsuleShape)(unsafe.Pointer(s.ptr))))
}

// TaperedCapsuleHalfHeight returns the tapered capsule half height.
func (s *Shape) TaperedCapsuleHalfHeight() float32 {
	return float32(C.JPH_TaperedCapsuleShape_GetHalfHeight((*C.JPH_TaperedCapsuleShape)(unsafe.Pointer(s.ptr))))
}

// TaperedCylinderTopRadius returns the tapered cylinder top radius.
func (s *Shape) TaperedCylinderTopRadius() float32 {
	return float32(C.JPH_TaperedCylinderShape_GetTopRadius((*C.JPH_TaperedCylinderShape)(unsafe.Pointer(s.ptr))))
}

// TaperedCylinderBottomRadius returns the tapered cylinder bottom radius.
func (s *Shape) TaperedCylinderBottomRadius() float32 {
	return float32(C.JPH_TaperedCylinderShape_GetBottomRadius((*C.JPH_TaperedCylinderShape)(unsafe.Pointer(s.ptr))))
}

// TaperedCylinderHalfHeight returns the tapered cylinder half height.
func (s *Shape) TaperedCylinderHalfHeight() float32 {
	return float32(C.JPH_TaperedCylinderShape_GetHalfHeight((*C.JPH_TaperedCylinderShape)(unsafe.Pointer(s.ptr))))
}

// RotatedTranslatedPosition returns the position of a rotated/translated shape.
func (s *Shape) RotatedTranslatedPosition() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_RotatedTranslatedShape_GetPosition((*C.JPH_RotatedTranslatedShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}

// RotatedTranslatedRotation returns the rotation of a rotated/translated shape.
func (s *Shape) RotatedTranslatedRotation() vecmath.Quat {
	var o C.JPH_Quat
	C.JPH_RotatedTranslatedShape_GetRotation((*C.JPH_RotatedTranslatedShape)(unsafe.Pointer(s.ptr)), &o)
	return quatFromC(o)
}

// ScaledScale returns the scale of a scaled shape.
func (s *Shape) ScaledScale() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_ScaledShape_GetScale((*C.JPH_ScaledShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}

// OffsetCenterOfMassOffset returns the offset of an offset-center-of-mass shape.
func (s *Shape) OffsetCenterOfMassOffset() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_OffsetCenterOfMassShape_GetOffset((*C.JPH_OffsetCenterOfMassShape)(unsafe.Pointer(s.ptr)), &o)
	return vec3FromC(o)
}
