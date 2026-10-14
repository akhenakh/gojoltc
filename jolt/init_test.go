package jolt

import (
	"strings"
	"sync"
	"testing"
)

func TestInitShutdown(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	Shutdown()

	// Reference counting: two Init calls need two Shutdown calls.
	if err := Init(); err != nil {
		t.Fatalf("Init (2nd): %v", err)
	}
	if err := Init(); err != nil {
		t.Fatalf("Init (3rd): %v", err)
	}
	Shutdown()
	Shutdown()
	// Extra Shutdown must be a harmless no-op.
	Shutdown()
}

func TestTraceHandler(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	defer Shutdown()

	var mu sync.Mutex
	var seen []string
	SetTraceHandler(func(msg string) {
		mu.Lock()
		seen = append(seen, msg)
		mu.Unlock()
	})
	defer SetTraceHandler(nil)

	// Jolt emits trace output on at least one of these paths; we only assert
	// that the handler mechanism is wired and never panics.
	if err := Init(); err != nil {
		t.Fatalf("Init (nested): %v", err)
	}
	Shutdown()

	mu.Lock()
	defer mu.Unlock()
	for _, m := range seen {
		if strings.Contains(m, "\x00") {
			t.Fatalf("trace message contains NUL: %q", m)
		}
	}
}
