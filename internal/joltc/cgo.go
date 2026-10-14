//go:build cgo

// Package joltc is the only package in this module that imports "C". It wraps
// the joltc C API (JoltPhysics) and must never leak C types to callers: every
// exported symbol returns Go handle types, vecmath values, primitives, or error.
//
// All bindings require the native static library built by scripts/build_native.sh
// (single precision by default, double with the jolt_double build tag).
package joltc

/*
#cgo CFLAGS: -I${SRCDIR}/../../third_party/joltc/include
#cgo linux LDFLAGS: -lm -lpthread -lstdc++
#cgo darwin LDFLAGS: -lm -lc++

#include <joltc.h>

extern void goJoltcTrace(char* message);
extern unsigned int goJoltcAssertFailure(char* expression, char* message, char* file, unsigned int line);

static void joltcTraceTrampoline(const char* message) {
	goJoltcTrace((char*)message);
}

static bool joltcAssertFailureTrampoline(const char* expression, const char* message, const char* file, uint32_t line) {
	return goJoltcAssertFailure((char*)expression, (char*)message, (char*)file, line) != 0;
}

static bool joltcInit(void) {
	if (!JPH_Init()) {
		return false;
	}
	JPH_SetTraceHandler(joltcTraceTrampoline);
	JPH_SetAssertFailureHandler(joltcAssertFailureTrampoline);
	return true;
}
*/
import "C"

import "sync"

var (
	initMu    sync.Mutex
	initCount int
)

// Initialize starts the Jolt engine. It is reference counted: each call must be
// paired with Shutdown, and the engine is shut down when the count reaches zero.
func Initialize() error {
	initMu.Lock()
	defer initMu.Unlock()
	if initCount == 0 {
		if !bool(C.joltcInit()) {
			return ErrInitFailed
		}
	}
	initCount++
	return nil
}

// Shutdown releases one reference taken by Initialize. It is a no-op when no
// reference is held.
func Shutdown() {
	initMu.Lock()
	defer initMu.Unlock()
	if initCount == 0 {
		return
	}
	initCount--
	if initCount == 0 {
		C.JPH_Shutdown()
	}
}

var (
	handlerMu     sync.RWMutex
	traceHandler  func(string)
	assertHandler func(expression, message, file string, line uint32) bool
)

// SetTraceHandler installs the callback that receives Jolt's trace messages.
// Pass nil to drop trace output.
func SetTraceHandler(fn func(string)) {
	handlerMu.Lock()
	traceHandler = fn
	handlerMu.Unlock()
}

// SetAssertFailureHandler installs the callback invoked when a Jolt assertion
// fails. It returns true to break into the debugger. Pass nil to ignore.
func SetAssertFailureHandler(fn func(expression, message, file string, line uint32) bool) {
	handlerMu.Lock()
	assertHandler = fn
	handlerMu.Unlock()
}

//export goJoltcTrace
func goJoltcTrace(message *C.char) {
	defer func() { _ = recover() }()
	handlerMu.RLock()
	fn := traceHandler
	handlerMu.RUnlock()
	if fn != nil {
		fn(C.GoString(message))
	}
}

//export goJoltcAssertFailure
func goJoltcAssertFailure(expression, message, file *C.char, line C.uint) C.uint {
	defer func() { _ = recover() }()
	handlerMu.RLock()
	fn := assertHandler
	handlerMu.RUnlock()
	if fn == nil {
		return 0
	}
	if fn(C.GoString(expression), C.GoString(message), C.GoString(file), uint32(line)) {
		return 1
	}
	return 0
}
