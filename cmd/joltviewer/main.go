// Command joltviewer is a raylib viewer for JoltPhysics simulations, in the
// spirit of JoltPhysics' JoltViewer: an orbit camera over a live simulation,
// rendering Jolt's debug geometry, with pause/step and scene switching.
//
// Difference from the C++ JoltViewer: it runs a live simulation rather than
// replaying a recorded .jor ObjectStream (that API is not wrapped here).
//
// Run from this directory: go run .
//
// Controls:
//
//	mouse      orbit / zoom
//	N          next scene
//	R          reset current scene
//	Space      pause/resume
//	Right      single step
//	F          wireframe / solid
//	P          save a screenshot (joltviewer.png)
//	Esc        quit
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/akhenakh/gojoltc/jolt"
)

const (
	layerNonMoving jolt.ObjectLayer = 0
	layerMoving    jolt.ObjectLayer = 1
	numLayers                       = 2
)

func joltColor(c jolt.Color) rl.Color {
	return rl.Color{
		R: uint8(c & 0xff),
		G: uint8((c >> 8) & 0xff),
		B: uint8((c >> 16) & 0xff),
		A: uint8((c >> 24) & 0xff),
	}
}

func toVec3(v jolt.RVec3) rl.Vector3 {
	return rl.Vector3{X: float32(v.X), Y: float32(v.Y), Z: float32(v.Z)}
}

type geometry struct {
	lines     []rl.Vector3
	lineColor []rl.Color
	tris      []rl.Vector3
	triColor  []rl.Color
}

func (g *geometry) reset() {
	g.lines = g.lines[:0]
	g.lineColor = g.lineColor[:0]
	g.tris = g.tris[:0]
	g.triColor = g.triColor[:0]
}

func (g *geometry) newRenderer() *jolt.DebugRenderer {
	return jolt.NewDebugRenderer(&jolt.DebugRendererCallbacks{
		DrawLine: func(from, to jolt.RVec3, color jolt.Color) {
			c := joltColor(color)
			g.lines = append(g.lines, toVec3(from), toVec3(to))
			g.lineColor = append(g.lineColor, c, c)
		},
		DrawTriangle: func(v1, v2, v3 jolt.RVec3, color jolt.Color, castShadow bool) {
			c := joltColor(color)
			g.tris = append(g.tris, toVec3(v1), toVec3(v2), toVec3(v3))
			g.triColor = append(g.triColor, c, c, c)
		},
	})
}

type world struct {
	system    *jolt.PhysicsSystem
	jobSystem *jolt.JobSystem
}

func newWorld() *world {
	jobSystem := jolt.NewJobSystemThreadPool(nil)

	bpLayer := jolt.NewBroadPhaseLayerInterfaceTable(numLayers, numLayers)
	bpLayer.MapObjectToBroadPhaseLayer(layerNonMoving, 0)
	bpLayer.MapObjectToBroadPhaseLayer(layerMoving, 1)
	pairFilter := jolt.NewObjectLayerPairFilterTable(numLayers)
	pairFilter.EnableCollision(layerNonMoving, layerMoving)
	pairFilter.EnableCollision(layerMoving, layerMoving)
	objBpFilter := jolt.NewObjectVsBroadPhaseLayerFilterTable(bpLayer, numLayers, pairFilter, numLayers)

	system := jolt.NewPhysicsSystem(&jolt.PhysicsSystemSettings{
		MaxBodies:                     16384,
		MaxBodyPairs:                  16384,
		MaxContactConstraints:         4096,
		BroadPhaseLayerInterface:      bpLayer,
		ObjectLayerPairFilter:         pairFilter,
		ObjectVsBroadPhaseLayerFilter: objBpFilter,
	})
	system.SetGravity(jolt.Vec3{0, -10, 0})
	return &world{system: system, jobSystem: jobSystem}
}

func (w *world) destroy() {
	w.system.Destroy()
	w.jobSystem.Destroy()
}

// --- scenes -----------------------------------------------------------------

var sceneNames = []string{"Pyramid", "Sphere Drop", "Soft Body Cloth", "Vehicle"}

func addFloor(bi *jolt.BodyInterface) {
	shape := jolt.NewBoxShape(jolt.Vec3{50, 1, 50}, jolt.DefaultConvexRadius)
	settings := jolt.NewBodyCreationSettings(shape, jolt.RVec3{0, -1, 0}, jolt.QuatIdentity(), jolt.MotionTypeStatic, layerNonMoving)
	bi.CreateAndAddBody(settings, jolt.ActivationDontActivate)
	settings.Destroy()
	shape.Destroy()
}

