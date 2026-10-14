package jolt

import "github.com/akhenakh/gojoltc/internal/joltc"

// Init starts the Jolt engine. It is reference counted: each call must be
// paired with a call to Shutdown.
func Init() error {
	return joltc.Initialize()
}

// Shutdown releases one reference taken by Init.
func Shutdown() {
	joltc.Shutdown()
}

// SetTraceHandler installs the callback that receives Jolt's trace messages.
// Pass nil to drop trace output.
func SetTraceHandler(fn func(string)) {
	joltc.SetTraceHandler(fn)
}

// SetAssertFailureHandler installs the callback invoked when a Jolt assertion
// fails. It returns true to break into the debugger. Pass nil to ignore.
func SetAssertFailureHandler(fn func(expression, message, file string, line uint32) bool) {
	joltc.SetAssertFailureHandler(fn)
}
