FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

VGO_PATH:=$(FIRST_GOPATH)/src/golang.org/x/vgo
VGO_VERSION:=master
VGO_BIN:=$(GOBIN)/vgo

GEN_CLAY_BIN:=$(CURDIR)/bin/protoc-gen-goclay
export GEN_CLAY_BIN

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
	-rm $(GEN_CLAY_BIN)
	cd ../cmd/protoc-gen-goclay && vgo build -o=$(GEN_CLAY_BIN)  .

test: $(VGO_BIN) build
		@ \
		success=0; \
		failure=0; \
		for i in `find */ -type f -name Makefile`; do \
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
		fi

clean:
	@find */ -type f -name Makefile -execdir sh -c "make clean; echo ;" \;
