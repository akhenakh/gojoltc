// Command rover is a small wheeled rover driving over rough terrain, with
// ground-clearance and forward-obstacle sensors.
//
// It is a worked example of the parts a driving robot needs:
//
//   - a procedural height field as the terrain
//   - a wheeled vehicle: chassis, wheels with tire friction curves, axles
//     connected by differentials, and a wheel collision tester
//   - filtered raycasts as rangefinders (a BodyFilter keeps the rover from
//     seeing its own chassis)
//   - a contact listener for surface feedback
//
// Run from this directory: go run .
//
// Controls:
//
//	W / up     drive forward
//	S / down   reverse
//	A / D, left / right   steer
//	Space      brake
//	R          respawn
//	F          wireframe / solid
//	mouse      orbit / zoom
//	Esc        quit
package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/akhenakh/gojoltc/jolt"
)

const (
	layerTerrain jolt.ObjectLayer = 0
	layerRover   jolt.ObjectLayer = 1
	numLayers                     = 2
)

// --- terrain -----------------------------------------------------------------

const (
	terrainSamples = 65    // samples per side (the height field API takes this, not len(samples))
	terrainSize    = 48.0  // metres across
	terrainSeed    = 0.123 // keeps the "noise" deterministic
)

// heightAt is a cheap smooth height field: a few sines at different scales, so
// the rover meets gentle slopes, ridges and small bumps.
func heightAt(x, z float64) float64 {
	h := 0.9 * math.Sin(x*0.32+terrainSeed) * math.Cos(z*0.27)
	h += 0.45 * math.Sin(x*0.85+1.3) * math.Sin(z*0.7-0.4)
	h += 0.18 * math.Sin(x*2.2) * math.Cos(z*1.8+0.7)
	return h
}

func buildTerrain(bi *jolt.BodyInterface) {
	step := float32(terrainSize) / float32(terrainSamples-1)
	samples := make([]float32, terrainSamples*terrainSamples)
	for j := range terrainSamples {
		for i := range terrainSamples {
			x := (float64(i) - terrainSamples/2) * float64(step)
			z := (float64(j) - terrainSamples/2) * float64(step)
			samples[j*terrainSamples+i] = float32(heightAt(x, z))
		}
	}

	shape := jolt.NewHeightFieldShape(
		samples,
		jolt.Vec3{-terrainSize / 2, 0, -terrainSize / 2}, // world offset of sample (0,0)
		jolt.Vec3{step, 1, step},                         // sample spacing and height scale
	)
	if shape == nil {
		panic("height field shape is nil (len(samples) must be a perfect square)")
	}
	defer shape.Destroy()

	settings := jolt.NewBodyCreationSettings(shape, jolt.RVec3{0, 0, 0}, jolt.QuatIdentity(), jolt.MotionTypeStatic, layerTerrain)
	settings.SetFriction(0.9)
	bi.CreateAndAddBody(settings, jolt.ActivationDontActivate)
	settings.Destroy()
}

// scatterRocks drops static rocks the rover has to drive around.
func scatterRocks(bi *jolt.BodyInterface) {
	for i := range 14 {
		a := float64(i) * 2.399963 // golden angle, so they do not line up
		r := 6.0 + float64(i)*1.15
		x, z := math.Cos(a)*r, math.Sin(a)*r
		size := 0.35 + 0.12*float64(i%5)
		y := heightAt(x, z)

		shape := jolt.NewSphereShape(float32(size))
		settings := jolt.NewBodyCreationSettings(shape, jolt.RVec3{jolt.Real(x), jolt.Real(y + size*0.4), jolt.Real(z)}, jolt.QuatIdentity(), jolt.MotionTypeStatic, layerTerrain)
		bi.CreateAndAddBody(settings, jolt.ActivationDontActivate)
		settings.Destroy()
		shape.Destroy()
	}
}

// --- rover -------------------------------------------------------------------

type rover struct {
	chassisID  jolt.BodyID
	vehicle    *jolt.VehicleConstraint
	controller *jolt.WheeledVehicleController

	// Sensors must not see the rover itself.
	sensorFilter *jolt.BodyFilter
}

