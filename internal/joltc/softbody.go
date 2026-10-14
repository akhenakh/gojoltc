//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

import (
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// Compile-time checks that the Go soft body value types are byte-compatible
// with the C structs we pass slices of.
var (
	_ [0]struct{} = [int(unsafe.Sizeof(SoftVertex{}) - unsafe.Sizeof(C.JPH_SoftVertex{}))]struct{}{}
	_ [0]struct{} = [int(unsafe.Sizeof(SoftFace{}) - unsafe.Sizeof(C.JPH_SoftFace{}))]struct{}{}
)

// SoftBodySharedSettings describes the mesh and constraints shared by soft
// bodies of the same type.
type SoftBodySharedSettings struct{ ptr *C.JPH_SoftBodySharedSettings }

// NewSoftBodySharedSettings creates empty shared settings.
func NewSoftBodySharedSettings() *SoftBodySharedSettings {
	return &SoftBodySharedSettings{ptr: C.JPH_SoftBodySharedSettings_Create()}
}

// Destroy releases the settings.
func (s *SoftBodySharedSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_SoftBodySharedSettings_Destroy(s.ptr)
		s.ptr = nil
	}
}

func (s *SoftBodySharedSettings) addVertices(vertices []SoftVertex) {
	n := len(vertices)
	if n == 0 {
		return
	}
	C.JPH_SoftBodySharedSettings_AddVertices(s.ptr, (*C.JPH_SoftVertex)(unsafe.Pointer(&vertices[0])), C.uint32_t(n))
}

// AddVertices adds vertices to the shared settings.
func (s *SoftBodySharedSettings) AddVertices(vertices []SoftVertex) { s.addVertices(vertices) }

// AddFaces adds triangle faces to the shared settings.
func (s *SoftBodySharedSettings) AddFaces(faces []SoftFace) {
	n := len(faces)
	if n == 0 {
		return
	}
	C.JPH_SoftBodySharedSettings_AddFaces(s.ptr, (*C.JPH_SoftFace)(unsafe.Pointer(&faces[0])), C.uint32_t(n))
}

// VertexCount returns the number of vertices.
func (s *SoftBodySharedSettings) VertexCount() uint32 {
	return uint32(C.JPH_SoftBodySharedSettings_GetVertexCount(s.ptr))
}

// FaceCount returns the number of faces.
func (s *SoftBodySharedSettings) FaceCount() uint32 {
	return uint32(C.JPH_SoftBodySharedSettings_GetFaceCount(s.ptr))
}

// CreateConstraints generates edge/bend constraints from the faces.
func (s *SoftBodySharedSettings) CreateConstraints(compliance float32, bendType SoftBodyBendType) {
	C.JPH_SoftBodySharedSettings_CreateConstraints(s.ptr, C.float(compliance), C.JPH_SoftBodyBendType(bendType))
}

// Optimize optimizes the constraint order.
func (s *SoftBodySharedSettings) Optimize() { C.JPH_SoftBodySharedSettings_Optimize(s.ptr) }

// SoftBodyCreationSettings describes a soft body to be created.
type SoftBodyCreationSettings struct {
	ptr *C.JPH_SoftBodyCreationSettings
}

// NewSoftBodyCreationSettings creates settings with Jolt defaults.
func NewSoftBodyCreationSettings() *SoftBodyCreationSettings {
	return &SoftBodyCreationSettings{ptr: C.JPH_SoftBodyCreationSettings_Create()}
}

// NewSoftBodyCreationSettingsFor creates settings for a shared mesh at a pose.
func NewSoftBodyCreationSettingsFor(shared *SoftBodySharedSettings, position vecmath.RVec3, rotation vecmath.Quat, layer ObjectLayer) *SoftBodyCreationSettings {
	p := rvec3ToC(position)
	r := quatToC(rotation)
	return &SoftBodyCreationSettings{ptr: C.JPH_SoftBodyCreationSettings_Create2(shared.ptr, &p, &r, C.JPH_ObjectLayer(layer))}
}

