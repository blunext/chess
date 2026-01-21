#!/bin/bash
# Build specific commits by offset from HEAD
# Usage: ./scripts/build_commits.sh [offsets...]
# Examples:
#   ./scripts/build_commits.sh 0 1      # Build HEAD and HEAD~1
#   ./scripts/build_commits.sh 0 1 2 3  # Build last 4 commits

set -e

# Default: build HEAD and HEAD~1
offsets=(${@:-0 1})

# Create builds directory
mkdir -p builds/commits

# Save current state
current_ref=$(git rev-parse HEAD)
echo "Current HEAD: $current_ref"
echo ""

# Resolve all hashes BEFORE any checkout (checkout changes HEAD!)
declare -a hashes names msgs
for offset in ${offsets[@]}; do
  hash=$(git rev-parse HEAD~$offset 2>/dev/null) || {
    echo "Error: HEAD~$offset does not exist"
    exit 1
  }
  short_hash=${hash:0:7}
  msg=$(git log -1 --format=%s $hash | head -c 40)
  hashes+=("$hash")
  names+=("v$offset-$short_hash")
  msgs+=("$msg")
done

built=()

# Now build each commit
for i in ${!offsets[@]}; do
  offset=${offsets[$i]}
  hash=${hashes[$i]}
  name=${names[$i]}
  msg=${msgs[$i]}

  echo "Building $name (HEAD~$offset): $msg"
  if ! git checkout "$hash" --quiet; then
    echo "  -> FAILED: git checkout failed (uncommitted changes?)"
    git checkout "$current_ref" --quiet 2>/dev/null || true
    exit 1
  fi

  if go build -o "builds/commits/$name" .; then
    echo "  -> builds/commits/$name"
    built+=("$name")
  else
    echo "  -> FAILED"
  fi
done

# Return to original state
git checkout "$current_ref" --quiet 2>/dev/null
echo ""
echo "Built ${#built[@]} versions:"
ls -lh builds/commits/ 2>/dev/null | grep -E "$(IFS=\|; echo "${built[*]}")" || true
