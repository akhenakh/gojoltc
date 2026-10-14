//go:build cgo

package joltc

/*
#include <joltc.h>
*/
import "C"

// Compile-time assertions that every hand-written Go enum value matches its
// joltc C enumerator. Each line is zero when the values agree ([0]struct{}
// matches) and non-zero otherwise, which fails to compile — so a drift between
// the Go values and the C header is caught at build time, not at runtime.
var (
	_ [0]struct{} = [int32(C.JPH_PhysicsUpdateError_None) - int32(PhysicsUpdateErrorNone)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_PhysicsUpdateError_ManifoldCacheFull) - int32(PhysicsUpdateErrorManifoldCacheFull)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_PhysicsUpdateError_BodyPairCacheFull) - int32(PhysicsUpdateErrorBodyPairCacheFull)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_PhysicsUpdateError_ContactConstraintsFull) - int32(PhysicsUpdateErrorContactConstraintsFull)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_BodyType_Rigid) - int32(BodyTypeRigid)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_BodyType_Soft) - int32(BodyTypeSoft)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_MotionType_Static) - int32(MotionTypeStatic)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_MotionType_Kinematic) - int32(MotionTypeKinematic)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_MotionType_Dynamic) - int32(MotionTypeDynamic)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_Activation_Activate) - int32(ActivationActivate)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_Activation_DontActivate) - int32(ActivationDontActivate)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ValidateResult_AcceptAllContactsForThisBodyPair) - int32(ValidateResultAcceptAllContactsForThisBodyPair)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ValidateResult_AcceptContact) - int32(ValidateResultAcceptContact)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ValidateResult_RejectContact) - int32(ValidateResultRejectContact)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ValidateResult_RejectAllContactsForThisBodyPair) - int32(ValidateResultRejectAllContactsForThisBodyPair)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ShapeType_Convex) - int32(ShapeTypeConvex)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_Compound) - int32(ShapeTypeCompound)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_Decorated) - int32(ShapeTypeDecorated)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_Mesh) - int32(ShapeTypeMesh)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_HeightField) - int32(ShapeTypeHeightField)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_SoftBody) - int32(ShapeTypeSoftBody)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_User1) - int32(ShapeTypeUser1)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_User2) - int32(ShapeTypeUser2)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_User3) - int32(ShapeTypeUser3)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeType_User4) - int32(ShapeTypeUser4)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Sphere) - int32(ShapeSubTypeSphere)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Box) - int32(ShapeSubTypeBox)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Triangle) - int32(ShapeSubTypeTriangle)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Capsule) - int32(ShapeSubTypeCapsule)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_TaperedCapsule) - int32(ShapeSubTypeTaperedCapsule)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Cylinder) - int32(ShapeSubTypeCylinder)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_ConvexHull) - int32(ShapeSubTypeConvexHull)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_StaticCompound) - int32(ShapeSubTypeStaticCompound)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_MutableCompound) - int32(ShapeSubTypeMutableCompound)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_RotatedTranslated) - int32(ShapeSubTypeRotatedTranslated)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Scaled) - int32(ShapeSubTypeScaled)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_OffsetCenterOfMass) - int32(ShapeSubTypeOffsetCenterOfMass)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_Mesh) - int32(ShapeSubTypeMesh)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_HeightField) - int32(ShapeSubTypeHeightField)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ShapeSubType_SoftBody) - int32(ShapeSubTypeSoftBody)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_MotionQuality_Discrete) - int32(MotionQualityDiscrete)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_MotionQuality_LinearCast) - int32(MotionQualityLinearCast)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_OverrideMassProperties_CalculateMassAndInertia) - int32(OverrideMassPropertiesCalculateMassAndInertia)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_OverrideMassProperties_CalculateInertia) - int32(OverrideMassPropertiesCalculateInertia)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_OverrideMassProperties_MassAndInertiaProvided) - int32(OverrideMassPropertiesMassAndInertiaProvided)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_All) - int32(AllowedDOFsAll)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_TranslationX) - int32(AllowedDOFsTranslationX)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_TranslationY) - int32(AllowedDOFsTranslationY)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_TranslationZ) - int32(AllowedDOFsTranslationZ)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_RotationX) - int32(AllowedDOFsRotationX)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_RotationY) - int32(AllowedDOFsRotationY)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_RotationZ) - int32(AllowedDOFsRotationZ)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_AllowedDOFs_Plane2D) - int32(AllowedDOFsPlane2D)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_GroundState_OnGround) - int32(GroundStateOnGround)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_GroundState_OnSteepGround) - int32(GroundStateOnSteepGround)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_GroundState_NotSupported) - int32(GroundStateNotSupported)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_GroundState_InAir) - int32(GroundStateInAir)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_BackFaceMode_IgnoreBackFaces) - int32(BackFaceModeIgnoreBackFaces)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_BackFaceMode_CollideWithBackFaces) - int32(BackFaceModeCollideWithBackFaces)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ActiveEdgeMode_CollideOnlyWithActive) - int32(ActiveEdgeModeCollideOnlyWithActive)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ActiveEdgeMode_CollideWithAll) - int32(ActiveEdgeModeCollideWithAll)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_CollectFacesMode_CollectFaces) - int32(CollectFacesModeCollectFaces)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_CollectFacesMode_NoFaces) - int32(CollectFacesModeNoFaces)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_CollisionCollectorType_AllHit) - int32(CollisionCollectorTypeAllHit)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_CollisionCollectorType_AllHitSorted) - int32(CollisionCollectorTypeAllHitSorted)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_CollisionCollectorType_ClosestHit) - int32(CollisionCollectorTypeClosestHit)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_CollisionCollectorType_AnyHit) - int32(CollisionCollectorTypeAnyHit)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ConstraintType_Constraint) - int32(ConstraintTypeConstraint)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintType_TwoBodyConstraint) - int32(ConstraintTypeTwoBodyConstraint)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Fixed) - int32(ConstraintSubTypeFixed)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Point) - int32(ConstraintSubTypePoint)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Hinge) - int32(ConstraintSubTypeHinge)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Slider) - int32(ConstraintSubTypeSlider)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Distance) - int32(ConstraintSubTypeDistance)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Cone) - int32(ConstraintSubTypeCone)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_SwingTwist) - int32(ConstraintSubTypeSwingTwist)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_SixDOF) - int32(ConstraintSubTypeSixDOF)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Path) - int32(ConstraintSubTypePath)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Vehicle) - int32(ConstraintSubTypeVehicle)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_RackAndPinion) - int32(ConstraintSubTypeRackAndPinion)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Gear) - int32(ConstraintSubTypeGear)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_Pulley) - int32(ConstraintSubTypePulley)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_User1) - int32(ConstraintSubTypeUser1)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_User2) - int32(ConstraintSubTypeUser2)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_User3) - int32(ConstraintSubTypeUser3)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSubType_User4) - int32(ConstraintSubTypeUser4)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_ConstraintSpace_LocalToBodyCOM) - int32(ConstraintSpaceLocalToBodyCOM)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_ConstraintSpace_WorldSpace) - int32(ConstraintSpaceWorldSpace)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_MotorState_Off) - int32(MotorStateOff)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_MotorState_Velocity) - int32(MotorStateVelocity)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_MotorState_Position) - int32(MotorStatePosition)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_SpringMode_FrequencyAndDamping) - int32(SpringModeFrequencyAndDamping)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SpringMode_StiffnessAndDamping) - int32(SpringModeStiffnessAndDamping)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_SwingType_Cone) - int32(SwingTypeCone)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SwingType_Pyramid) - int32(SwingTypePyramid)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_SixDOFConstraintAxis_TranslationX) - int32(SixDOFConstraintAxisTranslationX)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SixDOFConstraintAxis_TranslationY) - int32(SixDOFConstraintAxisTranslationY)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SixDOFConstraintAxis_TranslationZ) - int32(SixDOFConstraintAxisTranslationZ)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SixDOFConstraintAxis_RotationX) - int32(SixDOFConstraintAxisRotationX)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SixDOFConstraintAxis_RotationY) - int32(SixDOFConstraintAxisRotationY)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_SixDOFConstraintAxis_RotationZ) - int32(SixDOFConstraintAxisRotationZ)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_Mesh_Shape_BuildQuality_FavorRuntimePerformance) - int32(MeshShapeBuildQualityFavorRuntimePerformance)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_Mesh_Shape_BuildQuality_FavorBuildSpeed) - int32(MeshShapeBuildQualityFavorBuildSpeed)]struct{}{}

	_ [0]struct{} = [int32(C.JPH_TransmissionMode_Auto) - int32(TransmissionModeAuto)]struct{}{}
	_ [0]struct{} = [int32(C.JPH_TransmissionMode_Manual) - int32(TransmissionModeManual)]struct{}{}
)
