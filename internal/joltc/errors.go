package joltc

import "errors"

// ErrCgoRequired is returned by every binding when the package is built without
// cgo (CGO_ENABLED=0). joltc cannot function without the native library.
var ErrCgoRequired = errors.New("joltc: cgo is required; build with CGO_ENABLED=1 and a C toolchain")

// ErrInitFailed reports that JPH_Init returned false.
var ErrInitFailed = errors.New("joltc: JPH_Init failed")
