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

# Regression testing - build and test recent commits
# Usage:
#   make build-commits                    # Build HEAD and HEAD~1
#   make build-commits COMMITS="0 1 2"    # Build last 3 commits
#   make chain                            # Test HEAD vs HEAD~1
#   make chain COMMITS="0 1 2"            # Chain test last 3 commits
#   make regtest                          # Build + test HEAD vs HEAD~1
#   make regtest COMMITS="0 1 2 3"        # Build + chain test last 4

COMMITS ?= 0 1

build-commits:
	@echo ">> BUILD COMMITS: $(COMMITS)"
	@chmod +x scripts/build_commits.sh
	@./scripts/build_commits.sh $(COMMITS)

chain:
	@echo ">> CHAIN TEST: $(COMMITS)"
	@chmod +x scripts/chain_commits.sh
	@./scripts/chain_commits.sh $(COMMITS)

regtest: build-commits chain
	@echo ">> REGRESSION TEST COMPLETE"

# Build milestone versions (historical)
build-milestones:
	@echo ">> BUILD MILESTONES"
	@chmod +x build_milestones.sh
	@./build_milestones.sh

# Chain test milestones (historical)
chain-milestones:
	@echo ">> CHAIN TEST MILESTONES"
	@chmod +x chain_test.sh
	@./chain_test.sh

