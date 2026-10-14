package joltc

// PhysicsSettings holds the solver tuning parameters of a physics system
// (JPH_PhysicsSettings).
type PhysicsSettings struct {
	MaxInFlightBodyPairs                 int32
	StepListenersBatchSize               int32
	StepListenerBatchesPerJob            int32
	Baumgarte                            float32
	SpeculativeContactDistance           float32
	PenetrationSlop                      float32
	LinearCastThreshold                  float32
	LinearCastMaxPenetration             float32
	ManifoldTolerance                    float32
	MaxPenetrationDistance               float32
	BodyPairCacheMaxDeltaPositionSq      float32
	BodyPairCacheCosMaxDeltaRotationDiv2 float32
	ContactNormalCosMaxDeltaRotation     float32
	ContactPointPreserveLambdaMaxDistSq  float32
	NumVelocitySteps                     uint32
	NumPositionSteps                     uint32
	MinVelocityForRestitution            float32
	TimeBeforeSleep                      float32
	PointVelocitySleepThreshold          float32
	DeterministicSimulation              bool
	ConstraintWarmStart                  bool
	UseBodyPairContactCache              bool
	UseManifoldReduction                 bool
	UseLargeIslandSplitter               bool
	AllowSleeping                        bool
	CheckActiveEdges                     bool
}

// CollisionGroup is a collision group identifier plus an optional filter that
// decides whether two groups may collide (JPH_CollisionGroup).
type CollisionGroup struct {
	Filter     *GroupFilter
	GroupID    CollisionGroupID
	SubGroupID CollisionSubGroupID
}
