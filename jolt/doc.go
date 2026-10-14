// Package jolt is a Go wrapper around JoltPhysics, via the joltc C API.
//
// The package is pure Go: it exposes only Go types, functions, and errors.
// All cgo glue lives in internal/joltc, behind this facade. Callers therefore
// need no cgo knowledge to use it, though the module does require cgo and a C
// toolchain to build (see the README).
//
// # Lifecycle
//
// Call Init once before creating any engine object and Shutdown when done:
//
//	if err := jolt.Init(); err != nil {
//		log.Fatal(err)
//	}
//	defer jolt.Shutdown()
//
// # Precision
//
// The default build is single precision. Building with the jolt_double tag
// switches world-space types (Real, RVec3, RMat4) to float64.
package jolt
