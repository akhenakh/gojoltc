//go:build cgo

package joltc

/*
#include <joltc.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// CharacterSettings holds settings for a character.
type CharacterSettings struct {
	ptr *C.JPH_CharacterSettings
	// ownShape is true while base.shape is the default shape allocated by
	// JPH_CharacterSettings_Init; we must free it ourselves (joltc does not).
	ownShape bool
}

// NewCharacterSettings creates character settings with Jolt defaults.
func NewCharacterSettings() *CharacterSettings {
	p := (*C.JPH_CharacterSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_CharacterSettings{}))))
	C.JPH_CharacterSettings_Init(p)
	return &CharacterSettings{ptr: p, ownShape: true}
}

// Destroy releases the settings, including the default shape allocated by Init
// unless the caller replaced it.
func (s *CharacterSettings) Destroy() {
	if s.ptr != nil {
		if s.ownShape && s.ptr.base.shape != nil {
			C.JPH_Shape_Destroy(s.ptr.base.shape)
		}
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// SetShape sets the character shape, releasing the default shape allocated by
// Init if it was still owned.
func (s *CharacterSettings) SetShape(shape *Shape) {
	if s.ownShape && s.ptr.base.shape != nil {
		C.JPH_Shape_Destroy(s.ptr.base.shape)
	}
	s.ptr.base.shape = shape.ptr
	s.ownShape = false
}

// SetLayer sets the object layer.
func (s *CharacterSettings) SetLayer(layer ObjectLayer) { s.ptr.layer = C.JPH_ObjectLayer(layer) }

// SetMass sets the character mass.
func (s *CharacterSettings) SetMass(mass float32) { s.ptr.mass = C.float(mass) }

// Up returns the character up vector.
func (s *CharacterSettings) Up() vecmath.Vec3 { return vec3FromC(s.ptr.base.up) }

// Character is a rigid-body character.
type Character struct{ ptr *C.JPH_Character }

// NewCharacter creates a character with the given settings at the given pose.
func NewCharacter(settings *CharacterSettings, position vecmath.RVec3, rotation vecmath.Quat, userData uint64, system *PhysicsSystem) *Character {
	pos := rvec3ToC(position)
	rot := quatToC(rotation)
	return &Character{ptr: C.JPH_Character_Create(settings.ptr, &pos, &rot, C.uint64_t(userData), system.ptr)}
}

// Destroy releases the character.
func (c *Character) Destroy() {
	if c.ptr != nil {
		C.JPH_CharacterBase_Destroy((*C.JPH_CharacterBase)(unsafe.Pointer(c.ptr)))
		c.ptr = nil
	}
}

// AddToPhysicsSystem adds the character to the physics system.
func (c *Character) AddToPhysicsSystem(activation Activation, lockBodies bool) {
	C.JPH_Character_AddToPhysicsSystem(c.ptr, C.JPH_Activation(activation), C.bool(lockBodies))
}

// RemoveFromPhysicsSystem removes the character from the physics system.
func (c *Character) RemoveFromPhysicsSystem(lockBodies bool) {
	C.JPH_Character_RemoveFromPhysicsSystem(c.ptr, C.bool(lockBodies))
}

// BodyID returns the character's body identifier.
func (c *Character) BodyID() BodyID { return BodyID(C.JPH_Character_GetBodyID(c.ptr)) }

// Position returns the character's position.
func (c *Character) Position(lockBodies bool) vecmath.RVec3 {
	var o C.JPH_RVec3
	C.JPH_Character_GetPosition(c.ptr, &o, C.bool(lockBodies))
	return rvec3FromC(o)
}

// SetPosition sets the character's position.
func (c *Character) SetPosition(position vecmath.RVec3, activation Activation, lockBodies bool) {
	p := rvec3ToC(position)
	C.JPH_Character_SetPosition(c.ptr, &p, C.JPH_Activation(activation), C.bool(lockBodies))
}

// LinearVelocity returns the character's linear velocity.
func (c *Character) LinearVelocity() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_Character_GetLinearVelocity(c.ptr, &o)
	return vec3FromC(o)
}

// SetLinearVelocity sets the character's linear velocity.
func (c *Character) SetLinearVelocity(v vecmath.Vec3, lockBodies bool) {
	vel := vec3ToC(v)
	C.JPH_Character_SetLinearVelocity(c.ptr, &vel, C.bool(lockBodies))
}

// CharacterVirtualSettings holds settings for a virtual character.
type CharacterVirtualSettings struct {
	ptr      *C.JPH_CharacterVirtualSettings
	ownShape bool
}

// NewCharacterVirtualSettings creates virtual character settings with defaults.
func NewCharacterVirtualSettings() *CharacterVirtualSettings {
	p := (*C.JPH_CharacterVirtualSettings)(C.malloc(C.size_t(unsafe.Sizeof(C.JPH_CharacterVirtualSettings{}))))
	C.JPH_CharacterVirtualSettings_Init(p)
	return &CharacterVirtualSettings{ptr: p, ownShape: true}
}

// Destroy releases the settings, including the default shape allocated by Init
// unless the caller replaced it.
func (s *CharacterVirtualSettings) Destroy() {
	if s.ptr != nil {
		if s.ownShape && s.ptr.base.shape != nil {
			C.JPH_Shape_Destroy(s.ptr.base.shape)
		}
		C.free(unsafe.Pointer(s.ptr))
		s.ptr = nil
	}
}

// Mass returns the character mass.
func (s *CharacterVirtualSettings) Mass() float32 { return float32(s.ptr.mass) }

// SetShape sets the character shape, releasing the default shape allocated by
// Init if it was still owned.
func (s *CharacterVirtualSettings) SetShape(shape *Shape) {
	if s.ownShape && s.ptr.base.shape != nil {
		C.JPH_Shape_Destroy(s.ptr.base.shape)
	}
	s.ptr.base.shape = shape.ptr
	s.ownShape = false
}
