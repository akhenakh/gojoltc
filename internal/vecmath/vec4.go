package vecmath

// XYZ returns the first three components of v as a Vec3.
func (v Vec4) XYZ() Vec3 { return Vec3{v.X, v.Y, v.Z} }

// Add returns v + o.
func (v Vec4) Add(o Vec4) Vec4 { return Vec4{v.X + o.X, v.Y + o.Y, v.Z + o.Z, v.W + o.W} }

// Sub returns v - o.
func (v Vec4) Sub(o Vec4) Vec4 { return Vec4{v.X - o.X, v.Y - o.Y, v.Z - o.Z, v.W - o.W} }

// Mul returns the component-wise product v * o.
func (v Vec4) Mul(o Vec4) Vec4 { return Vec4{v.X * o.X, v.Y * o.Y, v.Z * o.Z, v.W * o.W} }

// MulScalar returns v * s.
func (v Vec4) MulScalar(s float32) Vec4 { return Vec4{v.X * s, v.Y * s, v.Z * s, v.W * s} }

// Dot returns the dot product v·o.
func (v Vec4) Dot(o Vec4) float32 { return v.X*o.X + v.Y*o.Y + v.Z*o.Z + v.W*o.W }

// LengthSquared returns the squared magnitude |v|².
func (v Vec4) LengthSquared() float32 { return v.Dot(v) }

// Length returns the magnitude |v|.
func (v Vec4) Length() float32 { return float32(sqrt32(v.LengthSquared())) }

// Normalized returns v scaled to unit length.
func (v Vec4) Normalized() Vec4 { return v.MulScalar(1 / v.Length()) }
