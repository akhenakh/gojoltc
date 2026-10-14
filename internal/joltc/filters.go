//go:build cgo

package joltc

/*
#include <joltc.h>

extern void joltcInstallBroadPhaseLayerFilterProcs(void);
extern void joltcInstallObjectLayerFilterProcs(void);
extern void joltcInstallBodyFilterProcs(void);
extern void joltcInstallShapeFilterProcs(void);
extern void* joltc_voidp(uintptr_t v);
*/
import "C"

import (
	"runtime/cgo"
	"sync"
	"unsafe"
)

// Filters restrict which bodies a query may touch. Every filter is a Go
// callback dispatched through a cgo.Handle; a nil callback accepts everything.
// Callbacks may be invoked from any thread while the simulation is stepping.

// BroadPhaseLayerFilter is a broad-phase layer filter.
type BroadPhaseLayerFilter struct {
	ptr    *C.JPH_BroadPhaseLayerFilter
	handle cgo.Handle
}

var bpLayerFilterProcsOnce sync.Once

// NewBroadPhaseLayerFilter creates a broad-phase layer filter.
func NewBroadPhaseLayerFilter(cb *BroadPhaseLayerFilterCallbacks) *BroadPhaseLayerFilter {
	bpLayerFilterProcsOnce.Do(func() { C.joltcInstallBroadPhaseLayerFilterProcs() })
	h := cgo.NewHandle(cb)
	return &BroadPhaseLayerFilter{ptr: C.JPH_BroadPhaseLayerFilter_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the filter.
func (f *BroadPhaseLayerFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_BroadPhaseLayerFilter_Destroy(f.ptr)
		f.ptr = nil
		f.handle.Delete()
	}
}

//export goJoltcBroadPhaseLayerFilterShouldCollide
func goJoltcBroadPhaseLayerFilterShouldCollide(userData unsafe.Pointer, layer C.JPH_BroadPhaseLayer) bool {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*BroadPhaseLayerFilterCallbacks)
	if !ok || cb.ShouldCollide == nil {
		return true
	}
	return cb.ShouldCollide(BroadPhaseLayer(layer))
}

// ObjectLayerFilter is an object layer filter.
type ObjectLayerFilter struct {
	ptr    *C.JPH_ObjectLayerFilter
	handle cgo.Handle
}

var objLayerFilterProcsOnce sync.Once

// NewObjectLayerFilter creates an object layer filter.
func NewObjectLayerFilter(cb *ObjectLayerFilterCallbacks) *ObjectLayerFilter {
	objLayerFilterProcsOnce.Do(func() { C.joltcInstallObjectLayerFilterProcs() })
	h := cgo.NewHandle(cb)
	return &ObjectLayerFilter{ptr: C.JPH_ObjectLayerFilter_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the filter.
func (f *ObjectLayerFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_ObjectLayerFilter_Destroy(f.ptr)
		f.ptr = nil
		f.handle.Delete()
	}
}

//export goJoltcObjectLayerFilterShouldCollide
func goJoltcObjectLayerFilterShouldCollide(userData unsafe.Pointer, layer C.JPH_ObjectLayer) bool {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ObjectLayerFilterCallbacks)
	if !ok || cb.ShouldCollide == nil {
		return true
	}
	return cb.ShouldCollide(ObjectLayer(layer))
}

// BodyFilter is a body filter.
type BodyFilter struct {
	ptr    *C.JPH_BodyFilter
	handle cgo.Handle
}

var bodyFilterProcsOnce sync.Once

// NewBodyFilter creates a body filter.
func NewBodyFilter(cb *BodyFilterCallbacks) *BodyFilter {
	bodyFilterProcsOnce.Do(func() { C.joltcInstallBodyFilterProcs() })
	h := cgo.NewHandle(cb)
	return &BodyFilter{ptr: C.JPH_BodyFilter_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the filter.
func (f *BodyFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_BodyFilter_Destroy(f.ptr)
		f.ptr = nil
		f.handle.Delete()
	}
}

//export goJoltcBodyFilterShouldCollide
func goJoltcBodyFilterShouldCollide(userData unsafe.Pointer, bodyID C.JPH_BodyID) bool {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*BodyFilterCallbacks)
	if !ok || cb.ShouldCollide == nil {
		return true
	}
	return cb.ShouldCollide(BodyID(bodyID))
}

//export goJoltcBodyFilterShouldCollideLocked
func goJoltcBodyFilterShouldCollideLocked(userData unsafe.Pointer, body *C.JPH_Body) bool {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*BodyFilterCallbacks)
	if !ok || cb.ShouldCollideLocked == nil {
		return true
	}
	return cb.ShouldCollideLocked(&Body{ptr: body})
}

// ShapeFilter is a shape filter.
type ShapeFilter struct {
	ptr    *C.JPH_ShapeFilter
	handle cgo.Handle
}

var shapeFilterProcsOnce sync.Once

// NewShapeFilter creates a shape filter.
func NewShapeFilter(cb *ShapeFilterCallbacks) *ShapeFilter {
	shapeFilterProcsOnce.Do(func() { C.joltcInstallShapeFilterProcs() })
	h := cgo.NewHandle(cb)
	return &ShapeFilter{ptr: C.JPH_ShapeFilter_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the filter.
func (f *ShapeFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_ShapeFilter_Destroy(f.ptr)
		f.ptr = nil
		f.handle.Delete()
	}
}

//export goJoltcShapeFilterShouldCollide
func goJoltcShapeFilterShouldCollide(userData unsafe.Pointer, shape2 *C.JPH_Shape, subShapeID2 *C.JPH_SubShapeID) bool {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ShapeFilterCallbacks)
	if !ok || cb.ShouldCollide == nil {
		return true
	}
	var id SubShapeID
	if subShapeID2 != nil {
		id = SubShapeID(*subShapeID2)
	}
	return cb.ShouldCollide(&Shape{ptr: shape2}, id)
}

//export goJoltcShapeFilterShouldCollide2
func goJoltcShapeFilterShouldCollide2(userData unsafe.Pointer, shape1 *C.JPH_Shape, subShapeID1 *C.JPH_SubShapeID, shape2 *C.JPH_Shape, subShapeID2 *C.JPH_SubShapeID) bool {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ShapeFilterCallbacks)
	if !ok || cb.ShouldCollide2 == nil {
		return true
	}
	var id1, id2 SubShapeID
	if subShapeID1 != nil {
		id1 = SubShapeID(*subShapeID1)
	}
	if subShapeID2 != nil {
		id2 = SubShapeID(*subShapeID2)
	}
	return cb.ShouldCollide2(&Shape{ptr: shape1}, id1, &Shape{ptr: shape2}, id2)
}
