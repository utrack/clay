GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m

test:
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