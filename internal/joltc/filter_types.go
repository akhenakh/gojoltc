package joltc

// Filter callback sets. They live here (not in the cgo file) so the cgo-less
// stub can reference them too; nothing here touches C.

// BroadPhaseLayerFilterCallbacks filters by broad-phase layer.
type BroadPhaseLayerFilterCallbacks struct {
	ShouldCollide func(layer BroadPhaseLayer) bool
}

// ObjectLayerFilterCallbacks filters by object layer.
type ObjectLayerFilterCallbacks struct {
	ShouldCollide func(layer ObjectLayer) bool
}

// BodyFilterCallbacks filters by body. ShouldCollide is used when the body is
// not locked; ShouldCollideLocked when it already is (e.g. during a query that
// holds a lock) and receives the body itself.
type BodyFilterCallbacks struct {
	ShouldCollide       func(bodyID BodyID) bool
	ShouldCollideLocked func(body *Body) bool
}

// ShapeFilterCallbacks filters by shape (and sub-shape).
type ShapeFilterCallbacks struct {
	ShouldCollide  func(shape2 *Shape, subShapeID2 SubShapeID) bool
	ShouldCollide2 func(shape1 *Shape, subShapeID1 SubShapeID, shape2 *Shape, subShapeID2 SubShapeID) bool
}
