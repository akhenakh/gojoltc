//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import "github.com/akhenakh/gojoltc/internal/vecmath"

// Conversions between pure-Go vecmath values and their C counterparts. These
// keep C types confined to this package.

func vec3ToC(v vecmath.Vec3) C.JPH_Vec3 {
	return C.JPH_Vec3{x: C.float(v.X), y: C.float(v.Y), z: C.float(v.Z)}
}

func vec3FromC(c C.JPH_Vec3) vecmath.Vec3 {
	return vecmath.Vec3{X: float32(c.x), Y: float32(c.y), Z: float32(c.z)}
}

func vec4ToC(v vecmath.Vec4) C.JPH_Vec4 {
	return C.JPH_Vec4{x: C.float(v.X), y: C.float(v.Y), z: C.float(v.Z), w: C.float(v.W)}
}

func vec4FromC(c C.JPH_Vec4) vecmath.Vec4 {
	return vecmath.Vec4{X: float32(c.x), Y: float32(c.y), Z: float32(c.z), W: float32(c.w)}
}

func quatToC(q vecmath.Quat) C.JPH_Quat {
	return C.JPH_Quat{x: C.float(q.X), y: C.float(q.Y), z: C.float(q.Z), w: C.float(q.W)}
}

func quatFromC(c C.JPH_Quat) vecmath.Quat {
	return vecmath.Quat{X: float32(c.x), Y: float32(c.y), Z: float32(c.z), W: float32(c.w)}
}

func mat4ToC(m vecmath.Mat4) C.JPH_Mat4 {
	var out C.JPH_Mat4
	for i := 0; i < 4; i++ {
		out.column[i] = vec4ToC(m.Col[i])
	}
	return out
}

func mat4FromC(c C.JPH_Mat4) vecmath.Mat4 {
	var out vecmath.Mat4
	for i := 0; i < 4; i++ {
		out.Col[i] = vec4FromC(c.column[i])
	}
	return out
}

func aaboxToC(b vecmath.AABox) C.JPH_AABox {
	return C.JPH_AABox{min: vec3ToC(b.Min), max: vec3ToC(b.Max)}
}

func aaboxFromC(c C.JPH_AABox) vecmath.AABox {
	return vecmath.AABox{Min: vec3FromC(c.min), Max: vec3FromC(c.max)}
}

func planeToC(p vecmath.Plane) C.JPH_Plane {
	return C.JPH_Plane{normal: vec3ToC(p.Normal), distance: C.float(p.Distance)}
}

func planeFromC(c C.JPH_Plane) vecmath.Plane {
	return vecmath.Plane{Normal: vec3FromC(c.normal), Distance: float32(c.distance)}
}

func colorToC(c vecmath.Color) C.JPH_Color { return C.JPH_Color(c) }

func colorFromC(c C.JPH_Color) vecmath.Color { return vecmath.Color(c) }
