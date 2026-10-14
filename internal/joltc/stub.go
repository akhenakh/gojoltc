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
