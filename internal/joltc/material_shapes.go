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

// Materials can be attached to shapes where joltc exposes it (plane and tapered
// cylinder settings). Mesh and height field material indices exist in their
// triangle structs, but joltc has no material-list setter, so per-triangle
// friction cannot be assigned through this API; Material() still reads back
// whatever a shape reports.

// NewPlaneShapeWithMaterial creates a plane shape with a physics material.
func NewPlaneShapeWithMaterial(plane vecmath.Plane, material *PhysicsMaterial, halfExtent float32) *Shape {
	p := planeToC(plane)
	var m *C.JPH_PhysicsMaterial
	if material != nil {
		m = material.ptr
	}
	return &Shape{ptr: (*C.JPH_Shape)(unsafe.Pointer(C.JPH_PlaneShape_Create(&p, m, C.float(halfExtent))))}
}

// NewTaperedCylinderShapeSettingsWithMaterial creates tapered cylinder settings
// with a physics material.
func NewTaperedCylinderShapeSettingsWithMaterial(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius float32, material *PhysicsMaterial) *TaperedCylinderShapeSettings {
	var m *C.JPH_PhysicsMaterial
	if material != nil {
		m = material.ptr
	}
	return &TaperedCylinderShapeSettings{shapeSettings{ptr: (*C.JPH_ShapeSettings)(unsafe.Pointer(
		C.JPH_TaperedCylinderShapeSettings_Create(C.float(halfHeightOfTaperedCylinder), C.float(topRadius), C.float(bottomRadius), C.float(convexRadius), m)))}}
}

// Material returns the material of a sub-shape, or nil.
func (s *Shape) Material(subShapeID SubShapeID) *PhysicsMaterial {
	ptr := C.JPH_Shape_GetMaterial(s.ptr, C.JPH_SubShapeID(subShapeID))
	if ptr == nil {
		return nil
	}
	return &PhysicsMaterial{ptr: ptr}
}

// HeightFieldMaterial returns the material of a height field cell, or nil.
func (s *Shape) HeightFieldMaterial(x, y uint32) *PhysicsMaterial {
	ptr := C.JPH_HeightFieldShape_GetMaterial((*C.JPH_HeightFieldShape)(unsafe.Pointer(s.ptr)), C.uint32_t(x), C.uint32_t(y))
	if ptr == nil {
		return nil
	}
	return &PhysicsMaterial{ptr: ptr}
}
