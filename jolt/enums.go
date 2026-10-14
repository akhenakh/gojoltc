package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Enum types mirror the joltc C enums. Their values are verified against the C
// header by the bindings' test suite.
type (
	// PhysicsUpdateError reports non-fatal issues from PhysicsSystem.Update.
	PhysicsUpdateError = joltc.PhysicsUpdateError
	// BodyType distinguishes rigid bodies from soft bodies.
	BodyType = joltc.BodyType
	// MotionType controls how a body is simulated.
	MotionType = joltc.MotionType
	// Activation controls whether a body starts active.
	Activation = joltc.Activation
	// ValidateResult is returned by contact validation callbacks.
	ValidateResult = joltc.ValidateResult
	// ShapeType is the top-level shape category.
	ShapeType = joltc.ShapeType
	// ShapeSubType is the specific shape category.
	ShapeSubType = joltc.ShapeSubType
	// MotionQuality selects collision detection quality.
	MotionQuality = joltc.MotionQuality
	// OverrideMassProperties controls mass/inertia computation.
	OverrideMassProperties = joltc.OverrideMassProperties
	// AllowedDOFs is a bitmask of a body's allowed motion.
	AllowedDOFs = joltc.AllowedDOFs
	// GroundState describes a character's ground contact state.
	GroundState = joltc.GroundState
	// BackFaceMode selects whether back faces collide.
	BackFaceMode = joltc.BackFaceMode
	// ActiveEdgeMode selects which triangle edges are active.
	ActiveEdgeMode = joltc.ActiveEdgeMode
	// CollectFacesMode selects whether supporting faces are collected.
	CollectFacesMode = joltc.CollectFacesMode
	// CollisionCollectorType selects how query results are collected.
	CollisionCollectorType = joltc.CollisionCollectorType
	// ConstraintType is the top-level constraint category.
	ConstraintType = joltc.ConstraintType
	// ConstraintSubType is the specific constraint category.
	ConstraintSubType = joltc.ConstraintSubType
	// ConstraintSpace selects a constraint's coordinate space.
	ConstraintSpace = joltc.ConstraintSpace
	// MotorState describes a constraint motor's mode.
	MotorState = joltc.MotorState
	// SpringMode selects how a spring is parameterized.
	SpringMode = joltc.SpringMode
	// SwingType selects a swing-twist constraint's swing shape.
	SwingType = joltc.SwingType
	// SixDOFConstraintAxis indexes a six DOF constraint's axes.
	SixDOFConstraintAxis = joltc.SixDOFConstraintAxis
	// MeshShapeBuildQuality trades mesh build time against runtime performance.
	MeshShapeBuildQuality = joltc.MeshShapeBuildQuality
	// TransmissionMode selects how a transmission shifts gears.
	TransmissionMode = joltc.TransmissionMode
)

// PhysicsUpdateError values.
const (
	PhysicsUpdateErrorNone                   = joltc.PhysicsUpdateErrorNone
	PhysicsUpdateErrorManifoldCacheFull      = joltc.PhysicsUpdateErrorManifoldCacheFull
	PhysicsUpdateErrorBodyPairCacheFull      = joltc.PhysicsUpdateErrorBodyPairCacheFull
	PhysicsUpdateErrorContactConstraintsFull = joltc.PhysicsUpdateErrorContactConstraintsFull
)

// BodyType values.
const (
	BodyTypeRigid = joltc.BodyTypeRigid
	BodyTypeSoft  = joltc.BodyTypeSoft
)

// MotionType values.
const (
	MotionTypeStatic    = joltc.MotionTypeStatic
	MotionTypeKinematic = joltc.MotionTypeKinematic
	MotionTypeDynamic   = joltc.MotionTypeDynamic
)

// Activation values.
const (
	ActivationActivate     = joltc.ActivationActivate
	ActivationDontActivate = joltc.ActivationDontActivate
)

// ValidateResult values.
const (
	ValidateResultAcceptAllContactsForThisBodyPair = joltc.ValidateResultAcceptAllContactsForThisBodyPair
	ValidateResultAcceptContact                    = joltc.ValidateResultAcceptContact
	ValidateResultRejectContact                    = joltc.ValidateResultRejectContact
	ValidateResultRejectAllContactsForThisBodyPair = joltc.ValidateResultRejectAllContactsForThisBodyPair
)

