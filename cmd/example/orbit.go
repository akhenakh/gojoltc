package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// orbitCamera is a turntable camera: drag with either mouse button to orbit and
// use the wheel to zoom.
//
// This is implemented here rather than with raylib's CameraOrbital mode because
// that mode ignores the mouse entirely (it slowly rotates around the target on
// its own), while the other modes move the camera on WASD instead of the mouse.
type orbitCamera struct {
	target     rl.Vector3
	yaw, pitch float32 // radians
	distance   float32
}

func newOrbitCamera(target rl.Vector3, yaw, pitch, distance float32) *orbitCamera {
	return &orbitCamera{target: target, yaw: yaw, pitch: pitch, distance: distance}
}

// position is the pure orbit math: the camera sits at distance from the target,
// at the current yaw/pitch. Kept free of input so it can be tested without a
// window.
func (o *orbitCamera) position() rl.Vector3 {
	cp, sp := float32(math.Cos(float64(o.pitch))), float32(math.Sin(float64(o.pitch)))
	cy, sy := float32(math.Cos(float64(o.yaw))), float32(math.Sin(float64(o.yaw)))
	return rl.Vector3{
		X: o.target.X + o.distance*cp*sy,
		Y: o.target.Y + o.distance*sp,
		Z: o.target.Z + o.distance*cp*cy,
	}
}

// updateInput applies mouse drag and wheel zoom.
func (o *orbitCamera) updateInput() {
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) || rl.IsMouseButtonDown(rl.MouseButtonRight) {
		d := rl.GetMouseDelta()
		o.yaw -= d.X * 0.006
		o.pitch = clampFloat(o.pitch+d.Y*0.006, 0.05, 1.5) // stop short of vertical
	}
	if w := rl.GetMouseWheelMove(); w != 0 {
		o.distance = clampFloat(o.distance*(1-w*0.1), 2, 150)
	}
}

// camera applies this frame's input and returns the raylib camera.
func (o *orbitCamera) camera() rl.Camera3D {
	o.updateInput()
	return rl.Camera3D{
		Position:   o.position(),
		Target:     o.target,
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}
}

func clampFloat(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
