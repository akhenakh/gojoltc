//go:build cgo

package joltc

/*
#include <joltc.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// GroupFilter decides, per collision group, whether two bodies may collide.
// The table variant (NewGroupFilterTable) is what joltc exposes.
type GroupFilter struct{ ptr *C.JPH_GroupFilter }

func (f *GroupFilter) table() *C.JPH_GroupFilterTable {
	return (*C.JPH_GroupFilterTable)(unsafe.Pointer(f.ptr))
}

// NewGroupFilterTable creates a group filter table with numSubGroups sub-groups
// (0 sizes it automatically).
func NewGroupFilterTable(numSubGroups uint32) *GroupFilter {
	return &GroupFilter{ptr: (*C.JPH_GroupFilter)(unsafe.Pointer(C.JPH_GroupFilterTable_Create(C.uint32_t(numSubGroups))))}
}

// Destroy releases the filter.
func (f *GroupFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_GroupFilter_Destroy(f.ptr)
		f.ptr = nil
	}
}

// EnableCollision enables collisions between two sub-groups.
func (f *GroupFilter) EnableCollision(subGroup1, subGroup2 CollisionSubGroupID) {
	C.JPH_GroupFilterTable_EnableCollision(f.table(), C.JPH_CollisionSubGroupID(subGroup1), C.JPH_CollisionSubGroupID(subGroup2))
}

// DisableCollision disables collisions between two sub-groups.
func (f *GroupFilter) DisableCollision(subGroup1, subGroup2 CollisionSubGroupID) {
	C.JPH_GroupFilterTable_DisableCollision(f.table(), C.JPH_CollisionSubGroupID(subGroup1), C.JPH_CollisionSubGroupID(subGroup2))
}

// IsCollisionEnabled reports whether two sub-groups collide.
func (f *GroupFilter) IsCollisionEnabled(subGroup1, subGroup2 CollisionSubGroupID) bool {
	return bool(C.JPH_GroupFilterTable_IsCollisionEnabled(f.table(), C.JPH_CollisionSubGroupID(subGroup1), C.JPH_CollisionSubGroupID(subGroup2)))
}

func collisionGroupToC(g CollisionGroup) C.JPH_CollisionGroup {
	var out C.JPH_CollisionGroup
	if g.Filter != nil {
		out.groupFilter = (*C.JPH_GroupFilter)(unsafe.Pointer(g.Filter.ptr))
	}
	out.groupID = C.JPH_CollisionGroupID(g.GroupID)
	out.subGroupID = C.JPH_CollisionSubGroupID(g.SubGroupID)
	return out
}

func collisionGroupFromC(c C.JPH_CollisionGroup) CollisionGroup {
	var filter *GroupFilter
	if c.groupFilter != nil {
		filter = &GroupFilter{ptr: (*C.JPH_GroupFilter)(unsafe.Pointer(c.groupFilter))}
	}
	return CollisionGroup{
		Filter:     filter,
		GroupID:    CollisionGroupID(c.groupID),
		SubGroupID: CollisionSubGroupID(c.subGroupID),
	}
}

// CanCollide reports whether two collision groups may collide.
func (f *GroupFilter) CanCollide(a, b CollisionGroup) bool {
	ca := collisionGroupToC(a)
	cb := collisionGroupToC(b)
	return bool(C.JPH_GroupFilter_CanCollide(f.ptr, &ca, &cb))
}
