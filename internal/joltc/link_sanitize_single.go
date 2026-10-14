//go:build cgo && joltc_asan && !jolt_double

// Single-precision native library built with sanitizers.
//
// Build with: JOLTC_SANITIZE=address,undefined scripts/build_native.sh single
// Run tests with: go test -tags joltc_asan -asan ./...
package joltc

/*
#cgo CFLAGS: -fsanitize=address,undefined -fno-omit-frame-pointer
#cgo LDFLAGS: -L${SRCDIR}/../native/single-sanitize/lib -ljoltc -lJolt -fsanitize=address,undefined
*/
import "C"
