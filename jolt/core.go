package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Core engine types. Handles are defined in the cgo bindings and re-exported
// here; their methods are part of this package's API.
type (
	// JobSystemConfig configures a thread-pool job system.
	JobSystemConfig = joltc.JobSystemConfig
	// JobSystem steps the simulation in parallel.
	JobSystem = joltc.JobSystem
	// BroadPhaseLayerInterface maps object layers to broad-phase layers.
	BroadPhaseLayerInterface = joltc.BroadPhaseLayerInterface
	// ObjectLayerPairFilter decides which object layers may collide.
	ObjectLayerPairFilter = joltc.ObjectLayerPairFilter
	// ObjectVsBroadPhaseLayerFilter decides which objects are tested against
	// which broad-phase layers.
	ObjectVsBroadPhaseLayerFilter = joltc.ObjectVsBroadPhaseLayerFilter
	// PhysicsSystemSettings configures a physics system.
	PhysicsSystemSettings = joltc.PhysicsSystemSettings
	// PhysicsSystem is the main simulation object.
	PhysicsSystem = joltc.PhysicsSystem
	// Shape is a collision shape.
	Shape = joltc.Shape
	// BodyCreationSettings describes a body to be created.
	BodyCreationSettings = joltc.BodyCreationSettings
	// Body is a rigid body instance.
	Body = joltc.Body
	// BodyInterface creates and manipulates bodies.
	BodyInterface = joltc.BodyInterface
)

// DefaultConvexRadius mirrors JPH_DEFAULT_CONVEX_RADIUS.
const DefaultConvexRadius = joltc.DefaultConvexRadius

// Constants mirroring joltc macros.
const (
	// MaxPhysicsJobs mirrors JPH_MAX_PHYSICS_JOBS.
	MaxPhysicsJobs = joltc.MaxPhysicsJobs
	// MaxPhysicsBarriers mirrors JPH_MAX_PHYSICS_BARRIERS.
	MaxPhysicsBarriers = joltc.MaxPhysicsBarriers
)

// NewJobSystemThreadPool creates a job system backed by a thread pool. A nil
// config uses Jolt's defaults.
func NewJobSystemThreadPool(cfg *JobSystemConfig) *JobSystem {
	return joltc.NewJobSystemThreadPool(cfg)
}

// NewBroadPhaseLayerInterfaceTable creates a table-based layer mapping.
func NewBroadPhaseLayerInterfaceTable(numObjectLayers, numBroadPhaseLayers uint32) *BroadPhaseLayerInterface {
	return joltc.NewBroadPhaseLayerInterfaceTable(numObjectLayers, numBroadPhaseLayers)
}

// NewObjectLayerPairFilterTable creates a table-based object layer pair filter.
func NewObjectLayerPairFilterTable(numObjectLayers uint32) *ObjectLayerPairFilter {
	return joltc.NewObjectLayerPairFilterTable(numObjectLayers)
}

// NewObjectVsBroadPhaseLayerFilterTable creates a table-based object vs
// broad-phase layer filter.
func NewObjectVsBroadPhaseLayerFilterTable(bp *BroadPhaseLayerInterface, numBroadPhaseLayers uint32, pair *ObjectLayerPairFilter, numObjectLayers uint32) *ObjectVsBroadPhaseLayerFilter {
	return joltc.NewObjectVsBroadPhaseLayerFilterTable(bp, numBroadPhaseLayers, pair, numObjectLayers)
}

// NewPhysicsSystem creates a physics system from settings.
func NewPhysicsSystem(settings *PhysicsSystemSettings) *PhysicsSystem {
	return joltc.NewPhysicsSystem(settings)
}

// NewBoxShape creates a box shape.
func NewBoxShape(halfExtent Vec3, convexRadius float32) *Shape {
	return joltc.NewBoxShape(halfExtent, convexRadius)
}

// NewSphereShape creates a sphere shape.
func NewSphereShape(radius float32) *Shape { return joltc.NewSphereShape(radius) }

// NewCapsuleShape creates a capsule shape.
func NewCapsuleShape(halfHeightOfCylinder, radius float32) *Shape {
	return joltc.NewCapsuleShape(halfHeightOfCylinder, radius)
}

// NewCylinderShape creates a cylinder shape.
func NewCylinderShape(halfHeight, radius float32) *Shape {
	return joltc.NewCylinderShape(halfHeight, radius)
}

// NewBodyCreationSettings creates settings for a body.
func NewBodyCreationSettings(shape *Shape, position RVec3, rotation Quat, motionType MotionType, layer ObjectLayer) *BodyCreationSettings {
	return joltc.NewBodyCreationSettings(shape, position, rotation, motionType, layer)
}
