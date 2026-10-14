package vecmath

import "math"

// Pi is the value of pi as float32 (Jolt's JPH_PI).
const Pi float32 = 3.14159265358979323846

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func sqrt32(v float32) float32 { return float32(math.Sqrt(float64(v))) }

// Vec3AxisX returns the unit X axis.
func Vec3AxisX() Vec3 { return Vec3{1, 0, 0} }

// Vec3AxisY returns the unit Y axis.
func Vec3AxisY() Vec3 { return Vec3{0, 1, 0} }

// Vec3AxisZ returns the unit Z axis.
func Vec3AxisZ() Vec3 { return Vec3{0, 0, 1} }

// Add returns v + o.
func (v Vec3) Add(o Vec3) Vec3 { return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z} }

// Sub returns v - o.
func (v Vec3) Sub(o Vec3) Vec3 { return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z} }

// Mul returns the component-wise product v * o.
func (v Vec3) Mul(o Vec3) Vec3 { return Vec3{v.X * o.X, v.Y * o.Y, v.Z * o.Z} }

// Div returns the component-wise quotient v / o.
func (v Vec3) Div(o Vec3) Vec3 { return Vec3{v.X / o.X, v.Y / o.Y, v.Z / o.Z} }

// MulScalar returns v * s.
func (v Vec3) MulScalar(s float32) Vec3 { return Vec3{v.X * s, v.Y * s, v.Z * s} }

// DivScalar returns v / s.
func (v Vec3) DivScalar(s float32) Vec3 { return v.MulScalar(1 / s) }

// Negate returns -v.
func (v Vec3) Negate() Vec3 { return Vec3{-v.X, -v.Y, -v.Z} }

// Abs returns the component-wise absolute value.
func (v Vec3) Abs() Vec3 { return Vec3{abs32(v.X), abs32(v.Y), abs32(v.Z)} }

// Dot returns the dot product v·o.
func (v Vec3) Dot(o Vec3) float32 { return v.X*o.X + v.Y*o.Y + v.Z*o.Z }

// Cross returns the cross product v × o.
func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}

// LengthSquared returns the squared magnitude |v|².
func (v Vec3) LengthSquared() float32 { return v.Dot(v) }

// Length returns the magnitude |v|.
func (v Vec3) Length() float32 { return float32(math.Sqrt(float64(v.LengthSquared()))) }

// Normalized returns v scaled to unit length.
func (v Vec3) Normalized() Vec3 { return v.DivScalar(v.Length()) }

// Normalize is an alias for Normalized.
func (v Vec3) Normalize() Vec3 { return v.Normalized() }

// NormalizedOr returns the unit vector, or fallback when v is (near) zero.
func (v Vec3) NormalizedOr(fallback Vec3) Vec3 {
	if v.LengthSquared() > 0 {
		return v.Normalized()
	}
	return fallback
}

// IsClose reports whether o is within maxDistSq squared distance of v.
func (v Vec3) IsClose(o Vec3, maxDistSq float32) bool {
	return o.Sub(v).LengthSquared() <= maxDistSq
}

// IsNearZero reports whether |v|² is within maxDistSq.
func (v Vec3) IsNearZero(maxDistSq float32) bool { return v.LengthSquared() <= maxDistSq }

// IsNormalized reports whether |v| is 1 within tolerance.
func (v Vec3) IsNormalized(tolerance float32) bool {
	return abs32(1-v.LengthSquared()) <= tolerance
}

// IsNaN reports whether any component is NaN.
func (v Vec3) IsNaN() bool { return v.X != v.X || v.Y != v.Y || v.Z != v.Z }

// GetNormalizedPerpendicular returns a unit vector perpendicular to v.
func (v Vec3) GetNormalizedPerpendicular() Vec3 {
	xx := v.X * v.X
	yy := v.Y * v.Y
	zz := v.Z * v.Z

	var perp Vec3
	if xx > yy {
		perp = Vec3{v.Z, 0, -v.X}
	} else {
		perp = Vec3{0, v.Z, -v.Y}
	}
	invLen := 1 / float32(math.Sqrt(float64(max32(xx, yy)+zz)))
	return perp.MulScalar(invLen)
}

// MulMatrix transforms v as a point by m (rotation + translation).
func (v Vec3) MulMatrix(m Mat4) Vec3 {
	return Vec3{
		m.Col[0].X*v.X + m.Col[1].X*v.Y + m.Col[2].X*v.Z + m.Col[3].X,
		m.Col[0].Y*v.X + m.Col[1].Y*v.Y + m.Col[2].Y*v.Z + m.Col[3].Y,
		m.Col[0].Z*v.X + m.Col[1].Z*v.Y + m.Col[2].Z*v.Z + m.Col[3].Z,
	}
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
