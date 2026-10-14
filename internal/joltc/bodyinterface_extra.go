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

// PhysicsMaterial assigns friction/restitution/color to a shape or sub-shape.
type PhysicsMaterial struct{ ptr *C.JPH_PhysicsMaterial }

// NewPhysicsMaterial creates a material with a debug name and packed color.
func NewPhysicsMaterial(name string, color vecmath.Color) *PhysicsMaterial {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return &PhysicsMaterial{ptr: C.JPH_PhysicsMaterial_Create(cname, C.uint32_t(color))}
}

// Destroy releases the material.
func (m *PhysicsMaterial) Destroy() {
	if m.ptr != nil {
		C.JPH_PhysicsMaterial_Destroy(m.ptr)
		m.ptr = nil
	}
}

// DebugName returns the material's debug name (copied into Go).
func (m *PhysicsMaterial) DebugName() string {
	return C.GoString(C.JPH_PhysicsMaterial_GetDebugName(m.ptr))
}

// DebugColor returns the material's packed debug color.
func (m *PhysicsMaterial) DebugColor() vecmath.Color {
	return vecmath.Color(C.JPH_PhysicsMaterial_GetDebugColor(m.ptr))
}

// GetMaterial returns the material of a sub-shape of a body, or nil.
func (bi *BodyInterface) GetMaterial(id BodyID, subShapeID SubShapeID) *PhysicsMaterial {
	ptr := C.JPH_BodyInterface_GetMaterial(bi.ptr, C.JPH_BodyID(id), C.JPH_SubShapeID(subShapeID))
	if ptr == nil {
		return nil
	}
	return &PhysicsMaterial{ptr: ptr}
}

// GetShape returns the body's shape. The shape is borrowed.
func (bi *BodyInterface) GetShape(id BodyID) *Shape {
	return &Shape{ptr: C.JPH_BodyInterface_GetShape(bi.ptr, C.JPH_BodyID(id))}
}

// SetShape replaces the body's shape.
func (bi *BodyInterface) SetShape(id BodyID, shape *Shape, updateMassProperties bool, activation Activation) {
	C.JPH_BodyInterface_SetShape(bi.ptr, C.JPH_BodyID(id), shape.ptr, C.bool(updateMassProperties), C.JPH_Activation(activation))
}

// MoveKinematic moves a kinematic body towards a target pose over deltaTime.
func (bi *BodyInterface) MoveKinematic(id BodyID, targetPosition vecmath.RVec3, targetRotation vecmath.Quat, deltaTime float32) {
	p := rvec3ToC(targetPosition)
	r := quatToC(targetRotation)
	C.JPH_BodyInterface_MoveKinematic(bi.ptr, C.JPH_BodyID(id), &p, &r, C.float(deltaTime))
}

// SetPositionAndRotationWhenChanged sets the pose only when it changed.
func (bi *BodyInterface) SetPositionAndRotationWhenChanged(id BodyID, position vecmath.RVec3, rotation vecmath.Quat, activation Activation) {
	p := rvec3ToC(position)
	r := quatToC(rotation)
	C.JPH_BodyInterface_SetPositionAndRotationWhenChanged(bi.ptr, C.JPH_BodyID(id), &p, &r, C.JPH_Activation(activation))
}

// SetPositionRotationAndVelocity sets pose and both velocities in one call.
func (bi *BodyInterface) SetPositionRotationAndVelocity(id BodyID, position vecmath.RVec3, rotation vecmath.Quat, linear, angular vecmath.Vec3) {
	p := rvec3ToC(position)
	r := quatToC(rotation)
	l := vec3ToC(linear)
	a := vec3ToC(angular)
	C.JPH_BodyInterface_SetPositionRotationAndVelocity(bi.ptr, C.JPH_BodyID(id), &p, &r, &l, &a)
}

// SetLinearAndAngularVelocity sets both velocities in one call.
func (bi *BodyInterface) SetLinearAndAngularVelocity(id BodyID, linear, angular vecmath.Vec3) {
	l := vec3ToC(linear)
	a := vec3ToC(angular)
	C.JPH_BodyInterface_SetLinearAndAngularVelocity(bi.ptr, C.JPH_BodyID(id), &l, &a)
}

// GetLinearAndAngularVelocity returns both velocities in one call.
func (bi *BodyInterface) GetLinearAndAngularVelocity(id BodyID) (vecmath.Vec3, vecmath.Vec3) {
	var l, a C.JPH_Vec3
	C.JPH_BodyInterface_GetLinearAndAngularVelocity(bi.ptr, C.JPH_BodyID(id), &l, &a)
	return vec3FromC(l), vec3FromC(a)
}

// AddForceAndTorque applies a force and a torque for the current step.
func (bi *BodyInterface) AddForceAndTorque(id BodyID, force, torque vecmath.Vec3) {
	f := vec3ToC(force)
	t := vec3ToC(torque)
	C.JPH_BodyInterface_AddForceAndTorque(bi.ptr, C.JPH_BodyID(id), &f, &t)
}

// AddForceAtPoint applies a force at a world point, producing torque.
func (bi *BodyInterface) AddForceAtPoint(id BodyID, force vecmath.Vec3, point vecmath.RVec3) {
	f := vec3ToC(force)
	p := rvec3ToC(point)
	C.JPH_BodyInterface_AddForce2(bi.ptr, C.JPH_BodyID(id), &f, &p)
}

// AddImpulseAtPoint applies an impulse at a world point, producing torque.
func (bi *BodyInterface) AddImpulseAtPoint(id BodyID, impulse vecmath.Vec3, point vecmath.RVec3) {
	im := vec3ToC(impulse)
	p := rvec3ToC(point)
	C.JPH_BodyInterface_AddImpulse2(bi.ptr, C.JPH_BodyID(id), &im, &p)
}

// PointVelocity returns the velocity of the body at a world point.
func (bi *BodyInterface) PointVelocity(id BodyID, point vecmath.RVec3) vecmath.Vec3 {
	p := rvec3ToC(point)
	var out C.JPH_Vec3
	C.JPH_BodyInterface_GetPointVelocity(bi.ptr, C.JPH_BodyID(id), &p, &out)
	return vec3FromC(out)
}

// NotifyShapeChanged informs the body that its shape's center of mass changed.
func (bi *BodyInterface) NotifyShapeChanged(id BodyID, previousCenterOfMass vecmath.Vec3, updateMassProperties bool, activation Activation) {
	c := vec3ToC(previousCenterOfMass)
	C.JPH_BodyInterface_NotifyShapeChanged(bi.ptr, C.JPH_BodyID(id), &c, C.bool(updateMassProperties), C.JPH_Activation(activation))
}

// InvalidateContactCache drops the body's cached contacts.
func (bi *BodyInterface) InvalidateContactCache(id BodyID) {
	C.JPH_BodyInterface_InvalidateContactCache(bi.ptr, C.JPH_BodyID(id))
}

// ResetSleepTimer prevents the body from sleeping this step.
func (bi *BodyInterface) ResetSleepTimer(id BodyID) {
	C.JPH_BodyInterface_ResetSleepTimer(bi.ptr, C.JPH_BodyID(id))
}
