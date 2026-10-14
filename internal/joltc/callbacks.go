//go:build cgo

package joltc

/*
#include <joltc.h>

extern void joltcInstallContactProcs(void);
extern void joltcInstallBodyActivationProcs(void);
extern void joltcInstallStepProcs(void);
extern void joltcInstallDebugRendererProcs(void);
extern void* joltc_voidp(uintptr_t v);
*/
import "C"

import (
	"runtime/cgo"
	"sync"
	"unsafe"

	"github.com/akhenakh/gojoltc/internal/vecmath"
)

// Callbacks are delivered from Jolt worker threads during PhysicsSystem.Update.
// A single process-wide procs table is installed per listener kind and
// multiplexed by the cgo.Handle passed as userData. Callbacks must not block,
// and must recover from panics (the trampolines do).

// --- Contact -----------------------------------------------------------------

// ContactManifold describes a contact between two bodies.
type ContactManifold struct{ ptr *C.JPH_ContactManifold }

// Normal returns the world-space contact normal.
func (m *ContactManifold) Normal() vecmath.Vec3 {
	var o C.JPH_Vec3
	C.JPH_ContactManifold_GetWorldSpaceNormal(m.ptr, &o)
	return vec3FromC(o)
}

// PenetrationDepth returns the penetration depth.
func (m *ContactManifold) PenetrationDepth() float32 {
	return float32(C.JPH_ContactManifold_GetPenetrationDepth(m.ptr))
}

// SubShapeID1 returns the sub-shape id on body 1.
func (m *ContactManifold) SubShapeID1() SubShapeID {
	return SubShapeID(C.JPH_ContactManifold_GetSubShapeID1(m.ptr))
}

// SubShapeID2 returns the sub-shape id on body 2.
func (m *ContactManifold) SubShapeID2() SubShapeID {
	return SubShapeID(C.JPH_ContactManifold_GetSubShapeID2(m.ptr))
}

// PointCount returns the number of contact points.
func (m *ContactManifold) PointCount() uint32 {
	return uint32(C.JPH_ContactManifold_GetPointCount(m.ptr))
}

// ContactPointOn1 returns the contact point on body 1.
func (m *ContactManifold) ContactPointOn1(index uint32) vecmath.RVec3 {
	var o C.JPH_RVec3
	C.JPH_ContactManifold_GetWorldSpaceContactPointOn1(m.ptr, C.uint32_t(index), &o)
	return rvec3FromC(o)
}

// ContactPointOn2 returns the contact point on body 2.
func (m *ContactManifold) ContactPointOn2(index uint32) vecmath.RVec3 {
	var o C.JPH_RVec3
	C.JPH_ContactManifold_GetWorldSpaceContactPointOn2(m.ptr, C.uint32_t(index), &o)
	return rvec3FromC(o)
}

// ContactSettings lets a contact callback tweak the contact. It is only valid
// for the duration of the callback.
type ContactSettings struct{ ptr *C.JPH_ContactSettings }

// SetFriction sets the combined friction.
func (s *ContactSettings) SetFriction(v float32) { s.ptr.combinedFriction = C.float(v) }

// SetRestitution sets the combined restitution.
func (s *ContactSettings) SetRestitution(v float32) { s.ptr.combinedRestitution = C.float(v) }

// SetIsSensor marks the contact as a sensor (no impulses applied).
func (s *ContactSettings) SetIsSensor(v bool) {
	if v {
		s.ptr.isSensor = 1
	} else {
		s.ptr.isSensor = 0
	}
}

// SetInvMassScale1 sets the inverse mass scale of body 1 (0 makes it immovable
// by this contact).
func (s *ContactSettings) SetInvMassScale1(v float32) { s.ptr.invMassScale1 = C.float(v) }

// SetInvMassScale2 sets the inverse mass scale of body 2.
func (s *ContactSettings) SetInvMassScale2(v float32) { s.ptr.invMassScale2 = C.float(v) }

// SetInvInertiaScale1 sets the inverse inertia scale of body 1.
func (s *ContactSettings) SetInvInertiaScale1(v float32) { s.ptr.invInertiaScale1 = C.float(v) }

