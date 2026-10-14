package joltc

// Identifier and layer types mirror the JPH_* integer typedefs.

// BodyID identifies a body within a physics system.
type BodyID uint32

// SubShapeID identifies a sub-shape within a compound shape.
type SubShapeID uint32

// ObjectLayer is a game-defined collision layer.
type ObjectLayer uint32

// BroadPhaseLayer is a broad-phase collision layer.
type BroadPhaseLayer uint8

// CollisionGroupID identifies a collision group.
type CollisionGroupID uint32

// CollisionSubGroupID identifies a collision sub-group.
type CollisionSubGroupID uint32

// CharacterID identifies a character.
type CharacterID uint32
