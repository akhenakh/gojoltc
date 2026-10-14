package jolt

import (
	"math"
	"testing"
)

// Headless version of cmd/rover: the same height field, wheel/friction setup and
// drive input, asserted instead of drawn.

func terrainAt(x, z float64) float32 {
	h := 0.9 * math.Sin(x*0.32+0.123) * math.Cos(z*0.27)
	h += 0.45 * math.Sin(x*0.85+1.3) * math.Sin(z*0.7-0.4)
	return float32(h)
}

func TestRoverDrivesOnTerrain(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	// --- terrain -----------------------------------------------------------------
	const samples = 33
	const size = 32.0
	step := float32(size) / float32(samples-1)
	heights := make([]float32, samples*samples)
	for j := range samples {
		for i := range samples {
			x := (float64(i) - samples/2) * float64(step)
			z := (float64(j) - samples/2) * float64(step)
			heights[j*samples+i] = terrainAt(x, z)
		}
	}
	terrain := NewHeightFieldShape(heights, Vec3{-size / 2, 0, -size / 2}, Vec3{step, 1, step})
	if terrain == nil {
		t.Fatal("terrain shape is nil")
	}
	defer terrain.Destroy()

	terrainSettings := NewBodyCreationSettings(terrain, RVec3{0, 0, 0}, QuatIdentity(), MotionTypeStatic, layerNonMoving)
	terrainSettings.SetFriction(0.9)
	terrainID := bi.CreateAndAddBody(terrainSettings, ActivationDontActivate)
	terrainSettings.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(terrainID) })

	// --- rover -------------------------------------------------------------------
	chassisShape := NewBoxShape(Vec3{0.6, 0.25, 0.45}, DefaultConvexRadius)
	chassisSettings := NewBodyCreationSettings(chassisShape, RVec3{0, 3, 0}, QuatIdentity(), MotionTypeDynamic, layerMoving)
	chassisSettings.SetFriction(0.4)
	chassisSettings.SetAllowSleeping(false)
	chassisID := bi.CreateAndAddBody(chassisSettings, ActivationActivate)
	chassisSettings.Destroy()
	chassisShape.Destroy()
	t.Cleanup(func() { bi.RemoveAndDestroyBody(chassisID) })

	curve := NewLinearCurve()
	curve.AddPoint(0, 0)
	curve.AddPoint(0.1, 1.0)
	curve.AddPoint(1.0, 0.55)
	curve.Sort()

	wheelPositions := []Vec3{
		{-0.55, -0.2, 0.5}, {0.55, -0.2, 0.5},
		{-0.55, -0.2, -0.5}, {0.55, -0.2, -0.5},
	}
	wheels := make([]*WheelSettings, 0, len(wheelPositions))
	for _, p := range wheelPositions {
		w := NewWheelSettingsWV()
		w.SetPosition(p)
		w.SetRadius(0.28)
		w.SetWidth(0.2)
		w.SetSuspensionMinLength(0.05)
		w.SetSuspensionMaxLength(0.35)
		w.SetMaxSteerAngle(0.55)
		w.SetMaxBrakeTorque(1200)
		w.SetInertia(0.4)
		w.SetLongitudinalFrictionCurve(curve)
		w.SetLateralFrictionCurve(curve)
		wheels = append(wheels, w)
		t.Cleanup(w.Destroy)
	}
	curve.Destroy() // the wheels copied it

	controller := NewWheeledVehicleControllerSettings()
	controller.AddDifferential(0, 1)
	controller.AddDifferential(2, 3)
	t.Cleanup(controller.Destroy)

	settings := NewVehicleConstraintSettings()
	settings.MaxPitchRollAngle = 0.6
	vehicle := NewVehicleConstraint(f.system.BodyPtr(chassisID), &settings, wheels, controller)
	tester := NewVehicleCollisionTesterRay(layerMoving, Vec3{0, 1, 0}, 0.9)
	t.Cleanup(tester.Destroy)
	vehicle.SetCollisionTester(tester)

	f.system.AddConstraint(vehicle.AsConstraint())
	stepListener := vehicle.StepListener()
	f.system.AddStepListener(stepListener)
	t.Cleanup(func() {
		// Remove before destroying, in the order the system expects.
		f.system.RemoveStepListener(stepListener)
		f.system.RemoveConstraint(vehicle.AsConstraint())
		vehicle.Destroy()
	})

	f.system.OptimizeBroadPhase()

	// Let the rover settle onto the terrain.
	for range 60 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	settled := bi.Position(chassisID)
	if settled.Y < 0 {
		t.Fatalf("rover fell through the terrain: y = %v", settled.Y)
	}

	// Drive forward for two seconds.
	vehicle.Controller().SetDriverInput(1, 0, 0, 0)
	for range 120 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	end := bi.Position(chassisID)

	travelled := math.Hypot(float64(end.X-settled.X), float64(end.Z-settled.Z))
	if travelled < 0.5 {
		t.Errorf("rover did not drive on the terrain: moved %.3f m", travelled)
	}
	if end.Y < -1 {
		t.Errorf("rover sank through the terrain: y = %v", end.Y)
	}
	// It should still be on the terrain, not flung into the air.
	if end.Y > settled.Y+3 {
		t.Errorf("rover launched off the terrain: y = %v (settled at %v)", end.Y, settled.Y)
	}
}
