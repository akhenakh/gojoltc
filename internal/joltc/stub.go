//go:build !cgo

// Package joltc is the only package in this module that imports "C".
//
// This file is compiled when cgo is disabled (CGO_ENABLED=0). It mirrors the
// exported surface of the bindings so that a cgo-less build of the public
// package compiles, and returns ErrCgoRequired at runtime instead of
// misbehaving. The handle types carry no state.
package joltc

import "github.com/akhenakh/gojoltc/internal/vecmath"

// DefaultConvexRadius mirrors JPH_DEFAULT_CONVEX_RADIUS.
const DefaultConvexRadius float32 = 0.05

// Initialize always returns ErrCgoRequired in a cgo-less build.
func Initialize() error { return ErrCgoRequired }

// Shutdown is a no-op in a cgo-less build.
func Shutdown() {}

// SetTraceHandler is a no-op in a cgo-less build.
func SetTraceHandler(fn func(string)) {}

// SetAssertFailureHandler is a no-op in a cgo-less build.
func SetAssertFailureHandler(fn func(expression, message, file string, line uint32) bool) {}

// JobSystemConfig configures a thread-pool job system.
type JobSystemConfig struct {
	MaxJobs     uint32
	MaxBarriers uint32
	NumThreads  int32
}

// JobSystem is a job system.
type JobSystem struct{}

// BroadPhaseLayerInterface maps object layers to broad-phase layers.
type BroadPhaseLayerInterface struct{}

// ObjectLayerPairFilter decides which object layers may collide.
type ObjectLayerPairFilter struct{}

// ObjectVsBroadPhaseLayerFilter decides which objects are tested against which
// broad-phase layers.
type ObjectVsBroadPhaseLayerFilter struct{}

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
type PhysicsSystem struct{}

// Shape is a collision shape.
type Shape struct{}

// BodyCreationSettings describes a body to be created.
type BodyCreationSettings struct{}

// Body is a rigid body instance.
type Body struct{}

// BodyInterface creates and manipulates bodies.
type BodyInterface struct{}

// NewJobSystemThreadPool returns nil in a cgo-less build.
func NewJobSystemThreadPool(cfg *JobSystemConfig) *JobSystem { return nil }

// NewBroadPhaseLayerInterfaceTable returns nil in a cgo-less build.
func NewBroadPhaseLayerInterfaceTable(numObjectLayers, numBroadPhaseLayers uint32) *BroadPhaseLayerInterface {
	return nil
}

// NewObjectLayerPairFilterTable returns nil in a cgo-less build.
func NewObjectLayerPairFilterTable(numObjectLayers uint32) *ObjectLayerPairFilter { return nil }

// NewObjectVsBroadPhaseLayerFilterTable returns nil in a cgo-less build.
func NewObjectVsBroadPhaseLayerFilterTable(bp *BroadPhaseLayerInterface, numBroadPhaseLayers uint32, pair *ObjectLayerPairFilter, numObjectLayers uint32) *ObjectVsBroadPhaseLayerFilter {
	return nil
}

// NewPhysicsSystem returns nil in a cgo-less build.
func NewPhysicsSystem(s *PhysicsSystemSettings) *PhysicsSystem { return nil }

// NewBoxShape returns nil in a cgo-less build.
func NewBoxShape(halfExtent vecmath.Vec3, convexRadius float32) *Shape { return nil }

// NewSphereShape returns nil in a cgo-less build.
func NewSphereShape(radius float32) *Shape { return nil }

// NewCapsuleShape returns nil in a cgo-less build.
func NewCapsuleShape(halfHeightOfCylinder, radius float32) *Shape { return nil }

// NewCylinderShape returns nil in a cgo-less build.
func NewCylinderShape(halfHeight, radius float32) *Shape { return nil }

// NewBodyCreationSettings returns nil in a cgo-less build.
func NewBodyCreationSettings(shape *Shape, position vecmath.RVec3, rotation vecmath.Quat, motionType MotionType, layer ObjectLayer) *BodyCreationSettings {
	return nil
}

// NewTriangleShape returns nil in a cgo-less build.
func NewTriangleShape(v1, v2, v3 vecmath.Vec3, convexRadius float32) *Shape { return nil }

