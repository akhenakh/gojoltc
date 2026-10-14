// Command example is a visual recreation of joltc's samples/01_HelloWorld.
//
// It builds a small scene (a large floor, a falling sphere, a stack of boxes
// and a player character), steps the simulation, and renders Jolt's debug
// geometry through raylib.
//
// Run from this directory: go run .
//
// Controls:
//
//	WASD   move the character
//	Space  pause/resume
//	F      toggle wireframe/solid
//	R      reset the scene
//	Esc    quit
package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/akhenakh/gojoltc/jolt"
)

const (
	layerNonMoving jolt.ObjectLayer = 0
	layerMoving    jolt.ObjectLayer = 1
	numLayers                       = 2
)

// joltColor decodes a Jolt packed RGBA color.
func joltColor(c jolt.Color) rl.Color {
	return rl.Color{
		R: uint8(c & 0xff),
		G: uint8((c >> 8) & 0xff),
		B: uint8((c >> 16) & 0xff),
		A: uint8((c >> 24) & 0xff),
	}
}

func vec3(v jolt.RVec3) rl.Vector3 {
	return rl.Vector3{X: float32(v.X), Y: float32(v.Y), Z: float32(v.Z)}
}

// debugGeometry accumulates the lines and triangles Jolt asks us to draw.
type debugGeometry struct {
	lines     []rl.Vector3 // pairs
	triangles []rl.Vector3 // triples
	colors    []rl.Color
}

func (g *debugGeometry) renderer() *jolt.DebugRenderer {
	return jolt.NewDebugRenderer(&jolt.DebugRendererCallbacks{
		DrawLine: func(from, to jolt.RVec3, color jolt.Color) {
			c := joltColor(color)
			g.lines = append(g.lines, vec3(from), vec3(to))
			g.colors = append(g.colors, c, c)
		},
		DrawTriangle: func(v1, v2, v3 jolt.RVec3, color jolt.Color, castShadow bool) {
			c := joltColor(color)
			g.triangles = append(g.triangles, vec3(v1), vec3(v2), vec3(v3))
			g.colors = append(g.colors, c, c, c)
		},
	})
}

func (g *debugGeometry) reset() {
	g.lines = g.lines[:0]
	g.triangles = g.triangles[:0]
	g.colors = g.colors[:0]
}

