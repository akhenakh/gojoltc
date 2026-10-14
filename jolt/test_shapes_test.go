package jolt

import (
	"math"
	"testing"
)

// Ported from third_party/joltc/tests/test_shapes.cpp.

func initShapes(t *testing.T) {
	t.Helper()
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	t.Cleanup(Shutdown)
}

func TestShapesSphereCreate(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	equalFloat(t, "radius", shape.SphereRadius(), 1)
}

func TestShapesSphereCreateWithSettings(t *testing.T) {
	initShapes(t)
	settings := NewSphereShapeSettings(2)
	defer settings.Destroy()
	equalFloat(t, "settings radius", settings.Radius(), 2)

	shape := settings.CreateShape()
	defer shape.Destroy()
	equalFloat(t, "shape radius", shape.SphereRadius(), 2)
}

func TestShapesSphereSettingsSetRadius(t *testing.T) {
	initShapes(t)
	settings := NewSphereShapeSettings(1)
	defer settings.Destroy()
	settings.SetRadius(3)
	equalFloat(t, "radius", settings.Radius(), 3)
}

func TestShapesSphereGetType(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	if shape.Type() != ShapeTypeConvex {
		t.Errorf("type = %d, want convex", shape.Type())
	}
	if shape.SubType() != ShapeSubTypeSphere {
		t.Errorf("subtype = %d, want sphere", shape.SubType())
	}
}

func TestShapesSphereGetVolume(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	expected := float32(4.0/3.0) * Pi
	nearFloat(t, "volume", shape.Volume(), expected)
}

func TestShapesBoxCreate(t *testing.T) {
	initShapes(t)
	shape := NewBoxShape(Vec3{1, 2, 3}, DefaultConvexRadius)
	defer shape.Destroy()
	equalVec3(t, "half extent", shape.BoxHalfExtent(), Vec3{1, 2, 3})
}

func TestShapesBoxCreateWithSettings(t *testing.T) {
	initShapes(t)
	settings := NewBoxShapeSettings(Vec3{2, 2, 2}, DefaultConvexRadius)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	equalVec3(t, "half extent", shape.BoxHalfExtent(), Vec3{2, 2, 2})
}

func TestShapesBoxGetConvexRadius(t *testing.T) {
	initShapes(t)
	shape := NewBoxShape(Vec3{1, 1, 1}, 0.1)
	defer shape.Destroy()
	equalFloat(t, "convex radius", shape.BoxConvexRadius(), 0.1)
}

func TestShapesBoxGetType(t *testing.T) {
	initShapes(t)
	shape := NewBoxShape(Vec3{1, 1, 1}, DefaultConvexRadius)
	defer shape.Destroy()
	if shape.Type() != ShapeTypeConvex {
		t.Errorf("type = %d, want convex", shape.Type())
	}
	if shape.SubType() != ShapeSubTypeBox {
		t.Errorf("subtype = %d, want box", shape.SubType())
	}
}

func TestShapesCapsuleCreate(t *testing.T) {
	initShapes(t)
	shape := NewCapsuleShape(1, 0.5)
	defer shape.Destroy()
	equalFloat(t, "half height", shape.CapsuleHalfHeightOfCylinder(), 1)
	equalFloat(t, "radius", shape.CapsuleRadius(), 0.5)
}

func TestShapesCapsuleCreateWithSettings(t *testing.T) {
	initShapes(t)
	settings := NewCapsuleShapeSettings(2, 1)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	equalFloat(t, "half height", shape.CapsuleHalfHeightOfCylinder(), 2)
	equalFloat(t, "radius", shape.CapsuleRadius(), 1)
}

func TestShapesCapsuleGetType(t *testing.T) {
	initShapes(t)
	shape := NewCapsuleShape(1, 0.5)
	defer shape.Destroy()
	if shape.SubType() != ShapeSubTypeCapsule {
		t.Errorf("subtype = %d, want capsule", shape.SubType())
	}
}

