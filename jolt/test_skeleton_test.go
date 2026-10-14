package jolt

import "testing"

// Ported from third_party/joltc/tests/test_skeleton.cpp.

func initSkeleton(t *testing.T) {
	t.Helper()
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	t.Cleanup(Shutdown)
}

func TestSkeletonCreate(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()
	if n := skeleton.JointCount(); n != 0 {
		t.Errorf("joint count = %d, want 0", n)
	}
}

func TestSkeletonAddJoint(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()

	if i := skeleton.AddJoint("root"); i != 0 {
		t.Errorf("root index = %d, want 0", i)
	}
	if n := skeleton.JointCount(); n != 1 {
		t.Errorf("joint count = %d, want 1", n)
	}
	if i := skeleton.AddJointWithParentIndex("child", 0); i != 1 {
		t.Errorf("child index = %d, want 1", i)
	}
	if n := skeleton.JointCount(); n != 2 {
		t.Errorf("joint count = %d, want 2", n)
	}
}

func TestSkeletonAddJointByName(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()

	skeleton.AddJoint("root")
	if i := skeleton.AddJointWithParentName("child", "root"); i != 1 {
		t.Errorf("child index = %d, want 1", i)
	}
}

func TestSkeletonGetJoint(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()

	skeleton.AddJoint("root")
	skeleton.AddJointWithParentIndex("spine", 0)

	name, parent := skeleton.Joint(0)
	if name != "root" || parent != -1 {
		t.Errorf("joint 0 = (%q, %d), want (root, -1)", name, parent)
	}
	name, parent = skeleton.Joint(1)
	if name != "spine" || parent != 0 {
		t.Errorf("joint 1 = (%q, %d), want (spine, 0)", name, parent)
	}
}

func TestSkeletonGetJointIndex(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()

	skeleton.AddJoint("root")
	skeleton.AddJointWithParentIndex("spine", 0)
	skeleton.AddJointWithParentIndex("head", 1)

	for name, want := range map[string]int{"root": 0, "spine": 1, "head": 2, "nonexistent": -1} {
		if got := skeleton.JointIndex(name); got != want {
			t.Errorf("index(%q) = %d, want %d", name, got, want)
		}
	}
}

func TestSkeletonJointsCorrectlyOrdered(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()

	skeleton.AddJoint("root")
	skeleton.AddJointWithParentIndex("child", 0)
	if !skeleton.AreJointsCorrectlyOrdered() {
		t.Error("joints should be correctly ordered")
	}
}

func TestSkeletonPoseCreate(t *testing.T) {
	initSkeleton(t)
	pose := NewSkeletonPose()
	defer pose.Destroy()
	if pose == nil {
		t.Fatal("pose is nil")
	}
}

func TestSkeletonPoseSetSkeleton(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()
	skeleton.AddJoint("root")
	skeleton.AddJointWithParentIndex("child", 0)

	pose := NewSkeletonPose()
	defer pose.Destroy()
	pose.SetSkeleton(skeleton)

	if !pose.Skeleton().Equal(skeleton) {
		t.Error("skeleton mismatch")
	}
	if n := pose.JointCount(); n != 2 {
		t.Errorf("joint count = %d, want 2", n)
	}
}

func TestSkeletonPoseRootOffset(t *testing.T) {
	initSkeleton(t)
	pose := NewSkeletonPose()
	defer pose.Destroy()

	pose.SetRootOffset(RVec3{1, 2, 3})
	equalVec3(t, "offset", pose.RootOffset().ToVec3(), Vec3{1, 2, 3})
}

func TestSkeletonPoseJointState(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()
	skeleton.AddJoint("root")

	pose := NewSkeletonPose()
	defer pose.Destroy()
	pose.SetSkeleton(skeleton)

	pose.SetJointState(0, Vec3{1, 2, 3}, QuatIdentity())
	tr, rot := pose.JointState(0)
	equalVec3(t, "translation", tr, Vec3{1, 2, 3})
	equalFloat(t, "rotation.w", rot.W, 1)
}

func TestSkeletonRagdollSettingsCreate(t *testing.T) {
	initSkeleton(t)
	settings := NewRagdollSettings()
	defer settings.Destroy()
	if settings == nil {
		t.Fatal("settings is nil")
	}
}

