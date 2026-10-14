package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Mesh and height field value types, re-exported from the bindings.
type (
	// Triangle is a standalone mesh triangle.
	Triangle = joltc.Triangle
	// IndexedTriangle references three vertices of a mesh.
	IndexedTriangle = joltc.IndexedTriangle
)

// NewMeshShape creates a mesh shape from standalone triangles.
func NewMeshShape(triangles []Triangle) *Shape { return joltc.NewMeshShapeSettings(triangles) }

// NewMeshShapeFromVertices creates a mesh shape from a shared vertex array and
// indexed triangles.
func NewMeshShapeFromVertices(vertices []Vec3, triangles []IndexedTriangle) *Shape {
	return joltc.MeshShapeSettingsFromVertices(vertices, triangles)
}

// NewHeightFieldShape creates a height field shape from a square sample grid
// (len(samples) must be n*n), with the given offset and scale.
func NewHeightFieldShape(samples []float32, offset, scale Vec3) *Shape {
	return joltc.NewHeightFieldShapeSettings(samples, offset, scale)
}