// ShapeType values.
const (
	ShapeTypeConvex      = joltc.ShapeTypeConvex
	ShapeTypeCompound    = joltc.ShapeTypeCompound
	ShapeTypeDecorated   = joltc.ShapeTypeDecorated
	ShapeTypeMesh        = joltc.ShapeTypeMesh
	ShapeTypeHeightField = joltc.ShapeTypeHeightField
	ShapeTypeSoftBody    = joltc.ShapeTypeSoftBody
	ShapeTypeUser1       = joltc.ShapeTypeUser1
	ShapeTypeUser2       = joltc.ShapeTypeUser2
	ShapeTypeUser3       = joltc.ShapeTypeUser3
	ShapeTypeUser4       = joltc.ShapeTypeUser4
)

// ShapeSubType values.
const (
	ShapeSubTypeSphere             = joltc.ShapeSubTypeSphere
	ShapeSubTypeBox                = joltc.ShapeSubTypeBox
	ShapeSubTypeTriangle           = joltc.ShapeSubTypeTriangle
	ShapeSubTypeCapsule            = joltc.ShapeSubTypeCapsule
	ShapeSubTypeTaperedCapsule     = joltc.ShapeSubTypeTaperedCapsule
	ShapeSubTypeCylinder           = joltc.ShapeSubTypeCylinder
	ShapeSubTypeConvexHull         = joltc.ShapeSubTypeConvexHull
	ShapeSubTypeStaticCompound     = joltc.ShapeSubTypeStaticCompound
	ShapeSubTypeMutableCompound    = joltc.ShapeSubTypeMutableCompound
	ShapeSubTypeRotatedTranslated  = joltc.ShapeSubTypeRotatedTranslated
	ShapeSubTypeScaled             = joltc.ShapeSubTypeScaled
	ShapeSubTypeOffsetCenterOfMass = joltc.ShapeSubTypeOffsetCenterOfMass
	ShapeSubTypeMesh               = joltc.ShapeSubTypeMesh
	ShapeSubTypeHeightField        = joltc.ShapeSubTypeHeightField
	ShapeSubTypeSoftBody           = joltc.ShapeSubTypeSoftBody
)

// MotionQuality values.
const (
	MotionQualityDiscrete   = joltc.MotionQualityDiscrete
	MotionQualityLinearCast = joltc.MotionQualityLinearCast
)

// OverrideMassProperties values.
const (
	OverrideMassPropertiesCalculateMassAndInertia = joltc.OverrideMassPropertiesCalculateMassAndInertia
	OverrideMassPropertiesCalculateInertia        = joltc.OverrideMassPropertiesCalculateInertia
	OverrideMassPropertiesMassAndInertiaProvided  = joltc.OverrideMassPropertiesMassAndInertiaProvided
)

// AllowedDOFs values.
const (
	AllowedDOFsAll          = joltc.AllowedDOFsAll
	AllowedDOFsTranslationX = joltc.AllowedDOFsTranslationX
	AllowedDOFsTranslationY = joltc.AllowedDOFsTranslationY
	AllowedDOFsTranslationZ = joltc.AllowedDOFsTranslationZ
	AllowedDOFsRotationX    = joltc.AllowedDOFsRotationX
	AllowedDOFsRotationY    = joltc.AllowedDOFsRotationY
	AllowedDOFsRotationZ    = joltc.AllowedDOFsRotationZ
	AllowedDOFsPlane2D      = joltc.AllowedDOFsPlane2D
)

// GroundState values.
const (
	GroundStateOnGround      = joltc.GroundStateOnGround
	GroundStateOnSteepGround = joltc.GroundStateOnSteepGround
	GroundStateNotSupported  = joltc.GroundStateNotSupported
	GroundStateInAir         = joltc.GroundStateInAir
)

// BackFaceMode values.
const (
	BackFaceModeIgnoreBackFaces      = joltc.BackFaceModeIgnoreBackFaces
	BackFaceModeCollideWithBackFaces = joltc.BackFaceModeCollideWithBackFaces
)

// ActiveEdgeMode values.
const (
	ActiveEdgeModeCollideOnlyWithActive = joltc.ActiveEdgeModeCollideOnlyWithActive
	ActiveEdgeModeCollideWithAll        = joltc.ActiveEdgeModeCollideWithAll
)

// CollectFacesMode values.
const (
	CollectFacesModeCollectFaces = joltc.CollectFacesModeCollectFaces
	CollectFacesModeNoFaces      = joltc.CollectFacesModeNoFaces
)