// NewPlaneShape returns nil in a cgo-less build.
func NewPlaneShape(plane vecmath.Plane, halfExtent float32) *Shape { return nil }

// NewRotatedTranslatedShape returns nil in a cgo-less build.
func NewRotatedTranslatedShape(position vecmath.Vec3, rotation vecmath.Quat, inner *Shape) *Shape {
	return nil
}

// NewScaledShape returns nil in a cgo-less build.
func NewScaledShape(inner *Shape, scale vecmath.Vec3) *Shape { return nil }

// NewOffsetCenterOfMassShape returns nil in a cgo-less build.
func NewOffsetCenterOfMassShape(offset vecmath.Vec3, inner *Shape) *Shape { return nil }

// Shape settings handle types.

// SphereShapeSettings holds settings for a sphere shape.
type SphereShapeSettings struct{}

// BoxShapeSettings holds settings for a box shape.
type BoxShapeSettings struct{}

// CapsuleShapeSettings holds settings for a capsule shape.
type CapsuleShapeSettings struct{}

// CylinderShapeSettings holds settings for a cylinder shape.
type CylinderShapeSettings struct{}

// TaperedCapsuleShapeSettings holds settings for a tapered capsule shape.
type TaperedCapsuleShapeSettings struct{}

// TaperedCylinderShapeSettings holds settings for a tapered cylinder shape.
type TaperedCylinderShapeSettings struct{}

// ConvexHullShapeSettings holds settings for a convex hull shape.
type ConvexHullShapeSettings struct{}

// CompoundShapeSettings is the base for compound shape settings.
type CompoundShapeSettings struct{}

// StaticCompoundShapeSettings holds settings for a static compound shape.
type StaticCompoundShapeSettings struct{ CompoundShapeSettings }

// MutableCompoundShapeSettings holds settings for a mutable compound shape.
type MutableCompoundShapeSettings struct{ CompoundShapeSettings }

// EmptyShapeSettings holds settings for an empty shape.
type EmptyShapeSettings struct{}

// BroadPhaseQuery is a broad-phase collision query interface.
type BroadPhaseQuery struct{}

// NarrowPhaseQuery is a narrow-phase collision query interface.
type NarrowPhaseQuery struct{}

// Constraint is a physics joint between bodies.
type Constraint struct{}

// FixedConstraintSettings holds settings for a fixed constraint.
type FixedConstraintSettings struct{}

// DistanceConstraintSettings holds settings for a distance constraint.
type DistanceConstraintSettings struct{}

// PointConstraintSettings holds settings for a point constraint.
type PointConstraintSettings struct{}

// HingeConstraintSettings holds settings for a hinge constraint.
type HingeConstraintSettings struct{}

// SliderConstraintSettings holds settings for a slider constraint.
type SliderConstraintSettings struct{}

// ConeConstraintSettings holds settings for a cone constraint.
type ConeConstraintSettings struct{}

// SwingTwistConstraintSettings holds settings for a swing-twist constraint.
type SwingTwistConstraintSettings struct{}

// SixDOFConstraintSettings holds settings for a six DOF constraint.
type SixDOFConstraintSettings struct{}

// GearConstraintSettings holds settings for a gear constraint.
type GearConstraintSettings struct{}

// NewFixedConstraintSettings returns nil in a cgo-less build.
func NewFixedConstraintSettings() *FixedConstraintSettings { return nil }

// NewDistanceConstraintSettings returns nil in a cgo-less build.
func NewDistanceConstraintSettings() *DistanceConstraintSettings { return nil }

// NewPointConstraintSettings returns nil in a cgo-less build.
func NewPointConstraintSettings() *PointConstraintSettings { return nil }

// NewHingeConstraintSettings returns nil in a cgo-less build.
func NewHingeConstraintSettings() *HingeConstraintSettings { return nil }

// NewSliderConstraintSettings returns nil in a cgo-less build.
func NewSliderConstraintSettings() *SliderConstraintSettings { return nil }

// NewConeConstraintSettings returns nil in a cgo-less build.
func NewConeConstraintSettings() *ConeConstraintSettings { return nil }

// NewSwingTwistConstraintSettings returns nil in a cgo-less build.
func NewSwingTwistConstraintSettings() *SwingTwistConstraintSettings { return nil }

