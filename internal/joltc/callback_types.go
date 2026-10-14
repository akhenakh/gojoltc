package joltc

import "github.com/akhenakh/gojoltc/internal/vecmath"

// Pure value types shared by the callbacks implementation and the cgo-less
// stub. Anything holding a C pointer lives in the cgo files.

// SubShapeIDPair identifies a contact between two sub-shapes.
type SubShapeIDPair struct {
	Body1ID     BodyID
	SubShapeID1 SubShapeID
	Body2ID     BodyID
	SubShapeID2 SubShapeID
}

// CollideShapeResult summarises a collision for contact validation.
type CollideShapeResult struct {
	ContactPointOn1  vecmath.Vec3
	ContactPointOn2  vecmath.Vec3
	PenetrationAxis  vecmath.Vec3
	PenetrationDepth float32
	SubShapeID1      SubShapeID
	SubShapeID2      SubShapeID
	BodyID2          BodyID
}

// ContactListenerCallbacks holds the Go callbacks for a contact listener. Nil
// fields are ignored.
type ContactListenerCallbacks struct {
	OnValidate  func(body1, body2 *Body, baseOffset vecmath.RVec3, result CollideShapeResult) ValidateResult
	OnAdded     func(body1, body2 *Body, manifold *ContactManifold, settings *ContactSettings)
	OnPersisted func(body1, body2 *Body, manifold *ContactManifold, settings *ContactSettings)
	OnRemoved   func(pair SubShapeIDPair)
}

// BodyActivationCallbacks holds the Go callbacks for a body activation listener.
type BodyActivationCallbacks struct {
	OnActivated   func(bodyID BodyID, bodyUserData uint64)
	OnDeactivated func(bodyID BodyID, bodyUserData uint64)
}

// StepContext describes a simulation step.
type StepContext struct {
	DeltaTime   float32
	IsFirstStep bool
	IsLastStep  bool
}

// StepListenerCallbacks holds the Go callbacks for a step listener.
type StepListenerCallbacks struct {
	OnStep func(ctx StepContext)
}

// DebugRendererCallbacks holds the Go callbacks for a debug renderer.
type DebugRendererCallbacks struct {
	DrawLine     func(from, to vecmath.RVec3, color vecmath.Color)
	DrawTriangle func(v1, v2, v3 vecmath.RVec3, color vecmath.Color, castShadow bool)
	DrawText3D   func(position vecmath.RVec3, text string, color vecmath.Color, height float32)
}

// DrawSettings selects what the debug renderer draws.
type DrawSettings struct {
	DrawShape                 bool
	DrawShapeWireframe        bool
	DrawBoundingBox           bool
	DrawCenterOfMassTransform bool
	DrawVelocity              bool
	DrawMassAndInertia        bool
	DrawSupportingFace        bool

	// Soft bodies are NOT drawn via DrawShape: Jolt only draws them through the
	// soft-body flags below, so a cloth is invisible without at least
	// DrawSoftBodyVertices or DrawSoftBodyEdgeConstraints.
	DrawSoftBodyVertices          bool
	DrawSoftBodyVertexVelocities  bool
	DrawSoftBodyEdgeConstraints   bool
	DrawSoftBodyBendConstraints   bool
	DrawSoftBodyVolumeConstraints bool
	DrawSoftBodySkinConstraints   bool
	DrawSoftBodyLRAConstraints    bool
	DrawSoftBodyPredictedBounds   bool
}

// DefaultDrawSettings returns settings that draw all body shapes as wireframe.
func DefaultDrawSettings() DrawSettings {
	return DrawSettings{DrawShape: true, DrawShapeWireframe: true}
}
