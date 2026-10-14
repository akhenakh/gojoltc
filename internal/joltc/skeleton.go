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

// Skeleton is a hierarchy of joints.
type Skeleton struct{ ptr *C.JPH_Skeleton }

// NewSkeleton creates an empty skeleton.
func NewSkeleton() *Skeleton { return &Skeleton{ptr: C.JPH_Skeleton_Create()} }

// Destroy releases the skeleton.
func (s *Skeleton) Destroy() {
	if s.ptr != nil {
		C.JPH_Skeleton_Destroy(s.ptr)
		s.ptr = nil
	}
}

// AddJoint appends a root joint and returns its index.
func (s *Skeleton) AddJoint(name string) uint32 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint32(C.JPH_Skeleton_AddJoint(s.ptr, cname))
}

// AddJointWithParentIndex appends a joint under parentIndex and returns its index.
func (s *Skeleton) AddJointWithParentIndex(name string, parentIndex int) uint32 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return uint32(C.JPH_Skeleton_AddJoint2(s.ptr, cname, C.int(parentIndex)))
}

// AddJointWithParentName appends a joint under the named parent and returns its index.
func (s *Skeleton) AddJointWithParentName(name, parentName string) uint32 {
	cname := C.CString(name)
	cparent := C.CString(parentName)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cparent))
	return uint32(C.JPH_Skeleton_AddJoint3(s.ptr, cname, cparent))
}

// JointCount returns the number of joints.
func (s *Skeleton) JointCount() int { return int(C.JPH_Skeleton_GetJointCount(s.ptr)) }

// Joint returns the name and parent index of the joint at index.
func (s *Skeleton) Joint(index int) (name string, parentJointIndex int) {
	var j C.JPH_SkeletonJoint
	C.JPH_Skeleton_GetJoint(s.ptr, C.int(index), &j)
	return C.GoString(j.name), int(j.parentJointIndex)
}

// JointIndex returns the index of the named joint, or -1 if absent.
func (s *Skeleton) JointIndex(name string) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return int(C.JPH_Skeleton_GetJointIndex(s.ptr, cname))
}

// AreJointsCorrectlyOrdered reports whether parents precede children.
func (s *Skeleton) AreJointsCorrectlyOrdered() bool {
	return bool(C.JPH_Skeleton_AreJointsCorrectlyOrdered(s.ptr))
}

// Equal reports whether s and o refer to the same underlying skeleton.
func (s *Skeleton) Equal(o *Skeleton) bool { return s.ptr == o.ptr }

// SkeletonPose is a pose for a skeleton.
type SkeletonPose struct{ ptr *C.JPH_SkeletonPose }

// NewSkeletonPose creates an empty pose.
func NewSkeletonPose() *SkeletonPose { return &SkeletonPose{ptr: C.JPH_SkeletonPose_Create()} }

// Destroy releases the pose.
func (p *SkeletonPose) Destroy() {
	if p.ptr != nil {
		C.JPH_SkeletonPose_Destroy(p.ptr)
		p.ptr = nil
	}
}

// SetSkeleton binds a skeleton to the pose.
func (p *SkeletonPose) SetSkeleton(s *Skeleton) { C.JPH_SkeletonPose_SetSkeleton(p.ptr, s.ptr) }

// Skeleton returns the pose's skeleton.
func (p *SkeletonPose) Skeleton() *Skeleton {
	return &Skeleton{ptr: C.JPH_SkeletonPose_GetSkeleton(p.ptr)}
}

// SetRootOffset sets the root offset.
func (p *SkeletonPose) SetRootOffset(offset vecmath.RVec3) {
	o := rvec3ToC(offset)
	C.JPH_SkeletonPose_SetRootOffset(p.ptr, &o)
}

// RootOffset returns the root offset.
func (p *SkeletonPose) RootOffset() vecmath.RVec3 {
	var o C.JPH_RVec3
	C.JPH_SkeletonPose_GetRootOffset(p.ptr, &o)
	return rvec3FromC(o)
}

// JointCount returns the number of joints in the pose.
func (p *SkeletonPose) JointCount() int { return int(C.JPH_SkeletonPose_GetJointCount(p.ptr)) }

// JointState returns the translation and rotation of the joint at index.
func (p *SkeletonPose) JointState(index int) (vecmath.Vec3, vecmath.Quat) {
	var t C.JPH_Vec3
	var r C.JPH_Quat
	C.JPH_SkeletonPose_GetJointState(p.ptr, C.int(index), &t, &r)
	return vec3FromC(t), quatFromC(r)
}

// SetJointState sets the translation and rotation of the joint at index.
func (p *SkeletonPose) SetJointState(index int, translation vecmath.Vec3, rotation vecmath.Quat) {
	t := vec3ToC(translation)
	r := quatToC(rotation)
	C.JPH_SkeletonPose_SetJointState(p.ptr, C.int(index), &t, &r)
}

// SkeletonMapper maps poses between two skeletons.
type SkeletonMapper struct{ ptr *C.JPH_SkeletonMapper }

// NewSkeletonMapper creates an empty mapper.
func NewSkeletonMapper() *SkeletonMapper { return &SkeletonMapper{ptr: C.JPH_SkeletonMapper_Create()} }