// SetInvInertiaScale2 sets the inverse inertia scale of body 2.
func (s *ContactSettings) SetInvInertiaScale2(v float32) { s.ptr.invInertiaScale2 = C.float(v) }

// SetRelativeLinearSurfaceVelocity sets the relative linear surface velocity.
func (s *ContactSettings) SetRelativeLinearSurfaceVelocity(v vecmath.Vec3) {
	s.ptr.relativeLinearSurfaceVelocity = vec3ToC(v)
}

func subShapeIDPairFromC(p *C.JPH_SubShapeIDPair) SubShapeIDPair {
	return SubShapeIDPair{
		Body1ID:     BodyID(p.Body1ID),
		SubShapeID1: SubShapeID(p.subShapeID1),
		Body2ID:     BodyID(p.Body2ID),
		SubShapeID2: SubShapeID(p.subShapeID2),
	}
}

func collideShapeResultFromC(r *C.JPH_CollideShapeResult) CollideShapeResult {
	return CollideShapeResult{
		ContactPointOn1:  vec3FromC(r.contactPointOn1),
		ContactPointOn2:  vec3FromC(r.contactPointOn2),
		PenetrationAxis:  vec3FromC(r.penetrationAxis),
		PenetrationDepth: float32(r.penetrationDepth),
		SubShapeID1:      SubShapeID(r.subShapeID1),
		SubShapeID2:      SubShapeID(r.subShapeID2),
		BodyID2:          BodyID(r.bodyID2),
	}
}

// ContactListener receives contact events from the physics system.
type ContactListener struct {
	ptr    *C.JPH_ContactListener
	handle cgo.Handle
}

var contactProcsOnce sync.Once

// NewContactListener creates a contact listener with the given callbacks.
func NewContactListener(cb *ContactListenerCallbacks) *ContactListener {
	contactProcsOnce.Do(func() { C.joltcInstallContactProcs() })
	h := cgo.NewHandle(cb)
	return &ContactListener{ptr: C.JPH_ContactListener_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the listener.
func (l *ContactListener) Destroy() {
	if l.ptr != nil {
		C.JPH_ContactListener_Destroy(l.ptr)
		l.ptr = nil
		l.handle.Delete()
	}
}

//export goJoltcContactValidate
func goJoltcContactValidate(userData unsafe.Pointer, b1, b2 *C.JPH_Body, baseOffset *C.JPH_RVec3, result *C.JPH_CollideShapeResult) C.uint32_t {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ContactListenerCallbacks)
	if !ok || cb.OnValidate == nil {
		return C.uint32_t(C.JPH_ValidateResult_AcceptAllContactsForThisBodyPair)
	}
	v := cb.OnValidate(&Body{ptr: b1}, &Body{ptr: b2}, rvec3FromC(*baseOffset), collideShapeResultFromC(result))
	return C.uint32_t(v)
}

//export goJoltcContactAdded
func goJoltcContactAdded(userData unsafe.Pointer, b1, b2 *C.JPH_Body, manifold *C.JPH_ContactManifold, settings *C.JPH_ContactSettings) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ContactListenerCallbacks)
	if !ok || cb.OnAdded == nil {
		return
	}
	cb.OnAdded(&Body{ptr: b1}, &Body{ptr: b2}, &ContactManifold{ptr: manifold}, &ContactSettings{ptr: settings})
}

//export goJoltcContactPersisted
func goJoltcContactPersisted(userData unsafe.Pointer, b1, b2 *C.JPH_Body, manifold *C.JPH_ContactManifold, settings *C.JPH_ContactSettings) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ContactListenerCallbacks)
	if !ok || cb.OnPersisted == nil {
		return
	}
	cb.OnPersisted(&Body{ptr: b1}, &Body{ptr: b2}, &ContactManifold{ptr: manifold}, &ContactSettings{ptr: settings})
}

//export goJoltcContactRemoved
func goJoltcContactRemoved(userData unsafe.Pointer, pair *C.JPH_SubShapeIDPair) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*ContactListenerCallbacks)
	if !ok || cb.OnRemoved == nil {
		return
	}
	cb.OnRemoved(subShapeIDPairFromC(pair))
}