func TestShapesCylinderCreate(t *testing.T) {
	initShapes(t)
	shape := NewCylinderShape(1, 0.5)
	defer shape.Destroy()
	equalFloat(t, "half height", shape.CylinderHalfHeight(), 1)
	equalFloat(t, "radius", shape.CylinderRadius(), 0.5)
}

func TestShapesCylinderCreateWithSettings(t *testing.T) {
	initShapes(t)
	settings := NewCylinderShapeSettings(2, 1, DefaultConvexRadius)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	equalFloat(t, "half height", shape.CylinderHalfHeight(), 2)
	equalFloat(t, "radius", shape.CylinderRadius(), 1)
}

func TestShapesCylinderGetType(t *testing.T) {
	initShapes(t)
	shape := NewCylinderShape(1, 0.5)
	defer shape.Destroy()
	if shape.SubType() != ShapeSubTypeCylinder {
		t.Errorf("subtype = %d, want cylinder", shape.SubType())
	}
}

func TestShapesTriangleCreate(t *testing.T) {
	initShapes(t)
	shape := NewTriangleShape(Vec3{0, 0, 0}, Vec3{1, 0, 0}, Vec3{0, 1, 0}, 0)
	defer shape.Destroy()
	equalVec3(t, "v1", shape.TriangleVertex1(), Vec3{0, 0, 0})
	equalVec3(t, "v2", shape.TriangleVertex2(), Vec3{1, 0, 0})
	equalVec3(t, "v3", shape.TriangleVertex3(), Vec3{0, 1, 0})
}

func TestShapesTriangleGetType(t *testing.T) {
	initShapes(t)
	shape := NewTriangleShape(Vec3{0, 0, 0}, Vec3{1, 0, 0}, Vec3{0, 1, 0}, 0)
	defer shape.Destroy()
	if shape.SubType() != ShapeSubTypeTriangle {
		t.Errorf("subtype = %d, want triangle", shape.SubType())
	}
}

func TestShapesTaperedCapsuleCreate(t *testing.T) {
	initShapes(t)
	settings := NewTaperedCapsuleShapeSettings(1, 0.5, 0.3)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	equalFloat(t, "top radius", shape.TaperedCapsuleTopRadius(), 0.5)
	equalFloat(t, "bottom radius", shape.TaperedCapsuleBottomRadius(), 0.3)
	equalFloat(t, "half height", shape.TaperedCapsuleHalfHeight(), 1)
}

func TestShapesTaperedCapsuleGetType(t *testing.T) {
	initShapes(t)
	settings := NewTaperedCapsuleShapeSettings(1, 0.5, 0.3)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	if shape.SubType() != ShapeSubTypeTaperedCapsule {
		t.Errorf("subtype = %d, want tapered capsule", shape.SubType())
	}
}

func TestShapesConvexHullCreate(t *testing.T) {
	initShapes(t)
	points := []Vec3{{0, 0, 0}, {1, 0, 0}, {0.5, 1, 0}, {0.5, 0.5, 1}}
	settings := NewConvexHullShapeSettings(points, DefaultConvexRadius)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	if n := shape.NumPoints(); n != 4 {
		t.Errorf("num points = %d, want 4", n)
	}
	if n := shape.NumFaces(); n != 4 {
		t.Errorf("num faces = %d, want 4", n)
	}
}

func TestShapesConvexHullGetPoint(t *testing.T) {
	initShapes(t)
	points := []Vec3{{0, 0, 0}, {1, 0, 0}, {0.5, 1, 0}, {0.5, 0.5, 1}}
	settings := NewConvexHullShapeSettings(points, DefaultConvexRadius)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()

	p := shape.Point(0)
	if math.IsNaN(float64(p.X)) || math.IsNaN(float64(p.Y)) || math.IsNaN(float64(p.Z)) {
		t.Errorf("point is NaN: %v", p)
	}
}

