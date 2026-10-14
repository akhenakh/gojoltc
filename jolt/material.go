package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// PhysicsMaterial assigns friction/restitution/color to a shape or sub-shape.
type PhysicsMaterial = joltc.PhysicsMaterial

// NewPhysicsMaterial creates a material with a debug name and packed color.
func NewPhysicsMaterial(name string, color Color) *PhysicsMaterial {
	return joltc.NewPhysicsMaterial(name, color)
}
