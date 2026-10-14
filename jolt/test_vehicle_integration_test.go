package jolt

import (
	"math"
	"testing"
)

// Behavioural test: build a real wheeled vehicle, drive it, and check it moves.
// The ported upstream tests only exercise vehicle *settings*, which proves the
// constraint can never actually be created or added.

func TestVehicleDrives(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	// Floor with its top at y = 0.
	floorShape := NewBoxShape(Vec3{100, 1, 100}, DefaultConvexRadius)
	floorSettings := NewBodyCreationSettings(floorShape, RVec3{0, -1, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	floorID := bi.CreateAndAddBody(floorSettings, ActivationDontActivate)
	floorSettings.Destroy()
	floorShape.Destroy()

	// Chassis.
	chassisShape := NewBoxShape(Vec3{1.5, 0.5, 1.0}, DefaultConvexRadius)
	chassisSettings := NewBodyCreationSettings(chassisShape, RVec3{0, 0.8, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	chassisID := bi.CreateAndAddBody(chassisSettings, ActivationActivate)
	chassisSettings.Destroy()
	chassisShape.Destroy()
	t.Cleanup(func() {
		bi.RemoveAndDestroyBody(chassisID)
		bi.RemoveAndDestroyBody(floorID)
	})

	// Four wheels near the corners.
	wheelPositions := []Vec3{
		{-1, -0.3, 0.8}, {1, -0.3, 0.8},
		{-1, -0.3, -0.8}, {1, -0.3, -0.8},
	}
	wheels := make([]*WheelSettings, 0, len(wheelPositions))
	for _, p := range wheelPositions {
		w := NewWheelSettingsWV()
		w.SetPosition(p)
		w.SetRadius(0.4)
		w.SetWidth(0.3)
		w.SetSuspensionMinLength(0.1)
		w.SetSuspensionMaxLength(0.5)
		w.SetMaxSteerAngle(0.5)
		w.SetMaxBrakeTorque(3000)
		w.SetMaxHandBrakeTorque(6000)
		wheels = append(wheels, w)
		t.Cleanup(w.Destroy)
	}

	controller := NewWheeledVehicleControllerSettings()
	// A differential per axle connects the engine to the wheels; without one
	// no torque reaches any wheel and the vehicle never moves.
	controller.AddDifferential(0, 1)
	controller.AddDifferential(2, 3)
	t.Cleanup(controller.Destroy)

	settings := NewVehicleConstraintSettings()
	vehicle := NewVehicleConstraint(f.system.BodyPtr(chassisID), &settings, wheels, controller)
	if vehicle == nil {
		t.Fatal("vehicle constraint is nil")
	}
	if n := vehicle.WheelsCount(); n != uint32(len(wheels)) {
		t.Fatalf("wheels = %d, want %d", n, len(wheels))
	}

	// The tester must use the vehicle's own object layer: the default layer
	// filter only permits layers that can collide with it (layerMoving x
	// layerNonMoving), which is what reaches the floor.
	tester := NewVehicleCollisionTesterRay(layerMoving, Vec3{0, 1, 0}, Pi/4)
	defer tester.Destroy()
	vehicle.SetCollisionTester(tester)

	// Add the constraint and its step listener to the system.
	f.system.AddConstraint(vehicle.AsConstraint())
	stepListener := vehicle.StepListener()
	f.system.AddStepListener(stepListener)

	t.Cleanup(func() {
		f.system.RemoveStepListener(stepListener)
		f.system.RemoveConstraint(vehicle.AsConstraint())
		vehicle.Destroy()
	})

	f.system.OptimizeBroadPhase()

	// Let the vehicle settle, then drive forward (+Z).
	for range 60 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	start := bi.Position(chassisID)

	vehicle.Controller().SetDriverInput(1, 0, 0, 0)
	for range 180 {
		// A settling vehicle can fall asleep; keep the chassis awake while driven.
		bi.ActivateBody(chassisID)
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	end := bi.Position(chassisID)

	travelled := math.Hypot(float64(end.X-start.X), float64(end.Z-start.Z))
	if travelled < 0.5 {
		t.Errorf("vehicle did not drive: moved %.3f m in 3s (start %v, end %v)", travelled, start, end)
	}
	if end.Y < -1 {
		t.Errorf("vehicle fell through the floor: y = %v", end.Y)
	}
}
