//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import "github.com/akhenakh/gojoltc/internal/vecmath"

// DefaultConvexRadius mirrors JPH_DEFAULT_CONVEX_RADIUS.
const DefaultConvexRadius float32 = 0.05

// JobSystemConfig configures a thread-pool job system.
type JobSystemConfig struct {
	MaxJobs     uint32
	MaxBarriers uint32
	NumThreads  int32
}

// JobSystem is a Jolt job system used to step the simulation in parallel.
type JobSystem struct {
	ptr *C.JPH_JobSystem
}

// NewJobSystemThreadPool creates a job system backed by a thread pool. A nil
// config uses Jolt's defaults.
func NewJobSystemThreadPool(cfg *JobSystemConfig) *JobSystem {
	if cfg == nil {
		return &JobSystem{ptr: C.JPH_JobSystemThreadPool_Create(nil)}
	}
	c := C.JobSystemThreadPoolConfig{
		maxJobs:     C.uint32_t(cfg.MaxJobs),
		maxBarriers: C.uint32_t(cfg.MaxBarriers),
		numThreads:  C.int32_t(cfg.NumThreads),
	}
	return &JobSystem{ptr: C.JPH_JobSystemThreadPool_Create(&c)}
}

// Destroy releases the job system.
func (j *JobSystem) Destroy() {
	if j.ptr != nil {
		C.JPH_JobSystem_Destroy(j.ptr)
		j.ptr = nil
	}
}

// BroadPhaseLayerInterface maps object layers to broad-phase layers.
type BroadPhaseLayerInterface struct {
	ptr *C.JPH_BroadPhaseLayerInterface
}

// NewBroadPhaseLayerInterfaceTable creates a table-based mapping.
func NewBroadPhaseLayerInterfaceTable(numObjectLayers, numBroadPhaseLayers uint32) *BroadPhaseLayerInterface {
	return &BroadPhaseLayerInterface{
		ptr: C.JPH_BroadPhaseLayerInterfaceTable_Create(C.uint32_t(numObjectLayers), C.uint32_t(numBroadPhaseLayers)),
	}
}

// MapObjectToBroadPhaseLayer maps an object layer to a broad-phase layer.
func (b *BroadPhaseLayerInterface) MapObjectToBroadPhaseLayer(objectLayer ObjectLayer, broadPhaseLayer BroadPhaseLayer) {
	C.JPH_BroadPhaseLayerInterfaceTable_MapObjectToBroadPhaseLayer(b.ptr, C.JPH_ObjectLayer(objectLayer), C.JPH_BroadPhaseLayer(broadPhaseLayer))
}

// Destroy releases the interface.
func (b *BroadPhaseLayerInterface) Destroy() {
	if b.ptr != nil {
		C.JPH_BroadPhaseLayerInterface_Destroy(b.ptr)
		b.ptr = nil
	}
}

// ObjectLayerPairFilter decides which object layers may collide.
type ObjectLayerPairFilter struct {
	ptr *C.JPH_ObjectLayerPairFilter
}

// NewObjectLayerPairFilterTable creates a table-based object layer pair filter.
func NewObjectLayerPairFilterTable(numObjectLayers uint32) *ObjectLayerPairFilter {
	return &ObjectLayerPairFilter{ptr: C.JPH_ObjectLayerPairFilterTable_Create(C.uint32_t(numObjectLayers))}
}

// EnableCollision enables collisions between two object layers.
func (f *ObjectLayerPairFilter) EnableCollision(layer1, layer2 ObjectLayer) {
	C.JPH_ObjectLayerPairFilterTable_EnableCollision(f.ptr, C.JPH_ObjectLayer(layer1), C.JPH_ObjectLayer(layer2))
}

// DisableCollision disables collisions between two object layers.
func (f *ObjectLayerPairFilter) DisableCollision(layer1, layer2 ObjectLayer) {
	C.JPH_ObjectLayerPairFilterTable_DisableCollision(f.ptr, C.JPH_ObjectLayer(layer1), C.JPH_ObjectLayer(layer2))
}

