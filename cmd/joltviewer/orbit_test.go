package main

import (
	"math"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// A regression guard for the orbit controls: raylib's CameraOrbital mode ignores
// the mouse, so the camera is driven by this type instead. Only the pure math is
// checked here (no window needed).
func TestOrbitCameraPosition(t *testing.T) {
	target := rl.Vector3{X: 0, Y: 1, Z: 0}
	o := newOrbitCamera(target, 0, 0.5, 10)

	p := o.position()
	if d := math.Hypot(math.Hypot(float64(p.X-target.X), float64(p.Y-target.Y)), float64(p.Z-target.Z)); math.Abs(d-10) > 1e-4 {
		t.Errorf("camera distance = %v, want 10", d)
	}
	if p.Y <= target.Y {
		t.Errorf("pitch 0.5 should put the camera above the target, got y = %v", p.Y)
	}

	// Yaw must actually move the camera around the target.
	o.yaw = math.Pi / 2
	q := o.position()
	if q.X == p.X && q.Z == p.Z {
		t.Error("yaw has no effect on the camera position")
	}

	// Zoom must actually change the distance.
	o.distance = 20
	if r := o.position(); math.Hypot(float64(r.X-target.X), float64(r.Z-target.Z)) == math.Hypot(float64(p.X-target.X), float64(p.Z-target.Z)) {
		t.Error("distance has no effect on the camera position")
	}

	if clampFloat(2, 0.05, 1.5) != 1.5 || clampFloat(0, 0.05, 1.5) != 0.05 {
		t.Error("pitch clamping is wrong")
	}
}
