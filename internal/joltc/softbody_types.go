package joltc

import "github.com/akhenakh/gojoltc/internal/vecmath"

// Soft body value types. These mirror the C structs in layout; the cgo side
// asserts the sizes before passing slices of them straight to C.

// SoftVertex is a soft body vertex (JPH_SoftVertex).
type SoftVertex struct {
	Position vecmath.Vec3
	Velocity vecmath.Vec3
	InvMass  float32
}

// SoftFace is a triangle of a soft body surface (JPH_SoftFace).
type SoftFace struct {
	Vertex1       uint32
	Vertex2       uint32
	Vertex3       uint32
	MaterialIndex uint32
}
