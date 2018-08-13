include env.mk

all: clean protoc-build test

test:
	@ \
	success=0; \
	failure=0; \
	for i in `find */ -type f -name Makefile -not -path "vendor/*"`; do \
		cd $${i%/*}; \
		if make -C . test; then \
			success=$$((success+1)); \
		else \
			failure=$$((failure+1)); \
			echo "${RED}TEST FAILED${NC}"; \
		fi; \
		cd ..; \
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
