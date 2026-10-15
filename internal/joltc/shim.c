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

/* --- query filters -------------------------------------------------------- */

extern bool goJoltcBroadPhaseLayerFilterShouldCollide(void* ud, JPH_BroadPhaseLayer layer);
extern bool goJoltcObjectLayerFilterShouldCollide(void* ud, JPH_ObjectLayer layer);
extern bool goJoltcBodyFilterShouldCollide(void* ud, JPH_BodyID bodyID);
extern bool goJoltcBodyFilterShouldCollideLocked(void* ud, const JPH_Body* body);
extern bool goJoltcShapeFilterShouldCollide(void* ud, const JPH_Shape* shape2, const JPH_SubShapeID* subShapeID2);
extern bool goJoltcShapeFilterShouldCollide2(void* ud, const JPH_Shape* shape1, const JPH_SubShapeID* subShapeID1, const JPH_Shape* shape2, const JPH_SubShapeID* subShapeID2);

static bool joltcBPLayerTramp(void* ud, JPH_BroadPhaseLayer l) { return goJoltcBroadPhaseLayerFilterShouldCollide(ud, l); }
static bool joltcObjLayerTramp(void* ud, JPH_ObjectLayer l) { return goJoltcObjectLayerFilterShouldCollide(ud, l); }
static bool joltcBodyFilterTramp(void* ud, JPH_BodyID id) { return goJoltcBodyFilterShouldCollide(ud, id); }
static bool joltcBodyFilterLockedTramp(void* ud, const JPH_Body* b) { return goJoltcBodyFilterShouldCollideLocked(ud, b); }
static bool joltcShapeFilterTramp(void* ud, const JPH_Shape* s2, const JPH_SubShapeID* id2) { return goJoltcShapeFilterShouldCollide(ud, s2, id2); }
static bool joltcShapeFilterTramp2(void* ud, const JPH_Shape* s1, const JPH_SubShapeID* id1, const JPH_Shape* s2, const JPH_SubShapeID* id2) { return goJoltcShapeFilterShouldCollide2(ud, s1, id1, s2, id2); }

/* joltc stores the procs pointer, so these need static storage duration. */
static JPH_BroadPhaseLayerFilter_Procs g_bpLayerFilterProcs;
static JPH_ObjectLayerFilter_Procs g_objLayerFilterProcs;
static JPH_BodyFilter_Procs g_bodyFilterProcs;
static JPH_ShapeFilter_Procs g_shapeFilterProcs;

void joltcInstallBroadPhaseLayerFilterProcs(void)
{
	g_bpLayerFilterProcs.ShouldCollide = joltcBPLayerTramp;
	JPH_BroadPhaseLayerFilter_SetProcs(&g_bpLayerFilterProcs);
}

void joltcInstallObjectLayerFilterProcs(void)
{
	g_objLayerFilterProcs.ShouldCollide = joltcObjLayerTramp;
	JPH_ObjectLayerFilter_SetProcs(&g_objLayerFilterProcs);
}

void joltcInstallBodyFilterProcs(void)
{
	g_bodyFilterProcs.ShouldCollide = joltcBodyFilterTramp;
	g_bodyFilterProcs.ShouldCollideLocked = joltcBodyFilterLockedTramp;
	JPH_BodyFilter_SetProcs(&g_bodyFilterProcs);
}

void joltcInstallShapeFilterProcs(void)
{
	g_shapeFilterProcs.ShouldCollide = joltcShapeFilterTramp;
	g_shapeFilterProcs.ShouldCollide2 = joltcShapeFilterTramp2;
	JPH_ShapeFilter_SetProcs(&g_shapeFilterProcs);
}

/* --- query collectors ----------------------------------------------------- */

extern float goJoltcCastRayCollector(void* ud, const JPH_RayCastResult* result);
extern float goJoltcCollidePointCollector(void* ud, const JPH_CollidePointResult* result);
extern float goJoltcCollideShapeCollector(void* ud, const JPH_CollideShapeResult* result);
extern float goJoltcCastShapeCollector(void* ud, const JPH_ShapeCastResult* result);

static float joltcCastRayCollTramp(void* ud, const JPH_RayCastResult* r) { return goJoltcCastRayCollector(ud, r); }
static float joltcCollidePointCollTramp(void* ud, const JPH_CollidePointResult* r) { return goJoltcCollidePointCollector(ud, r); }
static float joltcCollideShapeCollTramp(void* ud, const JPH_CollideShapeResult* r) { return goJoltcCollideShapeCollector(ud, r); }
static float joltcCastShapeCollTramp(void* ud, const JPH_ShapeCastResult* r) { return goJoltcCastShapeCollector(ud, r); }

bool joltcNarrowPhaseCastRay(const JPH_NarrowPhaseQuery* q, const JPH_RVec3* origin, const JPH_Vec3* dir, const JPH_RayCastSettings* s, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf)
{
	return JPH_NarrowPhaseQuery_CastRay2(q, origin, dir, s, joltcCastRayCollTramp, (void*)ud, bp, ol, bf, sf);
}

/* Sweep many rays with Jolt's native closest-hit cast (no collector callbacks, no
   per-ray cgo crossing). Misses are marked with Jolt's invalid body id. */
uint32_t joltcNarrowPhaseCastRays(const JPH_NarrowPhaseQuery* q, const JPH_RVec3* origins, const JPH_Vec3* dirs, uint32_t count,
	JPH_RayCastResult* out, const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf)
{
	uint32_t hits = 0;
	for (uint32_t i = 0; i < count; i++)
	{
		out[i].bodyID = (JPH_BodyID)0xffffffffu;
		JPH_RayCastResult r;
		if (JPH_NarrowPhaseQuery_CastRay(q, &origins[i], &dirs[i], &r, bp, ol, bf))
		{
			out[i] = r;
			hits++;
		}
	}
	return hits;
}

bool joltcNarrowPhaseCollidePoint(const JPH_NarrowPhaseQuery* q, const JPH_RVec3* point, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf)
{
	return JPH_NarrowPhaseQuery_CollidePoint(q, point, joltcCollidePointCollTramp, (void*)ud, bp, ol, bf, sf);
}

bool joltcNarrowPhaseCollideShape(const JPH_NarrowPhaseQuery* q, const JPH_Shape* shape, const JPH_Vec3* scale, const JPH_RMat4* comTransform,
	const JPH_CollideShapeSettings* s, JPH_RVec3* baseOffset, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf)
{
	return JPH_NarrowPhaseQuery_CollideShape(q, shape, scale, comTransform, s, baseOffset, joltcCollideShapeCollTramp, (void*)ud, bp, ol, bf, sf);
}

bool joltcNarrowPhaseCastShape(const JPH_NarrowPhaseQuery* q, const JPH_Shape* shape, const JPH_RMat4* worldTransform, const JPH_Vec3* dir,
	const JPH_ShapeCastSettings* s, JPH_RVec3* baseOffset, uintptr_t ud,
	const JPH_BroadPhaseLayerFilter* bp, const JPH_ObjectLayerFilter* ol, const JPH_BodyFilter* bf, const JPH_ShapeFilter* sf)
{
	return JPH_NarrowPhaseQuery_CastShape(q, shape, worldTransform, dir, s, baseOffset, joltcCastShapeCollTramp, (void*)ud, bp, ol, bf, sf);
}