// ShouldCollide reports whether two object layers collide.
func (f *ObjectLayerPairFilter) ShouldCollide(layer1, layer2 ObjectLayer) bool {
	return bool(C.JPH_ObjectLayerPairFilterTable_ShouldCollide(f.ptr, C.JPH_ObjectLayer(layer1), C.JPH_ObjectLayer(layer2)))
}

// Destroy releases the filter.
func (f *ObjectLayerPairFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_ObjectLayerPairFilter_Destroy(f.ptr)
		f.ptr = nil
	}
}

// ObjectVsBroadPhaseLayerFilter decides which objects are tested against which
// broad-phase layers.
type ObjectVsBroadPhaseLayerFilter struct {
	ptr *C.JPH_ObjectVsBroadPhaseLayerFilter
}

// NewObjectVsBroadPhaseLayerFilterTable creates a table-based filter from a
// broad-phase layer interface and an object layer pair filter.
func NewObjectVsBroadPhaseLayerFilterTable(bp *BroadPhaseLayerInterface, numBroadPhaseLayers uint32, pair *ObjectLayerPairFilter, numObjectLayers uint32) *ObjectVsBroadPhaseLayerFilter {
	return &ObjectVsBroadPhaseLayerFilter{
		ptr: C.JPH_ObjectVsBroadPhaseLayerFilterTable_Create(bp.ptr, C.uint32_t(numBroadPhaseLayers), pair.ptr, C.uint32_t(numObjectLayers)),
	}
}

// Destroy releases the filter.
func (f *ObjectVsBroadPhaseLayerFilter) Destroy() {
	if f.ptr != nil {
		C.JPH_ObjectVsBroadPhaseLayerFilter_Destroy(f.ptr)
		f.ptr = nil
	}
}

// PhysicsSystemSettings configures a physics system.
type PhysicsSystemSettings struct {
	MaxBodies                     uint32
	NumBodyMutexes                uint32
	MaxBodyPairs                  uint32
	MaxContactConstraints         uint32
	BroadPhaseLayerInterface      *BroadPhaseLayerInterface
	ObjectLayerPairFilter         *ObjectLayerPairFilter
	ObjectVsBroadPhaseLayerFilter *ObjectVsBroadPhaseLayerFilter
}

// PhysicsSystem is the main simulation object.
type PhysicsSystem struct {
	ptr *C.JPH_PhysicsSystem
}

// NewPhysicsSystem creates a physics system from settings.
//
// The system takes ownership of the broad-phase layer interface, the object
// layer pair filter, and the object vs broad-phase layer filter: the joltc
// destroy path deletes them. Those wrappers are cleared here so a later
// Destroy on them is a harmless no-op rather than a double free.
func NewPhysicsSystem(s *PhysicsSystemSettings) *PhysicsSystem {
	c := C.JPH_PhysicsSystemSettings{
		maxBodies:                     C.uint32_t(s.MaxBodies),
		numBodyMutexes:                C.uint32_t(s.NumBodyMutexes),
		maxBodyPairs:                  C.uint32_t(s.MaxBodyPairs),
		maxContactConstraints:         C.uint32_t(s.MaxContactConstraints),
		broadPhaseLayerInterface:      s.BroadPhaseLayerInterface.ptr,
		objectLayerPairFilter:         s.ObjectLayerPairFilter.ptr,
		objectVsBroadPhaseLayerFilter: s.ObjectVsBroadPhaseLayerFilter.ptr,
	}
	sys := &PhysicsSystem{ptr: C.JPH_PhysicsSystem_Create(&c)}
	if s.BroadPhaseLayerInterface != nil {
		s.BroadPhaseLayerInterface.ptr = nil
	}
	if s.ObjectLayerPairFilter != nil {
		s.ObjectLayerPairFilter.ptr = nil
	}
	if s.ObjectVsBroadPhaseLayerFilter != nil {
		s.ObjectVsBroadPhaseLayerFilter.ptr = nil
	}
	return sys
}

