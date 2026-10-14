package joltc

import "github.com/akhenakh/gojoltc/internal/vecmath"

// Mesh and height field value types. The cgo side asserts their sizes before
// passing slices straight to C.

// Triangle is a standalone mesh triangle (JPH_Triangle).
type Triangle struct {
	V1            vecmath.Vec3
	V2            vecmath.Vec3
	V3            vecmath.Vec3
	MaterialIndex uint32
}

// IndexedTriangle references three vertices (JPH_IndexedTriangle).
type IndexedTriangle struct {
	I1            uint32
	I2            uint32
	I3            uint32
	MaterialIndex uint32
	UserData      uint32
}
