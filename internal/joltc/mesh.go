//go:build cgo

package joltc

/*
#include <joltc.h>
#include <stdlib.h>
*/
import "C"

import (
	"math"
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// Compile-time checks that the Go mesh value types are byte-compatible with the
// C structs we pass slices of.
var (
	_ [0]struct{} = [int(unsafe.Sizeof(Triangle{}) - unsafe.Sizeof(C.JPH_Triangle{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(IndexedTriangle{}) - unsafe.Sizeof(C.JPH_IndexedTriangle{}))]struct{}{}
)

// NewMeshShapeSettings creates mesh shape settings from standalone triangles.
func NewMeshShapeSettings(triangles []Triangle) *Shape {
	if len(triangles) == 0 {
		return &Shape{ptr: (*C.JPH_Shape)(unsafe.Pointer(C.JPH_MeshShapeSettings_Create(nil, 0)))}
	}
	settings := C.JPH_MeshShapeSettings_Create(
		(*C.JPH_Triangle)(unsafe.Pointer(&triangles[0])), C.uint32_t(len(triangles)))
	return &Shape{ptr: (*C.JPH_Shape)(unsafe.Pointer(C.JPH_MeshShapeSettings_CreateShape(settings)))}
}

// MeshShapeSettingsFromVertices builds mesh settings from a shared vertex array
// and indexed triangles, then creates the shape.
func MeshShapeSettingsFromVertices(vertices []vecmath.Vec3, triangles []IndexedTriangle) *Shape {
	var vptr *C.JPH_Vec3
	if len(vertices) > 0 {
		vptr = (*C.JPH_Vec3)(unsafe.Pointer(&vertices[0]))
	}
	var tptr *C.JPH_IndexedTriangle
	if len(triangles) > 0 {
		tptr = (*C.JPH_IndexedTriangle)(unsafe.Pointer(&triangles[0]))
	}
	settings := C.JPH_MeshShapeSettings_Create2(vptr, C.uint32_t(len(vertices)), tptr, C.uint32_t(len(triangles)))
	C.JPH_MeshShapeSettings_Sanitize(settings)
	return &Shape{ptr: (*C.JPH_Shape)(unsafe.Pointer(C.JPH_MeshShapeSettings_CreateShape(settings)))}
}

// NewHeightFieldShapeSettings creates a height field shape from a square sample
// grid (len(samples) must be n*n) with the given offset and scale.
//
// Note: joltc's sampleCount parameter is the number of samples *per side* (n),
// not the total; passing the total makes Jolt read past the array.
func NewHeightFieldShapeSettings(samples []float32, offset, scale vecmath.Vec3) *Shape {
	if len(samples) == 0 {
		return nil
	}
	n := int(math.Sqrt(float64(len(samples))))
	if n*n != len(samples) {
		return nil
	}
	off := vec3ToC(offset)
	sc := vec3ToC(scale)
	settings := C.JPH_HeightFieldShapeSettings_Create(
		(*C.float)(unsafe.Pointer(&samples[0])), &off, &sc, C.uint32_t(n), nil)
	return &Shape{ptr: (*C.JPH_Shape)(unsafe.Pointer(C.JPH_HeightFieldShapeSettings_CreateShape(settings)))}
}
