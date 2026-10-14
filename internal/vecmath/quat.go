package vecmath

import "math"

// QuatIdentity returns the identity quaternion.
func QuatIdentity() Quat { return Quat{0, 0, 0, 1} }

// QuatRotation returns a quaternion for a rotation of angle radians around axis
// (which must be normalized).
func QuatRotation(axis Vec3, angle float32) Quat {
	half := 0.5 * angle
	s := float32(math.Sin(float64(half)))
	c := float32(math.Cos(float64(half)))
	return Quat{axis.X * s, axis.Y * s, axis.Z * s, c}
}

// QuatFromTo returns the rotation that takes from to to (Jolt's Quat::sFromTo).
func QuatFromTo(from, to Vec3) Quat {
	lenV1V2 := sqrt32(from.LengthSquared() * to.LengthSquared())
	w := lenV1V2 + from.Dot(to)
	if w == 0 {
		if lenV1V2 == 0 {
			return QuatIdentity()
		}
		p := from.GetNormalizedPerpendicular()
		return Quat{p.X, p.Y, p.Z, 0}
	}
	v := from.Cross(to)
	return Quat{v.X, v.Y, v.Z, w}.Normalized()
}

// QuatEulerAngles builds a quaternion from Euler angles in radians
// (Jolt's Quat::sEulerAngles).
func QuatEulerAngles(angles Vec3) Quat {
	hx := 0.5 * angles.X
	hy := 0.5 * angles.Y
	hz := 0.5 * angles.Z
	sx, cx := float32(math.Sin(float64(hx))), float32(math.Cos(float64(hx)))
	sy, cy := float32(math.Sin(float64(hy))), float32(math.Cos(float64(hy)))
	sz, cz := float32(math.Sin(float64(hz))), float32(math.Cos(float64(hz)))
	return Quat{
		cz*sx*cy - sz*cx*sy,
		cz*cx*sy + sz*sx*cy,
		sz*cx*cy - cz*sx*sy,
		cz*cx*cy + sz*sx*sy,
	}
}

// XYZ returns the vector part of the quaternion.
func (q Quat) XYZ() Vec3 { return Vec3{q.X, q.Y, q.Z} }

// Add returns q + o.
func (q Quat) Add(o Quat) Quat { return Quat{q.X + o.X, q.Y + o.Y, q.Z + o.Z, q.W + o.W} }

// Sub returns q - o.
func (q Quat) Sub(o Quat) Quat { return Quat{q.X - o.X, q.Y - o.Y, q.Z - o.Z, q.W - o.W} }

// Mul returns the Hamilton product q * o.
func (q Quat) Mul(o Quat) Quat {
	return Quat{
		q.W*o.X + q.X*o.W + q.Y*o.Z - q.Z*o.Y,
		q.W*o.Y - q.X*o.Z + q.Y*o.W + q.Z*o.X,
		q.W*o.Z + q.X*o.Y - q.Y*o.X + q.Z*o.W,
		q.W*o.W - q.X*o.X - q.Y*o.Y - q.Z*o.Z,
	}
}

// MulScalar returns q * s.
func (q Quat) MulScalar(s float32) Quat { return Quat{q.X * s, q.Y * s, q.Z * s, q.W * s} }

// DivScalar returns q / s.
func (q Quat) DivScalar(s float32) Quat { return q.MulScalar(1 / s) }

// Negated returns -q.
func (q Quat) Negated() Quat { return Quat{-q.X, -q.Y, -q.Z, -q.W} }

// Dot returns the dot product q·o.
func (q Quat) Dot(o Quat) float32 { return q.X*o.X + q.Y*o.Y + q.Z*o.Z + q.W*o.W }

// LengthSquared returns the squared magnitude |q|².
func (q Quat) LengthSquared() float32 { return q.Dot(q) }

// Length returns the magnitude |q|.
func (q Quat) Length() float32 { return sqrt32(q.LengthSquared()) }

// Normalized returns q scaled to unit length.
func (q Quat) Normalized() Quat { return q.DivScalar(q.Length()) }

// Conjugated returns the conjugate of q.
func (q Quat) Conjugated() Quat { return Quat{-q.X, -q.Y, -q.Z, q.W} }

// Inversed returns the inverse of q.
func (q Quat) Inversed() Quat { return q.Conjugated().DivScalar(q.LengthSquared()) }

// Rotate rotates v by q.
func (q Quat) Rotate(v Vec3) Vec3 {
	qv := q.XYZ()
	t := qv.Cross(v).MulScalar(2)
	return v.Add(t.MulScalar(q.W)).Add(qv.Cross(t))
}