func buildRover(system *jolt.PhysicsSystem) *rover {
	bi := system.BodyInterface()

	// Chassis: 1.2 m x 0.5 m x 0.9 m, spawned a little above the ground.
	chassisShape := jolt.NewBoxShape(jolt.Vec3{0.6, 0.25, 0.45}, jolt.DefaultConvexRadius)
	chassisSettings := jolt.NewBodyCreationSettings(chassisShape, jolt.RVec3{0, 2.5, 0}, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerRover)
	chassisSettings.SetFriction(0.4)
	chassisSettings.SetAllowSleeping(false) // a rover is never really asleep
	chassisSettings.SetLinearDamping(0.1)
	chassisSettings.SetAngularDamping(0.4)
	chassisID := bi.CreateAndAddBody(chassisSettings, jolt.ActivationActivate)
	chassisSettings.Destroy()
	chassisShape.Destroy()

	// Tire friction as a function of slip: grip peaks at ~10% slip and falls
	// away beyond it, which is what makes a rover slide on sand and loose rock.
	makeTireCurve := func(peak, tail float32) *jolt.LinearCurve {
		c := jolt.NewLinearCurve()
		c.AddPoint(0, 0)
		c.AddPoint(0.1, peak)
		c.AddPoint(1.0, tail)
		c.Sort()
		return c
	}
	longitudinal := makeTireCurve(1.0, 0.55)
	defer longitudinal.Destroy() // the wheel copies it
	lateral := makeTireCurve(0.9, 0.5)
	defer lateral.Destroy()

	wheelPositions := []jolt.Vec3{
		{-0.55, -0.2, 0.5}, {0.55, -0.2, 0.5},
		{-0.55, -0.2, -0.5}, {0.55, -0.2, -0.5},
	}
	wheels := make([]*jolt.WheelSettings, 0, len(wheelPositions))
	for _, p := range wheelPositions {
		w := jolt.NewWheelSettingsWV()
		w.SetPosition(p)
		w.SetRadius(0.28)
		w.SetWidth(0.2)
		w.SetSuspensionMinLength(0.05)
		w.SetSuspensionMaxLength(0.35)
		w.SetMaxSteerAngle(0.55) // ~31 degrees
		w.SetMaxBrakeTorque(1200)
		w.SetInertia(0.4)
		w.SetAngularDamping(0.05)
		w.SetLongitudinalFrictionCurve(longitudinal)
		w.SetLateralFrictionCurve(lateral)
		wheels = append(wheels, w)
	}

	controller := jolt.NewWheeledVehicleControllerSettings()
	engine := jolt.NewVehicleEngineSettings()
	engine.MaxTorque = 260
	engine.MaxRPM = 4000
	controller.SetEngine(engine)
	controller.AddDifferential(0, 1) // front axle: wheels 0 and 1
	controller.AddDifferential(2, 3) // rear axle: wheels 2 and 3

	settings := jolt.NewVehicleConstraintSettings()
	settings.MaxPitchRollAngle = 0.6 // rolls back upright beyond ~35 degrees
	vehicle := jolt.NewVehicleConstraint(system.BodyPtr(chassisID), &settings, wheels, controller)

	// The tester probes the ground under each wheel. It must use the rover's
	// own layer so the layer-pair filter lets it reach the terrain.
	tester := jolt.NewVehicleCollisionTesterRay(layerRover, jolt.Vec3{0, 1, 0}, 0.9)
	vehicle.SetCollisionTester(tester)

	system.AddConstraint(vehicle.AsConstraint())
	system.AddStepListener(vehicle.StepListener())

	return &rover{
		chassisID:  chassisID,
		vehicle:    vehicle,
		controller: vehicle.Controller(),
		sensorFilter: jolt.NewBodyFilter(&jolt.BodyFilterCallbacks{
			ShouldCollide: func(id jolt.BodyID) bool { return id != chassisID },
		}),
	}
}

func (r *rover) destroy() { r.sensorFilter.Destroy() }

// drive applies one frame of driver input: forward/reverse, steer, brake.
func (r *rover) drive(forward, steer, brake float32) {
	r.controller.SetDriverInput(forward, steer, brake, 0)
}

// respawn puts the rover back above the origin, at rest.
func (r *rover) respawn(bi *jolt.BodyInterface) {
	bi.SetPosition(r.chassisID, jolt.RVec3{0, 2.5, 0}, jolt.ActivationActivate)
	bi.SetRotation(r.chassisID, jolt.QuatIdentity(), jolt.ActivationActivate)
	bi.SetLinearVelocity(r.chassisID, jolt.Vec3{})
	bi.SetAngularVelocity(r.chassisID, jolt.Vec3{})
}

// sensorReading is what the rover's rangefinders see this frame.
type sensorReading struct {
	groundHit      bool
	groundDistance float32
	groundPoint    jolt.RVec3

	aheadHit      bool
	aheadDistance float32
	aheadPoint    jolt.RVec3
}

const (
	groundRange = float32(4) // metres the clearance sensor reaches
	aheadRange  = float32(6) // metres the obstacle sensor reaches
)