// --- Body activation ---------------------------------------------------------

// BodyActivationListener receives body activation events.
type BodyActivationListener struct {
	ptr    *C.JPH_BodyActivationListener
	handle cgo.Handle
}

var activationProcsOnce sync.Once

// NewBodyActivationListener creates a body activation listener.
func NewBodyActivationListener(cb *BodyActivationCallbacks) *BodyActivationListener {
	activationProcsOnce.Do(func() { C.joltcInstallBodyActivationProcs() })
	h := cgo.NewHandle(cb)
	return &BodyActivationListener{ptr: C.JPH_BodyActivationListener_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the listener.
func (l *BodyActivationListener) Destroy() {
	if l.ptr != nil {
		C.JPH_BodyActivationListener_Destroy(l.ptr)
		l.ptr = nil
		l.handle.Delete()
	}
}

//export goJoltcBodyActivated
func goJoltcBodyActivated(userData unsafe.Pointer, id C.JPH_BodyID, bodyUserData C.uint64_t) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*BodyActivationCallbacks)
	if !ok || cb.OnActivated == nil {
		return
	}
	cb.OnActivated(BodyID(id), uint64(bodyUserData))
}

//export goJoltcBodyDeactivated
func goJoltcBodyDeactivated(userData unsafe.Pointer, id C.JPH_BodyID, bodyUserData C.uint64_t) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*BodyActivationCallbacks)
	if !ok || cb.OnDeactivated == nil {
		return
	}
	cb.OnDeactivated(BodyID(id), uint64(bodyUserData))
}

// --- Step listener -----------------------------------------------------------

// StepListener receives a callback for each simulation step.
type StepListener struct {
	ptr    *C.JPH_PhysicsStepListener
	handle cgo.Handle
	// owned is false for listeners owned by another object (e.g. a vehicle),
	// which must not be destroyed here.
	owned bool
}

var stepProcsOnce sync.Once

// NewStepListener creates a step listener.
func NewStepListener(cb *StepListenerCallbacks) *StepListener {
	stepProcsOnce.Do(func() { C.joltcInstallStepProcs() })
	h := cgo.NewHandle(cb)
	return &StepListener{ptr: C.JPH_PhysicsStepListener_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h, owned: true}
}

// Destroy releases the listener, unless it is owned by another object.
func (l *StepListener) Destroy() {
	if l.ptr == nil {
		return
	}
	if l.owned {
		C.JPH_PhysicsStepListener_Destroy(l.ptr)
		l.handle.Delete()
	}
	l.ptr = nil
}

//export goJoltcStep
func goJoltcStep(userData unsafe.Pointer, ctx *C.JPH_PhysicsStepListenerContext) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*StepListenerCallbacks)
	if !ok || cb.OnStep == nil {
		return
	}
	cb.OnStep(StepContext{
		DeltaTime:   float32(ctx.deltaTime),
		IsFirstStep: ctx.isFirstStep != 0,
		IsLastStep:  ctx.isLastStep != 0,
	})
}

// --- Debug renderer ----------------------------------------------------------

// DebugRenderer turns Jolt's debug draw calls into Go callbacks.
type DebugRenderer struct {
	ptr    *C.JPH_DebugRenderer
	handle cgo.Handle
}

var debugProcsOnce sync.Once

// NewDebugRenderer creates a debug renderer with the given callbacks.
func NewDebugRenderer(cb *DebugRendererCallbacks) *DebugRenderer {
	debugProcsOnce.Do(func() { C.joltcInstallDebugRendererProcs() })
	h := cgo.NewHandle(cb)
	return &DebugRenderer{ptr: C.JPH_DebugRenderer_Create(C.joltc_voidp(C.uintptr_t(h))), handle: h}
}

// Destroy releases the renderer.
func (d *DebugRenderer) Destroy() {
	if d.ptr != nil {
		C.JPH_DebugRenderer_Destroy(d.ptr)
		d.ptr = nil
		d.handle.Delete()
	}
}

