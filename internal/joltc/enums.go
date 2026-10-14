package joltc

// Enum types mirror the JPH_* C enums in third_party/joltc/include/joltc.h.
// Their values are verified against the C header by TestEnumValuesMatchC, so a
// mismatch fails the test suite rather than silently mapping to the wrong value.

// PhysicsUpdateError reports non-fatal issues encountered during PhysicsSystem.Update.
type PhysicsUpdateError int32

const (
	PhysicsUpdateErrorNone                   PhysicsUpdateError = 0
	PhysicsUpdateErrorManifoldCacheFull      PhysicsUpdateError = 1 << 0
	PhysicsUpdateErrorBodyPairCacheFull      PhysicsUpdateError = 1 << 1
	PhysicsUpdateErrorContactConstraintsFull PhysicsUpdateError = 1 << 2
)

// BodyType distinguishes rigid bodies from soft bodies.
type BodyType int32

const (
	BodyTypeRigid BodyType = 0
	BodyTypeSoft  BodyType = 1
)

// MotionType controls how a body is simulated.
type MotionType int32

const (
	MotionTypeStatic    MotionType = 0
	MotionTypeKinematic MotionType = 1
	MotionTypeDynamic   MotionType = 2
)

// Activation controls whether a body starts active.
type Activation int32

const (
	ActivationActivate     Activation = 0
	ActivationDontActivate Activation = 1
)

// ValidateResult is returned by contact validation callbacks.
type ValidateResult int32

const (
	ValidateResultAcceptAllContactsForThisBodyPair ValidateResult = 0
	ValidateResultAcceptContact                    ValidateResult = 1
	ValidateResultRejectContact                    ValidateResult = 2
	ValidateResultRejectAllContactsForThisBodyPair ValidateResult = 3
)

// ShapeType is the top-level shape category.
type ShapeType int32

const (
	ShapeTypeConvex      ShapeType = 0
	ShapeTypeCompound    ShapeType = 1
	ShapeTypeDecorated   ShapeType = 2
	ShapeTypeMesh        ShapeType = 3
	ShapeTypeHeightField ShapeType = 4
	ShapeTypeSoftBody    ShapeType = 5
	ShapeTypeUser1       ShapeType = 6
	ShapeTypeUser2       ShapeType = 7
	ShapeTypeUser3       ShapeType = 8
	ShapeTypeUser4       ShapeType = 9
)

// ShapeSubType is the specific shape category.
type ShapeSubType int32

const (
	ShapeSubTypeSphere             ShapeSubType = 0
	ShapeSubTypeBox                ShapeSubType = 1
	ShapeSubTypeTriangle           ShapeSubType = 2
	ShapeSubTypeCapsule            ShapeSubType = 3
	ShapeSubTypeTaperedCapsule     ShapeSubType = 4
	ShapeSubTypeCylinder           ShapeSubType = 5
	ShapeSubTypeConvexHull         ShapeSubType = 6
	ShapeSubTypeStaticCompound     ShapeSubType = 7
	ShapeSubTypeMutableCompound    ShapeSubType = 8
	ShapeSubTypeRotatedTranslated  ShapeSubType = 9
	ShapeSubTypeScaled             ShapeSubType = 10
	ShapeSubTypeOffsetCenterOfMass ShapeSubType = 11
	ShapeSubTypeMesh               ShapeSubType = 12
	ShapeSubTypeHeightField        ShapeSubType = 13
	ShapeSubTypeSoftBody           ShapeSubType = 14
)

// MotionQuality selects the collision detection quality for a body.
type MotionQuality int32

const (
	MotionQualityDiscrete   MotionQuality = 0
	MotionQualityLinearCast MotionQuality = 1
)

// OverrideMassProperties controls how mass and inertia are computed.
type OverrideMassProperties int32

const (
	OverrideMassPropertiesCalculateMassAndInertia OverrideMassProperties = 0
	OverrideMassPropertiesCalculateInertia        OverrideMassProperties = 1
	OverrideMassPropertiesMassAndInertiaProvided  OverrideMassProperties = 2
)

// AllowedDOFs is a bitmask of the degrees of freedom a body may move in.
type AllowedDOFs int32

const (
	AllowedDOFsAll          AllowedDOFs = 0b111111
	AllowedDOFsTranslationX AllowedDOFs = 0b000001
	AllowedDOFsTranslationY AllowedDOFs = 0b000010
	AllowedDOFsTranslationZ AllowedDOFs = 0b000100
	AllowedDOFsRotationX    AllowedDOFs = 0b001000
	AllowedDOFsRotationY    AllowedDOFs = 0b010000
	AllowedDOFsRotationZ    AllowedDOFs = 0b100000
	AllowedDOFsPlane2D      AllowedDOFs = AllowedDOFsTranslationX | AllowedDOFsTranslationY | AllowedDOFsRotationZ
)

// GroundState describes a character's ground contact state.
type GroundState int32

