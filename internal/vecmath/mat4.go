package vecmath

// Mat4Identity returns the identity matrix.
func Mat4Identity() Mat4 {
	return Mat4{Col: [4]Vec4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}}
}

// Mat4Zero returns the zero matrix.
func Mat4Zero() Mat4 { return Mat4{} }

// Mat4Rotation returns a rotation matrix from quaternion q
// (Jolt's Mat44::sRotation).
func Mat4Rotation(q Quat) Mat4 {
	x, y, z, w := q.X, q.Y, q.Z, q.W

	x2, y2, z2 := x+x, y+y, z+z
	xx, xy, xz := x*x2, x*y2, x*z2
	yy, yz, zz := y*y2, y*z2, z*z2
	wx, wy, wz := w*x2, w*y2, w*z2

	return Mat4{Col: [4]Vec4{
		{1 - (yy + zz), xy + wz, xz - wy, 0},
		{xy - wz, 1 - (xx + zz), yz + wx, 0},
		{xz + wy, yz - wx, 1 - (xx + yy), 0},
		{0, 0, 0, 1},
	}}
}

// Mat4RotationAxisAngle returns a rotation matrix of angle radians around axis.
func Mat4RotationAxisAngle(axis Vec3, angle float32) Mat4 {
	return Mat4Rotation(QuatRotation(axis, angle))
}

// Mat4Translation returns a translation matrix.
func Mat4Translation(t Vec3) Mat4 {
	m := Mat4Identity()
	m.Col[3] = Vec4{t.X, t.Y, t.Z, 1}
	return m
}

// Mat4RotationTranslation returns a rotation followed by a translation.
func Mat4RotationTranslation(q Quat, t Vec3) Mat4 {
	m := Mat4Rotation(q)
	m.Col[3] = Vec4{t.X, t.Y, t.Z, 1}
	return m
}

// Mat4InverseRotationTranslation returns the inverse of
// Mat4RotationTranslation(q, t).
func Mat4InverseRotationTranslation(q Quat, t Vec3) Mat4 {
	m := Mat4Rotation(q.Conjugated())
	m.Col[3] = Vec4{
		-(m.Col[0].X*t.X + m.Col[1].X*t.Y + m.Col[2].X*t.Z),
		-(m.Col[0].Y*t.X + m.Col[1].Y*t.Y + m.Col[2].Y*t.Z),
		-(m.Col[0].Z*t.X + m.Col[1].Z*t.Y + m.Col[2].Z*t.Z),
		1,
	}
	return m
}

// Mat4Scale returns a scale matrix.
func Mat4Scale(s Vec3) Mat4 {
	return Mat4{Col: [4]Vec4{
		{s.X, 0, 0, 0},
		{0, s.Y, 0, 0},
		{0, 0, s.Z, 0},
		{0, 0, 0, 1},
	}}
}

// Add returns m + o.
func (m Mat4) Add(o Mat4) Mat4 {
	var r Mat4
	for i := 0; i < 4; i++ {
		r.Col[i] = m.Col[i].Add(o.Col[i])
	}
	return r
}

// Sub returns m - o.
func (m Mat4) Sub(o Mat4) Mat4 {
	var r Mat4
	for i := 0; i < 4; i++ {
		r.Col[i] = m.Col[i].Sub(o.Col[i])
	}
	return r
}

// Mul returns the matrix product m * o.
func (m Mat4) Mul(o Mat4) Mat4 {
	var r Mat4
	for i := 0; i < 4; i++ {
		r.Col[i] = m.Col[0].MulScalar(o.Col[i].X).
			Add(m.Col[1].MulScalar(o.Col[i].Y)).
			Add(m.Col[2].MulScalar(o.Col[i].Z)).
			Add(m.Col[3].MulScalar(o.Col[i].W))
	}
	return r
}

// MulScalar returns m * s.
func (m Mat4) MulScalar(s float32) Mat4 {
	var r Mat4
	for i := 0; i < 4; i++ {
		r.Col[i] = m.Col[i].MulScalar(s)
	}
	return r
}

