FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

VGO_PATH:=$(FIRST_GOPATH)/src/golang.org/x/vgo
VGO_VERSION:=master
VGO_BIN:=$(GOBIN)/vgo

LOCAL_BIN:=$(CURDIR)/bin
GEN_CLAY_BIN:=$(CURDIR)/bin/protoc-gen-goclay
export GEN_CLAY_BIN
GEN_GOFAST_BIN:=$(CURDIR)/bin/protoc-gen-gofast
export GEN_GOFAST_BIN

GRPC_GATEWAY_PKG:=$(shell vgo list -m all | grep github.com/grpc-ecosystem/grpc-gateway | awk '{print ($$4 != "" ? $$4 : $$1)}')
GRPC_GATEWAY_VERSION:=$(shell vgo list -m all | grep github.com/grpc-ecosystem/grpc-gateway | awk '{print ($$5 != "" ? $$5 : $$2)}')
GRPC_GATEWAY_PATH:=${FIRST_GOPATH}/src/mod/${GRPC_GATEWAY_PKG}@${GRPC_GATEWAY_VERSION}
export GRPC_GATEWAY_PATH

all: clean build test

# install vgo
$(VGO_BIN):
ifeq (${VGO_VERSION},master)
	$(info #Installing vgo version $(VGO_VERSION)...)
ifneq ($(wildcard $(VGO_PATH)),)
	rm -rf $(VGO_PATH)
endif
	go get -u golang.org/x/vgo

else
	$(info #Installing vgo version $(VGO_VERSION)...)
ifeq ($(wildcard $(VGO_PATH)),)
	mkdir -p $(VGO_PATH) && cd $(VGO_PATH) ;\
	git clone https://github.com/golang/vgo.git .
endif
	cd $(VGO_PATH) && git fetch --tags && git checkout $(VGO_VERSION) ;\
	git reset --hard && git clean -fd ;\
	go build -o=$(VGO_BIN) main.go
endif

GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m

build: $(VGO_BIN)
	$(info #Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) vgo install github.com/utrack/clay/v2/cmd/protoc-gen-goclay
	GOBIN=$(LOCAL_BIN) vgo install github.com/gogo/protobuf/protoc-gen-gofast
	vgo mod -vendor

test: $(VGO_BIN)
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
