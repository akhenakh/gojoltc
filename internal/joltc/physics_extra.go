//go:build cgo

package joltc

/*
#include <joltc.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// AddConstraint adds a constraint to the system. The system borrows it: keep it
// alive, and remove it before destroying it.
func (p *PhysicsSystem) AddConstraint(c *Constraint) {
	C.JPH_PhysicsSystem_AddConstraint(p.ptr, c.ptr)
}

// RemoveConstraint removes a constraint from the system.
func (p *PhysicsSystem) RemoveConstraint(c *Constraint) {
	C.JPH_PhysicsSystem_RemoveConstraint(p.ptr, c.ptr)
}

func constraintArray(cs []*Constraint) (*C.JPH_Constraint, func()) {
	if len(cs) == 0 {
		return nil, func() {}
	}
	buf := (*C.JPH_Constraint)(C.malloc(C.size_t(len(cs)) * C.size_t(unsafe.Sizeof(uintptr(0)))))
	arr := unsafe.Slice((**C.JPH_Constraint)(unsafe.Pointer(buf)), len(cs))
	for i, c := range cs {
		arr[i] = c.ptr
	}
	return buf, func() { C.free(unsafe.Pointer(buf)) }
}

// AddConstraints adds several constraints to the system in one call.
func (p *PhysicsSystem) AddConstraints(cs []*Constraint) {
	buf, free := constraintArray(cs)
	defer free()
	if buf == nil {
		return
	}
	C.JPH_PhysicsSystem_AddConstraints(p.ptr, (**C.JPH_Constraint)(unsafe.Pointer(buf)), C.uint32_t(len(cs)))
}

// RemoveConstraints removes several constraints from the system in one call.
func (p *PhysicsSystem) RemoveConstraints(cs []*Constraint) {
	buf, free := constraintArray(cs)
	defer free()
	if buf == nil {
		return
	}
	C.JPH_PhysicsSystem_RemoveConstraints(p.ptr, (**C.JPH_Constraint)(unsafe.Pointer(buf)), C.uint32_t(len(cs)))
}
