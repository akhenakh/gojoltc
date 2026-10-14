SHELL := /usr/bin/env bash

GO ?= go
GOFLAGS ?=
SINGLE_LIB := internal/native/single/lib/libjoltc.a
DOUBLE_LIB := internal/native/double/lib/libjoltc_double.a
GO_SOURCES := $(shell find . -name '*.go' -not -path './third_party/*' -not -path './internal/native/*')

.PHONY: all native native-double native-sanitize native-clean fmt fmt-check vet lint test test-double \
        build-nocgo cgo-boundary check check-double bench test-sanitize sbom examples clean

all: check

# --- Native library -------------------------------------------------------

native: $(SINGLE_LIB)

native-double: $(DOUBLE_LIB)

native-sanitize:
	JOLTC_SANITIZE=address,undefined scripts/build_native.sh single

sbom:
	scripts/gen_sbom.sh

$(SINGLE_LIB):
	scripts/build_native.sh single

$(DOUBLE_LIB):
	scripts/build_native.sh double

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

test: native
	$(GO) test $(GOFLAGS) ./...

test-double: native-double
	$(GO) test $(GOFLAGS) -tags jolt_double ./...

bench: native
	$(GO) test $(GOFLAGS) -run '^$$' -bench . -benchmem ./jolt/

test-sanitize: native-sanitize
	# Leak detection is off: third_party/joltc leaks a default EmptyShapeSettings
	# in JPH_CharacterBaseSettings_Init that cannot be freed from outside. ASan
	# still catches memory errors and UBSan catches undefined behaviour.
	ASAN_OPTIONS=detect_leaks=0 UBSAN_OPTIONS=halt_on_error=1:print_stacktrace=1 \
	$(GO) test $(GOFLAGS) -tags joltc_asan -asan ./...

build-nocgo:
	@pkgs="$$($(GO) list ./... | grep -v '/examples/')"; \
	CGO_ENABLED=0 $(GO) build $$pkgs

examples: native
	$(GO) build ./examples/...

# What must pass before a change is declared done.
check: fmt-check vet cgo-boundary lint test build-nocgo

check-double: fmt-check vet cgo-boundary test-double build-nocgo

clean: native-clean
	$(GO) clean -cache -testcache
