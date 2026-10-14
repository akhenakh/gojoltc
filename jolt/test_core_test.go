package jolt

import "testing"

// Ported from third_party/joltc/tests/test_core.cpp.

func TestCoreInitShutdown(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	Shutdown()
}

func TestCoreInitMultipleTimes(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	Shutdown()

	if err := Init(); err != nil {
		t.Fatalf("Init (2nd): %v", err)
	}
	Shutdown()
}

func TestCoreJobSystemDefault(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	js := NewJobSystemThreadPool(nil)
	if js == nil {
		t.Fatal("job system is nil")
	}
	js.Destroy()
}

func TestCoreJobSystemWithConfig(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	js := NewJobSystemThreadPool(&JobSystemConfig{
		MaxJobs:     MaxPhysicsJobs,
		MaxBarriers: MaxPhysicsBarriers,
		NumThreads:  2,
	})
	if js == nil {
		t.Fatal("job system is nil")
	}
	js.Destroy()
}

func TestCoreSetTraceHandler(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	SetTraceHandler(func(string) {})
	SetTraceHandler(nil)
}

func TestCoreSetAssertFailureHandler(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	SetAssertFailureHandler(func(expression, message, file string, line uint32) bool { return true })
	SetAssertFailureHandler(nil)
}
