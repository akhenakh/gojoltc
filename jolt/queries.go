package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Query interfaces, re-exported from the bindings.
type (
	// BroadPhaseQuery is a broad-phase collision query interface.
	BroadPhaseQuery = joltc.BroadPhaseQuery
	// NarrowPhaseQuery is a narrow-phase collision query interface.
	NarrowPhaseQuery = joltc.NarrowPhaseQuery
)
