//go:build cgo

package joltc

/*
#include <joltc.h>

extern bool joltcNarrowPhaseCastRay(const JPH_NarrowPhaseQuery* q, const JPH_RVec3* origin, const JPH_Vec3* dir, const JPH_RayCastSettings* s, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf);
extern bool joltcNarrowPhaseCollidePoint(const JPH_NarrowPhaseQuery* q, const JPH_RVec3* point, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf);
extern bool joltcNarrowPhaseCollideShape(const JPH_NarrowPhaseQuery* q, const JPH_Shape* shape, const JPH_Vec3* scale, const JPH_RMat4* comTransform,
	const JPH_CollideShapeSettings* s, JPH_RVec3* baseOffset, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf);
extern bool joltcNarrowPhaseCastShape(const JPH_NarrowPhaseQuery* q, const JPH_Shape* shape, const JPH_RMat4* worldTransform, const JPH_Vec3* dir,
	const JPH_ShapeCastSettings* s, JPH_RVec3* baseOffset, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf);
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

func (o *QueryOptions) cfilters() (bp *C.JPH_BroadPhaseLayerFilter, ol *C.JPH_ObjectLayerFilter, bf *C.JPH_BodyFilter, sf *C.JPH_ShapeFilter) {
	if o == nil {
		return
	}
	if o.BroadPhaseLayer != nil {
		bp = o.BroadPhaseLayer.ptr
	}
	if o.ObjectLayer != nil {
		ol = o.ObjectLayer.ptr
	}
	if o.Body != nil {
		bf = o.Body.ptr
	}
	if o.Shape != nil {
		sf = o.Shape.ptr
	}
	return
}

func (s RayCastSettings) toC() C.JPH_RayCastSettings {
	return C.JPH_RayCastSettings{
		backFaceModeTriangles: C.JPH_BackFaceMode(s.BackFaceModeTriangles),
		backFaceModeConvex:    C.JPH_BackFaceMode(s.BackFaceModeConvex),
		treatConvexAsSolid:    C.bool(s.TreatConvexAsSolid),
	}
}

func (s CollideShapeSettings) toC() C.JPH_CollideShapeSettings {
	var c C.JPH_CollideShapeSettings
	c.base.activeEdgeMode = C.JPH_ActiveEdgeMode(ActiveEdgeModeCollideOnlyWithActive)
	c.base.collectFacesMode = C.JPH_CollectFacesMode(CollectFacesModeNoFaces)
	c.base.collisionTolerance = C.float(1.0e-4)
	c.base.penetrationTolerance = C.float(1.0e-4)
	c.maxSeparationDistance = C.float(s.MaxSeparationDistance)
	c.backFaceMode = C.JPH_BackFaceMode(s.BackFaceMode)
	return c
}

func (s ShapeCastSettings) toC() C.JPH_ShapeCastSettings {
	var c C.JPH_ShapeCastSettings
	c.base.activeEdgeMode = C.JPH_ActiveEdgeMode(ActiveEdgeModeCollideOnlyWithActive)
	c.base.collectFacesMode = C.JPH_CollectFacesMode(CollectFacesModeNoFaces)
	c.base.collisionTolerance = C.float(1.0e-4)
	c.base.penetrationTolerance = C.float(1.0e-4)
	c.backFaceModeTriangles = C.JPH_BackFaceMode(s.BackFaceModeTriangles)
	c.backFaceModeConvex = C.JPH_BackFaceMode(s.BackFaceModeConvex)
	c.useShrunkenShapeAndConvexRadius = C.bool(s.UseShrunkenShapeAndConvexRadius)
	c.returnDeepestPoint = C.bool(s.ReturnDeepestPoint)
	return c
}

// --- collectors --------------------------------------------------------------

type rayCollector struct {
	all        []RayCastResult
	closest    RayCastResult
	hasClosest bool
	collectAll bool
}

//export goJoltcCastRayCollector
func goJoltcCastRayCollector(userData unsafe.Pointer, r *C.JPH_RayCastResult) C.float {
	defer func() { _ = recover() }()
	c, ok := cgo.Handle(uintptr(userData)).Value().(*rayCollector)
	if !ok {
		return 1
	}
	hit := RayCastResult{BodyID: BodyID(r.bodyID), Fraction: float32(r.fraction), SubShapeID2: SubShapeID(r.subShapeID2)}
	if !c.hasClosest || hit.Fraction < c.closest.Fraction {
		c.closest, c.hasClosest = hit, true
	}
	if c.collectAll {
		c.all = append(c.all, hit)
		return 1 // no early out: keep collecting
	}
	return C.float(hit.Fraction) // early out at the closest hit so far
}

type pointCollector struct{ hits []CollidePointResult }

//export goJoltcCollidePointCollector
func goJoltcCollidePointCollector(userData unsafe.Pointer, r *C.JPH_CollidePointResult) C.float {
	defer func() { _ = recover() }()
	c, ok := cgo.Handle(uintptr(userData)).Value().(*pointCollector)
	if !ok {
		return 1
	}
	c.hits = append(c.hits, CollidePointResult{BodyID: BodyID(r.bodyID), SubShapeID2: SubShapeID(r.subShapeID2)})
	return 1
}

type shapeCollector struct{ hits []CollideShapeResult }

//export goJoltcCollideShapeCollector
func goJoltcCollideShapeCollector(userData unsafe.Pointer, r *C.JPH_CollideShapeResult) C.float {
	defer func() { _ = recover() }()
	c, ok := cgo.Handle(uintptr(userData)).Value().(*shapeCollector)
	if !ok {
		return 1
	}
	// joltc frees the face arrays as soon as this returns, so only copy scalars.
	c.hits = append(c.hits, collideShapeResultFromC(r))
	return 1
}

type castCollector struct{ hits []ShapeCastResult }

//export goJoltcCastShapeCollector
func goJoltcCastShapeCollector(userData unsafe.Pointer, r *C.JPH_ShapeCastResult) C.float {
	defer func() { _ = recover() }()
	c, ok := cgo.Handle(uintptr(userData)).Value().(*castCollector)
	if !ok {
		return 1
	}
	c.hits = append(c.hits, ShapeCastResult{
		ContactPointOn1:  vec3FromC(r.contactPointOn1),
		ContactPointOn2:  vec3FromC(r.contactPointOn2),
		PenetrationAxis:  vec3FromC(r.penetrationAxis),
		PenetrationDepth: float32(r.penetrationDepth),
		SubShapeID1:      SubShapeID(r.subShapeID1),
		SubShapeID2:      SubShapeID(r.subShapeID2),
		BodyID2:          BodyID(r.bodyID2),
		Fraction:         float32(r.fraction),
		IsBackFaceHit:    bool(r.isBackFaceHit),
	})
	return 1
}

// --- queries -----------------------------------------------------------------

// CastRayAll returns every hit along the ray, in the order reported by Jolt.
func (q *NarrowPhaseQuery) CastRayAll(origin vecmath.RVec3, direction vecmath.Vec3, settings RayCastSettings, opts *QueryOptions) []RayCastResult {
	o := rvec3ToC(origin)
	d := vec3ToC(direction)
	s := settings.toC()
	bp, ol, bf, sf := opts.cfilters()
	c := &rayCollector{collectAll: true}
	h := cgo.NewHandle(c)
	defer h.Delete()
	C.joltcNarrowPhaseCastRay(q.ptr, &o, &d, &s, C.uintptr_t(h), bp, ol, bf, sf)
	return c.all
}

// CastRayClosest returns the first hit along the ray.
func (q *NarrowPhaseQuery) CastRayClosest(origin vecmath.RVec3, direction vecmath.Vec3, settings RayCastSettings, opts *QueryOptions) (RayCastResult, bool) {
	o := rvec3ToC(origin)
	d := vec3ToC(direction)
	s := settings.toC()
	bp, ol, bf, sf := opts.cfilters()
	c := &rayCollector{}
	h := cgo.NewHandle(c)
	defer h.Delete()
	C.joltcNarrowPhaseCastRay(q.ptr, &o, &d, &s, C.uintptr_t(h), bp, ol, bf, sf)
	return c.closest, c.hasClosest
}

// CollidePoint returns the bodies whose shape contains the point.
func (q *NarrowPhaseQuery) CollidePoint(point vecmath.RVec3, opts *QueryOptions) []CollidePointResult {
	p := rvec3ToC(point)
	bp, ol, bf, sf := opts.cfilters()
	c := &pointCollector{}
	h := cgo.NewHandle(c)
	defer h.Delete()
	C.joltcNarrowPhaseCollidePoint(q.ptr, &p, C.uintptr_t(h), bp, ol, bf, sf)
	return c.hits
}

// CollideShape returns the bodies overlapping shape, placed at transform and
// scaled by scale. settings.MaxSeparationDistance lets the query report
// near-misses (a proximity volume).
func (q *NarrowPhaseQuery) CollideShape(shape *Shape, scale vecmath.Vec3, transform vecmath.RMat4, settings CollideShapeSettings, opts *QueryOptions) []CollideShapeResult {
	sc := vec3ToC(scale)
	t := rmat4ToC(transform)
	s := settings.toC()
	var baseOffset C.JPH_RVec3
	bp, ol, bf, sf := opts.cfilters()
	c := &shapeCollector{}
	h := cgo.NewHandle(c)
	defer h.Delete()
	C.joltcNarrowPhaseCollideShape(q.ptr, shape.ptr, &sc, &t, &s, &baseOffset, C.uintptr_t(h), bp, ol, bf, sf)
	return c.hits
}

// CastShape sweeps shape (placed at transform) along direction and returns the
// hits.
func (q *NarrowPhaseQuery) CastShape(shape *Shape, transform vecmath.RMat4, direction vecmath.Vec3, settings ShapeCastSettings, opts *QueryOptions) []ShapeCastResult {
	t := rmat4ToC(transform)
	d := vec3ToC(direction)
	s := settings.toC()
	var baseOffset C.JPH_RVec3
	bp, ol, bf, sf := opts.cfilters()
	c := &castCollector{}
	h := cgo.NewHandle(c)
	defer h.Delete()
	C.joltcNarrowPhaseCastShape(q.ptr, shape.ptr, &t, &d, &s, &baseOffset, C.uintptr_t(h), bp, ol, bf, sf)
	return c.hits
}