// Transposed returns the transpose of m.
func (m Mat4) Transposed() Mat4 {
	var r Mat4
	for c := 0; c < 4; c++ {
		r.Col[c] = Vec4{m.Col[0].at(c), m.Col[1].at(c), m.Col[2].at(c), m.Col[3].at(c)}
	}
	return r
}

// Inversed returns the inverse of m.
func (m Mat4) Inversed() Mat4 {
	var a [16]float32
	for c := 0; c < 4; c++ {
		a[0*4+c] = m.Col[c].X
		a[1*4+c] = m.Col[c].Y
		a[2*4+c] = m.Col[c].Z
		a[3*4+c] = m.Col[c].W
	}
	inv := invert4(a)
	var r Mat4
	for c := 0; c < 4; c++ {
		r.Col[c] = Vec4{inv[0*4+c], inv[1*4+c], inv[2*4+c], inv[3*4+c]}
	}
	return r
}

// AxisX returns the X (first) column as a vector.
func (m Mat4) AxisX() Vec3 { return m.Col[0].XYZ() }

// AxisY returns the Y (second) column as a vector.
func (m Mat4) AxisY() Vec3 { return m.Col[1].XYZ() }

// AxisZ returns the Z (third) column as a vector.
func (m Mat4) AxisZ() Vec3 { return m.Col[2].XYZ() }

// Translation returns the translation column as a vector.
func (m Mat4) Translation() Vec3 { return m.Col[3].XYZ() }

// Quaternion returns the rotation encoded in m (Jolt's Mat44::GetQuaternion).
func (m Mat4) Quaternion() Quat {
	tr := m.Col[0].X + m.Col[1].Y + m.Col[2].Z
	if tr >= 0 {
		s := sqrt32(tr + 1)
		is := 0.5 / s
		return Quat{
			(m.Col[1].Z - m.Col[2].Y) * is,
			(m.Col[2].X - m.Col[0].Z) * is,
			(m.Col[0].Y - m.Col[1].X) * is,
			0.5 * s,
		}
	}

	i := 0
	if m.Col[1].Y > m.Col[0].X {
		i = 1
	}
	if m.at(i, i) < m.at(2, 2) {
		i = 2
	}

	switch i {
	case 0:
		s := sqrt32(m.Col[0].X - (m.Col[1].Y + m.Col[2].Z) + 1)
		is := 0.5 / s
		return Quat{
			0.5 * s,
			(m.Col[1].X + m.Col[0].Y) * is,
			(m.Col[0].Z + m.Col[2].X) * is,
			(m.Col[1].Z - m.Col[2].Y) * is,
		}
	case 1:
		s := sqrt32(m.Col[1].Y - (m.Col[2].Z + m.Col[0].X) + 1)
		is := 0.5 / s
		return Quat{
			(m.Col[1].X + m.Col[0].Y) * is,
			0.5 * s,
			(m.Col[2].Y + m.Col[1].Z) * is,
			(m.Col[2].X - m.Col[0].Z) * is,
		}
	default:
		s := sqrt32(m.Col[2].Z - (m.Col[0].X + m.Col[1].Y) + 1)
		is := 0.5 / s
		return Quat{
			(m.Col[0].Z + m.Col[2].X) * is,
			(m.Col[2].Y + m.Col[1].Z) * is,
			0.5 * s,
			(m.Col[0].Y - m.Col[1].X) * is,
		}
	}
}

// at returns element (row r, column c).
func (m Mat4) at(r, c int) float32 { return m.Col[c].at(r) }

// at returns component i (0=X, 1=Y, 2=Z, 3=W).
func (v Vec4) at(i int) float32 {
	switch i {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	default:
		return v.W
	}
}

