// Batch helpers that keep per-body work inside C, so hot loops make one cgo
// crossing for the whole sweep instead of one per body. C only: no Go pointers
// are stored or dereferenced on the Go side of these calls.
#include <joltc.h>

void joltc_body_positions(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_RVec3* out)
{
	for (uint32_t i = 0; i < count; ++i)
		JPH_BodyInterface_GetPosition(bi, ids[i], &out[i]);
}

void joltc_body_rotations(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_Quat* out)
{
	for (uint32_t i = 0; i < count; ++i)
		JPH_BodyInterface_GetRotation(bi, ids[i], &out[i]);
}

void joltc_body_linear_velocities(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_Vec3* out)
{
	for (uint32_t i = 0; i < count; ++i)
		JPH_BodyInterface_GetLinearVelocity(bi, ids[i], &out[i]);
}

void joltc_body_position_and_rotation(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_RVec3* outPos, JPH_Quat* outRot)
{
	for (uint32_t i = 0; i < count; ++i)
		JPH_BodyInterface_GetPositionAndRotation(bi, ids[i], &outPos[i], &outRot[i]);
}

void joltc_body_set_positions(JPH_BodyInterface* bi, const JPH_BodyID* ids, uint32_t count, JPH_RVec3* in, JPH_Activation activation)
{
	for (uint32_t i = 0; i < count; ++i)
		JPH_BodyInterface_SetPosition(bi, ids[i], &in[i], activation);
}

/* --- callbacks: handle helpers and procs installation --------------------- */

/* Reinterpret a cgo.Handle (uintptr) as the void* userData C expects. */
void* joltc_voidp(uintptr_t v)
{
	return (void*)v;
}

extern uint32_t goJoltcContactValidate(void* ud, const JPH_Body* b1, const JPH_Body* b2, const JPH_RVec3* baseOffset, const JPH_CollideShapeResult* collisionResult);
extern void goJoltcContactAdded(void* ud, const JPH_Body* b1, const JPH_Body* b2, const JPH_ContactManifold* manifold, JPH_ContactSettings* settings);
extern void goJoltcContactPersisted(void* ud, const JPH_Body* b1, const JPH_Body* b2, const JPH_ContactManifold* manifold, JPH_ContactSettings* settings);
extern void goJoltcContactRemoved(void* ud, const JPH_SubShapeIDPair* pair);

static JPH_ValidateResult joltcContactValidateTramp(void* ud, const JPH_Body* b1, const JPH_Body* b2, const JPH_RVec3* o, const JPH_CollideShapeResult* cr)
{
	return (JPH_ValidateResult)goJoltcContactValidate(ud, b1, b2, o, cr);
}
static void joltcContactAddedTramp(void* ud, const JPH_Body* b1, const JPH_Body* b2, const JPH_ContactManifold* m, JPH_ContactSettings* s) { goJoltcContactAdded(ud, b1, b2, m, s); }
static void joltcContactPersistedTramp(void* ud, const JPH_Body* b1, const JPH_Body* b2, const JPH_ContactManifold* m, JPH_ContactSettings* s) { goJoltcContactPersisted(ud, b1, b2, m, s); }
static void joltcContactRemovedTramp(void* ud, const JPH_SubShapeIDPair* p) { goJoltcContactRemoved(ud, p); }

/* joltc stores the procs pointer, so these must have static storage duration. */
static JPH_ContactListener_Procs g_contactProcs;

void joltcInstallContactProcs(void)
{
	g_contactProcs.OnContactValidate = joltcContactValidateTramp;
	g_contactProcs.OnContactAdded = joltcContactAddedTramp;
	g_contactProcs.OnContactPersisted = joltcContactPersistedTramp;
	g_contactProcs.OnContactRemoved = joltcContactRemovedTramp;
	JPH_ContactListener_SetProcs(&g_contactProcs);
}

extern void goJoltcBodyActivated(void* ud, JPH_BodyID id, uint64_t userData);
extern void goJoltcBodyDeactivated(void* ud, JPH_BodyID id, uint64_t userData);

static void joltcBodyActivatedTramp(void* ud, JPH_BodyID id, uint64_t u) { goJoltcBodyActivated(ud, id, u); }
static void joltcBodyDeactivatedTramp(void* ud, JPH_BodyID id, uint64_t u) { goJoltcBodyDeactivated(ud, id, u); }

static JPH_BodyActivationListener_Procs g_activationProcs;

void joltcInstallBodyActivationProcs(void)
{
	g_activationProcs.OnBodyActivated = joltcBodyActivatedTramp;
	g_activationProcs.OnBodyDeactivated = joltcBodyDeactivatedTramp;
	JPH_BodyActivationListener_SetProcs(&g_activationProcs);
}

extern void goJoltcStep(void* ud, const JPH_PhysicsStepListenerContext* ctx);

static void joltcStepTramp(void* ud, const JPH_PhysicsStepListenerContext* ctx) { goJoltcStep(ud, ctx); }

static JPH_PhysicsStepListener_Procs g_stepProcs;

void joltcInstallStepProcs(void)
{
	g_stepProcs.OnStep = joltcStepTramp;
	JPH_PhysicsStepListener_SetProcs(&g_stepProcs);
}

extern void goJoltcDrawLine(void* ud, const JPH_RVec3* from, const JPH_RVec3* to, JPH_Color color);
extern void goJoltcDrawTriangle(void* ud, const JPH_RVec3* v1, const JPH_RVec3* v2, const JPH_RVec3* v3, JPH_Color color, int castShadow);
extern void goJoltcDrawText3D(void* ud, const JPH_RVec3* position, const char* str, JPH_Color color, float height);

static void joltcDrawLineTramp(void* ud, const JPH_RVec3* f, const JPH_RVec3* t, JPH_Color c) { goJoltcDrawLine(ud, f, t, c); }
static void joltcDrawTriangleTramp(void* ud, const JPH_RVec3* a, const JPH_RVec3* b, const JPH_RVec3* c, JPH_Color col, JPH_DebugRenderer_CastShadow cs) { goJoltcDrawTriangle(ud, a, b, c, col, (int)cs); }
static void joltcDrawText3DTramp(void* ud, const JPH_RVec3* p, const char* s, JPH_Color c, float h) { goJoltcDrawText3D(ud, p, s, c, h); }

static JPH_DebugRenderer_Procs g_debugProcs;

void joltcInstallDebugRendererProcs(void)
{
	g_debugProcs.DrawLine = joltcDrawLineTramp;
	g_debugProcs.DrawTriangle = joltcDrawTriangleTramp;
	g_debugProcs.DrawText3D = joltcDrawText3DTramp;
	JPH_DebugRenderer_SetProcs(&g_debugProcs);
}