func TestShapesPlaneCreate(t *testing.T) {
	initShapes(t)
	shape := NewPlaneShape(Plane{Normal: Vec3{0, 1, 0}}, 1000)
	defer shape.Destroy()

	plane := shape.Plane()
	equalFloat(t, "normal.x", plane.Normal.X, 0)
	equalFloat(t, "normal.y", plane.Normal.Y, 1)
	equalFloat(t, "normal.z", plane.Normal.Z, 0)
	equalFloat(t, "distance", plane.Distance, 0)
	equalFloat(t, "half extent", shape.PlaneHalfExtent(), 1000)
}

func TestShapesShapeUserData(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	shape.SetUserData(12345)
	if got := shape.UserData(); got != 12345 {
		t.Errorf("user data = %d, want 12345", got)
	}
}

func TestShapesShapeGetLocalBounds(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	b := shape.LocalBounds()
	equalVec3(t, "min", b.Min, Vec3{-1, -1, -1})
	equalVec3(t, "max", b.Max, Vec3{1, 1, 1})
}

func TestShapesShapeGetCenterOfMass(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	equalVec3(t, "com", shape.CenterOfMass(), Vec3{0, 0, 0})
}

func TestShapesShapeGetInnerRadius(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	equalFloat(t, "inner radius", shape.InnerRadius(), 1)
}

func TestShapesShapeGetMassProperties(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	mp := shape.MassProperties()
	if mp.Mass <= 0 {
		t.Errorf("mass = %v, want > 0", mp.Mass)
	}
}

func TestShapesShapeMustBeStatic(t *testing.T) {
	initShapes(t)
	sphere := NewSphereShape(1)
	defer sphere.Destroy()
	if sphere.MustBeStatic() {
		t.Error("sphere must not be static")
	}

	planeShape := NewPlaneShape(Plane{Normal: Vec3{0, 1, 0}}, 1000)
	defer planeShape.Destroy()
	if !planeShape.MustBeStatic() {
		t.Error("plane must be static")
	}
}

func TestShapesConvexShapeDensity(t *testing.T) {
	initShapes(t)
	settings := NewSphereShapeSettings(1)
	defer settings.Destroy()

	if d := settings.Density(); d <= 0 {
		t.Errorf("default density = %v, want > 0", d)
	}
	settings.SetDensity(2000)
	equalFloat(t, "density", settings.Density(), 2000)
}

func TestShapesConvexShapeSetDensity(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()

	if d := shape.Density(); d <= 0 {
		t.Errorf("default density = %v, want > 0", d)
	}
	shape.SetDensity(5000)
	equalFloat(t, "density", shape.Density(), 5000)
}

func TestShapesShapeSettingsUserData(t *testing.T) {
	initShapes(t)
	settings := NewSphereShapeSettings(1)
	defer settings.Destroy()
	settings.SetUserData(99999)
	if got := settings.UserData(); got != 99999 {
		t.Errorf("user data = %d, want 99999", got)
	}
}

func TestShapesRotatedTranslatedShapeCreate(t *testing.T) {
	initShapes(t)
	inner := NewSphereShape(1)

	shape := NewRotatedTranslatedShape(Vec3{1, 2, 3}, QuatIdentity(), inner)
	defer shape.Destroy()

	equalVec3(t, "position", shape.RotatedTranslatedPosition(), Vec3{1, 2, 3})
	equalFloat(t, "rotation.w", shape.RotatedTranslatedRotation().W, 1)
}

func TestShapesScaledShapeCreate(t *testing.T) {
	initShapes(t)
	inner := NewSphereShape(1)
	shape := NewScaledShape(inner, Vec3{2, 2, 2})
	defer shape.Destroy()
	equalVec3(t, "scale", shape.ScaledScale(), Vec3{2, 2, 2})
}

func TestShapesOffsetCenterOfMassShapeCreate(t *testing.T) {
	initShapes(t)
	inner := NewSphereShape(1)
	shape := NewOffsetCenterOfMassShape(Vec3{0, 1, 0}, inner)
	defer shape.Destroy()
	equalVec3(t, "offset", shape.OffsetCenterOfMassOffset(), Vec3{0, 1, 0})
}