// Destroy releases the settings.
func (s *SoftBodyCreationSettings) Destroy() {
	if s.ptr != nil {
		C.JPH_SoftBodyCreationSettings_Destroy(s.ptr)
		s.ptr = nil
	}
}

// SharedSettings returns the shared mesh settings.
func (s *SoftBodyCreationSettings) SharedSettings() *SoftBodySharedSettings {
	return &SoftBodySharedSettings{ptr: (*C.JPH_SoftBodySharedSettings)(unsafe.Pointer(C.JPH_SoftBodyCreationSettings_GetSettings(s.ptr)))}
}

// SetObjectLayer sets the object layer.
func (s *SoftBodyCreationSettings) SetObjectLayer(layer ObjectLayer) {
	C.JPH_SoftBodyCreationSettings_SetObjectLayer(s.ptr, C.JPH_ObjectLayer(layer))
}

// SetNumIterations sets the number of solver iterations.
func (s *SoftBodyCreationSettings) SetNumIterations(v uint32) {
	C.JPH_SoftBodyCreationSettings_SetNumIterations(s.ptr, C.uint32_t(v))
}

// SetFriction sets the friction.
func (s *SoftBodyCreationSettings) SetFriction(v float32) {
	C.JPH_SoftBodyCreationSettings_SetFriction(s.ptr, C.float(v))
}

// SetRestitution sets the restitution.
func (s *SoftBodyCreationSettings) SetRestitution(v float32) {
	C.JPH_SoftBodyCreationSettings_SetRestitution(s.ptr, C.float(v))
}

// SetPressure sets the pressure (inflates the body).
func (s *SoftBodyCreationSettings) SetPressure(v float32) {
	C.JPH_SoftBodyCreationSettings_SetPressure(s.ptr, C.float(v))
}

// SetGravityFactor sets the gravity factor.
func (s *SoftBodyCreationSettings) SetGravityFactor(v float32) {
	C.JPH_SoftBodyCreationSettings_SetGravityFactor(s.ptr, C.float(v))
}

// SetVertexRadius sets the collision radius of the vertices.
func (s *SoftBodyCreationSettings) SetVertexRadius(v float32) {
	C.JPH_SoftBodyCreationSettings_SetVertexRadius(s.ptr, C.float(v))
}

// --- BodyInterface -----------------------------------------------------------

// CreateSoftBody creates a soft body without adding it to the system.
func (bi *BodyInterface) CreateSoftBody(settings *SoftBodyCreationSettings) *Body {
	return &Body{ptr: C.JPH_BodyInterface_CreateSoftBody(bi.ptr, settings.ptr)}
}

// CreateAndAddSoftBody creates a soft body and adds it to the system.
func (bi *BodyInterface) CreateAndAddSoftBody(settings *SoftBodyCreationSettings, activation Activation) BodyID {
	return BodyID(C.JPH_BodyInterface_CreateAndAddSoftBody(bi.ptr, settings.ptr, C.JPH_Activation(activation)))
}

// --- Body soft body accessors ------------------------------------------------

// SoftBodyVertexPosition returns the world position of a soft body vertex.
func (b *Body) SoftBodyVertexPosition(index uint32) vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_Body_GetSoftBodyVertexPosition(b.ptr, C.uint32_t(index), &o)
	return vec3FromC(o)
}

// SoftBodyVertexPositions writes the world positions of all soft body vertices
// into dst and returns the number written.
func (b *Body) SoftBodyVertexPositions(dst []vecmath.Vec3) int {
	n := len(dst)
	if n == 0 {
		return 0
	}
	var count C.uint32_t
	C.JPH_Body_GetSoftBodyVertexPositions(b.ptr, (*C.JPH_Vec3)(unsafe.Pointer(&dst[0])), C.uint32_t(n), &count)
	return int(count)
}
