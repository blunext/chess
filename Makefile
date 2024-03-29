PKGS             := $(shell go list ./...)

all: codequality test

test:
	@echo ">> TEST, \"full-mode\": race detector on"
	@$(foreach pkg, $(PKGS),\
		echo -n "     ";\
		go test -run '(Test|Example)' -race $(pkg) || exit 1;\
		)

codequality:
	@echo ">> CODE QUALITY"

	@echo -n "     GOLANGCI-LINTERS \n"
	@golangci-lint -v run ./...
	@$(call ok)

	@echo -n "     REVIVE"
	@revive -config revive.toml -formatter friendly -exclude vendor/... ./...
	@$(call ok)
