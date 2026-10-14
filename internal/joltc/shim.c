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
