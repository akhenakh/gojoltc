package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Listener and debug-render types, re-exported from the bindings.
type (
	// ContactManifold describes a contact between two bodies.
	ContactManifold = joltc.ContactManifold
	// ContactSettings lets a contact callback tweak the contact.
	ContactSettings = joltc.ContactSettings
	// SubShapeIDPair identifies a contact between two sub-shapes.
	SubShapeIDPair = joltc.SubShapeIDPair
	// CollideShapeResult summarises a collision for contact validation.
	CollideShapeResult = joltc.CollideShapeResult
	// ContactListenerCallbacks holds the Go callbacks for a contact listener.
	ContactListenerCallbacks = joltc.ContactListenerCallbacks
	// ContactListener receives contact events.
	ContactListener = joltc.ContactListener
	// BodyActivationCallbacks holds the Go callbacks for a body activation listener.
	BodyActivationCallbacks = joltc.BodyActivationCallbacks
	// BodyActivationListener receives body activation events.
	BodyActivationListener = joltc.BodyActivationListener
	// StepContext describes a simulation step.
	StepContext = joltc.StepContext
	// StepListenerCallbacks holds the Go callbacks for a step listener.
	StepListenerCallbacks = joltc.StepListenerCallbacks
	// StepListener receives a callback for each simulation step.
	StepListener = joltc.StepListener
	// DebugRendererCallbacks holds the Go callbacks for a debug renderer.
	DebugRendererCallbacks = joltc.DebugRendererCallbacks
	// DebugRenderer turns Jolt's debug draw calls into Go callbacks.
	DebugRenderer = joltc.DebugRenderer
	// DrawSettings selects what the debug renderer draws.
	DrawSettings = joltc.DrawSettings
)

// NewContactListener creates a contact listener with the given callbacks.
func NewContactListener(cb *ContactListenerCallbacks) *ContactListener {
	return joltc.NewContactListener(cb)
}

// NewBodyActivationListener creates a body activation listener.
func NewBodyActivationListener(cb *BodyActivationCallbacks) *BodyActivationListener {
	return joltc.NewBodyActivationListener(cb)
}

// NewStepListener creates a step listener.
func NewStepListener(cb *StepListenerCallbacks) *StepListener {
	return joltc.NewStepListener(cb)
}

// NewDebugRenderer creates a debug renderer with the given callbacks.
func NewDebugRenderer(cb *DebugRendererCallbacks) *DebugRenderer {
	return joltc.NewDebugRenderer(cb)
}

// DefaultDrawSettings returns settings that draw all body shapes as wireframe.
func DefaultDrawSettings() DrawSettings { return joltc.DefaultDrawSettings() }
