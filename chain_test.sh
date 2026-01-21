#!/bin/bash
# Chain testing in stages: v1 vs v2, v2 vs v3, etc.
# Results shown live - press Ctrl+C anytime to stop if you see problems

set -e

# macOS needs higher file descriptor limit for fastchess
ulimit -n 65536 2>/dev/null || true

# All versions (regressions removed, v01-quiescence skipped - no time mgmt)
versions=(
  "v01-iterative-deepening"
  "v02-tt"
  "v03-nmp-disabled"
  "v04-check-ext"
  "v05-king-safety"
  "v06-pawn-struct"
  "v07-better-book"
  "v08-mate-threat-disabled"
  "v09-tt-probe-fix"
  "v10-pesto"
  "v11-space-bonus"
  "v12-mobility"
  "v13-piece-dev"
  "v14-single-reply-ext"
  "v15-killers"
  "v16-delta-history"
)

# Test configuration
TC="10+0.1"
ROUNDS=50  # 50 rounds = 100 games (faster feedback)
CONCURRENCY=5

# Check if fastchess exists
if [ ! -f "./tools/tournament/fastchess" ]; then
  echo "Error: fastchess not found at ./tools/tournament/fastchess"
  exit 1
fi

# Parse command line arg for stage number
STAGE=${1:-"all"}

echo "==================================================="
echo "CHAIN REGRESSION TEST"
echo "==================================================="
echo "Time control: $TC"
echo "Rounds: $ROUNDS (= $(($ROUNDS * 2)) games per match)"
echo "Concurrency: $CONCURRENCY"
echo "Stage: $STAGE"
echo ""
echo "NOTE: Results update LIVE after each game!"
echo "Press Ctrl+C anytime if you see a problem."
echo "==================================================="
echo ""

get_stage_range() {
  case $1 in
    1) echo "0:4" ;;   # v01-v04: TT + NMP disable + check ext (3 matches)
    2) echo "4:7" ;;   # v05-v07: King safety + pawn eval (2 matches)
    3) echo "7:10" ;;  # v08-v10: Bug fixes + PeSTO (2 matches)
    4) echo "10:16" ;; # v11-v16: Advanced eval + killers (5 matches)
  esac
}

run_stage() {
  local stage_num=$1
  local range=$(get_stage_range $stage_num)
  local start=${range%:*}
  local end=${range#*:}

  echo ""
  echo "###################################################"
  echo "STAGE $stage_num: ${versions[$start]} → ${versions[$((end-1))]}"
  echo "###################################################"

  for i in $(seq $start $((end - 2))); do
    old="${versions[$i]}"
    new="${versions[$((i + 1))]}"
    match_num=$((i - start + 1))
    total_matches=$((end - start - 1))

    echo ""
    echo "==================================================="
    echo "Stage $stage_num, Match $match_num/$total_matches: $old vs $new"
    echo "==================================================="

    ./tools/tournament/fastchess \
      -engine cmd=./builds/milestones/$old name=$old \
      -engine cmd=./builds/milestones/$new name=$new \
      -each tc=$TC \
      -rounds $ROUNDS \
      -concurrency $CONCURRENCY 2>&1 | grep --line-buffered -v "Warning; Last info string with score not found"

    echo ""
    # Auto-continue (comment out for manual pause)
    # echo "Press Enter to continue (or Ctrl+C to stop)..."
    # read
  done
}

if [ "$STAGE" = "all" ]; then
  for s in 1 2 3 4; do
    run_stage $s
  done
else
  run_stage $STAGE
fi

echo ""
echo "==================================================="
echo "TESTS COMPLETE"
echo "=================================================="
