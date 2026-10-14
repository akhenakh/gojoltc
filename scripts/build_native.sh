#!/usr/bin/env bash
# Build the native joltc static library used by the cgo bindings.
#
# Usage: scripts/build_native.sh [single|double]
#
# Environment overrides:
#   JOLTC_BUILD_TYPE       CMake build type            (default: Release)
#   JOLTC_DEBUG_RENDERER   Debug renderer on/off        (default: OFF)
#   JOLTC_PROFILER         Profiler on/off              (default: OFF)
#   JOLTC_ASSERTS          Jolt asserts on/off          (default: OFF)
#   JOLTC_AVX512           AVX512 on/off                (default: OFF)
#   JOLTC_LTO              Interprocedural opt (LTO)     (default: OFF)
#   JOLTC_JOBS             Parallel build jobs          (default: nproc)
set -euo pipefail

variant="${1:-single}"
case "$variant" in
  single) double=OFF; libname="libjoltc.a" ;;
  double) double=ON;  libname="libjoltc_double.a" ;;
  *) echo "usage: $0 [single|double]" >&2; exit 2 ;;
esac

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
joltc_src="$root/third_party/joltc"
jolt_src="$root/third_party/JoltPhysics"
build_dir="$root/internal/native/build-$variant"
out_dir="$root/internal/native/$variant/lib"

if [[ ! -f "$joltc_src/include/joltc.h" || ! -f "$jolt_src/Build/CMakeLists.txt" ]]; then
  echo "error: submodules missing; run: git submodule update --init --recursive" >&2
  exit 1
fi

: "${JOLTC_BUILD_TYPE:=Release}"
: "${JOLTC_DEBUG_RENDERER:=OFF}"
: "${JOLTC_PROFILER:=OFF}"
: "${JOLTC_ASSERTS:=OFF}"
: "${JOLTC_AVX512:=OFF}"
: "${JOLTC_LTO:=OFF}"
: "${JOLTC_JOBS:=$(nproc 2>/dev/null || echo 4)}"

echo ">> configuring joltc ($variant, $JOLTC_BUILD_TYPE)"
cmake -S "$joltc_src" -B "$build_dir" \
  -DCMAKE_BUILD_TYPE="$JOLTC_BUILD_TYPE" \
  -DCMAKE_POSITION_INDEPENDENT_CODE=ON \
  -DJOLT_PHYSICS_ROOT="$jolt_src" \
  -DJPH_BUILD_SHARED=OFF \
  -DJPH_SAMPLES=OFF \
  -DJPH_TESTS=OFF \
  -DJPH_INSTALL=OFF \
  -DDOUBLE_PRECISION="$double" \
  -DUSE_ASSERTS="$JOLTC_ASSERTS" \
  -DCPP_EXCEPTIONS_ENABLED=OFF \
  -DCPP_RTTI_ENABLED=OFF \
  -DUSE_AVX512="$JOLTC_AVX512" \
  -DINTERPROCEDURAL_OPTIMIZATION="$JOLTC_LTO" \
  -DDEBUG_RENDERER_IN_DEBUG_AND_RELEASE="$JOLTC_DEBUG_RENDERER" \
  -DDEBUG_RENDERER_IN_DISTRIBUTION="$JOLTC_DEBUG_RENDERER" \
  -DPROFILER_IN_DEBUG_AND_RELEASE="$JOLTC_PROFILER" \
  -DPROFILER_IN_DISTRIBUTION="$JOLTC_PROFILER" \
  -DJPH_USE_DX12=OFF \
  -DJPH_USE_VK=OFF \
  -DJPH_USE_MTL=OFF \
  -DJPH_USE_CPU_COMPUTE=ON \
  -DENABLE_INSTALL=OFF \
  -DTARGET_UNIT_TESTS=OFF \
  -DTARGET_HELLO_WORLD=OFF \
  -DTARGET_PERFORMANCE_TEST=OFF \
  -DTARGET_SAMPLES=OFF \
  -DTARGET_VIEWER=OFF

echo ">> building joltc ($JOLTC_JOBS jobs)"
cmake --build "$build_dir" --config "$JOLTC_BUILD_TYPE" -j "$JOLTC_JOBS"

found="$(find "$build_dir" -name "$libname" -type f | head -1)"
if [[ -z "$found" ]]; then
  echo "error: $libname not found under $build_dir" >&2
  exit 1
fi

jolt_lib="$(find "$build_dir" -name 'libJolt.a' -type f | head -1)"
if [[ -z "$jolt_lib" ]]; then
  echo "error: libJolt.a not found under $build_dir" >&2
  exit 1
fi

mkdir -p "$out_dir"
cp -f "$found" "$out_dir/$libname"
cp -f "$jolt_lib" "$out_dir/libJolt.a"
echo ">> built $out_dir/$libname (+ libJolt.a)"
