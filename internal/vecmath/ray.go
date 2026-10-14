package vecmath

// RayCast is a ray defined by an origin point and a direction vector.
type RayCast struct {
	Origin    Vec3
	Direction Vec3
}

// GetPointOnRay returns the point at fraction along the ray.
func (r RayCast) GetPointOnRay(fraction float32) Vec3 {
	return r.Origin.Add(r.Direction.MulScalar(fraction))
}

// RayCastGetPointOnRay returns the point at fraction along the ray defined by
// origin and direction (Jolt's JPH_RayCast_GetPointOnRay).
func RayCastGetPointOnRay(origin, direction Vec3, fraction float32) Vec3 {
	return origin.Add(direction.MulScalar(fraction))
}
