FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

LOCAL_BIN:=$(CURDIR)/bin
GEN_CLAY_BIN:=$(CURDIR)/bin/protoc-gen-goclay
export GEN_CLAY_BIN
GEN_GOFAST_BIN:=$(CURDIR)/bin/protoc-gen-gofast
export GEN_GOFAST_BIN
GEN_GOGOFAST_BIN:=$(CURDIR)/bin/protoc-gen-gogofast
export GEN_GOGOFAST_BIN

GRPC_GATEWAY_PKG:=$(shell vgo list -m all | grep github.com/grpc-ecosystem/grpc-gateway | awk '{print ($$4 != "" ? $$4 : $$1)}')
GRPC_GATEWAY_VERSION:=$(shell vgo list -m all | grep github.com/grpc-ecosystem/grpc-gateway | awk '{print ($$5 != "" ? $$5 : $$2)}')
GRPC_GATEWAY_PATH:=${FIRST_GOPATH}/src/mod/${GRPC_GATEWAY_PKG}@${GRPC_GATEWAY_VERSION}
export GRPC_GATEWAY_PATH

all: clean build test

GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m

build:
	$(info #Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) vgo install github.com/utrack/clay/v2/cmd/protoc-gen-goclay
	GOBIN=$(LOCAL_BIN) vgo install github.com/gogo/protobuf/protoc-gen-gofast
	GOBIN=$(LOCAL_BIN) vgo install github.com/gogo/protobuf/protoc-gen-gogofast

test:
	@ \
	success=0; \
	failure=0; \
	for i in `find */ -type f -name Makefile -not -path "vendor/*"`; do \
		if make -C $${i%/*} test; then \
			success=$$((success+1)); \
		else \
			failure=$$((failure+1)); \
			echo "${RED}TEST FAILED${NC}"; \
		fi; \
		echo; \
	done; \
	echo "== RESULTS =="; \
	echo "${GREEN}$$success TEST PASSED${NC}"; \
	if [ "$$failure" -gt "0" ]; then \
		echo "${RED}$$failure TEST FAILED${NC}"; \
		exit 1; \
	fi

clean:
	@find */ -type f -name Makefile -execdir sh -c "make clean; echo ;" \;