// Destroy releases the physics system.
func (p *PhysicsSystem) Destroy() {
	if p.ptr != nil {
		C.JPH_PhysicsSystem_Destroy(p.ptr)
		p.ptr = nil
	}
}

// BodyInterface returns the locking body interface.
func (p *PhysicsSystem) BodyInterface() *BodyInterface {
	return &BodyInterface{ptr: C.JPH_PhysicsSystem_GetBodyInterface(p.ptr)}
}

// BodyInterfaceNoLock returns the non-locking body interface.
func (p *PhysicsSystem) BodyInterfaceNoLock() *BodyInterface {
	return &BodyInterface{ptr: C.JPH_PhysicsSystem_GetBodyInterfaceNoLock(p.ptr)}
}

// SetGravity sets the world gravity.
func (p *PhysicsSystem) SetGravity(g vecmath.Vec3) {
	c := vec3ToC(g)
	C.JPH_PhysicsSystem_SetGravity(p.ptr, &c)
}

// Gravity returns the world gravity.
func (p *PhysicsSystem) Gravity() vecmath.Vec3 {
	var out C.JPH_Vec3
	C.JPH_PhysicsSystem_GetGravity(p.ptr, &out)
	return vec3FromC(out)
}

// NumBodies returns the number of bodies in the system.
func (p *PhysicsSystem) NumBodies() uint32 { return uint32(C.JPH_PhysicsSystem_GetNumBodies(p.ptr)) }

// NumActiveBodies returns the number of active bodies of the given type.
func (p *PhysicsSystem) NumActiveBodies(t BodyType) uint32 {
	return uint32(C.JPH_PhysicsSystem_GetNumActiveBodies(p.ptr, C.JPH_BodyType(t)))
}

// MaxBodies returns the configured body capacity.
func (p *PhysicsSystem) MaxBodies() uint32 { return uint32(C.JPH_PhysicsSystem_GetMaxBodies(p.ptr)) }

// OptimizeBroadPhase optimizes the broad-phase tree. Call after adding static
// geometry and before simulating.
func (p *PhysicsSystem) OptimizeBroadPhase() { C.JPH_PhysicsSystem_OptimizeBroadPhase(p.ptr) }

// Update advances the simulation by deltaTime seconds.
func (p *PhysicsSystem) Update(deltaTime float32, collisionSteps int, jobSystem *JobSystem) PhysicsUpdateError {
	var js *C.JPH_JobSystem
	if jobSystem != nil {
		js = jobSystem.ptr
	}
	return PhysicsUpdateError(C.JPH_PhysicsSystem_Update(p.ptr, C.float(deltaTime), C.int(collisionSteps), js))
}

// Shape is a collision shape (reference counted).
type Shape struct {
	ptr *C.JPH_Shape
}

// Destroy releases one reference to the shape.
func (s *Shape) Destroy() {
	if s.ptr != nil {
		C.JPH_Shape_Destroy(s.ptr)
		s.ptr = nil
	}
}

// Type returns the top-level shape category.
func (s *Shape) Type() ShapeType { return ShapeType(C.JPH_Shape_GetType(s.ptr)) }

// SubType returns the specific shape category.
func (s *Shape) SubType() ShapeSubType { return ShapeSubType(C.JPH_Shape_GetSubType(s.ptr)) }

// NewBoxShape creates a box shape.
func NewBoxShape(halfExtent vecmath.Vec3, convexRadius float32) *Shape {
	c := vec3ToC(halfExtent)
	return &Shape{ptr: (*C.JPH_Shape)(C.JPH_BoxShape_Create(&c, C.float(convexRadius)))}
}

// NewSphereShape creates a sphere shape.
func NewSphereShape(radius float32) *Shape {
	return &Shape{ptr: (*C.JPH_Shape)(C.JPH_SphereShape_Create(C.float(radius)))}
}

