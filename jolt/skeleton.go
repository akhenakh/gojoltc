package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Skeleton and ragdoll types, re-exported from the bindings.
type (
	// Skeleton is a hierarchy of joints.
	Skeleton = joltc.Skeleton
	// SkeletonPose is a pose for a skeleton.
	SkeletonPose = joltc.SkeletonPose
	// SkeletonMapper maps poses between two skeletons.
	SkeletonMapper = joltc.SkeletonMapper
	// RagdollSettings holds settings for a ragdoll.
	RagdollSettings = joltc.RagdollSettings
	// Ragdoll is a physics-driven skeleton.
	Ragdoll = joltc.Ragdoll
)

// NewSkeleton creates an empty skeleton.
func NewSkeleton() *Skeleton { return joltc.NewSkeleton() }

// NewSkeletonPose creates an empty pose.
func NewSkeletonPose() *SkeletonPose { return joltc.NewSkeletonPose() }

// NewSkeletonMapper creates an empty mapper.
func NewSkeletonMapper() *SkeletonMapper { return joltc.NewSkeletonMapper() }

// NewRagdollSettings creates empty ragdoll settings.
func NewRagdollSettings() *RagdollSettings { return joltc.NewRagdollSettings() }
