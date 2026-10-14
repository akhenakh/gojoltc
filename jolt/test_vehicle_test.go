package jolt

import "testing"

// Ported from third_party/joltc/tests/test_vehicle.cpp.

func initVehicle(t *testing.T) {
	t.Helper()
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	t.Cleanup(Shutdown)
}

func TestVehicleWheelSettingsCreate(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()
	if settings == nil {
		t.Fatal("settings is nil")
	}
}

func TestVehicleWheelSettingsPosition(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()

	settings.SetPosition(Vec3{1, 0, 2})
	p := settings.Position()
	equalFloat(t, "x", p.X, 1)
	equalFloat(t, "z", p.Z, 2)
}

func TestVehicleWheelSettingsRadius(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()

	settings.SetRadius(0.5)
	equalFloat(t, "radius", settings.Radius(), 0.5)
}

func TestVehicleWheelSettingsWidth(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()

	settings.SetWidth(0.3)
	equalFloat(t, "width", settings.Width(), 0.3)
}

func TestVehicleWheelSettingsSuspension(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()

	settings.SetSuspensionMinLength(0.1)
	settings.SetSuspensionMaxLength(0.5)
	equalFloat(t, "min", settings.SuspensionMinLength(), 0.1)
	equalFloat(t, "max", settings.SuspensionMaxLength(), 0.5)
}

func TestVehicleWheelCreate(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()
	wheel := NewWheel(settings)
	defer wheel.Destroy()
	if wheel == nil {
		t.Fatal("wheel is nil")
	}
}

func TestVehicleWheelAngularVelocity(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()
	wheel := NewWheel(settings)
	defer wheel.Destroy()

	wheel.SetAngularVelocity(10)
	equalFloat(t, "angular velocity", wheel.AngularVelocity(), 10)
}

func TestVehicleWheelRotationAngle(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()
	wheel := NewWheel(settings)
	defer wheel.Destroy()

	wheel.SetRotationAngle(Pi / 4)
	nearFloat(t, "rotation angle", wheel.RotationAngle(), Pi/4)
}

func TestVehicleWheelSteerAngle(t *testing.T) {
	initVehicle(t)
	settings := NewWheelSettings()
	defer settings.Destroy()
	wheel := NewWheel(settings)
	defer wheel.Destroy()

	wheel.SetSteerAngle(0.1)
	equalFloat(t, "steer angle", wheel.SteerAngle(), 0.1)
}

func TestVehicleCollisionTesterRayCreate(t *testing.T) {
	initVehicle(t)
	tester := NewVehicleCollisionTesterRay(layerMoving, Vec3{0, 1, 0}, Pi/4)
	defer tester.Destroy()
	if tester == nil {
		t.Fatal("tester is nil")
	}
	if tester.ObjectLayer() != layerMoving {
		t.Errorf("object layer = %d, want %d", tester.ObjectLayer(), layerMoving)
	}
}

func TestVehicleCollisionTesterCastSphereCreate(t *testing.T) {
	initVehicle(t)
	tester := NewVehicleCollisionTesterCastSphere(layerMoving, 0.1, Vec3{0, 1, 0}, Pi/4)
	defer tester.Destroy()
	if tester == nil {
		t.Fatal("tester is nil")
	}
}

func TestVehicleCollisionTesterCastCylinderCreate(t *testing.T) {
	initVehicle(t)
	tester := NewVehicleCollisionTesterCastCylinder(layerMoving, 0.05)
	defer tester.Destroy()
	if tester == nil {
		t.Fatal("tester is nil")
	}
}

func TestVehicleEngineSettingsInit(t *testing.T) {
	settings := NewVehicleEngineSettings()
	if settings.MaxRPM <= 0 {
		t.Errorf("maxRPM = %v, want > 0", settings.MaxRPM)
	}
	if settings.MaxTorque <= 0 {
		t.Errorf("maxTorque = %v, want > 0", settings.MaxTorque)
	}
}

func TestVehicleDifferentialSettingsInit(t *testing.T) {
	settings := NewVehicleDifferentialSettings()
	if settings.LeftWheel < -1 {
		t.Errorf("leftWheel = %d, want >= -1", settings.LeftWheel)
	}
}

func TestVehicleTransmissionSettingsCreate(t *testing.T) {
	initVehicle(t)
	settings := NewVehicleTransmissionSettings()
	defer settings.Destroy()
	if settings == nil {
		t.Fatal("settings is nil")
	}
}

func TestVehicleTransmissionSettingsMode(t *testing.T) {
	initVehicle(t)
	settings := NewVehicleTransmissionSettings()
	defer settings.Destroy()

	settings.SetMode(TransmissionModeManual)
	if settings.Mode() != TransmissionModeManual {
		t.Error("mode should be manual")
	}
	settings.SetMode(TransmissionModeAuto)
	if settings.Mode() != TransmissionModeAuto {
		t.Error("mode should be auto")
	}
}

func TestVehicleTransmissionSettingsGearRatios(t *testing.T) {
	initVehicle(t)
	settings := NewVehicleTransmissionSettings()
	defer settings.Destroy()

	settings.SetGearRatios([]float32{3, 2, 1.5, 1, 0.8})
	if n := settings.GearRatioCount(); n != 5 {
		t.Errorf("gear ratio count = %d, want 5", n)
	}
	equalFloat(t, "gear 0", settings.GearRatio(0), 3)
	equalFloat(t, "gear 4", settings.GearRatio(4), 0.8)
}

func TestVehicleTransmissionSettingsTiming(t *testing.T) {
	initVehicle(t)
	settings := NewVehicleTransmissionSettings()
	defer settings.Destroy()

	settings.SetSwitchTime(0.5)
	equalFloat(t, "switch time", settings.SwitchTime(), 0.5)
	settings.SetClutchReleaseTime(0.3)
	equalFloat(t, "clutch release time", settings.ClutchReleaseTime(), 0.3)
	settings.SetShiftUpRPM(5000)
	equalFloat(t, "shift up rpm", settings.ShiftUpRPM(), 5000)
	settings.SetShiftDownRPM(2000)
	equalFloat(t, "shift down rpm", settings.ShiftDownRPM(), 2000)
}

func TestVehicleWheeledControllerSettingsCreate(t *testing.T) {
	initVehicle(t)
	settings := NewWheeledVehicleControllerSettings()
	defer settings.Destroy()
	if settings == nil {
		t.Fatal("settings is nil")
	}
}

func TestVehicleWheeledControllerSettingsEngine(t *testing.T) {
	initVehicle(t)
	settings := NewWheeledVehicleControllerSettings()
	defer settings.Destroy()

	engine := NewVehicleEngineSettings()
	engine.MaxRPM = 8000
	engine.MaxTorque = 500
	settings.SetEngine(engine)

	result := settings.Engine()
	equalFloat(t, "maxRPM", result.MaxRPM, 8000)
	equalFloat(t, "maxTorque", result.MaxTorque, 500)
}

func TestVehicleAntiRollBarInit(t *testing.T) {
	bar := NewVehicleAntiRollBar()
	if bar.LeftWheel < 0 {
		t.Errorf("leftWheel = %d, want >= 0", bar.LeftWheel)
	}
}

func TestVehicleConstraintSettingsInit(t *testing.T) {
	settings := NewVehicleConstraintSettings()
	f := settings.Forward
	if f.X == 0 && f.Y == 0 && f.Z == 0 {
		t.Error("forward should be non-zero after Init")
	}
}
