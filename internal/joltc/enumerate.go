//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import "unsafe"

// The *Into variants batch a full body enumeration into a single cgo call and
// write into a caller-owned slice, so hot loops stay allocation-free. The
// allocating variants are convenience wrappers for cold paths.

// Bodies returns the IDs of all bodies in the system.
func (p *PhysicsSystem) Bodies() []BodyID {
	n := int(C.JPH_PhysicsSystem_GetNumBodies(p.ptr))
	if n == 0 {
		return nil
	}
	ids := make([]BodyID, n)
	C.JPH_PhysicsSystem_GetBodies(p.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n))
	return ids
}

// BodyIDsInto fills dst with body IDs and returns the number written. It reads
// at most len(dst) IDs.
func (p *PhysicsSystem) BodyIDsInto(dst []BodyID) int {
	if len(dst) == 0 {
		return 0
	}
	n := int(C.JPH_PhysicsSystem_GetNumBodies(p.ptr))
	n = min(n, len(dst))
	if n == 0 {
		return 0
	}
	C.JPH_PhysicsSystem_GetBodies(p.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&dst[0])), C.uint32_t(n))
	return n
}

// ActiveBodies returns the active body IDs of the given type.
func (p *PhysicsSystem) ActiveBodies(t BodyType) []BodyID {
	n := int(C.JPH_PhysicsSystem_GetNumActiveBodies(p.ptr, C.JPH_BodyType(t)))
	if n == 0 {
		return nil
	}
	ids := make([]BodyID, n)
	C.JPH_PhysicsSystem_GetActiveBodies(p.ptr, C.JPH_BodyType(t), (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(n))
	return ids
}

// ActiveBodyIDsInto fills dst with the active body IDs of the given type and
// returns the number written. It reads at most len(dst) IDs.
func (p *PhysicsSystem) ActiveBodyIDsInto(t BodyType, dst []BodyID) int {
	if len(dst) == 0 {
		return 0
	}
	n := int(C.JPH_PhysicsSystem_GetNumActiveBodies(p.ptr, C.JPH_BodyType(t)))
	n = min(n, len(dst))
	if n == 0 {
		return 0
	}
	C.JPH_PhysicsSystem_GetActiveBodies(p.ptr, C.JPH_BodyType(t), (*C.JPH_BodyID)(unsafe.Pointer(&dst[0])), C.uint32_t(n))
	return n
}
