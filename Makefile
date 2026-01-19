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

build:
	@echo ">> BUILD ENGINE"
	go build -o chessengine .

tournament:
	@echo ">> BUILD TOURNAMENT TOOL"
	go build -o tournament ./tools/tournament

match: build tournament
	@echo ">> RUNNING SELF-PLAY MATCH"
	./tournament -engine1 ./chessengine -engine2 ./chessengine -games 10 -tc "1+0.01"