func main() {
	if err := jolt.Init(); err != nil {
		panic(err)
	}
	defer jolt.Shutdown()

	jobSystem := jolt.NewJobSystemThreadPool(nil)
	defer jobSystem.Destroy()

	bpLayer := jolt.NewBroadPhaseLayerInterfaceTable(numLayers, numLayers)
	bpLayer.MapObjectToBroadPhaseLayer(layerNonMoving, 0)
	bpLayer.MapObjectToBroadPhaseLayer(layerMoving, 1)
	pairFilter := jolt.NewObjectLayerPairFilterTable(numLayers)
	pairFilter.EnableCollision(layerNonMoving, layerMoving)
	pairFilter.EnableCollision(layerMoving, layerMoving)
	objBpFilter := jolt.NewObjectVsBroadPhaseLayerFilterTable(bpLayer, numLayers, pairFilter, numLayers)

	system := jolt.NewPhysicsSystem(&jolt.PhysicsSystemSettings{
		MaxBodies:                     8192,
		MaxBodyPairs:                  8192,
		MaxContactConstraints:         2048,
		BroadPhaseLayerInterface:      bpLayer,
		ObjectLayerPairFilter:         pairFilter,
		ObjectVsBroadPhaseLayerFilter: objBpFilter,
	})
	defer system.Destroy()
	system.SetGravity(jolt.Vec3{0, -10, 0})

	geometry := &debugGeometry{}
	renderer := geometry.renderer()
	defer renderer.Destroy()

	buildScene := func() *jolt.Character {
		bi := system.BodyInterface()

		// Floor.
		floorShape := jolt.NewBoxShape(jolt.Vec3{50, 1, 50}, jolt.DefaultConvexRadius)
		floorSettings := jolt.NewBodyCreationSettings(floorShape, jolt.RVec3{0, -1, 0}, jolt.QuatIdentity(), jolt.MotionTypeStatic, layerNonMoving)
		bi.CreateAndAddBody(floorSettings, jolt.ActivationDontActivate)
		floorSettings.Destroy()
		floorShape.Destroy()

		// A falling sphere with an initial velocity.
		sphereShape := jolt.NewSphereShape(0.5)
		sphereSettings := jolt.NewBodyCreationSettings(sphereShape, jolt.RVec3{0, 8, 0}, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerMoving)
		sphereID := bi.CreateAndAddBody(sphereSettings, jolt.ActivationActivate)
		sphereSettings.Destroy()
		sphereShape.Destroy()
		bi.SetLinearVelocity(sphereID, jolt.Vec3{0, -5, 0})

		// A stack of boxes.
		boxShape := jolt.NewBoxShape(jolt.Vec3{0.5, 0.5, 0.5}, jolt.DefaultConvexRadius)
		for i := range 6 {
			pos := jolt.RVec3{-4, 0.5 + float32(i), 0}
			bs := jolt.NewBodyCreationSettings(boxShape, pos, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerMoving)
			bi.CreateAndAddBody(bs, jolt.ActivationActivate)
			bs.Destroy()
		}
		boxShape.Destroy()

		// Player character (a rigid capsule character).
		capsule := jolt.NewCapsuleShape(0.5, 0.3)
		characterSettings := jolt.NewCharacterSettings()
		characterSettings.SetShape(capsule)
		character := jolt.NewCharacter(characterSettings, jolt.RVec3{4, 1.5, 0}, jolt.QuatIdentity(), 0, system)
		character.AddToPhysicsSystem(jolt.ActivationActivate, true)
		characterSettings.Destroy()
		capsule.Destroy()

		system.OptimizeBroadPhase()
		return character
	}

	character := buildScene()
	defer character.Destroy()
	characterID := character.BodyID()

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(1280, 720, "gojoltc — joltc 01_HelloWorld")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	cam := newOrbitCamera(rl.Vector3{X: 0, Y: 1, Z: 0}, 0.8, 0.5, 18)

	wireframe := true
	paused := false

	for !rl.WindowShouldClose() {
		switch {
		case rl.IsKeyPressed(rl.KeySpace):
			paused = !paused
		case rl.IsKeyPressed(rl.KeyF):
			wireframe = !wireframe
		}

		// Drive the character with WASD.
		vel := jolt.Vec3{}
		if rl.IsKeyDown(rl.KeyW) {
			vel.Z -= 4
		}
		if rl.IsKeyDown(rl.KeyS) {
			vel.Z += 4
		}
		if rl.IsKeyDown(rl.KeyA) {
			vel.X -= 4
		}
		if rl.IsKeyDown(rl.KeyD) {
			vel.X += 4
		}
		character.SetLinearVelocity(vel, true)

		if !paused {
			system.Update(rl.GetFrameTime(), 1, jobSystem)
		}

		geometry.reset()
		settings := jolt.DrawSettings{DrawShape: true, DrawShapeWireframe: wireframe}
		system.DrawBodies(settings, renderer)

		camera := cam.camera()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		if !wireframe {
			for i := 0; i+2 < len(geometry.triangles); i += 3 {
				rl.DrawTriangle3D(geometry.triangles[i], geometry.triangles[i+1], geometry.triangles[i+2], geometry.colors[i])
			}
		}
		for i := 0; i+1 < len(geometry.lines); i += 2 {
			rl.DrawLine3D(geometry.lines[i], geometry.lines[i+1], geometry.colors[i])
		}
		rl.DrawGrid(50, 1)

		rl.EndMode3D()

		pos := character.Position(true)
		rl.DrawText(fmt.Sprintf("character: (%.1f, %.1f, %.1f)", pos.X, pos.Y, pos.Z), 10, 10, 20, rl.DarkGray)
		rl.DrawText("WASD move  Space pause  F wireframe  drag to orbit, wheel to zoom", 10, 34, 20, rl.Gray)

		rl.EndDrawing()
	}

	_ = characterID
}