// scan runs a downward clearance probe and a forward obstacle probe. Both use a
// body filter that excludes the rover, and the unfiltered ray would otherwise
// hit the chassis at fraction 0.
func (r *rover) scan(system *jolt.PhysicsSystem, bi *jolt.BodyInterface) sensorReading {
	query := system.NarrowPhaseQuery()
	pos := bi.Position(r.chassisID)
	rot := bi.Rotation(r.chassisID)
	opts := &jolt.QueryOptions{Body: r.sensorFilter}

	var out sensorReading

	// Clearance: straight down from just below the chassis.
	groundOrigin := jolt.RVec3{pos.X, pos.Y - 0.25, pos.Z}
	if hit, ok := query.CastRayClosest(groundOrigin, jolt.Vec3{0, -groundRange, 0}, jolt.DefaultRayCastSettings(), opts); ok {
		out.groundHit = true
		out.groundDistance = hit.Fraction * groundRange
		out.groundPoint = jolt.RVec3{groundOrigin.X, groundOrigin.Y - jolt.Real(out.groundDistance), groundOrigin.Z}
	}

	// Obstacle: along the chassis' own forward axis (+Z in body space).
	forward := rot.Rotate(jolt.Vec3{0, 0, 1})
	if hit, ok := query.CastRayClosest(pos, forward.MulScalar(aheadRange), jolt.DefaultRayCastSettings(), opts); ok {
		out.aheadHit = true
		out.aheadDistance = hit.Fraction * aheadRange
		out.aheadPoint = jolt.RVec3{
			pos.X + jolt.Real(forward.X*out.aheadDistance),
			pos.Y + jolt.Real(forward.Y*out.aheadDistance),
			pos.Z + jolt.Real(forward.Z*out.aheadDistance),
		}
	}
	return out
}

// --- render helpers ----------------------------------------------------------

func joltColor(c jolt.Color) rl.Color {
	return rl.Color{R: uint8(c & 0xff), G: uint8((c >> 8) & 0xff), B: uint8((c >> 16) & 0xff), A: uint8((c >> 24) & 0xff)}
}

func vec3(v jolt.RVec3) rl.Vector3 {
	return rl.Vector3{X: float32(v.X), Y: float32(v.Y), Z: float32(v.Z)}
}

// geometry accumulates Jolt's debug draw calls for one frame.
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

func (g *geometry) renderer() *jolt.DebugRenderer {
	return jolt.NewDebugRenderer(&jolt.DebugRendererCallbacks{
		// One colour per vertex, so the draw loops can index positions directly.
		DrawLine: func(from, to jolt.RVec3, color jolt.Color) {
			c := joltColor(color)
			g.lines = append(g.lines, vec3(from), vec3(to))
			g.lineColor = append(g.lineColor, c, c)
		},
		DrawTriangle: func(v1, v2, v3 jolt.RVec3, color jolt.Color, castShadow bool) {
			c := joltColor(color)
			g.tris = append(g.tris, vec3(v1), vec3(v2), vec3(v3))
			g.triColor = append(g.triColor, c, c, c)
		},
	})
}

// --- main --------------------------------------------------------------------

