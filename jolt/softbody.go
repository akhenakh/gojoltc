package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Soft body types, re-exported from the bindings.
type (
	// SoftVertex is a soft body vertex.
	SoftVertex = joltc.SoftVertex
	// SoftFace is a triangle of a soft body surface.
	SoftFace = joltc.SoftFace
	// SoftBodySharedSettings describes the mesh shared by soft bodies.
	SoftBodySharedSettings = joltc.SoftBodySharedSettings
	// SoftBodyCreationSettings describes a soft body to be created.
	SoftBodyCreationSettings = joltc.SoftBodyCreationSettings
)

// NewSoftBodySharedSettings creates empty shared settings.
func NewSoftBodySharedSettings() *SoftBodySharedSettings { return joltc.NewSoftBodySharedSettings() }

// NewSoftBodyCreationSettings creates soft body settings with Jolt defaults.
func NewSoftBodyCreationSettings() *SoftBodyCreationSettings {
	return joltc.NewSoftBodyCreationSettings()
}

// NewSoftBodyCreationSettingsFor creates settings for a shared mesh at a pose.
func NewSoftBodyCreationSettingsFor(shared *SoftBodySharedSettings, position RVec3, rotation Quat, layer ObjectLayer) *SoftBodyCreationSettings {
	return joltc.NewSoftBodyCreationSettingsFor(shared, position, rotation, layer)
}