// NewSixDOFConstraintSettings returns nil in a cgo-less build.
func NewSixDOFConstraintSettings() *SixDOFConstraintSettings { return nil }

// NewGearConstraintSettings returns nil in a cgo-less build.
func NewGearConstraintSettings() *GearConstraintSettings { return nil }

// CharacterSettings holds settings for a character.
type CharacterSettings struct{}

// Character is a rigid-body character.
type Character struct{}

// CharacterVirtualSettings holds settings for a virtual character.
type CharacterVirtualSettings struct{}

// NewCharacterSettings returns nil in a cgo-less build.
func NewCharacterSettings() *CharacterSettings { return nil }

// NewCharacterVirtualSettings returns nil in a cgo-less build.
func NewCharacterVirtualSettings() *CharacterVirtualSettings { return nil }

// NewCharacter returns nil in a cgo-less build.
func NewCharacter(settings *CharacterSettings, position vecmath.RVec3, rotation vecmath.Quat, userData uint64, system *PhysicsSystem) *Character {
	return nil
}

// Skeleton is a hierarchy of joints.
type Skeleton struct{}

// SkeletonPose is a pose for a skeleton.
type SkeletonPose struct{}

// SkeletonMapper maps poses between two skeletons.
type SkeletonMapper struct{}

// RagdollSettings holds settings for a ragdoll.
type RagdollSettings struct{}

// Ragdoll is a physics-driven skeleton.
type Ragdoll struct{}

// NewSkeleton returns nil in a cgo-less build.
func NewSkeleton() *Skeleton { return nil }

// NewSkeletonPose returns nil in a cgo-less build.
func NewSkeletonPose() *SkeletonPose { return nil }

// NewSkeletonMapper returns nil in a cgo-less build.
func NewSkeletonMapper() *SkeletonMapper { return nil }

// NewRagdollSettings returns nil in a cgo-less build.
func NewRagdollSettings() *RagdollSettings { return nil }

// WheelSettings holds settings for a wheel.
type WheelSettings struct{}

// Wheel is a runtime wheel instance.
type Wheel struct{}

// VehicleCollisionTester is a vehicle ground collision tester.
type VehicleCollisionTester struct{}

// VehicleTransmissionSettings holds transmission settings.
type VehicleTransmissionSettings struct{}

// WheeledVehicleControllerSettings holds settings for a wheeled vehicle controller.
type WheeledVehicleControllerSettings struct{}

// NewWheelSettings returns nil in a cgo-less build.
func NewWheelSettings() *WheelSettings { return nil }

// NewWheel returns nil in a cgo-less build.
func NewWheel(settings *WheelSettings) *Wheel { return nil }

// NewVehicleCollisionTesterRay returns nil in a cgo-less build.
func NewVehicleCollisionTesterRay(layer ObjectLayer, up vecmath.Vec3, maxSlopeAngle float32) *VehicleCollisionTester {
	return nil
}

// NewVehicleCollisionTesterCastSphere returns nil in a cgo-less build.
func NewVehicleCollisionTesterCastSphere(layer ObjectLayer, radius float32, up vecmath.Vec3, maxSlopeAngle float32) *VehicleCollisionTester {
	return nil
}

// NewVehicleCollisionTesterCastCylinder returns nil in a cgo-less build.
func NewVehicleCollisionTesterCastCylinder(layer ObjectLayer, convexRadiusFraction float32) *VehicleCollisionTester {
	return nil
}

// NewVehicleEngineSettings returns zero settings in a cgo-less build.
func NewVehicleEngineSettings() VehicleEngineSettings { return VehicleEngineSettings{} }

// NewVehicleDifferentialSettings returns zero settings in a cgo-less build.
func NewVehicleDifferentialSettings() VehicleDifferentialSettings {
	return VehicleDifferentialSettings{}
}

// NewVehicleAntiRollBar returns zero settings in a cgo-less build.
func NewVehicleAntiRollBar() VehicleAntiRollBar { return VehicleAntiRollBar{} }

// NewVehicleConstraintSettings returns zero settings in a cgo-less build.
func NewVehicleConstraintSettings() VehicleConstraintSettings { return VehicleConstraintSettings{} }

