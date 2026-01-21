#!/bin/bash
# Chain test commits by offset from HEAD
# Usage: ./scripts/chain_commits.sh [offsets...]
# Examples:
#   ./scripts/chain_commits.sh 0 1      # Test HEAD vs HEAD~1
#   ./scripts/chain_commits.sh 0 1 2    # Test HEAD vs HEAD~1, HEAD~1 vs HEAD~2

set -e

# macOS needs higher file descriptor limit for fastchess
ulimit -n 65536 2>/dev/null || true

# Default: test HEAD vs HEAD~1
offsets=(${@:-0 1})

# Test configuration
TC="${TC:-10+0.1}"
ROUNDS="${ROUNDS:-50}"
CONCURRENCY="${CONCURRENCY:-5}"

# Check fastchess
FASTCHESS="./tools/tournament/fastchess"
if [ ! -f "$FASTCHESS" ]; then
  echo "Error: fastchess not found at $FASTCHESS"
  exit 1
fi

echo "==================================================="
echo "CHAIN TEST: Recent Commits"
echo "==================================================="
echo "Offsets: ${offsets[*]}"
echo "Time control: $TC"
echo "Rounds: $ROUNDS (= $(($ROUNDS * 2)) games per match)"
echo "Concurrency: $CONCURRENCY"
echo "==================================================="
echo ""

# Build version names (same format as build_commits.sh)
versions=()
for offset in ${offsets[@]}; do
  hash=$(git rev-parse HEAD~$offset 2>/dev/null)
  short_hash=${hash:0:7}
  name="v$offset-$short_hash"

  if [ ! -f "builds/commits/$name" ]; then
    echo "Error: builds/commits/$name not found"
    echo "Run: make build-commits COMMITS='${offsets[*]}'"
    exit 1
  fi
  versions+=("$name")
done

# Run chain tests (newer vs older, so v0 vs v1, v1 vs v2, etc.)
num_versions=${#versions[@]}
for ((i = 0; i < num_versions - 1; i++)); do
  new="${versions[$i]}"
  old="${versions[$((i + 1))]}"

  echo ""
  echo "==================================================="
  echo "Match $((i + 1))/$((num_versions - 1)): $new vs $old"
  echo "==================================================="

  $FASTCHESS \
    -engine cmd=./builds/commits/$new name=$new \
    -engine cmd=./builds/commits/$old name=$old \
    -each tc=$TC \
    -rounds $ROUNDS \
    -concurrency $CONCURRENCY 2>&1 | grep --line-buffered -v "Warning; Last info string"

  echo ""
done

echo "==================================================="
echo "CHAIN TEST COMPLETE"
echo "==================================================="