// Destroy releases the mapper.
func (m *SkeletonMapper) Destroy() {
	if m.ptr != nil {
		C.JPH_SkeletonMapper_Destroy(m.ptr)
		m.ptr = nil
	}
}

// RagdollSettings holds settings for a ragdoll.
type RagdollSettings struct{ ptr *C.JPH_RagdollSettings }

// NewRagdollSettings creates empty ragdoll settings.
func NewRagdollSettings() *RagdollSettings {
	return &RagdollSettings{ptr: C.JPH_RagdollSettings_Create()}
}

// Destroy releases the settings.
func (s *RagdollSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_RagdollSettings_Destroy(s.ptr)
		s.ptr = nil
	}
}

// SetSkeleton binds a skeleton to the settings.
func (s *RagdollSettings) SetSkeleton(skeleton *Skeleton) {
	C.JPH_RagdollSettings_SetSkeleton(s.ptr, skeleton.ptr)
}

// Skeleton returns the settings' skeleton.
func (s *RagdollSettings) Skeleton() *Skeleton {
	return &Skeleton{ptr: (*C.JPH_Skeleton)(unsafe.Pointer(C.JPH_RagdollSettings_GetSkeleton(s.ptr)))}
}

// ResizeParts resizes the part list.
func (s *RagdollSettings) ResizeParts(count int) {
	C.JPH_RagdollSettings_ResizeParts(s.ptr, C.int(count))
}

// PartCount returns the number of parts.
func (s *RagdollSettings) PartCount() int { return int(C.JPH_RagdollSettings_GetPartCount(s.ptr)) }

// SetPartShape sets the shape of a part.
func (s *RagdollSettings) SetPartShape(index int, shape *Shape) {
	C.JPH_RagdollSettings_SetPartShape(s.ptr, C.int(index), shape.ptr)
}

// SetPartPosition sets the position of a part.
func (s *RagdollSettings) SetPartPosition(index int, position vecmath.RVec3) {
	p := rvec3ToC(position)
	C.JPH_RagdollSettings_SetPartPosition(s.ptr, C.int(index), &p)
}

// SetPartRotation sets the rotation of a part.
func (s *RagdollSettings) SetPartRotation(index int, rotation vecmath.Quat) {
	r := quatToC(rotation)
	C.JPH_RagdollSettings_SetPartRotation(s.ptr, C.int(index), &r)
}

// SetPartMotionType sets the motion type of a part.
func (s *RagdollSettings) SetPartMotionType(index int, motionType MotionType) {
	C.JPH_RagdollSettings_SetPartMotionType(s.ptr, C.int(index), C.JPH_MotionType(motionType))
}

// SetPartObjectLayer sets the object layer of a part.
func (s *RagdollSettings) SetPartObjectLayer(index int, layer ObjectLayer) {
	C.JPH_RagdollSettings_SetPartObjectLayer(s.ptr, C.int(index), C.JPH_ObjectLayer(layer))
}

// CreateRagdoll creates a ragdoll in the given physics system.
func (s *RagdollSettings) CreateRagdoll(system *PhysicsSystem, collisionGroup CollisionGroupID, userData uint64) *Ragdoll {
	return &Ragdoll{ptr: C.JPH_RagdollSettings_CreateRagdoll(s.ptr, system.ptr, C.JPH_CollisionGroupID(collisionGroup), C.uint64_t(userData))}
}

// Ragdoll is a physics-driven skeleton.
type Ragdoll struct{ ptr *C.JPH_Ragdoll }

// Destroy releases the ragdoll.
func (r *Ragdoll) Destroy() {
	if r.ptr != nil {
		C.JPH_Ragdoll_Destroy(r.ptr)
		r.ptr = nil
	}
}

// AddToPhysicsSystem adds the ragdoll to the physics system.
func (r *Ragdoll) AddToPhysicsSystem(activation Activation, lockBodies bool) {
	C.JPH_Ragdoll_AddToPhysicsSystem(r.ptr, C.JPH_Activation(activation), C.bool(lockBodies))
}

// RemoveFromPhysicsSystem removes the ragdoll from the physics system.
func (r *Ragdoll) RemoveFromPhysicsSystem(lockBodies bool) {
	C.JPH_Ragdoll_RemoveFromPhysicsSystem(r.ptr, C.bool(lockBodies))
}

// IsActive reports whether the ragdoll is active.
func (r *Ragdoll) IsActive(lockBodies bool) bool {
	return bool(C.JPH_Ragdoll_IsActive(r.ptr, C.bool(lockBodies)))
}

// BodyCount returns the number of bodies in the ragdoll.
func (r *Ragdoll) BodyCount() int { return int(C.JPH_Ragdoll_GetBodyCount(r.ptr)) }

// BodyID returns the body identifier of the body at index.
func (r *Ragdoll) BodyID(index int) BodyID {
	return BodyID(C.JPH_Ragdoll_GetBodyID(r.ptr, C.int(index)))
}

// RootTransform returns the ragdoll root position and rotation.
func (r *Ragdoll) RootTransform(lockBodies bool) (vecmath.RVec3, vecmath.Quat) {
	var p C.JPH_RVec3
	var q C.JPH_Quat
	C.JPH_Ragdoll_GetRootTransform(r.ptr, &p, &q, C.bool(lockBodies))
	return rvec3FromC(p), quatFromC(q)
}
