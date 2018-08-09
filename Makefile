FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

VGO_PATH:=$(FIRST_GOPATH)/src/golang.org/x/vgo
VGO_VERSION:=master
VGO_BIN:=$(GOBIN)/vgo

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

.PHONY: install
install: $(VGO_BIN)
	vgo mod vendor

.PHONY: integration
integration: install
	$(MAKE) -C ./integration test
