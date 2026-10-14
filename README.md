# gojoltc

Go bindings for [JoltPhysics](https://github.com/jrouwe/JoltPhysics) via the
[joltc](https://github.com/amerkoleci/joltc) C API, for simulating 3D physics in
games.

The public API is pure Go: `import "C"` is confined to `internal/joltc`, so
callers never deal with cgo types. The native library is linked statically, so
the resulting Go binary has no runtime Jolt dependency.

## Status

The full joltc C API surface exercised by the upstream test suite is wrapped:
engine lifetimes, job system, layer filters, physics system, bodies, all shape
types and settings, constraints, broad/narrow phase queries, characters,
vehicles, skeletons and ragdolls. All 188 upstream joltc test cases are ported
to Go and passing (`make test`).

Not yet wrapped: contact/activation/step listeners, debug renderer, soft bodies,
and the long tail of accessors with no upstream test. See `.operator/plans/implementation.md`.

## Requirements

- Go with cgo enabled (`CGO_ENABLED=1`) and a C/C++ toolchain (gcc/clang, MSVC).
- CMake ≥ 3.16.
- Git submodules checked out.

```sh
git submodule update --init --recursive
```

## Building

The native static library (`libjoltc.a` + `libJolt.a`) is built once with CMake
and linked by cgo:

```sh
make native          # single precision (default)
make native-double   # double precision (jolt_double build tag)
```

Configurable via environment: `JOLTC_BUILD_TYPE` (default `Release`),
`JOLTC_DEBUG_RENDERER`, `JOLTC_PROFILER`, `JOLTC_ASSERTS` (default `OFF`),
`JOLTC_AVX512` (default `OFF`), `JOLTC_LTO` (default `OFF`), `JOLTC_JOBS`.

`make check` runs the definition-of-done gate: gofmt, `go vet`, the cgo-boundary
check, `go test ./...`, and a `CGO_ENABLED=0` build.

## Usage

```go
package main

import (
	"log"

	"github.com/akhenakh/gojoltc/jolt"
)

func main() {
	if err := jolt.Init(); err != nil {
		log.Fatal(err)
	}
	defer jolt.Shutdown()

	// ... create a physics system, bodies, and step the simulation.
}
```

## Precision

The default build is single precision. Building with `-tags jolt_double` switches
world-space types (`jolt.Real`, `jolt.RVec3`, `jolt.RMat4`) to `float64` and links
the double-precision library. Single and double are separate builds.

## Platform support

| OS      | Architectures   | Status        |
| ------- | --------------- | ------------- |
| Linux   | amd64, arm64    | supported     |
| macOS   | arm64, amd64    | supported     |
| Windows | amd64           | supported     |

## Native dependencies

Vendored as git submodules under `third_party/`:

- [joltc](https://github.com/amerkoleci/joltc) — MIT
- [JoltPhysics](https://github.com/jrouwe/JoltPhysics) — MIT, pinned at `v5.6.0`

GPU compute backends (DX12/Vulkan/Metal) are disabled; the build uses Jolt's CPU
compute only.

## License

MIT.