func TestSkeletonRagdollSettingsSetSkeleton(t *testing.T) {
	initSkeleton(t)
	skeleton := NewSkeleton()
	defer skeleton.Destroy()
	skeleton.AddJoint("root")

	settings := NewRagdollSettings()
	defer settings.Destroy()
	settings.SetSkeleton(skeleton)

	if !settings.Skeleton().Equal(skeleton) {
		t.Error("skeleton mismatch")
	}
}

func TestSkeletonRagdollSettingsParts(t *testing.T) {
	initSkeleton(t)
	settings := NewRagdollSettings()
	defer settings.Destroy()

	settings.ResizeParts(3)
	if n := settings.PartCount(); n != 3 {
		t.Errorf("part count = %d, want 3", n)
	}

	shape := NewSphereShape(0.5)
	defer shape.Destroy()

	settings.SetPartShape(0, shape)
	settings.SetPartPosition(0, RVec3{0, 1, 0})
	settings.SetPartRotation(0, QuatIdentity())
	settings.SetPartMotionType(0, MotionTypeDynamic)
	settings.SetPartObjectLayer(0, layerMoving)
}

func buildSimpleRagdoll(t *testing.T, f *physicsFixture, position RVec3) (*Skeleton, *RagdollSettings, *Shape, *Ragdoll) {
	t.Helper()
	skeleton := NewSkeleton()
	skeleton.AddJoint("root")

	settings := NewRagdollSettings()
	settings.SetSkeleton(skeleton)
	settings.ResizeParts(1)

	shape := NewSphereShape(0.5)
	settings.SetPartShape(0, shape)
	settings.SetPartPosition(0, position)
	settings.SetPartRotation(0, QuatIdentity())
	settings.SetPartMotionType(0, MotionTypeDynamic)
	settings.SetPartObjectLayer(0, layerMoving)

	ragdoll := settings.CreateRagdoll(f.system, 0, 0)
	return skeleton, settings, shape, ragdoll
}

func TestSkeletonRagdollCreateSimple(t *testing.T) {
	f := newPhysicsFixture(t)
	skeleton, settings, shape, ragdoll := buildSimpleRagdoll(t, f, RVec3{0, 1, 0})
	defer skeleton.Destroy()
	defer settings.Destroy()
	defer shape.Destroy()
	defer ragdoll.Destroy()

	if ragdoll == nil {
		t.Fatal("ragdoll is nil")
	}
	if n := ragdoll.BodyCount(); n != 1 {
		t.Errorf("body count = %d, want 1", n)
	}
	if id := ragdoll.BodyID(0); id == invalidBodyID {
		t.Errorf("body id is invalid (%#x)", uint32(id))
	}
}

func TestSkeletonRagdollAddRemoveFromPhysicsSystem(t *testing.T) {
	f := newPhysicsFixture(t)
	skeleton, settings, shape, ragdoll := buildSimpleRagdoll(t, f, RVec3{0, 1, 0})
	defer skeleton.Destroy()
	defer settings.Destroy()
	defer shape.Destroy()
	defer ragdoll.Destroy()

	ragdoll.AddToPhysicsSystem(ActivationActivate, true)
	if !ragdoll.IsActive(true) {
		t.Error("ragdoll should be active")
	}
	ragdoll.RemoveFromPhysicsSystem(true)
}

func TestSkeletonRagdollGetRootTransform(t *testing.T) {
	f := newPhysicsFixture(t)
	skeleton, settings, shape, ragdoll := buildSimpleRagdoll(t, f, RVec3{1, 2, 3})
	defer skeleton.Destroy()
	defer settings.Destroy()
	defer shape.Destroy()
	defer ragdoll.Destroy()

	ragdoll.AddToPhysicsSystem(ActivationActivate, true)

	pos, _ := ragdoll.RootTransform(true)
	equalReal(t, "root.x", pos.X, 1)
	equalReal(t, "root.y", pos.Y, 2)
	equalReal(t, "root.z", pos.Z, 3)

	ragdoll.RemoveFromPhysicsSystem(true)
}

func TestSkeletonMapperCreate(t *testing.T) {
	initSkeleton(t)
	mapper := NewSkeletonMapper()
	defer mapper.Destroy()
	if mapper == nil {
		t.Fatal("mapper is nil")
	}
}