func main() {
	if err := jolt.Init(); err != nil {
		panic(err)
	}
	defer jolt.Shutdown()

	jobSystem := jolt.NewJobSystemThreadPool(nil)
	defer jobSystem.Destroy()

	bpLayer := jolt.NewBroadPhaseLayerInterfaceTable(numLayers, numLayers)
	bpLayer.MapObjectToBroadPhaseLayer(layerTerrain, 0)
	bpLayer.MapObjectToBroadPhaseLayer(layerRover, 1)
	pairFilter := jolt.NewObjectLayerPairFilterTable(numLayers)
	pairFilter.EnableCollision(layerTerrain, layerRover)
	pairFilter.EnableCollision(layerRover, layerRover)
	objBpFilter := jolt.NewObjectVsBroadPhaseLayerFilterTable(bpLayer, numLayers, pairFilter, numLayers)

	system := jolt.NewPhysicsSystem(&jolt.PhysicsSystemSettings{
		MaxBodies:                     4096,
		MaxBodyPairs:                  4096,
		MaxContactConstraints:         2048,
		BroadPhaseLayerInterface:      bpLayer,
		ObjectLayerPairFilter:         pairFilter,
		ObjectVsBroadPhaseLayerFilter: objBpFilter,
	})
	defer system.Destroy()
	system.SetGravity(jolt.Vec3{0, -9.81, 0})

	bi := system.BodyInterface()

	// Surface feedback: count contacts (a real rover would use the sub-shape or
	// the material to tell sand from rock and change its tire behaviour).
	var contacts int32
	contactListener := jolt.NewContactListener(&jolt.ContactListenerCallbacks{
		OnAdded: func(a, b *jolt.Body, manifold *jolt.ContactManifold, settings *jolt.ContactSettings) {
			contacts++
		},
	})
	defer contactListener.Destroy()
	system.SetContactListener(contactListener)

	buildTerrain(bi)
	scatterRocks(bi)
	r := buildRover(system)
	defer r.destroy()

	system.OptimizeBroadPhase()

	geom := &geometry{}
	renderer := geom.renderer()
	defer renderer.Destroy()

	// FlagWindowHighdpi: on a scaled display (e.g. 125%) the OS hands back a
	// larger framebuffer than the logical window size. Without this flag raylib
	// keeps the render surface at the logical size, so the 3D viewport covers
	// only part of the window and the rest stays blank. 2D drawing keeps its
	// logical coordinates (raylib applies the display scale), so the HUD is
	// unaffected - just rendered at the display's real resolution.
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowHighdpi)
	rl.InitWindow(1280, 720, "gojoltc — rover on rough terrain")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	// Report the real surface sizes: on a scaled display the render size is larger
	// than the logical one, which is what makes the 3D viewport fill the window.
	fmt.Printf("window: %dx%d logical, %dx%d render, display scale %.2f\n",
		rl.GetScreenWidth(), rl.GetScreenHeight(), rl.GetRenderWidth(), rl.GetRenderHeight(), rl.GetWindowScaleDPI().X)

	cam := newOrbitCamera(rl.Vector3{X: 0, Y: 1, Z: 0}, 2.4, 0.35, 16)

	wireframe := true

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyR) {
			r.respawn(bi)
		}
		if rl.IsKeyPressed(rl.KeyF) {
			wireframe = !wireframe
		}

		forward, steer, brake := float32(0), float32(0), float32(0)
		if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
			forward += 1
		}
		if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
			forward -= 1
		}
		if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
			steer += 1
		}
		if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
			steer -= 1
		}
		if rl.IsKeyDown(rl.KeySpace) {
			brake = 1
		}
		r.drive(forward, steer, brake)

		system.Update(rl.GetFrameTime(), 1, jobSystem)
		camera := cam.camera()

		reading := r.scan(system, bi)
		chassis := r.vehicle.VehicleBody()

		geom.reset()
		system.DrawBodies(jolt.DrawSettings{DrawShape: true, DrawShapeWireframe: wireframe}, renderer)

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(206, 222, 238, 255))
		rl.BeginMode3D(camera)

		if !wireframe {
			for i := 0; i+2 < len(geom.tris); i += 3 {
				rl.DrawTriangle3D(geom.tris[i], geom.tris[i+1], geom.tris[i+2], geom.triColor[i])
			}
		}
		for i := 0; i+1 < len(geom.lines); i += 2 {
			rl.DrawLine3D(geom.lines[i], geom.lines[i+1], geom.lineColor[i])
		}

		// Draw the sensors.
		pos := chassis.Position()
		if reading.groundHit {
			rl.DrawLine3D(vec3(jolt.RVec3{pos.X, pos.Y - 0.25, pos.Z}), vec3(reading.groundPoint), rl.Green)
			rl.DrawSphere(vec3(reading.groundPoint), 0.07, rl.Green)
		}
		if reading.aheadHit {
			rl.DrawLine3D(vec3(pos), vec3(reading.aheadPoint), rl.Red)
			rl.DrawSphere(vec3(reading.aheadPoint), 0.09, rl.Red)
		}

		rl.EndMode3D()

		velocity := chassis.LinearVelocity()
		kmh := float32(math.Hypot(float64(velocity.X), float64(velocity.Z))) * 3.6

		rl.DrawRectangle(0, 0, 470, 132, rl.NewColor(255, 255, 255, 215))
		rl.DrawText("gojoltc rover — rough terrain", 10, 8, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("speed %.1f km/h    contacts %d", kmh, contacts), 10, 32, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("clearance %.2f m%s", reading.groundDistance, mark(reading.groundHit)), 10, 56, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("obstacle ahead %.2f m%s", reading.aheadDistance, mark(reading.aheadHit)), 10, 80, 20, rl.DarkGray)
		rl.DrawText("WASD/arrows drive  Space brake  R respawn  F wireframe", 10, 106, 18, rl.Gray)
		rl.DrawText("drag a mouse button to orbit, wheel to zoom", 10, 128, 18, rl.Gray)

		rl.EndDrawing()
	}
}

func mark(hit bool) string {
	if hit {
		return ""
	}
	return " (none)"
}
