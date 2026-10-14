# gojoltc

Go bindings for [JoltPhysics](https://github.com/jrouwe/JoltPhysics) via the
[joltc](https://github.com/amerkoleci/joltc) C API, for simulating 3D physics in
games.

The public API is pure Go: `import "C"` is confined to `internal/joltc`, so
callers never deal with cgo types. The native library is linked statically, so
the resulting Go binary has no runtime Jolt dependency.

## Status

The API surface exercised by the upstream joltc test suite is wrapped, and a
usability pass makes it actually usable in games:

- **Constraints** can be attached to a world (`AddConstraint`/`RemoveConstraint`).
- **Wheeled vehicles** drive: chassis + `NewWheelSettingsWV` wheels + a controller
  (with at least one differential), added as both a constraint and a step listener.
- **Mesh and height-field** shapes for level and terrain geometry.
- **Sensing queries**: filtered raycasts (closest hit and all hits), point and
  shape overlap (with proximity volumes), and shape sweeps —
  `BroadPhaseLayerFilter`/`ObjectLayerFilter`/`BodyFilter`/`ShapeFilter` as Go callbacks.
- **Listeners** (contact/activation/step) and a **debug renderer**.
- **Soft bodies**, characters, skeletons/ragdolls.
- **Solver tuning** (`PhysicsSettings`), **collision groups**, physics materials.
- The gameplay slice of `BodyInterface` (transforms, layers, user data, kinematic
  moves, torque/forces) and `PhysicsMaterial`.

Not wrapped: `CharacterVirtual` (the rigid `Character` covers this today) and the
remaining composite functions — see `docs/coverage.md`. See `.operator/plans/implementation.md`.

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

## Performance

The hot paths stay on the C heap where it matters:

- `PhysicsSystem.Update` is a single cgo crossing; stepping parallelism comes
  from Jolt's own job system, not from Go.
- Sub-interfaces (`BodyInterface`, `NarrowPhaseQuery`, `BroadPhaseQuery`, ...)
  are cached on the system, so repeated access neither allocates nor calls into C.
- Per-frame loops should use the allocation-free batched accessors
  (`PositionsInto`, `PositionsAndRotationsInto`, `ActiveBodyIDsInto`, ...), which
  write into caller-owned slices. The per-call accessors allocate a small
  C out-parameter on the Go heap each call.

`make bench` reports allocations per op to keep this honest. On a 12-core
machine with 1024 dynamic bodies: a full position sweep is ~11µs with 0
allocations via `PositionsInto` (one cgo crossing for the whole sweep), versus
~49µs with 1024 allocations via the per-call `Position`; `Update` with 512
bodies is allocation-free.

## Examples

- `examples/helloworld` — console port of joltc's `samples/01_HelloWorld`.
- `cmd/example` — the same scene rendered with
  [raylib-go](https://github.com/gen2brain/raylib-go), driving a character with
  WASD and drawing Jolt's debug geometry.
- `cmd/joltviewer` — a JoltViewer-style viewer: orbit camera, pause/step,
  wireframe toggle and switchable scenes (pyramid, sphere drop, soft-body cloth).

`cmd/example` and `cmd/joltviewer` are **separate Go modules** (each with its own
`go.mod` and a `replace` to the parent) so the library does not pull in raylib:

```sh
cd cmd/example && go run .
cd cmd/joltviewer && go run .
```

They need an OpenGL/X11 (or Wayland) toolchain to build raylib.

## Precision

The default build is single precision. Building with `-tags jolt_double` switches
world-space types (`jolt.Real`, `jolt.RVec3`, `jolt.RMat4`) to `float64` and links
the double-precision library. Single and double are separate builds.

## Platform support

| OS      | Architectures | Status                                   |
| ------- | ------------- | ---------------------------------------- |
| Linux   | amd64, arm64  | CI-tested                                 |
| macOS   | arm64         | CI-tested                                 |
| Windows | amd64         | not yet (needs MinGW link flags)          |

Both precisions are exercised in CI, along with an AddressSanitizer/UBSan run.

## Development

```sh
make native          # build the single-precision native static library
make check           # gofmt, vet, lint, generate-check, tests, no-cgo build (the done gate)
make bench           # benchmarks with allocation reporting
make test-sanitize   # native library + tests under ASan/UBSan
make generate        # regenerate generated accessors + docs/coverage.md
make sbom            # regenerate sbom.json
```

## Code generation

`gen/` parses `third_party/joltc/include/joltc.h` and emits thin accessor
wrappers (getters/setters returning scalars, enums and vectors) for the
functions not already hand-wrapped, into `internal/joltc/gen_accessors.go`
(committed; `make generate-check` fails if it is stale). The generated methods
hang off the existing handle types, so the public `jolt` API picks them up and
the `!cgo` stub needs no change. Composite signatures (create/destroy,
callbacks, POD structs) stay hand-written. See
[`docs/coverage.md`](docs/coverage.md) for current coverage.

Design decisions are recorded in [`docs/adr/`](docs/adr/) and the native
dependency versions in [`sbom.json`](sbom.json).

## Native dependencies

Vendored as git submodules under `third_party/`:

- [joltc](https://github.com/amerkoleci/joltc) — MIT
- [JoltPhysics](https://github.com/jrouwe/JoltPhysics) — MIT, pinned at `v5.6.0`

Versions and licenses are tracked in [`sbom.json`](sbom.json) (`make sbom`).
GPU compute backends (DX12/Vulkan/Metal) are disabled; the build uses Jolt's CPU
compute only. Both are linked statically, so there is no runtime shared object.

## License

MIT.
