# ADR-0001: Wrap joltc with cgo rather than purego or an RPC sidecar

- Status: accepted
- Date: 2026-10-14

## Context

`gojoltc` exposes JoltPhysics to Go. The upstream native interface is
[joltc](https://github.com/amerkoleci/joltc), a flat C API (`~1,269` functions)
over JoltPhysics, built as a static library and linked into the binary.

Alternatives considered:

- **purego** (no cgo, dynamic symbol lookup), as used by
  `NilssonCreative/jolt-purego`. Requires a shared library at runtime and routes
  every call through reflection-based dynamic lookup.
- **RPC sidecar**: run physics in a separate process and talk over IPC.
- **Pure-Go port**: not realistic for a mature SIMD physics engine.

## Decision

Use cgo, with the FFI confined to a single internal package (`internal/joltc`)
and the public surface kept pure Go (`jolt`).

## Rationale

- **Latency and allocation control.** The hot path is `PhysicsSystem.Update`
  (one cgo crossing per step) plus per-frame body sweeps. A C shim keeps a whole
  sweep to a single crossing with zero allocations; a per-call, reflection-based
  dynamic lookup cannot match this.
- **Static linking.** cgo over the vendored static library yields a single
  self-contained binary with no runtime shared-object or discovery step.
- **No reflection on the call path.** purego's `RegisterLibFunc` indirection adds
  overhead per call and complicates the ownership story.
- **Build reproducibility.** Sources are pinned as submodules; the module builds
  the library with CMake once and then compiles quickly.

Trade-offs accepted:

- A C toolchain is required to build (cgo is mandatory; a `!cgo` stub returns
  `ErrCgoRequired`).
- Cross-compilation needs a target toolchain.
- All FFI hazards (pointer lifetimes, callbacks) must be handled explicitly; the
  invariants are recorded in `.operator/specs/ffi-and-build.md`.

## Consequences

- `import "C"` is allowed only under `internal/joltc` (enforced by a CI gate).
- The public `jolt` API never exposes `C.*` types.
- A native build step (`make native`) precedes `go test`.