// NextFrame resets the renderer for a new frame.
func (d *DebugRenderer) NextFrame() { C.JPH_DebugRenderer_NextFrame(d.ptr) }

// SetCameraPos informs the renderer of the camera position (used for LOD).
func (d *DebugRenderer) SetCameraPos(p vecmath.RVec3) {
	c := rvec3ToC(p)
	C.JPH_DebugRenderer_SetCameraPos(d.ptr, &c)
}

//export goJoltcDrawLine
func goJoltcDrawLine(userData unsafe.Pointer, from, to *C.JPH_RVec3, color C.JPH_Color) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*DebugRendererCallbacks)
	if !ok || cb.DrawLine == nil {
		return
	}
	cb.DrawLine(rvec3FromC(*from), rvec3FromC(*to), vecmath.Color(color))
}

//export goJoltcDrawTriangle
func goJoltcDrawTriangle(userData unsafe.Pointer, v1, v2, v3 *C.JPH_RVec3, color C.JPH_Color, castShadow C.int) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*DebugRendererCallbacks)
	if !ok || cb.DrawTriangle == nil {
		return
	}
	cb.DrawTriangle(rvec3FromC(*v1), rvec3FromC(*v2), rvec3FromC(*v3), vecmath.Color(color), castShadow != 0)
}

//export goJoltcDrawText3D
func goJoltcDrawText3D(userData unsafe.Pointer, position *C.JPH_RVec3, str *C.char, color C.JPH_Color, height C.float) {
	defer func() { _ = recover() }()
	cb, ok := cgo.Handle(uintptr(userData)).Value().(*DebugRendererCallbacks)
	if !ok || cb.DrawText3D == nil {
		return
	}
	cb.DrawText3D(rvec3FromC(*position), C.GoString(str), vecmath.Color(color), float32(height))
}

// --- PhysicsSystem wiring ----------------------------------------------------

// SetContactListener installs (or clears, with nil) the contact listener.
func (p *PhysicsSystem) SetContactListener(l *ContactListener) {
	if l == nil {
		C.JPH_PhysicsSystem_SetContactListener(p.ptr, nil)
		return
	}
	C.JPH_PhysicsSystem_SetContactListener(p.ptr, l.ptr)
}

// SetBodyActivationListener installs (or clears, with nil) the activation listener.
func (p *PhysicsSystem) SetBodyActivationListener(l *BodyActivationListener) {
	if l == nil {
		C.JPH_PhysicsSystem_SetBodyActivationListener(p.ptr, nil)
		return
	}
	C.JPH_PhysicsSystem_SetBodyActivationListener(p.ptr, l.ptr)
}

// AddStepListener registers a step listener.
func (p *PhysicsSystem) AddStepListener(l *StepListener) {
	C.JPH_PhysicsSystem_AddStepListener(p.ptr, l.ptr)
}

// RemoveStepListener unregisters a step listener.
func (p *PhysicsSystem) RemoveStepListener(l *StepListener) {
	C.JPH_PhysicsSystem_RemoveStepListener(p.ptr, l.ptr)
}

func (s DrawSettings) toC() C.JPH_DrawSettings {
	var c C.JPH_DrawSettings
	C.JPH_DrawSettings_InitDefault(&c)
	c.drawShape = C.bool(s.DrawShape)
	c.drawShapeWireframe = C.bool(s.DrawShapeWireframe)
	c.drawBoundingBox = C.bool(s.DrawBoundingBox)
	c.drawCenterOfMassTransform = C.bool(s.DrawCenterOfMassTransform)
	c.drawVelocity = C.bool(s.DrawVelocity)
	c.drawMassAndInertia = C.bool(s.DrawMassAndInertia)
	c.drawGetSupportingFace = C.bool(s.DrawSupportingFace)
	c.drawSoftBodyVertices = C.bool(s.DrawSoftBodyVertices)
	return c
}

// DrawBodies draws all bodies through the debug renderer.
func (p *PhysicsSystem) DrawBodies(settings DrawSettings, renderer *DebugRenderer) {
	c := settings.toC()
	C.JPH_PhysicsSystem_DrawBodies(p.ptr, &c, renderer.ptr, nil)
}
