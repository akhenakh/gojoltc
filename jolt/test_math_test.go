package jolt

import (
	"math"
	"testing"
)

// Ported from third_party/joltc/tests/test_math.cpp.

const mathEpsilon = 1e-5

func nearFloat(t *testing.T, name string, got, want float32) {
	t.Helper()
	if d := got - want; d > mathEpsilon || d < -mathEpsilon {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func equalFloat(t *testing.T, name string, got, want float32) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func nearVec3(t *testing.T, name string, got, want Vec3) {
	t.Helper()
	nearFloat(t, name+".X", got.X, want.X)
	nearFloat(t, name+".Y", got.Y, want.Y)
	nearFloat(t, name+".Z", got.Z, want.Z)
}

func equalVec3(t *testing.T, name string, got, want Vec3) {
	t.Helper()
	equalFloat(t, name+".X", got.X, want.X)
	equalFloat(t, name+".Y", got.Y, want.Y)
	equalFloat(t, name+".Z", got.Z, want.Z)
}

func nearQuat(t *testing.T, name string, got, want Quat) {
	t.Helper()
	nearFloat(t, name+".X", got.X, want.X)
	nearFloat(t, name+".Y", got.Y, want.Y)
	nearFloat(t, name+".Z", got.Z, want.Z)
	nearFloat(t, name+".W", got.W, want.W)
}

func equalQuat(t *testing.T, name string, got, want Quat) {
	t.Helper()
	equalFloat(t, name+".X", got.X, want.X)
	equalFloat(t, name+".Y", got.Y, want.Y)
	equalFloat(t, name+".Z", got.Z, want.Z)
	equalFloat(t, name+".W", got.W, want.W)
}

func TestMathVec3(t *testing.T) {
	t.Run("AxisX", func(t *testing.T) { equalVec3(t, "axis", Vec3AxisX(), Vec3{1, 0, 0}) })
	t.Run("AxisY", func(t *testing.T) { equalVec3(t, "axis", Vec3AxisY(), Vec3{0, 1, 0}) })
	t.Run("AxisZ", func(t *testing.T) { equalVec3(t, "axis", Vec3AxisZ(), Vec3{0, 0, 1}) })

	t.Run("Add", func(t *testing.T) {
		got := (Vec3{1, 2, 3}).Add(Vec3{4, 5, 6})
		equalVec3(t, "add", got, Vec3{5, 7, 9})
	})
	t.Run("Subtract", func(t *testing.T) {
		got := (Vec3{4, 5, 6}).Sub(Vec3{1, 2, 3})
		equalVec3(t, "sub", got, Vec3{3, 3, 3})
	})
	t.Run("Multiply", func(t *testing.T) {
		got := (Vec3{1, 2, 3}).Mul(Vec3{2, 3, 4})
		equalVec3(t, "mul", got, Vec3{2, 6, 12})
	})
	t.Run("MultiplyScalar", func(t *testing.T) {
		got := (Vec3{1, 2, 3}).MulScalar(2)
		equalVec3(t, "mulScalar", got, Vec3{2, 4, 6})
	})
	t.Run("Divide", func(t *testing.T) {
		got := (Vec3{4, 6, 8}).Div(Vec3{2, 3, 4})
		equalVec3(t, "div", got, Vec3{2, 2, 2})
	})
	t.Run("DivideScalar", func(t *testing.T) {
		got := (Vec3{2, 4, 6}).DivScalar(2)
		equalVec3(t, "divScalar", got, Vec3{1, 2, 3})
	})
	t.Run("Length", func(t *testing.T) {
		equalFloat(t, "length", (Vec3{3, 4, 0}).Length(), 5)
	})
	t.Run("LengthSquared", func(t *testing.T) {
		equalFloat(t, "lengthSq", (Vec3{3, 4, 0}).LengthSquared(), 25)
	})
	t.Run("Normalized", func(t *testing.T) {
		got := (Vec3{3, 0, 0}).Normalized()
		equalVec3(t, "normalized", got, Vec3{1, 0, 0})
	})
	t.Run("DotProduct", func(t *testing.T) {
		equalFloat(t, "dot", (Vec3{1, 2, 3}).Dot(Vec3{4, 5, 6}), 32)
	})
	t.Run("Cross", func(t *testing.T) {
		got := (Vec3{1, 0, 0}).Cross(Vec3{0, 1, 0})
		equalVec3(t, "cross", got, Vec3{0, 0, 1})
	})
	t.Run("Negate", func(t *testing.T) {
		got := (Vec3{1, -2, 3}).Negate()
		equalVec3(t, "negate", got, Vec3{-1, 2, -3})
	})
	t.Run("Abs", func(t *testing.T) {
		got := (Vec3{-1, -2, 3}).Abs()
		equalVec3(t, "abs", got, Vec3{1, 2, 3})
	})
	t.Run("IsClose", func(t *testing.T) {
		if !(Vec3{1, 2, 3}).IsClose(Vec3{1, 2, 3}, 0.001) {
			t.Error("expected close")
		}
		if (Vec3{1, 2, 3}).IsClose(Vec3{1.1, 2.1, 3.1}, 0.001) {
			t.Error("expected not close")
		}
	})
	t.Run("IsNearZero", func(t *testing.T) {
		if !(Vec3{0, 0, 0}).IsNearZero(0.001) {
			t.Error("expected near zero")
		}
		if !(Vec3{0.0001, 0.0001, 0.0001}).IsNearZero(0.001) {
			t.Error("expected near zero")
		}
		if (Vec3{1, 0, 0}).IsNearZero(0.001) {
			t.Error("expected not near zero")
		}
	})
	t.Run("IsNormalized", func(t *testing.T) {
		if !(Vec3{1, 0, 0}).IsNormalized(0.001) {
			t.Error("expected normalized")
		}
		if (Vec3{2, 0, 0}).IsNormalized(0.001) {
			t.Error("expected not normalized")
		}
	})
	t.Run("IsNaN", func(t *testing.T) {
		if (Vec3{1, 2, 3}).IsNaN() {
			t.Error("expected not NaN")
		}
		if !(Vec3{float32(math.NaN()), 2, 3}).IsNaN() {
			t.Error("expected NaN")
		}
	})
}

func TestMathQuat(t *testing.T) {
	t.Run("FromEulerAngles", func(t *testing.T) {
		got := QuatEulerAngles(Vec3{0, 0, 0})
		nearQuat(t, "quat", got, Quat{0, 0, 0, 1})
	})
	t.Run("GetEulerAngles", func(t *testing.T) {
		got := QuatIdentity().GetEulerAngles()
		nearVec3(t, "euler", got, Vec3{0, 0, 0})
	})
	t.Run("Multiply", func(t *testing.T) {
		got := QuatIdentity().Mul(QuatIdentity())
		nearQuat(t, "quat", got, Quat{0, 0, 0, 1})
	})
	t.Run("Conjugated", func(t *testing.T) {
		got := (Quat{1, 2, 3, 4}).Conjugated()
		equalQuat(t, "quat", got, Quat{-1, -2, -3, 4})
	})
	t.Run("Inversed", func(t *testing.T) {
		got := QuatIdentity().Inversed()
		nearQuat(t, "quat", got, Quat{0, 0, 0, 1})
	})
	t.Run("Dot", func(t *testing.T) {
		equalFloat(t, "dot", (Quat{1, 0, 0, 0}).Dot(Quat{1, 0, 0, 0}), 1)
	})
	t.Run("Add", func(t *testing.T) {
		got := (Quat{1, 2, 3, 4}).Add(Quat{5, 6, 7, 8})
		equalQuat(t, "quat", got, Quat{6, 8, 10, 12})
	})
	t.Run("Subtract", func(t *testing.T) {
		got := (Quat{5, 6, 7, 8}).Sub(Quat{1, 2, 3, 4})
		equalQuat(t, "quat", got, Quat{4, 4, 4, 4})
	})
	t.Run("MultiplyScalar", func(t *testing.T) {
		got := (Quat{1, 2, 3, 4}).MulScalar(2)
		equalQuat(t, "quat", got, Quat{2, 4, 6, 8})
	})
	t.Run("Rotate", func(t *testing.T) {
		angle := Pi / 2
		q := Quat{0, float32(math.Sin(float64(angle / 2))), 0, float32(math.Cos(float64(angle / 2)))}
		got := q.Rotate(Vec3{1, 0, 0})
		nearVec3(t, "rotated", got, Vec3{0, 0, -1})
	})
	t.Run("Slerp", func(t *testing.T) {
		got := QuatIdentity().Slerp(QuatIdentity(), 0.5)
		nearQuat(t, "quat", got, Quat{0, 0, 0, 1})
	})
	t.Run("Lerp", func(t *testing.T) {
		got := QuatIdentity().Lerp(QuatIdentity(), 0.5)
		nearQuat(t, "quat", got, Quat{0, 0, 0, 1})
	})
}

func TestMathMat4(t *testing.T) {
	t.Run("Identity", func(t *testing.T) {
		m := Mat4Identity()
		equalFloat(t, "m00", m.Col[0].X, 1)
		equalFloat(t, "m11", m.Col[1].Y, 1)
		equalFloat(t, "m22", m.Col[2].Z, 1)
		equalFloat(t, "m33", m.Col[3].W, 1)
		equalFloat(t, "m01", m.Col[0].Y, 0)
		equalFloat(t, "m02", m.Col[0].Z, 0)
		equalFloat(t, "m03", m.Col[0].W, 0)
	})
	t.Run("Zero", func(t *testing.T) {
		m := Mat4Zero()
		for i := range 4 {
			equalFloat(t, "x", m.Col[i].X, 0)
			equalFloat(t, "y", m.Col[i].Y, 0)
			equalFloat(t, "z", m.Col[i].Z, 0)
			equalFloat(t, "w", m.Col[i].W, 0)
		}
	})
	t.Run("Translation", func(t *testing.T) {
		m := Mat4Translation(Vec3{1, 2, 3})
		equalFloat(t, "t.x", m.Col[3].X, 1)
		equalFloat(t, "t.y", m.Col[3].Y, 2)
		equalFloat(t, "t.z", m.Col[3].Z, 3)
		equalFloat(t, "t.w", m.Col[3].W, 1)
		equalFloat(t, "m00", m.Col[0].X, 1)
		equalFloat(t, "m11", m.Col[1].Y, 1)
		equalFloat(t, "m22", m.Col[2].Z, 1)
	})
	t.Run("Scale", func(t *testing.T) {
		m := Mat4Scale(Vec3{2, 3, 4})
		equalFloat(t, "m00", m.Col[0].X, 2)
		equalFloat(t, "m11", m.Col[1].Y, 3)
		equalFloat(t, "m22", m.Col[2].Z, 4)
		equalFloat(t, "m33", m.Col[3].W, 1)
	})
	t.Run("Rotation", func(t *testing.T) {
		m := Mat4Rotation(QuatIdentity())
		nearFloat(t, "m00", m.Col[0].X, 1)
		nearFloat(t, "m11", m.Col[1].Y, 1)
		nearFloat(t, "m22", m.Col[2].Z, 1)
		nearFloat(t, "m33", m.Col[3].W, 1)
	})
	t.Run("Multiply", func(t *testing.T) {
		m := Mat4Identity().Mul(Mat4Identity())
		nearFloat(t, "m00", m.Col[0].X, 1)
		nearFloat(t, "m11", m.Col[1].Y, 1)
		nearFloat(t, "m22", m.Col[2].Z, 1)
		nearFloat(t, "m33", m.Col[3].W, 1)
	})
	t.Run("Add", func(t *testing.T) {
		m := Mat4Identity().Add(Mat4Identity())
		equalFloat(t, "m00", m.Col[0].X, 2)
		equalFloat(t, "m11", m.Col[1].Y, 2)
		equalFloat(t, "m22", m.Col[2].Z, 2)
		equalFloat(t, "m33", m.Col[3].W, 2)
	})
	t.Run("Subtract", func(t *testing.T) {
		m := Mat4Identity().Sub(Mat4Identity())
		for i := range 4 {
			equalFloat(t, "x", m.Col[i].X, 0)
			equalFloat(t, "y", m.Col[i].Y, 0)
			equalFloat(t, "z", m.Col[i].Z, 0)
			equalFloat(t, "w", m.Col[i].W, 0)
		}
	})
	t.Run("MultiplyScalar", func(t *testing.T) {
		m := Mat4Identity().MulScalar(2)
		equalFloat(t, "m00", m.Col[0].X, 2)
		equalFloat(t, "m11", m.Col[1].Y, 2)
		equalFloat(t, "m22", m.Col[2].Z, 2)
		equalFloat(t, "m33", m.Col[3].W, 2)
	})
	t.Run("Transposed", func(t *testing.T) {
		m := Mat4Identity()
		m.Col[0].Y = 1 // [row 1, col 0] = 1
		got := m.Transposed()
		equalFloat(t, "m10", got.Col[1].X, 1) // [row 0, col 1] = 1
	})
	t.Run("GetAxisX", func(t *testing.T) { equalVec3(t, "axis", Mat4Identity().AxisX(), Vec3{1, 0, 0}) })
	t.Run("GetAxisY", func(t *testing.T) { equalVec3(t, "axis", Mat4Identity().AxisY(), Vec3{0, 1, 0}) })
	t.Run("GetAxisZ", func(t *testing.T) { equalVec3(t, "axis", Mat4Identity().AxisZ(), Vec3{0, 0, 1}) })
	t.Run("GetTranslation", func(t *testing.T) {
		m := Mat4Translation(Vec3{1, 2, 3})
		equalVec3(t, "translation", m.Translation(), Vec3{1, 2, 3})
	})
	t.Run("GetQuaternion", func(t *testing.T) {
		q := Mat4Rotation(QuatIdentity()).Quaternion()
		nearFloat(t, "w", float32(math.Abs(float64(q.W))), 1)
	})
	t.Run("RotationTranslation", func(t *testing.T) {
		m := Mat4RotationTranslation(QuatIdentity(), Vec3{1, 2, 3})
		equalFloat(t, "t.x", m.Col[3].X, 1)
		equalFloat(t, "t.y", m.Col[3].Y, 2)
		equalFloat(t, "t.z", m.Col[3].Z, 3)
		nearFloat(t, "m00", m.Col[0].X, 1)
		nearFloat(t, "m11", m.Col[1].Y, 1)
		nearFloat(t, "m22", m.Col[2].Z, 1)
	})
}

func TestMathUtility(t *testing.T) {
	t.Run("Sin", func(t *testing.T) {
		nearFloat(t, "sin0", float32(math.Sin(0)), 0)
		nearFloat(t, "sinHalfPi", float32(math.Sin(float64(Pi/2))), 1)
	})
	t.Run("Cos", func(t *testing.T) {
		nearFloat(t, "cos0", float32(math.Cos(0)), 1)
		nearFloat(t, "cosPi", float32(math.Cos(float64(Pi))), -1)
	})
}

func TestMathRayCast(t *testing.T) {
	t.Run("GetPointOnRay", func(t *testing.T) {
		got := RayCastGetPointOnRay(Vec3{0, 0, 0}, Vec3{1, 0, 0}, 0.5)
		equalVec3(t, "point", got, Vec3{0.5, 0, 0})
	})
}
