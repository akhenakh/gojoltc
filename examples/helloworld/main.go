// Command helloworld is a port of joltc's samples/01_HelloWorld.
package main

import (
	"fmt"
	"log"

	"github.com/akhenakh/gojoltc/jolt"
)

const (
	layerNonMoving jolt.ObjectLayer = 0
	layerMoving    jolt.ObjectLayer = 1
	numLayers                       = 2
)

func main() {
	if err := jolt.Init(); err != nil {
		log.Fatalf("init: %v", err)
	}
	defer jolt.Shutdown()

	jobSystem := jolt.NewJobSystemThreadPool(nil)
	defer jobSystem.Destroy()

	// The physics system takes ownership of the layer interface and filters.
	bpLayer := jolt.NewBroadPhaseLayerInterfaceTable(numLayers, numLayers)
	bpLayer.MapObjectToBroadPhaseLayer(layerNonMoving, 0)
	bpLayer.MapObjectToBroadPhaseLayer(layerMoving, 1)

	pairFilter := jolt.NewObjectLayerPairFilterTable(numLayers)
	pairFilter.EnableCollision(layerNonMoving, layerMoving)
	pairFilter.EnableCollision(layerMoving, layerMoving)

	objBpFilter := jolt.NewObjectVsBroadPhaseLayerFilterTable(bpLayer, numLayers, pairFilter, numLayers)

	system := jolt.NewPhysicsSystem(&jolt.PhysicsSystemSettings{
		MaxBodies:                     65536,
		MaxBodyPairs:                  65536,
		MaxContactConstraints:         65536,
		BroadPhaseLayerInterface:      bpLayer,
		ObjectLayerPairFilter:         pairFilter,
		ObjectVsBroadPhaseLayerFilter: objBpFilter,
	})
	defer system.Destroy()

	bi := system.BodyInterface()

	// Floor.
	floorShape := jolt.NewBoxShape(jolt.Vec3{100, 1, 100}, jolt.DefaultConvexRadius)
	floorSettings := jolt.NewBodyCreationSettings(floorShape, jolt.Vec3{0, -1, 0}, jolt.QuatIdentity(), jolt.MotionTypeStatic, layerNonMoving)
	bi.CreateAndAddBody(floorSettings, jolt.ActivationDontActivate)
	floorSettings.Destroy()
	floorShape.Destroy()

	// Falling sphere.
	sphereShape := jolt.NewSphereShape(0.5)
	sphereSettings := jolt.NewBodyCreationSettings(sphereShape, jolt.Vec3{0, 2, 0}, jolt.QuatIdentity(), jolt.MotionTypeDynamic, layerMoving)
	sphereID := bi.CreateAndAddBody(sphereSettings, jolt.ActivationActivate)
	sphereSettings.Destroy()
	sphereShape.Destroy()

	bi.SetLinearVelocity(sphereID, jolt.Vec3{0, -5, 0})

	system.OptimizeBroadPhase()

	for range 60 {
		system.Update(1.0/60.0, 1, jobSystem)
	}

	pos := bi.Position(sphereID)
	fmt.Printf("Sphere position after one second: (%.2f, %.2f, %.2f)\n", pos.X, pos.Y, pos.Z)
}