// NewCapsuleShape creates a capsule shape.
func NewCapsuleShape(halfHeightOfCylinder, radius float32) *Shape {
	return &Shape{ptr: (*C.JPH_Shape)(C.JPH_CapsuleShape_Create(C.float(halfHeightOfCylinder), C.float(radius)))}
}

// NewCylinderShape creates a cylinder shape.
func NewCylinderShape(halfHeight, radius float32) *Shape {
	return &Shape{ptr: (*C.JPH_Shape)(C.JPH_CylinderShape_Create(C.float(halfHeight), C.float(radius)))}
}

// BodyCreationSettings describes a body to be created.
type BodyCreationSettings struct {
	ptr *C.JPH_BodyCreationSettings
}

// NewBodyCreationSettings creates settings for a body with the given shape,
// position, rotation, motion type and object layer.
func NewBodyCreationSettings(shape *Shape, position vecmath.RVec3, rotation vecmath.Quat, motionType MotionType, layer ObjectLayer) *BodyCreationSettings {
	pos := rvec3ToC(position)
	rot := quatToC(rotation)
	return &BodyCreationSettings{
		ptr: C.JPH_BodyCreationSettings_Create3(shape.ptr, &pos, &rot, C.JPH_MotionType(motionType), C.JPH_ObjectLayer(layer)),
	}
}

// Destroy releases the settings.
func (b *BodyCreationSettings) Destroy() {
	if b.ptr != nil {
		C.JPH_BodyCreationSettings_Destroy(b.ptr)
		b.ptr = nil
	}
}

// SetFriction sets the friction coefficient.
func (b *BodyCreationSettings) SetFriction(v float32) {
	C.JPH_BodyCreationSettings_SetFriction(b.ptr, C.float(v))
}

// Friction returns the friction coefficient.
func (b *BodyCreationSettings) Friction() float32 {
	return float32(C.JPH_BodyCreationSettings_GetFriction(b.ptr))
}

// SetRestitution sets the restitution coefficient.
func (b *BodyCreationSettings) SetRestitution(v float32) {
	C.JPH_BodyCreationSettings_SetRestitution(b.ptr, C.float(v))
}

// Restitution returns the restitution coefficient.
func (b *BodyCreationSettings) Restitution() float32 {
	return float32(C.JPH_BodyCreationSettings_GetRestitution(b.ptr))
}

// SetLinearVelocity sets the initial linear velocity.
func (b *BodyCreationSettings) SetLinearVelocity(v vecmath.Vec3) {
	c := vec3ToC(v)
	C.JPH_BodyCreationSettings_SetLinearVelocity(b.ptr, &c)
}

// SetAngularVelocity sets the initial angular velocity.
func (b *BodyCreationSettings) SetAngularVelocity(v vecmath.Vec3) {
	c := vec3ToC(v)
	C.JPH_BodyCreationSettings_SetAngularVelocity(b.ptr, &c)
}

// Body is a rigid body instance.
type Body struct {
	ptr *C.JPH_Body
}

// ID returns the body's identifier.
func (b *Body) ID() BodyID { return BodyID(C.JPH_Body_GetID(b.ptr)) }

// BodyInterface creates and manipulates bodies in a physics system.
type BodyInterface struct {
	ptr *C.JPH_BodyInterface
}

// CreateBody creates a body without adding it to the system.
func (bi *BodyInterface) CreateBody(settings *BodyCreationSettings) *Body {
	return &Body{ptr: C.JPH_BodyInterface_CreateBody(bi.ptr, settings.ptr)}
}

// CreateAndAddBody creates a body and immediately adds it to the system.
func (bi *BodyInterface) CreateAndAddBody(settings *BodyCreationSettings, activation Activation) BodyID {
	return BodyID(C.JPH_BodyInterface_CreateAndAddBody(bi.ptr, settings.ptr, C.JPH_Activation(activation)))
}

