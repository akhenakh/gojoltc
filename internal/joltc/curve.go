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

// Compile-time check that vecmath.Point matches JPH_Point (we pass slices).
var _ [0]struct{} = [int(unsafe.Sizeof(vecmath.Point{}) - unsafe.Sizeof(C.JPH_Point{}))]struct{}{}

// LinearCurve is a piecewise-linear function. Wheeled wheels use it for their
// longitudinal and lateral friction-vs-slip curves.
type LinearCurve struct {
	ptr *C.JPH_LinearCurve
	// borrowed marks curves owned by another object (a wheel's friction curve);
	// Destroy is a no-op for those.
	borrowed bool
}

// NewLinearCurve creates an empty curve.
func NewLinearCurve() *LinearCurve { return &LinearCurve{ptr: C.JPH_LinearCurve_Create()} }

// Destroy releases the curve (no-op for a borrowed curve).
func (c *LinearCurve) Destroy() {
	if c.ptr != nil && !c.borrowed {
		C.JPH_LinearCurve_Destroy(c.ptr)
	}
	c.ptr = nil
}

// Clear removes all points.
func (c *LinearCurve) Clear() { C.JPH_LinearCurve_Clear(c.ptr) }

// Reserve preallocates space for n points.
func (c *LinearCurve) Reserve(n uint32) { C.JPH_LinearCurve_Reserve(c.ptr, C.uint32_t(n)) }

// AddPoint appends a point.
func (c *LinearCurve) AddPoint(x, y float32) {
	C.JPH_LinearCurve_AddPoint(c.ptr, C.float(x), C.float(y))
}

// Sort orders the points by x.
func (c *LinearCurve) Sort() { C.JPH_LinearCurve_Sort(c.ptr) }

// MinX returns the smallest x.
func (c *LinearCurve) MinX() float32 { return float32(C.JPH_LinearCurve_GetMinX(c.ptr)) }

// MaxX returns the largest x.
func (c *LinearCurve) MaxX() float32 { return float32(C.JPH_LinearCurve_GetMaxX(c.ptr)) }

// Value evaluates the curve at x.
func (c *LinearCurve) Value(x float32) float32 {
	return float32(C.JPH_LinearCurve_GetValue(c.ptr, C.float(x)))
}

// PointCount returns the number of points.
func (c *LinearCurve) PointCount() uint32 { return uint32(C.JPH_LinearCurve_GetPointCount(c.ptr)) }

// Point returns the point at index.
func (c *LinearCurve) Point(index uint32) vecmath.Point {
	var p C.JPH_Point
	C.JPH_LinearCurve_GetPoint(c.ptr, C.uint32_t(index), &p)
	return vecmath.Point{X: float32(p.x), Y: float32(p.y)}
}

// Points copies up to len(dst) points into dst and returns the number written.
func (c *LinearCurve) Points(dst []vecmath.Point) int {
	if len(dst) == 0 {
		return 0
	}
	count := C.uint32_t(len(dst))
	C.JPH_LinearCurve_GetPoints(c.ptr, (*C.JPH_Point)(unsafe.Pointer(&dst[0])), &count)
	return int(count)
}

// SetLongitudinalFrictionCurve sets the tire's longitudinal friction curve. The
// curve is copied, so the caller may destroy it afterwards.
func (s *WheelSettings) SetLongitudinalFrictionCurve(c *LinearCurve) {
	C.JPH_WheelSettingsWV_SetLongitudinalFriction(s.wv(), c.ptr)
}

// LongitudinalFrictionCurve returns a borrowed view of the wheel's longitudinal
// friction curve. It is owned by the wheel settings; do not destroy it.
func (s *WheelSettings) LongitudinalFrictionCurve() *LinearCurve {
	return &LinearCurve{ptr: C.JPH_WheelSettingsWV_GetLongitudinalFriction(s.wv()), borrowed: true}
}

// SetLateralFrictionCurve sets the tire's lateral friction curve (copied).
func (s *WheelSettings) SetLateralFrictionCurve(c *LinearCurve) {
	C.JPH_WheelSettingsWV_SetLateralFriction(s.wv(), c.ptr)
}

// LateralFrictionCurve returns a borrowed view of the wheel's lateral friction
// curve. It is owned by the wheel settings; do not destroy it.
func (s *WheelSettings) LateralFrictionCurve() *LinearCurve {
	return &LinearCurve{ptr: C.JPH_WheelSettingsWV_GetLateralFriction(s.wv()), borrowed: true}
}