// InverseRotate rotates v by the inverse of q.
func (q Quat) InverseRotate(v Vec3) Vec3 { return q.Conjugated().Rotate(v) }

// RotateAxisX returns the X axis rotated by q.
func (q Quat) RotateAxisX() Vec3 { return q.Rotate(Vec3AxisX()) }

// RotateAxisY returns the Y axis rotated by q.
func (q Quat) RotateAxisY() Vec3 { return q.Rotate(Vec3AxisY()) }

// RotateAxisZ returns the Z axis rotated by q.
func (q Quat) RotateAxisZ() Vec3 { return q.Rotate(Vec3AxisZ()) }

// GetRotationAngle returns the rotation angle around axis.
func (q Quat) GetRotationAngle(axis Vec3) float32 {
	if q.W == 0 {
		return Pi
	}
	return 2 * float32(math.Atan(float64(q.XYZ().Dot(axis)/q.W)))
}

// GetAxisAngle returns the rotation axis and angle.
func (q Quat) GetAxisAngle() (Vec3, float32) {
	wPos := q
	if wPos.W < 0 {
		wPos = wPos.Negated()
	}
	absW := wPos.W
	if absW >= 1 {
		return Vec3{}, 0
	}
	angle := 2 * float32(math.Acos(float64(absW)))
	return wPos.XYZ().NormalizedOr(Vec3{}), angle
}

// GetTwist returns the twist component of q around axis.
func (q Quat) GetTwist(axis Vec3) Quat {
	d := q.XYZ().Dot(axis)
	twist := Quat{d * axis.X, d * axis.Y, d * axis.Z, q.W}
	twistLen := twist.LengthSquared()
	if twistLen != 0 {
		return twist.DivScalar(sqrt32(twistLen))
	}
	return QuatIdentity()
}

// GetSwingTwist decomposes q into a swing and a twist (both around X).
func (q Quat) GetSwingTwist() (swing, twist Quat) {
	s := sqrt32(q.W*q.W + q.X*q.X)
	if s != 0 {
		twist = Quat{q.X / s, 0, 0, q.W / s}
		swing = Quat{0, (q.W*q.Y - q.X*q.Z) / s, (q.W*q.Z + q.X*q.Y) / s, s}
	} else {
		twist = QuatIdentity()
		swing = q
	}
	return swing, twist
}

// GetPerpendicular returns a quaternion perpendicular to q (Jolt's
// Quat::GetPerpendicular).
func (q Quat) GetPerpendicular() Quat { return Quat{q.Y, -q.X, q.W, -q.Z} }

// Slerp returns the spherical linear interpolation between q and to.
func (q Quat) Slerp(to Quat, fraction float32) Quat {
	cos := q.Dot(to)
	target := to
	if cos < 0 {
		target = to.Negated()
		cos = -cos
	}
	if cos > 0.9995 {
		return q.MulScalar(1 - fraction).Add(target.MulScalar(fraction)).Normalized()
	}
	theta := float32(math.Acos(float64(cos)))
	sinTheta := float32(math.Sin(float64(theta)))
	wa := float32(math.Sin(float64((1-fraction)*theta))) / sinTheta
	wb := float32(math.Sin(float64(fraction*theta))) / sinTheta
	return q.MulScalar(wa).Add(target.MulScalar(wb))
}

// Lerp returns the linear interpolation between q and to, normalized.
func (q Quat) Lerp(to Quat, fraction float32) Quat {
	return q.MulScalar(1 - fraction).Add(to.MulScalar(fraction)).Normalized()
}

// GetEulerAngles returns the Euler angles (radians) of q
// (Jolt's Quat::GetEulerAngles).
func (q Quat) GetEulerAngles() Vec3 {
	ySq := q.Y * q.Y

	t0 := 2 * (q.W*q.X + q.Y*q.Z)
	t1 := 1 - 2*(q.X*q.X+ySq)

	t2 := 2 * (q.W*q.Y - q.Z*q.X)
	t2 = clamp(t2, -1, 1)

	t3 := 2 * (q.W*q.Z + q.X*q.Y)
	t4 := 1 - 2*(ySq+q.Z*q.Z)

	return Vec3{
		float32(math.Atan2(float64(t0), float64(t1))),
		float32(math.Asin(float64(t2))),
		float32(math.Atan2(float64(t3), float64(t4))),
	}
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