func TestShapesStaticCompoundShapeCreate(t *testing.T) {
	initShapes(t)
	settings := NewStaticCompoundShapeSettings()
	defer settings.Destroy()

	sphere1 := NewSphereShape(1)
	defer sphere1.Destroy()
	sphere2 := NewSphereShape(0.5)
	defer sphere2.Destroy()
	settings.AddShape(Vec3{-1, 0, 0}, QuatIdentity(), sphere1, 0)
	settings.AddShape(Vec3{1, 0, 0}, QuatIdentity(), sphere2, 0)

	shape := settings.CreateShape()
	defer shape.Destroy()

	if n := shape.NumSubShapes(); n != 2 {
		t.Errorf("num sub shapes = %d, want 2", n)
	}
}

func TestShapesMutableCompoundShapeCreate(t *testing.T) {
	initShapes(t)
	settings := NewMutableCompoundShapeSettings()
	defer settings.Destroy()

	shape := settings.CreateShape()
	defer shape.Destroy()

	sphere := NewSphereShape(1)
	defer sphere.Destroy()
	index := shape.AddShape(Vec3{0, 0, 0}, QuatIdentity(), sphere, 0, ^uint32(0))
	if index != 0 {
		t.Errorf("index = %d, want 0", index)
	}
	if n := shape.NumSubShapes(); n != 1 {
		t.Errorf("num sub shapes = %d, want 1", n)
	}

	shape.RemoveShape(0)
	if n := shape.NumSubShapes(); n != 0 {
		t.Errorf("num sub shapes = %d, want 0", n)
	}
}

func TestShapesEmptyShapeCreate(t *testing.T) {
	initShapes(t)
	settings := NewEmptyShapeSettings(Vec3{0, 0, 0})
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()
	equalFloat(t, "volume", shape.Volume(), 0)
}

func TestShapesShapeCastRay(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()

	hit, ok := shape.CastRay(Vec3{-5, 0, 0}, Vec3{10, 0, 0})
	if !ok {
		t.Fatal("expected hit")
	}
	nearFloat(t, "fraction", hit.Fraction, 0.4)
}

func TestShapesShapeCastRayMiss(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()

	if _, ok := shape.CastRay(Vec3{0, 5, 0}, Vec3{1, 0, 0}); ok {
		t.Error("expected miss")
	}
}

func TestShapesShapeCollidePointInside(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	if !shape.CollidePoint(Vec3{0, 0, 0}) {
		t.Error("point should be inside")
	}
}

func TestShapesShapeCollidePointOutside(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	if shape.CollidePoint(Vec3{5, 0, 0}) {
		t.Error("point should be outside")
	}
}

func TestShapesTaperedCylinderShapeCreate(t *testing.T) {
	initShapes(t)
	settings := NewTaperedCylinderShapeSettings(1, 0.5, 1, DefaultConvexRadius)
	defer settings.Destroy()
	shape := settings.CreateShape()
	defer shape.Destroy()

	equalFloat(t, "top radius", shape.TaperedCylinderTopRadius(), 0.5)
	equalFloat(t, "bottom radius", shape.TaperedCylinderBottomRadius(), 1)
	equalFloat(t, "half height", shape.TaperedCylinderHalfHeight(), 1)
}

func TestShapesShapeIsValidScale(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	if !shape.IsValidScale(Vec3{2, 2, 2}) {
		t.Error("uniform scale should be valid")
	}
	if shape.IsValidScale(Vec3{2, 1, 1}) {
		t.Error("non-uniform scale should be invalid for a sphere")
	}
}

func TestShapesShapeMakeScaleValid(t *testing.T) {
	initShapes(t)
	shape := NewSphereShape(1)
	defer shape.Destroy()
	result := shape.MakeScaleValid(Vec3{2, 1, 1})
	if !shape.IsValidScale(result) {
		t.Errorf("result scale %v should be valid", result)
	}
}
