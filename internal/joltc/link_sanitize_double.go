//go:build cgo && joltc_asan && jolt_double

// Double-precision native library built with sanitizers.
//
// Build with: JOLTC_SANITIZE=address,undefined scripts/build_native.sh double
// Run tests with: go test -tags 'joltc_asan jolt_double' -asan ./...
package joltc

/*
#cgo CFLAGS: -DJPH_DOUBLE_PRECISION -fsanitize=address,undefined -fno-omit-frame-pointer
#cgo LDFLAGS: -L${SRCDIR}/../native/double-sanitize/lib -ljoltc_double -lJolt -fsanitize=address,undefined
*/
import "C"
