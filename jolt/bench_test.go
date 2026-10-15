package jolt

import "testing"

// Benchmarks use b.Loop (Go 1.24+) and report allocations so regressions in the
// cgo hot paths are visible.

func newBenchSystem(b *testing.B, numBodies int) (*PhysicsSystem, *JobSystem, []BodyID) {
	b.Helper()

	if err := Init(); err != nil {
		b.Fatalf("Init: %v", err)
	}

	jobSystem := NewJobSystemThreadPool(nil)

	bpLayer := NewBroadPhaseLayerInterfaceTable(2, 2)
	bpLayer.MapObjectToBroadPhaseLayer(layerNonMoving, 0)
	bpLayer.MapObjectToBroadPhaseLayer(layerMoving, 1)

	pairFilter := NewObjectLayerPairFilterTable(2)
	pairFilter.EnableCollision(layerNonMoving, layerMoving)
	pairFilter.EnableCollision(layerMoving, layerMoving)

	objBpFilter := NewObjectVsBroadPhaseLayerFilterTable(bpLayer, 2, pairFilter, 2)

	system := NewPhysicsSystem(&PhysicsSystemSettings{
		MaxBodies:                     uint32(numBodies) + 1,
		MaxBodyPairs:                  65536,
		MaxContactConstraints:         65536,
		BroadPhaseLayerInterface:      bpLayer,
		ObjectLayerPairFilter:         pairFilter,
		ObjectVsBroadPhaseLayerFilter: objBpFilter,
	})

	bi := system.BodyInterface()
	shape := NewSphereShape(0.5)
	ids := make([]BodyID, numBodies)
	for i := range numBodies {
		settings := NewBodyCreationSettings(shape, RVec3{Real(i % 64), 10, Real(i / 64)}, QuatIdentity(), MotionTypeDynamic, layerMoving)
		ids[i] = bi.CreateAndAddBody(settings, ActivationActivate)
		settings.Destroy()
	}
	shape.Destroy()
	system.OptimizeBroadPhase()

	b.Cleanup(func() {
		for _, id := range ids {
			bi.RemoveAndDestroyBody(id)
		}
		system.Destroy()
		jobSystem.Destroy()
		Shutdown()
	})
	return system, jobSystem, ids
}

func BenchmarkPhysicsUpdate(b *testing.B) {
	system, jobSystem, _ := newBenchSystem(b, 512)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		system.Update(1.0/60.0, 1, jobSystem)
	}
}

// BenchmarkBodyPositions exercises the per-call path for reading every body
// position (value returns, no allocation).
func BenchmarkBodyPositions(b *testing.B) {
	system, _, ids := newBenchSystem(b, 1024)
	bi := system.BodyInterface()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		for _, id := range ids {
			_ = bi.Position(id)
		}
	}
}

// BenchmarkBodyPositionsInto is the allocation-free sweep: results go into a
// caller-owned slice reused across iterations.
func BenchmarkBodyPositionsInto(b *testing.B) {
	system, _, ids := newBenchSystem(b, 1024)
	bi := system.BodyInterface()
	dst := make([]RVec3, len(ids))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		bi.PositionsInto(ids, dst)
	}
}

// BenchmarkBodyPosesInto reads positions and rotations in one sweep.
func BenchmarkBodyPosesInto(b *testing.B) {
	system, _, ids := newBenchSystem(b, 1024)
	bi := system.BodyInterface()
	positions := make([]RVec3, len(ids))
	rotations := make([]Quat, len(ids))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		bi.PositionsAndRotationsInto(ids, positions, rotations)
	}
}

func BenchmarkEnumerateBodiesAlloc(b *testing.B) {
	system, _, _ := newBenchSystem(b, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = system.Bodies()
	}
}

func BenchmarkEnumerateBodiesInto(b *testing.B) {
	system, _, _ := newBenchSystem(b, 1024)
	buf := make([]BodyID, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = system.BodyIDsInto(buf)
	}
}

// BenchmarkRaySweep compares sweeping a sensor bank of rays (a LIDAR pattern)
// three ways: per-ray through the collector, per-ray through Jolt's native cast,
// and the batched native sweep. It reuses the helper's body grid as targets, so
// the body count stays within the system it was sized for.
func BenchmarkRaySweep(b *testing.B) {
	system, _, ids := newBenchSystem(b, 256)
	system.OptimizeBroadPhase()

	query := system.NarrowPhaseQuery()
	bi := system.BodyInterface()

	// Rays straight down at the grid of spheres sitting at y = 10.
	const rays = 256
	origins := make([]RVec3, rays)
	dirs := make([]Vec3, rays)
	out := make([]RayCastResult, rays)
	for i := range rays {
		origins[i] = bi.Position(ids[i])
		origins[i].Y += 2
		dirs[i] = Vec3{0, -20, 0}
	}

	b.Run("collector-per-ray", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			for i := range rays {
				query.CastRayClosest(origins[i], dirs[i], DefaultRayCastSettings(), nil)
			}
		}
	})
	b.Run("native-per-ray", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			for i := range rays {
				query.CastRayFiltered(origins[i], dirs[i], nil)
			}
		}
	})
	b.Run("native-batch", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			query.CastRaysFiltered(origins, dirs, out, nil)
		}
	})
}
