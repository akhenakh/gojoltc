package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Query interfaces and filters, re-exported from the bindings.
type (
	// BroadPhaseQuery is a broad-phase collision query interface.
	BroadPhaseQuery = joltc.BroadPhaseQuery
	// NarrowPhaseQuery is a narrow-phase collision query interface.
	NarrowPhaseQuery = joltc.NarrowPhaseQuery
	// QueryOptions bundles the optional filters for a query.
	QueryOptions = joltc.QueryOptions
	// RayCastSettings tunes a ray cast.
	RayCastSettings = joltc.RayCastSettings
	// CollideShapeSettings tunes a shape overlap query.
	CollideShapeSettings = joltc.CollideShapeSettings
	// ShapeCastSettings tunes a shape sweep.
	ShapeCastSettings = joltc.ShapeCastSettings
	// CollidePointResult reports a body containing a point.
	CollidePointResult = joltc.CollidePointResult
	// ShapeCastResult reports a shape sweep hit.
	ShapeCastResult = joltc.ShapeCastResult

	// BroadPhaseLayerFilter restricts a query by broad-phase layer.
	BroadPhaseLayerFilter = joltc.BroadPhaseLayerFilter
	// BroadPhaseLayerFilterCallbacks holds the Go callbacks for a broad-phase layer filter.
	BroadPhaseLayerFilterCallbacks = joltc.BroadPhaseLayerFilterCallbacks
	// ObjectLayerFilter restricts a query by object layer.
	ObjectLayerFilter = joltc.ObjectLayerFilter
	// ObjectLayerFilterCallbacks holds the Go callbacks for an object layer filter.
	ObjectLayerFilterCallbacks = joltc.ObjectLayerFilterCallbacks
	// BodyFilter restricts a query by body.
	BodyFilter = joltc.BodyFilter
	// BodyFilterCallbacks holds the Go callbacks for a body filter.
	BodyFilterCallbacks = joltc.BodyFilterCallbacks
	// ShapeFilter restricts a query by shape.
	ShapeFilter = joltc.ShapeFilter
	// ShapeFilterCallbacks holds the Go callbacks for a shape filter.
	ShapeFilterCallbacks = joltc.ShapeFilterCallbacks
)

// DefaultRayCastSettings returns Jolt's ray cast defaults.
func DefaultRayCastSettings() RayCastSettings { return joltc.DefaultRayCastSettings() }

// DefaultCollideShapeSettings returns Jolt's collide shape defaults.
func DefaultCollideShapeSettings() CollideShapeSettings { return joltc.DefaultCollideShapeSettings() }

// DefaultShapeCastSettings returns Jolt's shape cast defaults.
func DefaultShapeCastSettings() ShapeCastSettings { return joltc.DefaultShapeCastSettings() }

// NewBroadPhaseLayerFilter creates a broad-phase layer filter.
func NewBroadPhaseLayerFilter(cb *BroadPhaseLayerFilterCallbacks) *BroadPhaseLayerFilter {
	return joltc.NewBroadPhaseLayerFilter(cb)
}

// NewObjectLayerFilter creates an object layer filter.
func NewObjectLayerFilter(cb *ObjectLayerFilterCallbacks) *ObjectLayerFilter {
	return joltc.NewObjectLayerFilter(cb)
}

// NewBodyFilter creates a body filter.
func NewBodyFilter(cb *BodyFilterCallbacks) *BodyFilter { return joltc.NewBodyFilter(cb) }

// NewShapeFilter creates a shape filter.
func NewShapeFilter(cb *ShapeFilterCallbacks) *ShapeFilter { return joltc.NewShapeFilter(cb) }
