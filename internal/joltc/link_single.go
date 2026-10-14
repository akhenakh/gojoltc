//go:build cgo && !jolt_double

// Single-precision native library (built by: make native).
package joltc

/*
#cgo LDFLAGS: -L${SRCDIR}/../native/single/lib -ljoltc -lJolt
*/
import "C"
