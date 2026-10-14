package jolt

import "testing"

func TestSoftBody(t *testing.T) {
	f := newPhysicsFixture(t)
	bi := f.system.BodyInterface()

	shared := NewSoftBodySharedSettings()
	defer shared.Destroy()

	// A 2x2 quad (4 vertices, 2 triangles) as a cloth-like sheet.
	vertices := []SoftVertex{
		{Position: Vec3{-1, 0, -1}, InvMass: 1},
		{Position: Vec3{1, 0, -1}, InvMass: 1},
		{Position: Vec3{1, 0, 1}, InvMass: 1},
		{Position: Vec3{-1, 0, 1}, InvMass: 1},
	}
	faces := []SoftFace{
		{Vertex1: 0, Vertex2: 1, Vertex3: 2},
		{Vertex1: 0, Vertex2: 2, Vertex3: 3},
	}
	shared.AddVertices(vertices)
	shared.AddFaces(faces)
	shared.CreateConstraints(1.0e-4, SoftBodyBendTypeDistance)
	shared.Optimize()

	if n := shared.VertexCount(); n != 4 {
		t.Fatalf("vertex count = %d, want 4", n)
	}
	if n := shared.FaceCount(); n != 2 {
		t.Fatalf("face count = %d, want 2", n)
	}

	settings := NewSoftBodyCreationSettingsFor(shared, RVec3{0, 5, 0}, QuatIdentity(), layerMoving)
	defer settings.Destroy()
	settings.SetFriction(0.5)
	settings.SetRestitution(0.1)
	settings.SetVertexRadius(0.05)
	settings.SetNumIterations(5)

	if settings.SharedSettings() == nil {
		t.Fatal("shared settings not set")
	}

	id := bi.CreateAndAddSoftBody(settings, ActivationActivate)
	if id == invalidBodyID {
		t.Fatal("invalid soft body id")
	}
	t.Cleanup(func() {
		if bi.IsAdded(id) {
			bi.RemoveAndDestroyBody(id)
		}
	})

	body := f.system.BodyPtr(id)
	if body == nil {
		t.Fatal("body pointer is nil")
	}
	if n := body.SoftBodyVertexCount(); n != 4 {
		t.Fatalf("body vertex count = %d, want 4", n)
	}

	positions := make([]Vec3, 4)
	if n := body.SoftBodyVertexPositions(positions); n != 4 {
		t.Fatalf("read %d vertex positions, want 4", n)
	}

	f.system.SetGravity(Vec3{0, -10, 0})
	start := body.SoftBodyVertexPosition(0)
	for range 60 {
		f.system.Update(1.0/60.0, 1, f.jobSystem)
	}
	after := body.SoftBodyVertexPosition(0)
	if after.Y >= start.Y {
		t.Errorf("soft body did not fall: y %v -> %v", start.Y, after.Y)
	}
}
