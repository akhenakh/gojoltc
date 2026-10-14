//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import "github.com/akhenakh/gojoltc/internal/vecmath"

// BroadPhaseQuery is a broad-phase collision query interface.
type BroadPhaseQuery struct {
	ptr *C.JPH_BroadPhaseQuery
}

// NarrowPhaseQuery is a narrow-phase collision query interface.
type NarrowPhaseQuery struct {
	ptr *C.JPH_NarrowPhaseQuery
}

// BroadPhaseQuery returns the broad-phase query interface.
func (p *PhysicsSystem) BroadPhaseQuery() *BroadPhaseQuery {
	return &BroadPhaseQuery{ptr: C.JPH_PhysicsSystem_GetBroadPhaseQuery(p.ptr)}
}

// NarrowPhaseQuery returns the locking narrow-phase query interface.
func (p *PhysicsSystem) NarrowPhaseQuery() *NarrowPhaseQuery {
	return &NarrowPhaseQuery{ptr: C.JPH_PhysicsSystem_GetNarrowPhaseQuery(p.ptr)}
}

// NarrowPhaseQueryNoLock returns the non-locking narrow-phase query interface.
func (p *PhysicsSystem) NarrowPhaseQueryNoLock() *NarrowPhaseQuery {
	return &NarrowPhaseQuery{ptr: C.JPH_PhysicsSystem_GetNarrowPhaseQueryNoLock(p.ptr)}
}

// CastRay casts a ray against the world. direction's length is the ray length.
// It returns the closest hit and whether anything was hit.
func (q *NarrowPhaseQuery) CastRay(origin vecmath.RVec3, direction vecmath.Vec3) (RayCastResult, bool) {
	o := rvec3ToC(origin)
	d := vec3ToC(direction)
	var hit C.JPH_RayCastResult
	ok := bool(C.JPH_NarrowPhaseQuery_CastRay(q.ptr, &o, &d, &hit, nil, nil, nil))
	if !ok {
		return RayCastResult{}, false
	}
	return RayCastResult{
		BodyID:      BodyID(hit.bodyID),
		Fraction:    float32(hit.fraction),
		SubShapeID2: SubShapeID(hit.subShapeID2),
	}, true
}
