package jolt

import "testing"

// Ported from third_party/joltc/tests/test_character.cpp.

func TestCharacterSettingsInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	settings := NewCharacterSettings()
	defer settings.Destroy()
	if settings.Up().Y == 0 {
		t.Error("up.y should be non-zero after Init")
	}
}

func TestCharacterCreate(t *testing.T) {
	f := newPhysicsFixture(t)

	shape := NewCapsuleShape(0.5, 0.3)
	defer shape.Destroy()

	settings := NewCharacterSettings()
	defer settings.Destroy()
	settings.SetShape(shape)

	character := NewCharacter(settings, RVec3{0, 1, 0}, QuatIdentity(), 0, f.system)
	if character == nil {
		t.Fatal("character is nil")
	}
	defer character.Destroy()

	character.AddToPhysicsSystem(ActivationActivate, true)
	if id := character.BodyID(); id == invalidBodyID {
		t.Errorf("body id is invalid (%#x)", uint32(id))
	}
	character.RemoveFromPhysicsSystem(true)
}

func TestCharacterPositionAndRotation(t *testing.T) {
	f := newPhysicsFixture(t)

	shape := NewCapsuleShape(0.5, 0.3)
	defer shape.Destroy()

	settings := NewCharacterSettings()
	defer settings.Destroy()
	settings.SetShape(shape)

	character := NewCharacter(settings, RVec3{1, 2, 3}, QuatIdentity(), 0, f.system)
	defer character.Destroy()
	character.AddToPhysicsSystem(ActivationActivate, true)

	p := character.Position(true)
	equalReal(t, "pos.x", p.X, 1)
	equalReal(t, "pos.y", p.Y, 2)
	equalReal(t, "pos.z", p.Z, 3)

	character.SetPosition(RVec3{5, 6, 7}, ActivationActivate, true)
	p = character.Position(true)
	equalReal(t, "new pos.x", p.X, 5)
	equalReal(t, "new pos.y", p.Y, 6)
	equalReal(t, "new pos.z", p.Z, 7)

	character.RemoveFromPhysicsSystem(true)
}

func TestCharacterVelocity(t *testing.T) {
	f := newPhysicsFixture(t)

	shape := NewCapsuleShape(0.5, 0.3)
	defer shape.Destroy()

	settings := NewCharacterSettings()
	defer settings.Destroy()
	settings.SetShape(shape)

	character := NewCharacter(settings, RVec3{0, 1, 0}, QuatIdentity(), 0, f.system)
	defer character.Destroy()
	character.AddToPhysicsSystem(ActivationActivate, true)

	character.SetLinearVelocity(Vec3{1, 0, 0}, true)
	equalFloat(t, "vel.x", character.LinearVelocity().X, 1)

	character.RemoveFromPhysicsSystem(true)
}

func TestCharacterVirtualSettingsInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	settings := NewCharacterVirtualSettings()
	defer settings.Destroy()
	if settings.Mass() <= 0 {
		t.Errorf("mass = %v, want > 0", settings.Mass())
	}
}