// invert4 inverts a row-major 4x4 matrix (a[row*4+col]).
func invert4(a [16]float32) [16]float32 {
	inv := [16]float32{
		a[5]*a[10]*a[15] - a[5]*a[11]*a[14] - a[9]*a[6]*a[15] + a[9]*a[7]*a[14] + a[13]*a[6]*a[11] - a[13]*a[7]*a[10],
		-a[4]*a[10]*a[15] + a[4]*a[11]*a[14] + a[8]*a[6]*a[15] - a[8]*a[7]*a[14] - a[12]*a[6]*a[11] + a[12]*a[7]*a[10],
		a[4]*a[9]*a[15] - a[4]*a[11]*a[13] - a[8]*a[5]*a[15] + a[8]*a[7]*a[13] + a[12]*a[5]*a[11] - a[12]*a[7]*a[9],
		-a[4]*a[9]*a[14] + a[4]*a[10]*a[13] + a[8]*a[5]*a[14] - a[8]*a[6]*a[13] - a[12]*a[5]*a[10] + a[12]*a[6]*a[9],
		-a[1]*a[10]*a[15] + a[1]*a[11]*a[14] + a[9]*a[2]*a[15] - a[9]*a[3]*a[14] - a[13]*a[2]*a[11] + a[13]*a[3]*a[10],
		a[0]*a[10]*a[15] - a[0]*a[11]*a[14] - a[8]*a[2]*a[15] + a[8]*a[3]*a[14] + a[12]*a[2]*a[11] - a[12]*a[3]*a[10],
		-a[0]*a[9]*a[15] + a[0]*a[11]*a[13] + a[8]*a[1]*a[15] - a[8]*a[3]*a[13] - a[12]*a[1]*a[11] + a[12]*a[3]*a[9],
		a[0]*a[9]*a[14] - a[0]*a[10]*a[13] - a[8]*a[1]*a[14] + a[8]*a[2]*a[13] + a[12]*a[1]*a[10] - a[12]*a[2]*a[9],
		a[1]*a[6]*a[15] - a[1]*a[7]*a[14] - a[5]*a[2]*a[15] + a[5]*a[3]*a[14] + a[13]*a[2]*a[7] - a[13]*a[3]*a[6],
		-a[0]*a[6]*a[15] + a[0]*a[7]*a[14] + a[4]*a[2]*a[15] - a[4]*a[3]*a[14] - a[12]*a[2]*a[7] + a[12]*a[3]*a[6],
		a[0]*a[5]*a[15] - a[0]*a[7]*a[13] - a[4]*a[1]*a[15] + a[4]*a[3]*a[13] + a[12]*a[1]*a[7] - a[12]*a[3]*a[5],
		-a[0]*a[5]*a[14] + a[0]*a[6]*a[13] + a[4]*a[1]*a[14] - a[4]*a[2]*a[13] - a[12]*a[1]*a[6] + a[12]*a[2]*a[5],
		-a[1]*a[6]*a[11] + a[1]*a[7]*a[10] + a[5]*a[2]*a[11] - a[5]*a[3]*a[10] - a[9]*a[2]*a[7] + a[9]*a[3]*a[6],
		a[0]*a[6]*a[11] - a[0]*a[7]*a[10] - a[4]*a[2]*a[11] + a[4]*a[3]*a[10] + a[8]*a[2]*a[7] - a[8]*a[3]*a[6],
		-a[0]*a[5]*a[11] + a[0]*a[7]*a[9] + a[4]*a[1]*a[11] - a[4]*a[3]*a[9] - a[8]*a[1]*a[7] + a[8]*a[3]*a[5],
		a[0]*a[5]*a[10] - a[0]*a[6]*a[9] - a[4]*a[1]*a[10] + a[4]*a[2]*a[9] + a[8]*a[1]*a[6] - a[8]*a[2]*a[5],
	}

	det := a[0]*inv[0] + a[1]*inv[4] + a[2]*inv[8] + a[3]*inv[12]
	if det == 0 {
		return [16]float32{}
	}
	invDet := 1 / det
	for i := range inv {
		inv[i] *= invDet
	}
	return inv
}
