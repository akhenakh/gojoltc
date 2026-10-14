//go:build cgo && jolt_double

// Double-precision native library (built by: make native-double).
//
// The JPH_DOUBLE_PRECISION define must reach the C side so JPH_RVec3 and
// JPH_RMat4 use their double-precision layout.
package joltc

/*
#cgo CFLAGS: -DJPH_DOUBLE_PRECISION
#cgo LDFLAGS: -L${SRCDIR}/../native/double/lib -ljoltc_double -lJolt
*/
import "C"
