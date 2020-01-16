all: test

include env.mk

SUBDIRS := $(wildcard */.)
SUBDIRS := $(filter-out bin/., $(SUBDIRS))
SUBCLEAN = $(addsuffix .clean,$(SUBDIRS))


test: protoc-build $(SUBDIRS)
clean: $(SUBCLEAN)


.PHONY: $(SUBDIRS)
$(SUBDIRS):
	cd $@ && $(MAKE) -C . test

.PHONY: $(SUBCLEAN)
$(SUBCLEAN): %.clean:
	cd $* && $(MAKE) -C . clean