const (
	GroundStateOnGround      GroundState = 0
	GroundStateOnSteepGround GroundState = 1
	GroundStateNotSupported  GroundState = 2
	GroundStateInAir         GroundState = 3
)

// BackFaceMode selects whether back faces of triangles collide.
type BackFaceMode int32

const (
	BackFaceModeIgnoreBackFaces      BackFaceMode = 0
	BackFaceModeCollideWithBackFaces BackFaceMode = 1
)

// ActiveEdgeMode selects which triangle edges are considered active.
type ActiveEdgeMode int32

const (
	ActiveEdgeModeCollideOnlyWithActive ActiveEdgeMode = 0
	ActiveEdgeModeCollideWithAll        ActiveEdgeMode = 1
)

// CollectFacesMode selects whether supporting faces are collected.
type CollectFacesMode int32

const (
	CollectFacesModeCollectFaces CollectFacesMode = 0
	CollectFacesModeNoFaces      CollectFacesMode = 1
)

// CollisionCollectorType selects how collision query results are collected.
type CollisionCollectorType int32

const (
	CollisionCollectorTypeAllHit       CollisionCollectorType = 0
	CollisionCollectorTypeAllHitSorted CollisionCollectorType = 1
	CollisionCollectorTypeClosestHit   CollisionCollectorType = 2
	CollisionCollectorTypeAnyHit       CollisionCollectorType = 3
)

// ConstraintType is the top-level constraint category.
type ConstraintType int32

const (
	ConstraintTypeConstraint        ConstraintType = 0
	ConstraintTypeTwoBodyConstraint ConstraintType = 1
)

// ConstraintSubType is the specific constraint category.
type ConstraintSubType int32

const (
	ConstraintSubTypeFixed         ConstraintSubType = 0
	ConstraintSubTypePoint         ConstraintSubType = 1
	ConstraintSubTypeHinge         ConstraintSubType = 2
	ConstraintSubTypeSlider        ConstraintSubType = 3
	ConstraintSubTypeDistance      ConstraintSubType = 4
	ConstraintSubTypeCone          ConstraintSubType = 5
	ConstraintSubTypeSwingTwist    ConstraintSubType = 6
	ConstraintSubTypeSixDOF        ConstraintSubType = 7
	ConstraintSubTypePath          ConstraintSubType = 8
	ConstraintSubTypeVehicle       ConstraintSubType = 9
	ConstraintSubTypeRackAndPinion ConstraintSubType = 10
	ConstraintSubTypeGear          ConstraintSubType = 11
	ConstraintSubTypePulley        ConstraintSubType = 12
	ConstraintSubTypeUser1         ConstraintSubType = 13
	ConstraintSubTypeUser2         ConstraintSubType = 14
	ConstraintSubTypeUser3         ConstraintSubType = 15
	ConstraintSubTypeUser4         ConstraintSubType = 16
)

// ConstraintSpace selects the coordinate space of a constraint.
type ConstraintSpace int32

const (
	ConstraintSpaceLocalToBodyCOM ConstraintSpace = 0
	ConstraintSpaceWorldSpace     ConstraintSpace = 1
)

// MotorState describes a constraint motor's mode.
type MotorState int32

const (
	MotorStateOff      MotorState = 0
	MotorStateVelocity MotorState = 1
	MotorStatePosition MotorState = 2
)

// SpringMode selects how a spring is parameterized.
type SpringMode int32

const (
	SpringModeFrequencyAndDamping SpringMode = 0
	SpringModeStiffnessAndDamping SpringMode = 1
)

// SwingType selects how a swing-twist constraint's swing cone is shaped.
type SwingType int32

const (
	SwingTypeCone    SwingType = 0
	SwingTypePyramid SwingType = 1
)

// SixDOFConstraintAxis indexes the six axes of a six DOF constraint.
type SixDOFConstraintAxis int32

const (
	SixDOFConstraintAxisTranslationX SixDOFConstraintAxis = 0
	SixDOFConstraintAxisTranslationY SixDOFConstraintAxis = 1
	SixDOFConstraintAxisTranslationZ SixDOFConstraintAxis = 2
	SixDOFConstraintAxisRotationX    SixDOFConstraintAxis = 3
	SixDOFConstraintAxisRotationY    SixDOFConstraintAxis = 4
	SixDOFConstraintAxisRotationZ    SixDOFConstraintAxis = 5
)

// MeshShapeBuildQuality trades mesh build time against runtime performance.
type MeshShapeBuildQuality int32

const (
	MeshShapeBuildQualityFavorRuntimePerformance MeshShapeBuildQuality = 0
	MeshShapeBuildQualityFavorBuildSpeed         MeshShapeBuildQuality = 1
)

// TransmissionMode selects how a transmission shifts gears.
type TransmissionMode int32

const (
	TransmissionModeAuto   TransmissionMode = 0
	TransmissionModeManual TransmissionMode = 1
)
