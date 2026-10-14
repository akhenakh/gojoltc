//go:build !cgo

// Package joltc is the only package in this module that imports "C".
//
// This file is compiled when cgo is disabled (CGO_ENABLED=0). It mirrors the
// exported surface of the bindings and returns ErrCgoRequired so that a
// cgo-less build compiles but fails loudly at runtime instead of misbehaving.
package joltc

// Initialize always returns ErrCgoRequired in a cgo-less build.
func Initialize() error { return ErrCgoRequired }

// Shutdown is a no-op in a cgo-less build.
func Shutdown() {}

// SetTraceHandler is a no-op in a cgo-less build.
func SetTraceHandler(fn func(string)) {}

// SetAssertFailureHandler is a no-op in a cgo-less build.
func SetAssertFailureHandler(fn func(expression, message, file string, line uint32) bool) {}
