//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

// PhysicsSettings returns the system's solver tuning parameters.
func (p *PhysicsSystem) PhysicsSettings() PhysicsSettings {
	var c C.JPH_PhysicsSettings
	C.JPH_PhysicsSystem_GetPhysicsSettings(p.ptr, &c)
	return PhysicsSettings{
		MaxInFlightBodyPairs:                 int32(c.maxInFlightBodyPairs),
		StepListenersBatchSize:               int32(c.stepListenersBatchSize),
		StepListenerBatchesPerJob:            int32(c.stepListenerBatchesPerJob),
		Baumgarte:                            float32(c.baumgarte),
		SpeculativeContactDistance:           float32(c.speculativeContactDistance),
		PenetrationSlop:                      float32(c.penetrationSlop),
		LinearCastThreshold:                  float32(c.linearCastThreshold),
		LinearCastMaxPenetration:             float32(c.linearCastMaxPenetration),
		ManifoldTolerance:                    float32(c.manifoldTolerance),
		MaxPenetrationDistance:               float32(c.maxPenetrationDistance),
		BodyPairCacheMaxDeltaPositionSq:      float32(c.bodyPairCacheMaxDeltaPositionSq),
		BodyPairCacheCosMaxDeltaRotationDiv2: float32(c.bodyPairCacheCosMaxDeltaRotationDiv2),
		ContactNormalCosMaxDeltaRotation:     float32(c.contactNormalCosMaxDeltaRotation),
		ContactPointPreserveLambdaMaxDistSq:  float32(c.contactPointPreserveLambdaMaxDistSq),
		NumVelocitySteps:                     uint32(c.numVelocitySteps),
		NumPositionSteps:                     uint32(c.numPositionSteps),
		MinVelocityForRestitution:            float32(c.minVelocityForRestitution),
		TimeBeforeSleep:                      float32(c.timeBeforeSleep),
		PointVelocitySleepThreshold:          float32(c.pointVelocitySleepThreshold),
		DeterministicSimulation:              bool(c.deterministicSimulation),
		ConstraintWarmStart:                  bool(c.constraintWarmStart),
		UseBodyPairContactCache:              bool(c.useBodyPairContactCache),
		UseManifoldReduction:                 bool(c.useManifoldReduction),
		UseLargeIslandSplitter:               bool(c.useLargeIslandSplitter),
		AllowSleeping:                        bool(c.allowSleeping),
		CheckActiveEdges:                     bool(c.checkActiveEdges),
	}
}

// SetPhysicsSettings applies solver tuning parameters to the system.
func (p *PhysicsSystem) SetPhysicsSettings(s PhysicsSettings) {
	c := C.JPH_PhysicsSettings{
		maxInFlightBodyPairs:                 C.int(s.MaxInFlightBodyPairs),
		stepListenersBatchSize:               C.int(s.StepListenersBatchSize),
		stepListenerBatchesPerJob:            C.int(s.StepListenerBatchesPerJob),
		baumgarte:                            C.float(s.Baumgarte),
		speculativeContactDistance:           C.float(s.SpeculativeContactDistance),
		penetrationSlop:                      C.float(s.PenetrationSlop),
		linearCastThreshold:                  C.float(s.LinearCastThreshold),
		linearCastMaxPenetration:             C.float(s.LinearCastMaxPenetration),
		manifoldTolerance:                    C.float(s.ManifoldTolerance),
		maxPenetrationDistance:               C.float(s.MaxPenetrationDistance),
		bodyPairCacheMaxDeltaPositionSq:      C.float(s.BodyPairCacheMaxDeltaPositionSq),
		bodyPairCacheCosMaxDeltaRotationDiv2: C.float(s.BodyPairCacheCosMaxDeltaRotationDiv2),
		contactNormalCosMaxDeltaRotation:     C.float(s.ContactNormalCosMaxDeltaRotation),
		contactPointPreserveLambdaMaxDistSq:  C.float(s.ContactPointPreserveLambdaMaxDistSq),
		numVelocitySteps:                     C.uint32_t(s.NumVelocitySteps),
		numPositionSteps:                     C.uint32_t(s.NumPositionSteps),
		minVelocityForRestitution:            C.float(s.MinVelocityForRestitution),
		timeBeforeSleep:                      C.float(s.TimeBeforeSleep),
		pointVelocitySleepThreshold:          C.float(s.PointVelocitySleepThreshold),
		deterministicSimulation:              C.bool(s.DeterministicSimulation),
		constraintWarmStart:                  C.bool(s.ConstraintWarmStart),
		useBodyPairContactCache:              C.bool(s.UseBodyPairContactCache),
		useManifoldReduction:                 C.bool(s.UseManifoldReduction),
		useLargeIslandSplitter:               C.bool(s.UseLargeIslandSplitter),
		allowSleeping:                        C.bool(s.AllowSleeping),
		checkActiveEdges:                     C.bool(s.CheckActiveEdges),
	}
	C.JPH_PhysicsSystem_SetPhysicsSettings(p.ptr, &c)
}