func buildPyramid(bi *jolt.BodyInterface) {
	addFloor(bi)
	shape := jolt.NewBoxShape(jolt.Vec3{0.5, 0.5, 0.5}, jolt.DefaultConvexRadius)
	defer shape.Destroy()
	for row := range 8 {
		for col := range 8 - row {
			pos := jolt.RVec3{float32(col)*1.05 + float32(row)*0.525 - 3.5, float32(row) + 0.5, 0}
			settings := jolt.NewBodyCreationSettings(shape, pos, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerMoving)
			bi.CreateAndAddBody(settings, jolt.ActivationActivate)
			settings.Destroy()
		}
	}
}

func buildSphereDrop(bi *jolt.BodyInterface) {
	addFloor(bi)
	shape := jolt.NewSphereShape(0.5)
	defer shape.Destroy()
	for i := range 120 {
		x := float32(i%10)*1.2 - 5.4
		z := float32(i/10)*1.2 - 3
		pos := jolt.RVec3{x, float32(2 + i%7), z}
		settings := jolt.NewBodyCreationSettings(shape, pos, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerMoving)
		bi.CreateAndAddBody(settings, jolt.ActivationActivate)
		settings.Destroy()
	}
}

func buildSoftCloth(bi *jolt.BodyInterface) {
	addFloor(bi)

	const n = 12
	shared := jolt.NewSoftBodySharedSettings()
	defer shared.Destroy()

	vertices := make([]jolt.SoftVertex, 0, n*n)
	for j := range n {
		for i := range n {
			// The first row is pinned (zero inverse mass), so the cloth hangs
			// like a curtain instead of falling to the floor.
			invMass := float32(1)
			if j == 0 {
				invMass = 0
			}
			vertices = append(vertices, jolt.SoftVertex{
				Position: jolt.Vec3{float32(i)*0.25 - 1.5, 3, float32(j)*0.25 - 1.5},
				InvMass:  invMass,
			})
		}
	}
	shared.AddVertices(vertices)

	faces := make([]jolt.SoftFace, 0, (n-1)*(n-1)*2)
	for j := range n - 1 {
		for i := range n - 1 {
			a := uint32(j*n + i)
			faces = append(faces,
				jolt.SoftFace{Vertex1: a, Vertex2: a + 1, Vertex3: a + n},
				jolt.SoftFace{Vertex1: a + 1, Vertex2: a + n + 1, Vertex3: a + n},
			)
		}
	}
	shared.AddFaces(faces)
	shared.CreateConstraints(1.0e-4, jolt.SoftBodyBendTypeDistance)
	shared.Optimize()

	settings := jolt.NewSoftBodyCreationSettingsFor(shared, jolt.RVec3{0, 3, 0}, jolt.QuatIdentity(), layerMoving)
	defer settings.Destroy()
	settings.SetFriction(0.5)
	settings.SetVertexRadius(0.05)
	bi.CreateAndAddSoftBody(settings, jolt.ActivationActivate)
}

