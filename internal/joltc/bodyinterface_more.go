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

// ActivateBodies activates several bodies at once.
func (bi *BodyInterface) ActivateBodies(ids []BodyID) {
	if len(ids) == 0 {
		return
	}
	C.JPH_BodyInterface_ActivateBodies(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(len(ids)))
}

// DeactivateBodies deactivates several bodies at once.
func (bi *BodyInterface) DeactivateBodies(ids []BodyID) {
	if len(ids) == 0 {
		return
	}
	C.JPH_BodyInterface_DeactivateBodies(bi.ptr, (*C.JPH_BodyID)(unsafe.Pointer(&ids[0])), C.uint32_t(len(ids)))
}

// AddLinearAndAngularVelocity adds to both velocities in one call.
func (bi *BodyInterface) AddLinearAndAngularVelocity(id BodyID, linear, angular vecmath.Vec3) {
	l := vec3ToC(linear)
	a := vec3ToC(angular)
	C.JPH_BodyInterface_AddLinearAndAngularVelocity(bi.ptr, C.JPH_BodyID(id), &l, &a)
}

// ApplyBuoyancyImpulse applies a buoyancy impulse and fluid drag, and reports
// whether the surface plane was hit.
func (bi *BodyInterface) ApplyBuoyancyImpulse(id BodyID, surfacePosition vecmath.RVec3, surfaceNormal vecmath.Vec3, buoyancy, linearDrag, angularDrag float32, fluidVelocity, gravity vecmath.Vec3, deltaTime float32) bool {
	p := rvec3ToC(surfacePosition)
	n := vec3ToC(surfaceNormal)
	fv := vec3ToC(fluidVelocity)
	g := vec3ToC(gravity)
	return bool(C.JPH_BodyInterface_ApplyBuoyancyImpulse(bi.ptr, C.JPH_BodyID(id), &p, &n,
		C.float(buoyancy), C.float(linearDrag), C.float(angularDrag), &fv, &g, C.float(deltaTime)))
}

// CollisionGroup returns the body's collision group.
func (bi *BodyInterface) CollisionGroup(id BodyID) CollisionGroup {
	var c C.JPH_CollisionGroup
	C.JPH_BodyInterface_GetCollisionGroup(bi.ptr, C.JPH_BodyID(id), &c)
	return collisionGroupFromC(c)
}

// SetCollisionGroup sets the body's collision group.
func (bi *BodyInterface) SetCollisionGroup(id BodyID, g CollisionGroup) {
	c := collisionGroupToC(g)
	C.JPH_BodyInterface_SetCollisionGroup(bi.ptr, C.JPH_BodyID(id), &c)
}

// CollisionGroup returns the settings' collision group.
func (s *BodyCreationSettings) CollisionGroup() CollisionGroup {
	var c C.JPH_CollisionGroup
	C.JPH_BodyCreationSettings_GetCollisionGroup(s.ptr, &c)
	return collisionGroupFromC(c)
}

// SetCollisionGroup sets the settings' collision group.
func (s *BodyCreationSettings) SetCollisionGroup(g CollisionGroup) {
	c := collisionGroupToC(g)
	C.JPH_BodyCreationSettings_SetCollisionGroup(s.ptr, &c)
}

// CreateBodyWithID creates a body with a caller-chosen id, without adding it.
func (bi *BodyInterface) CreateBodyWithID(id BodyID, settings *BodyCreationSettings) *Body {
	return &Body{ptr: C.JPH_BodyInterface_CreateBodyWithID(bi.ptr, C.JPH_BodyID(id), settings.ptr)}
}

// CreateBodyWithoutID creates a body with no id, without adding it. Give it an
// id later with AssignBodyID.
func (bi *BodyInterface) CreateBodyWithoutID(settings *BodyCreationSettings) *Body {
	return &Body{ptr: C.JPH_BodyInterface_CreateBodyWithoutID(bi.ptr, settings.ptr)}
}

// DestroyBodyWithoutID destroys a body created without an id.
func (bi *BodyInterface) DestroyBodyWithoutID(body *Body) {
	C.JPH_BodyInterface_DestroyBodyWithoutID(bi.ptr, body.ptr)
}

// AssignBodyID gives a previously unassigned body the next free id.
func (bi *BodyInterface) AssignBodyID(body *Body) bool {
	return bool(C.JPH_BodyInterface_AssignBodyID(bi.ptr, body.ptr))
}

// AssignBodyIDTo gives a previously unassigned body a specific id.
func (bi *BodyInterface) AssignBodyIDTo(body *Body, id BodyID) bool {
	return bool(C.JPH_BodyInterface_AssignBodyID2(bi.ptr, body.ptr, C.JPH_BodyID(id)))
}

// UnassignBodyID removes the body's id and returns the body.
func (bi *BodyInterface) UnassignBodyID(id BodyID) *Body {
	return &Body{ptr: C.JPH_BodyInterface_UnassignBodyID(bi.ptr, C.JPH_BodyID(id))}
}

// ActivateBodiesInAABox wakes all bodies in an AABB on the given object layer.
func (p *PhysicsSystem) ActivateBodiesInAABox(box vecmath.AABox, layer ObjectLayer) {
	c := aaboxToC(box)
	C.JPH_PhysicsSystem_ActivateBodiesInAABox(p.ptr, &c, C.JPH_ObjectLayer(layer))
}
