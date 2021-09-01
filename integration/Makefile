all: test

include env.mk

SUBDIRS := $(wildcard */)
SUBDIRS := $(filter-out bin/ third_party/, $(SUBDIRS))
SUBCLEAN = $(addsuffix .clean,$(SUBDIRS))


test: protoc-build $(SUBDIRS)
clean: $(SUBCLEAN)
	rm -rf ./bin

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(info #Running $@...)
	cd $@ && $(MAKE) -C . test

.PHONY: $(SUBCLEAN)
$(SUBCLEAN): %.clean:
	cd $* && $(MAKE) -C . clean
