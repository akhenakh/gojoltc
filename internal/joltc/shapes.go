//go:build cgo

package joltc

/*
#include <joltc.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// shapeSettings is the common base for the concrete shape settings types. Its
// methods are promoted to those types.
type shapeSettings struct {
	ptr *C.JPH_ShapeSettings
}

// Destroy releases the settings.
func (s *shapeSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_ShapeSettings_Destroy(s.ptr)
		s.ptr = nil
	}
}

// UserData returns the user data attached to the settings.
func (s *shapeSettings) UserData() uint64 { return uint64(C.JPH_ShapeSettings_GetUserData(s.ptr)) }

// SetUserData attaches user data to the settings.
func (s *shapeSettings) SetUserData(v uint64) { C.JPH_ShapeSettings_SetUserData(s.ptr, C.uint64_t(v)) }

// Density returns the density of convex shape settings.
func (s *shapeSettings) Density() float32 {
	return float32(C.JPH_ConvexShapeSettings_GetDensity((*C.JPH_ConvexShapeSettings)(unsafe.Pointer(s.ptr))))
}

// SetDensity sets the density of convex shape settings.
func (s *shapeSettings) SetDensity(v float32) {
	C.JPH_ConvexShapeSettings_SetDensity((*C.JPH_ConvexShapeSettings)(unsafe.Pointer(s.ptr)), C.float(v))
}

func shapeFrom(ptr *C.JPH_Shape) *Shape { return &Shape{ptr: ptr} }

// --- Shapes -----------------------------------------------------------------

// UserData returns the user data attached to the shape.
func (s *Shape) UserData() uint64 { return uint64(C.JPH_Shape_GetUserData(s.ptr)) }

// SetUserData attaches user data to the shape.
func (s *Shape) SetUserData(v uint64) { C.JPH_Shape_SetUserData(s.ptr, C.uint64_t(v)) }

// MustBeStatic reports whether the shape can only be used on static bodies.
func (s *Shape) MustBeStatic() bool { return bool(C.JPH_Shape_MustBeStatic(s.ptr)) }

// CenterOfMass returns the shape's center of mass.
func (s *Shape) CenterOfMass() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_Shape_GetCenterOfMass(s.ptr, &o)
	return vec3FromC(o)
}

// LocalBounds returns the shape's local-space bounding box.
func (s *Shape) LocalBounds() vecmath.AABox {
	var o C.JPH_AABox
	C.JPH_Shape_GetLocalBounds(s.ptr, &o)
	return aaboxFromC(o)
}

// InnerRadius returns the shape's inner radius.
func (s *Shape) InnerRadius() float32 { return float32(C.JPH_Shape_GetInnerRadius(s.ptr)) }

// MassProperties returns the shape's mass properties.
func (s *Shape) MassProperties() MassProperties {
	var o C.JPH_MassProperties
	C.JPH_Shape_GetMassProperties(s.ptr, &o)
	return MassProperties{Mass: float32(o.mass), Inertia: mat4FromC(o.inertia)}
}

// Volume returns the shape's volume.
func (s *Shape) Volume() float32 { return float32(C.JPH_Shape_GetVolume(s.ptr)) }

// IsValidScale reports whether scale is a valid scale for the shape.
func (s *Shape) IsValidScale(scale vecmath.Vec3) bool {
	c := vec3ToC(scale)
	return bool(C.JPH_Shape_IsValidScale(s.ptr, &c))
}

// MakeScaleValid returns a scale that is valid for the shape.
func (s *Shape) MakeScaleValid(scale vecmath.Vec3) vecmath.Vec3 {
	c := vec3ToC(scale)
	var o C.JPH_Vec3
	C.JPH_Shape_MakeScaleValid(s.ptr, &c, &o)
	return vec3FromC(o)
}

// CastRay casts a ray (origin, direction, with direction length = ray length)
// against the shape. It returns whether the ray hit and the hit details.
func (s *Shape) CastRay(origin, direction vecmath.Vec3) (RayCastResult, bool) {
	co := vec3ToC(origin)
	cd := vec3ToC(direction)
	var hit C.JPH_RayCastResult
	ok := bool(C.JPH_Shape_CastRay(s.ptr, &co, &cd, &hit))
	if !ok {
		return RayCastResult{}, false
	}
	return RayCastResult{
		BodyID:      BodyID(hit.bodyID),
		Fraction:    float32(hit.fraction),
		SubShapeID2: SubShapeID(hit.subShapeID2),
	}, true
}

// CollidePoint reports whether point is inside the shape.
func (s *Shape) CollidePoint(point vecmath.Vec3) bool {
	c := vec3ToC(point)
	return bool(C.JPH_Shape_CollidePoint(s.ptr, &c, nil))
}

// ConvexShape density accessors operate on the base Shape.

// Density returns the density of a convex shape.
func (s *Shape) Density() float32 {
	return float32(C.JPH_ConvexShape_GetDensity((*C.JPH_ConvexShape)(unsafe.Pointer(s.ptr))))
}

// SetDensity sets the density of a convex shape.
func (s *Shape) SetDensity(v float32) {
	C.JPH_ConvexShape_SetDensity((*C.JPH_ConvexShape)(unsafe.Pointer(s.ptr)), C.float(v))
}

// --- Convex shape constructors ----------------------------------------------

// NewBoxShape creates a box shape.
func NewTriangleShape(v1, v2, v3 vecmath.Vec3, convexRadius float32) *Shape {
	c1, c2, c3 := vec3ToC(v1), vec3ToC(v2), vec3ToC(v3)
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_TriangleShape_Create(&c1, &c2, &c3, C.float(convexRadius)))))
}

// NewPlaneShape creates a plane shape. material may be nil.
func NewPlaneShape(plane vecmath.Plane, halfExtent float32) *Shape {
	c := planeToC(plane)
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_PlaneShape_Create(&c, nil, C.float(halfExtent)))))
}

// NewRotatedTranslatedShape creates a shape decorated with a translation and rotation.
func NewRotatedTranslatedShape(position vecmath.Vec3, rotation vecmath.Quat, inner *Shape) *Shape {
	p := vec3ToC(position)
	r := quatToC(rotation)
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_RotatedTranslatedShape_Create(&p, &r, inner.ptr))))
}

// NewScaledShape creates a shape decorated with a scale.
func NewScaledShape(inner *Shape, scale vecmath.Vec3) *Shape {
	s := vec3ToC(scale)
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_ScaledShape_Create(inner.ptr, &s))))
}

// NewOffsetCenterOfMassShape creates a shape decorated with a center of mass offset.
func NewOffsetCenterOfMassShape(offset vecmath.Vec3, inner *Shape) *Shape {
	o := vec3ToC(offset)
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_OffsetCenterOfMassShape_Create(&o, inner.ptr))))
}

// --- Settings ---------------------------------------------------------------

// SphereShapeSettings holds settings for a sphere shape.
type SphereShapeSettings struct{ shapeSettings }

// NewSphereShapeSettings creates sphere shape settings.
func NewSphereShapeSettings(radius float32) *SphereShapeSettings {
	return &SphereShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_SphereShapeSettings_Create(C.float(radius))))}}
}

// Radius returns the sphere radius.
func (s *SphereShapeSettings) Radius() float32 {
	return float32(C.JPH_SphereShapeSettings_GetRadius((*C.JPH_SphereShapeSettings)(unsafe.Pointer(s.ptr))))
}

// SetRadius sets the sphere radius.
func (s *SphereShapeSettings) SetRadius(r float32) {
	C.JPH_SphereShapeSettings_SetRadius((*C.JPH_SphereShapeSettings)(unsafe.Pointer(s.ptr)), C.float(r))
}

// CreateShape builds the sphere shape.
func (s *SphereShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_SphereShapeSettings_CreateShape((*C.JPH_SphereShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// BoxShapeSettings holds settings for a box shape.
type BoxShapeSettings struct{ shapeSettings }

// NewBoxShapeSettings creates box shape settings.
func NewBoxShapeSettings(halfExtent vecmath.Vec3, convexRadius float32) *BoxShapeSettings {
	c := vec3ToC(halfExtent)
	return &BoxShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_BoxShapeSettings_Create(&c, C.float(convexRadius))))}}
}

// CreateShape builds the box shape.
func (s *BoxShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_BoxShapeSettings_CreateShape((*C.JPH_BoxShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// CapsuleShapeSettings holds settings for a capsule shape.
type CapsuleShapeSettings struct{ shapeSettings }

// NewCapsuleShapeSettings creates capsule shape settings.
func NewCapsuleShapeSettings(halfHeightOfCylinder, radius float32) *CapsuleShapeSettings {
	return &CapsuleShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_CapsuleShapeSettings_Create(C.float(halfHeightOfCylinder), C.float(radius))))}}
}

// CreateShape builds the capsule shape.
func (s *CapsuleShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_CapsuleShapeSettings_CreateShape((*C.JPH_CapsuleShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// CylinderShapeSettings holds settings for a cylinder shape.
type CylinderShapeSettings struct{ shapeSettings }

// NewCylinderShapeSettings creates cylinder shape settings.
func NewCylinderShapeSettings(halfHeight, radius, convexRadius float32) *CylinderShapeSettings {
	return &CylinderShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_CylinderShapeSettings_Create(C.float(halfHeight), C.float(radius), C.float(convexRadius))))}}
}

// CreateShape builds the cylinder shape.
func (s *CylinderShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_CylinderShapeSettings_CreateShape((*C.JPH_CylinderShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// TaperedCapsuleShapeSettings holds settings for a tapered capsule shape.
type TaperedCapsuleShapeSettings struct{ shapeSettings }

// NewTaperedCapsuleShapeSettings creates tapered capsule shape settings.
func NewTaperedCapsuleShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius float32) *TaperedCapsuleShapeSettings {
	return &TaperedCapsuleShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_TaperedCapsuleShapeSettings_Create(C.float(halfHeightOfTaperedCylinder), C.float(topRadius), C.float(bottomRadius))))}}
}

// CreateShape builds the tapered capsule shape.
func (s *TaperedCapsuleShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_TaperedCapsuleShapeSettings_CreateShape((*C.JPH_TaperedCapsuleShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// TaperedCylinderShapeSettings holds settings for a tapered cylinder shape.
type TaperedCylinderShapeSettings struct{ shapeSettings }

// NewTaperedCylinderShapeSettings creates tapered cylinder shape settings.
func NewTaperedCylinderShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius float32) *TaperedCylinderShapeSettings {
	return &TaperedCylinderShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_TaperedCylinderShapeSettings_Create(C.float(halfHeightOfTaperedCylinder), C.float(topRadius), C.float(bottomRadius), C.float(convexRadius), nil)))}}
}

// CreateShape builds the tapered cylinder shape.
func (s *TaperedCylinderShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_TaperedCylinderShapeSettings_CreateShape((*C.JPH_TaperedCylinderShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// ConvexHullShapeSettings holds settings for a convex hull shape.
type ConvexHullShapeSettings struct{ shapeSettings }

// NewConvexHullShapeSettings creates convex hull shape settings from points.
func NewConvexHullShapeSettings(points []vecmath.Vec3, maxConvexRadius float32) *ConvexHullShapeSettings {
	var pptr *C.JPH_Vec3
	if n := len(points); n > 0 {
		pptr = (*C.JPH_Vec3)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.JPH_Vec3{}))))
		defer C.free(unsafe.Pointer(pptr))
		arr := unsafe.Slice(pptr, n)
		for i, p := range points {
			arr[i] = vec3ToC(p)
		}
	}
	return &ConvexHullShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_ConvexHullShapeSettings_Create(pptr, C.uint32_t(len(points)), C.float(maxConvexRadius))))}}
}

// CreateShape builds the convex hull shape.
func (s *ConvexHullShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_ConvexHullShapeSettings_CreateShape((*C.JPH_ConvexHullShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// NumPoints returns the number of points in the hull.
func (s *Shape) NumPoints() uint32 {
	return uint32(C.JPH_ConvexHullShape_GetNumPoints((*C.JPH_ConvexHullShape)(unsafe.Pointer(s.ptr))))
}

// Point returns the hull point at index.
func (s *Shape) Point(index uint32) vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_ConvexHullShape_GetPoint((*C.JPH_ConvexHullShape)(unsafe.Pointer(s.ptr)), C.uint32_t(index), &o)
	return vec3FromC(o)
}

// NumFaces returns the number of faces in the hull.
func (s *Shape) NumFaces() uint32 {
	return uint32(C.JPH_ConvexHullShape_GetNumFaces((*C.JPH_ConvexHullShape)(unsafe.Pointer(s.ptr))))
}

// CompoundShapeSettings is the base for compound shape settings.
type CompoundShapeSettings struct{ shapeSettings }

// AddShape adds a child shape (as a shape) to the compound settings.
func (s *CompoundShapeSettings) AddShape(position vecmath.Vec3, rotation vecmath.Quat, child *Shape, userData uint32) {
	p := vec3ToC(position)
	r := quatToC(rotation)
	C.JPH_CompoundShapeSettings_AddShape2((*C.JPH_CompoundShapeSettings)(unsafe.Pointer(s.ptr)), &p, &r, child.ptr, C.uint32_t(userData))
}

// NumSubShapes returns the number of sub-shapes in a compound shape.
func (s *Shape) NumSubShapes() uint32 {
	return uint32(C.JPH_CompoundShape_GetNumSubShapes((*C.JPH_CompoundShape)(unsafe.Pointer(s.ptr))))
}

// StaticCompoundShapeSettings holds settings for a static compound shape.
type StaticCompoundShapeSettings struct{ CompoundShapeSettings }

// NewStaticCompoundShapeSettings creates empty static compound shape settings.
func NewStaticCompoundShapeSettings() *StaticCompoundShapeSettings {
	return &StaticCompoundShapeSettings{CompoundShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_StaticCompoundShapeSettings_Create()))}}}
}

// CreateShape builds the static compound shape.
func (s *StaticCompoundShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_StaticCompoundShape_Create((*C.JPH_StaticCompoundShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// MutableCompoundShapeSettings holds settings for a mutable compound shape.
type MutableCompoundShapeSettings struct{ CompoundShapeSettings }

// NewMutableCompoundShapeSettings creates empty mutable compound shape settings.
func NewMutableCompoundShapeSettings() *MutableCompoundShapeSettings {
	return &MutableCompoundShapeSettings{CompoundShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_MutableCompoundShapeSettings_Create()))}}}
}

// CreateShape builds the mutable compound shape.
func (s *MutableCompoundShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_MutableCompoundShape_Create((*C.JPH_MutableCompoundShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// AddShape adds a child shape to a mutable compound shape at index (or
// ^uint32(0) to append). It returns the index of the added shape.
func (s *Shape) AddShape(position vecmath.Vec3, rotation vecmath.Quat, child *Shape, userData uint32, index uint32) uint32 {
	p := vec3ToC(position)
	r := quatToC(rotation)
	return uint32(C.JPH_MutableCompoundShape_AddShape((*C.JPH_MutableCompoundShape)(unsafe.Pointer(s.ptr)), &p, &r, child.ptr, C.uint32_t(userData), C.uint32_t(index)))
}

// RemoveShape removes the child shape at index from a mutable compound shape.
func (s *Shape) RemoveShape(index uint32) {
	C.JPH_MutableCompoundShape_RemoveShape((*C.JPH_MutableCompoundShape)(unsafe.Pointer(s.ptr)), C.uint32_t(index))
}

// EmptyShapeSettings holds settings for an empty shape.
type EmptyShapeSettings struct{ shapeSettings }

// NewEmptyShapeSettings creates empty shape settings with the given center of mass.
func NewEmptyShapeSettings(centerOfMass vecmath.Vec3) *EmptyShapeSettings {
	c := vec3ToC(centerOfMass)
	return &EmptyShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(C.JPH_EmptyShapeSettings_Create(&c)))}}
}

// CreateShape builds the empty shape.
func (s *EmptyShapeSettings) CreateShape() *Shape {
	return shapeFrom((*C.JPH_Shape)(unsafe.Pointer(C.JPH_EmptyShapeSettings_CreateShape((*C.JPH_EmptyShapeSettings)(unsafe.Pointer(s.ptr))))))
}

// GetInnerShape returns the inner shape of a decorated shape. The returned
// shape is owned by the decorated shape.
func (s *Shape) GetInnerShape() *Shape {
	return &Shape{ptr: C.JPH_DecoratedShape_GetInnerShape((*C.JPH_DecoratedShape)(unsafe.Pointer(s.ptr)))}
}