// AddBody adds a previously created body to the system.
func (bi *BodyInterface) AddBody(id BodyID, activation Activation) {
	C.JPH_BodyInterface_AddBody(bi.ptr, C.JPH_BodyID(id), C.JPH_Activation(activation))
}

// RemoveBody removes a body from the system without destroying it.
func (bi *BodyInterface) RemoveBody(id BodyID) {
	C.JPH_BodyInterface_RemoveBody(bi.ptr, C.JPH_BodyID(id))
}

// DestroyBody destroys a body that is not in the system.
func (bi *BodyInterface) DestroyBody(id BodyID) {
	C.JPH_BodyInterface_DestroyBody(bi.ptr, C.JPH_BodyID(id))
}

// RemoveAndDestroyBody removes a body from the system and destroys it.
func (bi *BodyInterface) RemoveAndDestroyBody(id BodyID) {
	C.JPH_BodyInterface_RemoveAndDestroyBody(bi.ptr, C.JPH_BodyID(id))
}

// IsAdded reports whether the body is in the system.
func (bi *BodyInterface) IsAdded(id BodyID) bool {
	return bool(C.JPH_BodyInterface_IsAdded(bi.ptr, C.JPH_BodyID(id)))
}

// IsActive reports whether the body is active.
func (bi *BodyInterface) IsActive(id BodyID) bool {
	return bool(C.JPH_BodyInterface_IsActive(bi.ptr, C.JPH_BodyID(id)))
}

// ActivateBody activates the body.
func (bi *BodyInterface) ActivateBody(id BodyID) {
	C.JPH_BodyInterface_ActivateBody(bi.ptr, C.JPH_BodyID(id))
}

// DeactivateBody deactivates the body.
func (bi *BodyInterface) DeactivateBody(id BodyID) {
	C.JPH_BodyInterface_DeactivateBody(bi.ptr, C.JPH_BodyID(id))
}

// Position returns the body's position.
func (bi *BodyInterface) Position(id BodyID) vecmath.RVec3 {
	var out C.JPH_RVec3
	C.JPH_BodyInterface_GetPosition(bi.ptr, C.JPH_BodyID(id), &out)
	return rvec3FromC(out)
}

// SetPosition sets the body's position.
func (bi *BodyInterface) SetPosition(id BodyID, p vecmath.RVec3, activation Activation) {
	c := rvec3ToC(p)
	C.JPH_BodyInterface_SetPosition(bi.ptr, C.JPH_BodyID(id), &c, C.JPH_Activation(activation))
}

// Rotation returns the body's rotation.
func (bi *BodyInterface) Rotation(id BodyID) vecmath.Quat {
	var out C.JPH_Quat
	C.JPH_BodyInterface_GetRotation(bi.ptr, C.JPH_BodyID(id), &out)
	return quatFromC(out)
}

// SetRotation sets the body's rotation.
func (bi *BodyInterface) SetRotation(id BodyID, q vecmath.Quat, activation Activation) {
	c := quatToC(q)
	C.JPH_BodyInterface_SetRotation(bi.ptr, C.JPH_BodyID(id), &c, C.JPH_Activation(activation))
}

// PositionAndRotation returns the body's position and rotation.
func (bi *BodyInterface) PositionAndRotation(id BodyID) (vecmath.RVec3, vecmath.Quat) {
	var p C.JPH_RVec3
	var q C.JPH_Quat
	C.JPH_BodyInterface_GetPositionAndRotation(bi.ptr, C.JPH_BodyID(id), &p, &q)
	return rvec3FromC(p), quatFromC(q)
}

// SetPositionAndRotation sets the body's position and rotation.
func (bi *BodyInterface) SetPositionAndRotation(id BodyID, p vecmath.RVec3, q vecmath.Quat, activation Activation) {
	cp := rvec3ToC(p)
	cq := quatToC(q)
	C.JPH_BodyInterface_SetPositionAndRotation(bi.ptr, C.JPH_BodyID(id), &cp, &cq, C.JPH_Activation(activation))
}