// buildVehicle builds a drivable vehicle and returns a function that applies
// driver input each frame.
func buildVehicle(system *jolt.PhysicsSystem) func() {
	bi := system.BodyInterface()

	chassisShape := jolt.NewBoxShape(jolt.Vec3{1.5, 0.5, 1.0}, jolt.DefaultConvexRadius)
	chassisSettings := jolt.NewBodyCreationSettings(chassisShape, jolt.RVec3{0, 0.8, 0}, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerMoving)
	chassisID := bi.CreateAndAddBody(chassisSettings, jolt.ActivationActivate)
	chassisSettings.Destroy()
	chassisShape.Destroy()

	wheels := make([]*jolt.WheelSettings, 0, 4)
	for _, p := range []jolt.Vec3{{-1, -0.3, 0.8}, {1, -0.3, 0.8}, {-1, -0.3, -0.8}, {1, -0.3, -0.8}} {
		w := jolt.NewWheelSettingsWV()
		w.SetPosition(p)
		w.SetRadius(0.4)
		w.SetWidth(0.3)
		w.SetSuspensionMinLength(0.1)
		w.SetSuspensionMaxLength(0.5)
		w.SetMaxSteerAngle(0.5)
		w.SetMaxBrakeTorque(3000)
		wheels = append(wheels, w)
	}

	controller := jolt.NewWheeledVehicleControllerSettings()
	controller.AddDifferential(0, 1)
	controller.AddDifferential(2, 3)

	settings := jolt.NewVehicleConstraintSettings()
	vehicle := jolt.NewVehicleConstraint(system.BodyPtr(chassisID), &settings, wheels, controller)

	tester := jolt.NewVehicleCollisionTesterRay(layerMoving, jolt.Vec3{0, 1, 0}, 0.5)
	vehicle.SetCollisionTester(tester)

	system.AddConstraint(vehicle.AsConstraint())
	system.AddStepListener(vehicle.StepListener())

	ctrl := vehicle.Controller()
	return func() {
		fwd, right, brake := float32(0), float32(0), float32(0)
		if rl.IsKeyDown(rl.KeyUp) {
			fwd = 1
		}
		if rl.IsKeyDown(rl.KeyDown) {
			fwd = -1
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			right = 0.5
		}
		if rl.IsKeyDown(rl.KeyRight) {
			right = -0.5
		}
		if rl.IsKeyDown(rl.KeyB) {
			brake = 1
		}
		ctrl.SetDriverInput(fwd, right, brake, 0)
	}
}

func buildScene(system *jolt.PhysicsSystem, index int) func() {
	bi := system.BodyInterface()
	switch index % len(sceneNames) {
	case 0:
		buildPyramid(bi)
	case 1:
		buildSphereDrop(bi)
	case 2:
		buildSoftCloth(bi)
	default:
		addFloor(bi)
		drive := buildVehicle(system)
		system.OptimizeBroadPhase()
		return drive
	}
	system.OptimizeBroadPhase()
	return nil
}

// --- main -------------------------------------------------------------------

func main() {
	if err := jolt.Init(); err != nil {
		panic(err)
	}
	defer jolt.Shutdown()

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(1280, 720, "gojoltc — joltviewer")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	geom := &geometry{}
	renderer := geom.newRenderer()
	defer renderer.Destroy()

	scene := 0
	w := newWorld()
	drive := buildScene(w.system, scene)
	defer func() { w.destroy() }()

	cam := newOrbitCamera(rl.Vector3{X: 0, Y: 1, Z: 0}, 0.8, 0.55, 22)

	wireframe := true
	paused := false

	for !rl.WindowShouldClose() {
		stepOnce := false

		if rl.IsKeyPressed(rl.KeyN) {
			w.destroy()
			scene++
			w = newWorld()
			drive = buildScene(w.system, scene)
			paused = false
		}
		if rl.IsKeyPressed(rl.KeyR) {
			w.destroy()
			w = newWorld()
			drive = buildScene(w.system, scene)
			paused = false
		}
		if rl.IsKeyPressed(rl.KeySpace) {
			paused = !paused
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			stepOnce = true
		}
		if rl.IsKeyPressed(rl.KeyF) {
			wireframe = !wireframe
		}
		if rl.IsKeyPressed(rl.KeyP) {
			rl.TakeScreenshot("joltviewer.png")
		}

		if drive != nil {
			drive()
		}
		if !paused || stepOnce {
			w.system.Update(rl.GetFrameTime(), 1, w.jobSystem)
		}

		geom.reset()
		// Soft bodies are only drawn through the soft-body flags, so enable them
		// or the cloth scene shows nothing.
		w.system.DrawBodies(jolt.DrawSettings{
			DrawShape:                   true,
			DrawShapeWireframe:          wireframe,
			DrawSoftBodyVertices:        true,
			DrawSoftBodyEdgeConstraints: true,
			DrawSoftBodyBendConstraints: true,
		}, renderer)

		camera := cam.camera()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode3D(camera)

		if !wireframe {
			for i := 0; i+2 < len(geom.tris); i += 3 {
				rl.DrawTriangle3D(geom.tris[i], geom.tris[i+1], geom.tris[i+2], geom.triColor[i])
			}
		}
		for i := 0; i+1 < len(geom.lines); i += 2 {
			rl.DrawLine3D(geom.lines[i], geom.lines[i+1], geom.lineColor[i])
		}
		rl.DrawGrid(50, 1)
		rl.EndMode3D()

		drawHUD(hudRows(sceneNames[scene%len(sceneNames)], int(w.system.NumBodies()), paused, drive != nil))
		rl.EndDrawing()
	}
}
