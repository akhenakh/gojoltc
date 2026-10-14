package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Character types, re-exported from the bindings.
type (
	// CharacterSettings holds settings for a character.
	CharacterSettings = joltc.CharacterSettings
	// Character is a rigid-body character.
	Character = joltc.Character
	// CharacterVirtualSettings holds settings for a virtual character.
	CharacterVirtualSettings = joltc.CharacterVirtualSettings
)

// NewCharacterSettings creates character settings with Jolt defaults.
func NewCharacterSettings() *CharacterSettings { return joltc.NewCharacterSettings() }

// NewCharacterVirtualSettings creates virtual character settings with defaults.
func NewCharacterVirtualSettings() *CharacterVirtualSettings {
	return joltc.NewCharacterVirtualSettings()
}

// NewCharacter creates a character with the given settings at the given pose.
func NewCharacter(settings *CharacterSettings, position RVec3, rotation Quat, userData uint64, system *PhysicsSystem) *Character {
	return joltc.NewCharacter(settings, position, rotation, userData, system)
}
