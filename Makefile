FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

.PHONY: install
install: $(GO_BIN)
	go get

.PHONY: integration
integration: install
	$(MAKE) -C ./integration test
