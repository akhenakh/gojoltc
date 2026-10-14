package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Identifier and layer types, re-exported from the cgo bindings.
type (
	// BodyID identifies a body within a physics system.
	BodyID = joltc.BodyID
	// SubShapeID identifies a sub-shape within a compound shape.
	SubShapeID = joltc.SubShapeID
	// ObjectLayer is a game-defined collision layer.
	ObjectLayer = joltc.ObjectLayer
	// BroadPhaseLayer is a broad-phase collision layer.
	BroadPhaseLayer = joltc.BroadPhaseLayer
	// CollisionGroupID identifies a collision group.
	CollisionGroupID = joltc.CollisionGroupID
	// CollisionSubGroupID identifies a collision sub-group.
	CollisionSubGroupID = joltc.CollisionSubGroupID
	// CharacterID identifies a character.
	CharacterID = joltc.CharacterID
)