// CenterOfMassPosition returns the body's center of mass position.
func (bi *BodyInterface) CenterOfMassPosition(id BodyID) vecmath.RVec3 {
	var out C.JPH_RVec3
	C.JPH_BodyInterface_GetCenterOfMassPosition(bi.ptr, C.JPH_BodyID(id), &out)
	return rvec3FromC(out)
}

// LinearVelocity returns the body's linear velocity.
func (bi *BodyInterface) LinearVelocity(id BodyID) vecmath.Vec3 {
	var out C.JPH_Vec3
	C.JPH_BodyInterface_GetLinearVelocity(bi.ptr, C.JPH_BodyID(id), &out)
	return vec3FromC(out)
}

// SetLinearVelocity sets the body's linear velocity.
func (bi *BodyInterface) SetLinearVelocity(id BodyID, v vecmath.Vec3) {
	c := vec3ToC(v)
	C.JPH_BodyInterface_SetLinearVelocity(bi.ptr, C.JPH_BodyID(id), &c)
}

// AngularVelocity returns the body's angular velocity.
func (bi *BodyInterface) AngularVelocity(id BodyID) vecmath.Vec3 {
	var out C.JPH_Vec3
	C.JPH_BodyInterface_GetAngularVelocity(bi.ptr, C.JPH_BodyID(id), &out)
	return vec3FromC(out)
}

// SetAngularVelocity sets the body's angular velocity.
func (bi *BodyInterface) SetAngularVelocity(id BodyID, v vecmath.Vec3) {
	c := vec3ToC(v)
	C.JPH_BodyInterface_SetAngularVelocity(bi.ptr, C.JPH_BodyID(id), &c)
}

// AddForce applies a force to the body for the current step.
func (bi *BodyInterface) AddForce(id BodyID, f vecmath.Vec3) {
	c := vec3ToC(f)
	C.JPH_BodyInterface_AddForce(bi.ptr, C.JPH_BodyID(id), &c)
}

// AddImpulse applies an impulse to the body.
func (bi *BodyInterface) AddImpulse(id BodyID, impulse vecmath.Vec3) {
	c := vec3ToC(impulse)
	C.JPH_BodyInterface_AddImpulse(bi.ptr, C.JPH_BodyID(id), &c)
}

// MotionType returns the body's motion type.
func (bi *BodyInterface) MotionType(id BodyID) MotionType {
	return MotionType(C.JPH_BodyInterface_GetMotionType(bi.ptr, C.JPH_BodyID(id)))
}

// SetMotionType sets the body's motion type.
func (bi *BodyInterface) SetMotionType(id BodyID, m MotionType, activation Activation) {
	C.JPH_BodyInterface_SetMotionType(bi.ptr, C.JPH_BodyID(id), C.JPH_MotionType(m), C.JPH_Activation(activation))
}

// Friction returns the body's friction.
func (bi *BodyInterface) Friction(id BodyID) float32 {
	return float32(C.JPH_BodyInterface_GetFriction(bi.ptr, C.JPH_BodyID(id)))
}

// SetFriction sets the body's friction.
func (bi *BodyInterface) SetFriction(id BodyID, v float32) {
	C.JPH_BodyInterface_SetFriction(bi.ptr, C.JPH_BodyID(id), C.float(v))
}

// Restitution returns the body's restitution.
func (bi *BodyInterface) Restitution(id BodyID) float32 {
	return float32(C.JPH_BodyInterface_GetRestitution(bi.ptr, C.JPH_BodyID(id)))
}

// SetRestitution sets the body's restitution.
func (bi *BodyInterface) SetRestitution(id BodyID, v float32) {
	C.JPH_BodyInterface_SetRestitution(bi.ptr, C.JPH_BodyID(id), C.float(v))
}
