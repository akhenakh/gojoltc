SHELL := /usr/bin/env bash

GO ?= go
GOFLAGS ?=
SINGLE_LIB := internal/native/single/lib/libjoltc.a
DOUBLE_LIB := internal/native/double/lib/libjoltc_double.a
GO_SOURCES := $(shell find . -name '*.go' -not -path './third_party/*' -not -path './internal/native/*')

.PHONY: all native native-double native-sanitize native-clean fmt fmt-check vet lint generate generate-check \
        test test-double build-nocgo cgo-boundary check check-double bench test-sanitize sbom examples clean

all: check

# --- Native library -------------------------------------------------------

# The scripts decide whether a rebuild is needed (submodule revisions and build
# options are recorded next to the archive), so these always run.
native:
	scripts/build_native.sh single

native-double:
	scripts/build_native.sh double

native-sanitize:
	JOLTC_SANITIZE=address,undefined scripts/build_native.sh single

sbom:
	scripts/gen_sbom.sh

native-clean:
	rm -rf internal/native

# --- Go quality gates -----------------------------------------------------
#
# vet runs with -composites=false: vecmath value types (Vec3, Quat, Mat4, ...)
# are constructed positionally (Vec3{1, 2, 3}), which the composites analyzer
# rejects for types defined in another package. This is intentional.

fmt:
	gofmt -w $(GO_SOURCES)

fmt-check:
	@out="$$(gofmt -l $(GO_SOURCES))"; \
	if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

vet:
	$(GO) vet -composites=false $(GOFLAGS) ./...

cgo-boundary:
	scripts/check_cgo_boundary.sh

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 5m; \
	else \
		echo "golangci-lint not found; skipping (install for full lint)"; \
	fi

generate:
	$(GO) run ./gen

generate-check:
	$(GO) run ./gen -check

# -count=1: the native archive is not part of Go's test cache key, so a rebuilt
# library would otherwise report a cached (stale) pass.
test: native
	$(GO) test -count=1 $(GOFLAGS) ./...

test-double: native-double
	$(GO) test -count=1 $(GOFLAGS) -tags jolt_double ./...

bench: native
	$(GO) test $(GOFLAGS) -run '^$$' -bench . -benchmem ./jolt/

test-sanitize: native-sanitize
	# Leak detection is off: third_party/joltc leaks a default EmptyShapeSettings
	# in JPH_CharacterBaseSettings_Init that cannot be freed from outside, the
	# suite has a handful of its own leaks to triage, and the frames are
	# unsymbolized so no leak: suppression can be written by name. ASan still
	# catches memory errors and UBSan catches undefined behaviour.
	ASAN_OPTIONS=detect_leaks=0 UBSAN_OPTIONS=halt_on_error=1:print_stacktrace=1 \
	$(GO) test -count=1 $(GOFLAGS) -tags joltc_asan -asan ./...

build-nocgo:
	@pkgs="$$($(GO) list ./... | grep -v '/examples/')"; \
	CGO_ENABLED=0 $(GO) build $$pkgs

examples: native
	$(GO) build ./examples/...

# What must pass before a change is declared done.
check: fmt-check vet cgo-boundary lint generate-check test build-nocgo

check-double: fmt-check vet cgo-boundary test-double build-nocgo

clean: native-clean
	$(GO) clean -cache -testcache