// CollisionCollectorType values.
const (
	CollisionCollectorTypeAllHit       = joltc.CollisionCollectorTypeAllHit
	CollisionCollectorTypeAllHitSorted = joltc.CollisionCollectorTypeAllHitSorted
	CollisionCollectorTypeClosestHit   = joltc.CollisionCollectorTypeClosestHit
	CollisionCollectorTypeAnyHit       = joltc.CollisionCollectorTypeAnyHit
)

// ConstraintType values.
const (
	ConstraintTypeConstraint        = joltc.ConstraintTypeConstraint
	ConstraintTypeTwoBodyConstraint = joltc.ConstraintTypeTwoBodyConstraint
)

// ConstraintSubType values.
const (
	ConstraintSubTypeFixed         = joltc.ConstraintSubTypeFixed
	ConstraintSubTypePoint         = joltc.ConstraintSubTypePoint
	ConstraintSubTypeHinge         = joltc.ConstraintSubTypeHinge
	ConstraintSubTypeSlider        = joltc.ConstraintSubTypeSlider
	ConstraintSubTypeDistance      = joltc.ConstraintSubTypeDistance
	ConstraintSubTypeCone          = joltc.ConstraintSubTypeCone
	ConstraintSubTypeSwingTwist    = joltc.ConstraintSubTypeSwingTwist
	ConstraintSubTypeSixDOF        = joltc.ConstraintSubTypeSixDOF
	ConstraintSubTypePath          = joltc.ConstraintSubTypePath
	ConstraintSubTypeVehicle       = joltc.ConstraintSubTypeVehicle
	ConstraintSubTypeRackAndPinion = joltc.ConstraintSubTypeRackAndPinion
	ConstraintSubTypeGear          = joltc.ConstraintSubTypeGear
	ConstraintSubTypePulley        = joltc.ConstraintSubTypePulley
	ConstraintSubTypeUser1         = joltc.ConstraintSubTypeUser1
	ConstraintSubTypeUser2         = joltc.ConstraintSubTypeUser2
	ConstraintSubTypeUser3         = joltc.ConstraintSubTypeUser3
	ConstraintSubTypeUser4         = joltc.ConstraintSubTypeUser4
)

// ConstraintSpace values.
const (
	ConstraintSpaceLocalToBodyCOM = joltc.ConstraintSpaceLocalToBodyCOM
	ConstraintSpaceWorldSpace     = joltc.ConstraintSpaceWorldSpace
)

// MotorState values.
const (
	MotorStateOff      = joltc.MotorStateOff
	MotorStateVelocity = joltc.MotorStateVelocity
	MotorStatePosition = joltc.MotorStatePosition
)

// SpringMode values.
const (
	SpringModeFrequencyAndDamping = joltc.SpringModeFrequencyAndDamping
	SpringModeStiffnessAndDamping = joltc.SpringModeStiffnessAndDamping
)

// SwingType values.
const (
	SwingTypeCone    = joltc.SwingTypeCone
	SwingTypePyramid = joltc.SwingTypePyramid
)

// SixDOFConstraintAxis values.
const (
	SixDOFConstraintAxisTranslationX = joltc.SixDOFConstraintAxisTranslationX
	SixDOFConstraintAxisTranslationY = joltc.SixDOFConstraintAxisTranslationY
	SixDOFConstraintAxisTranslationZ = joltc.SixDOFConstraintAxisTranslationZ
	SixDOFConstraintAxisRotationX    = joltc.SixDOFConstraintAxisRotationX
	SixDOFConstraintAxisRotationY    = joltc.SixDOFConstraintAxisRotationY
	SixDOFConstraintAxisRotationZ    = joltc.SixDOFConstraintAxisRotationZ
)

// MeshShapeBuildQuality values.
const (
	MeshShapeBuildQualityFavorRuntimePerformance = joltc.MeshShapeBuildQualityFavorRuntimePerformance
	MeshShapeBuildQualityFavorBuildSpeed         = joltc.MeshShapeBuildQualityFavorBuildSpeed
)

// TransmissionMode values.
const (
	TransmissionModeAuto   = joltc.TransmissionModeAuto
	TransmissionModeManual = joltc.TransmissionModeManual
)

// SoftBodyBendType values.
const (
	SoftBodyBendTypeNone     = joltc.SoftBodyBendTypeNone
	SoftBodyBendTypeDistance = joltc.SoftBodyBendTypeDistance
	SoftBodyBendTypeDihedral = joltc.SoftBodyBendTypeDihedral
)