// NewVehicleTransmissionSettings returns nil in a cgo-less build.
func NewVehicleTransmissionSettings() *VehicleTransmissionSettings { return nil }

// NewWheeledVehicleControllerSettings returns nil in a cgo-less build.
func NewWheeledVehicleControllerSettings() *WheeledVehicleControllerSettings { return nil }

// ContactManifold describes a contact between two bodies.
type ContactManifold struct{}

// ContactSettings lets a contact callback tweak the contact.
type ContactSettings struct{}

// ContactListener receives contact events from the physics system.
type ContactListener struct{}

// BodyActivationListener receives body activation events.
type BodyActivationListener struct{}

// StepListener receives a callback for each simulation step.
type StepListener struct{}

// DebugRenderer turns Jolt's debug draw calls into Go callbacks.
type DebugRenderer struct{}

// NewContactListener returns nil in a cgo-less build.
func NewContactListener(cb *ContactListenerCallbacks) *ContactListener { return nil }

// NewBodyActivationListener returns nil in a cgo-less build.
func NewBodyActivationListener(cb *BodyActivationCallbacks) *BodyActivationListener { return nil }

// NewStepListener returns nil in a cgo-less build.
func NewStepListener(cb *StepListenerCallbacks) *StepListener { return nil }

// NewDebugRenderer returns nil in a cgo-less build.
func NewDebugRenderer(cb *DebugRendererCallbacks) *DebugRenderer { return nil }

// SoftBodySharedSettings describes a soft body mesh.
type SoftBodySharedSettings struct{}

// SoftBodyCreationSettings describes a soft body to be created.
type SoftBodyCreationSettings struct{}

// NewSoftBodySharedSettings returns nil in a cgo-less build.
func NewSoftBodySharedSettings() *SoftBodySharedSettings { return nil }

// NewSoftBodyCreationSettings returns nil in a cgo-less build.
func NewSoftBodyCreationSettings() *SoftBodyCreationSettings { return nil }

// NewSoftBodyCreationSettingsFor returns nil in a cgo-less build.
func NewSoftBodyCreationSettingsFor(shared *SoftBodySharedSettings, position vecmath.RVec3, rotation vecmath.Quat, layer ObjectLayer) *SoftBodyCreationSettings {
	return nil
}

// NewSphereShapeSettings returns nil in a cgo-less build.
func NewSphereShapeSettings(radius float32) *SphereShapeSettings { return nil }

// NewBoxShapeSettings returns nil in a cgo-less build.
func NewBoxShapeSettings(halfExtent vecmath.Vec3, convexRadius float32) *BoxShapeSettings {
	return nil
}

// NewCapsuleShapeSettings returns nil in a cgo-less build.
func NewCapsuleShapeSettings(halfHeightOfCylinder, radius float32) *CapsuleShapeSettings {
	return nil
}

// NewCylinderShapeSettings returns nil in a cgo-less build.
func NewCylinderShapeSettings(halfHeight, radius, convexRadius float32) *CylinderShapeSettings {
	return nil
}

// NewTaperedCapsuleShapeSettings returns nil in a cgo-less build.
func NewTaperedCapsuleShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius float32) *TaperedCapsuleShapeSettings {
	return nil
}

// NewTaperedCylinderShapeSettings returns nil in a cgo-less build.
func NewTaperedCylinderShapeSettings(halfHeightOfTaperedCylinder, topRadius, bottomRadius, convexRadius float32) *TaperedCylinderShapeSettings {
	return nil
}

// NewConvexHullShapeSettings returns nil in a cgo-less build.
func NewConvexHullShapeSettings(points []vecmath.Vec3, maxConvexRadius float32) *ConvexHullShapeSettings {
	return nil
}

// NewStaticCompoundShapeSettings returns nil in a cgo-less build.
func NewStaticCompoundShapeSettings() *StaticCompoundShapeSettings { return nil }

// NewMutableCompoundShapeSettings returns nil in a cgo-less build.
func NewMutableCompoundShapeSettings() *MutableCompoundShapeSettings { return nil }

// NewEmptyShapeSettings returns nil in a cgo-less build.
func NewEmptyShapeSettings(centerOfMass vecmath.Vec3) *EmptyShapeSettings { return nil }
