package joltc

import "github.com/akhenakh/gojoltc/internal/vecmath"

// MassProperties holds a shape's mass and inertia tensor (JPH_MassProperties).
type MassProperties struct {
	Mass    float32
	Inertia vecmath.Mat4
}

// RayCastResult describes a ray hit against a shape (JPH_RayCastResult).
type RayCastResult struct {
	BodyID      BodyID
	Fraction    float32
	SubShapeID2 SubShapeID
}
