package joltc

// Query value types. They live here (not in the cgo file) so the cgo-less stub
// and the public package can reference them; the C conversion methods stay in
// the cgo file.

import "github.com/akhenakh/gojoltc/internal/vecmath"

// QueryOptions bundles the optional filters for a query. A nil field (or a nil
// *QueryOptions) disables that filter.
type QueryOptions struct {
	BroadPhaseLayer *BroadPhaseLayerFilter
	ObjectLayer     *ObjectLayerFilter
	Body            *BodyFilter
	Shape           *ShapeFilter
}

// RayCastSettings tunes a ray cast.
type RayCastSettings struct {
	BackFaceModeTriangles BackFaceMode
	BackFaceModeConvex    BackFaceMode
	TreatConvexAsSolid    bool
}

// DefaultRayCastSettings returns Jolt's ray cast defaults.
func DefaultRayCastSettings() RayCastSettings {
	return RayCastSettings{TreatConvexAsSolid: true}
}

// CollideShapeSettings tunes a shape overlap query.
type CollideShapeSettings struct {
	MaxSeparationDistance float32
	BackFaceMode          BackFaceMode
}

// DefaultCollideShapeSettings returns Jolt's collide shape defaults.
func DefaultCollideShapeSettings() CollideShapeSettings { return CollideShapeSettings{} }

// ShapeCastSettings tunes a shape sweep.
type ShapeCastSettings struct {
	BackFaceModeTriangles           BackFaceMode
	BackFaceModeConvex              BackFaceMode
	UseShrunkenShapeAndConvexRadius bool
	ReturnDeepestPoint              bool
}

// DefaultShapeCastSettings returns Jolt's shape cast defaults.
func DefaultShapeCastSettings() ShapeCastSettings { return ShapeCastSettings{} }

// CollidePointResult reports a body containing a point.
type CollidePointResult struct {
	BodyID      BodyID
	SubShapeID2 SubShapeID
}

// ShapeCastResult reports a shape sweep hit.
type ShapeCastResult struct {
	ContactPointOn1  vecmath.Vec3
	ContactPointOn2  vecmath.Vec3
	PenetrationAxis  vecmath.Vec3
	PenetrationDepth float32
	SubShapeID1      SubShapeID
	SubShapeID2      SubShapeID
	BodyID2          BodyID
	Fraction         float32
	IsBackFaceHit    bool
}
