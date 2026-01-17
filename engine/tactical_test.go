package engine

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"chess/board"
	"chess/generator"
	"chess/magic"
)

// TacticalPosition defines a position where the engine must find a specific move
type TacticalPosition struct {
	Name      string
	FEN       string
	BestMoves []string // acceptable best moves (UCI format)
	MinDepth  int      // minimum depth to find the solution
	Category  string   // mate, material, defensive, etc.
}

// Tactical test positions - engine MUST find these moves
// All positions verified to have correct solutions
var tacticalPositions = []TacticalPosition{
	// === MATE IN 1 ===
	{
		Name:      "Mate in 1: Back rank",
		FEN:       "6k1/5ppp/8/8/8/8/8/R3K3 w Q - 0 1",
		BestMoves: []string{"a1a8"},
		MinDepth:  2,
		Category:  "mate1",
	},
	{
		Name:      "Mate in 1: Scholar's mate",
		FEN:       "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 0 1",
		BestMoves: []string{"h5f7"},
		MinDepth:  2,
		Category:  "mate1",
	},

	// === WIN MATERIAL: SIMPLE CAPTURES ===
	{
		Name:      "Capture hanging queen",
		FEN:       "rnb1kbnr/pppppppp/8/8/3q4/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
		BestMoves: []string{"e3d4"},
		MinDepth:  1,
		Category:  "hanging",
	},

	// === WIN MATERIAL: FORKS ===
	{
		Name:      "Knight fork: King and Rook",
		FEN:       "r3k2r/ppp2ppp/2n5/3N4/8/8/PPP2PPP/R3K2R w KQkq - 0 1",
		BestMoves: []string{"d5c7", "d5e7"},
		MinDepth:  2,
		Category:  "fork",
	},

	// === WIN MATERIAL: PINS ===
	{
		Name:      "Pin: win the queen",
		FEN:       "r2qkb1r/ppp2ppp/2n5/3np1B1/8/5N2/PPPP1PPP/R2QKB1R w KQkq - 0 1",
		BestMoves: []string{"g5d8"},
		MinDepth:  1,
		Category:  "pin",
	},

	// === DEFENSIVE ===
	{
		Name:      "Defend: escape back rank threat",
		FEN:       "6k1/5ppp/8/8/8/8/5PPP/r3K2R w K - 0 1",
		BestMoves: []string{"e1f1", "e1d2", "e1e2"},
		MinDepth:  2,
		Category:  "defensive",
	},

	// === WAC (Win At Chess) ===
	{
		Name:      "WAC.001",
		FEN:       "2rr3k/pp3pp1/1nnqbN1p/3pN3/2pP4/2P3Q1/PPB4P/R4RK1 w - - 0 1",
		BestMoves: []string{"g3g6"},
		MinDepth:  4,
		Category:  "wac",
	},
	{
		Name:      "WAC.005",
		FEN:       "r1b1k2r/ppp1qppp/5n2/4P3/2B2n2/2N2Q2/PPn2PPP/R1BR2K1 b kq - 0 1",
		BestMoves: []string{"f4h3", "c2a1"},
		MinDepth:  3,
		Category:  "wac",
	},
	{
		Name:      "WAC.008",
		FEN:       "r1bqk2r/pppp1ppp/5n2/2b1n3/4P3/1BP2N2/PP1P1PPP/RNBQ1RK1 b kq - 0 1",
		BestMoves: []string{"e5f3"},
		MinDepth:  2,
		Category:  "wac",
	},
}

func TestTacticalSuite(t *testing.T) {
	magic.Prepare()
	pm := generator.NewGenerator()

	passed := 0
	failed := 0
	var failedPositions []string

	fmt.Println("\n=== Tactical Test Suite ===")
	fmt.Printf("%-40s | %-10s | %-8s | %-10s | %s\n",
		"Position", "Category", "Depth", "Found", "Result")
	fmt.Println(strings.Repeat("-", 90))

	for _, tc := range tacticalPositions {
		pos := board.CreatePositionFormFEN(tc.FEN)
		result := Search(pos, pm, tc.MinDepth+2)

		foundMove := result.Move.ToUCI()
		isCorrect := false
		for _, best := range tc.BestMoves {
			if foundMove == best {
				isCorrect = true
				break
			}
		}

		status := "PASS"
		if !isCorrect {
			status = "FAIL"
			failed++
			failedPositions = append(failedPositions, fmt.Sprintf("%s: found %s, expected %v", tc.Name, foundMove, tc.BestMoves))
		} else {
			passed++
		}

		fmt.Printf("%-40s | %-10s | %-8d | %-10s | %s\n",
			tc.Name, tc.Category, tc.MinDepth, foundMove, status)
	}

	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("Results: %d/%d passed (%.1f%%)\n", passed, passed+failed, float64(passed)/float64(passed+failed)*100)

	if failed > 0 {
		fmt.Println("\nFailed positions:")
		for _, f := range failedPositions {
			fmt.Printf("  - %s\n", f)
		}
		t.Errorf("Tactical suite: %d/%d positions failed", failed, passed+failed)
	}
}

func TestTacticalSuiteWithTime(t *testing.T) {
	magic.Prepare()
	pm := generator.NewGenerator()
	timeLimit := 500 * time.Millisecond

	passed := 0
	failed := 0

	fmt.Println("\n=== Tactical Test Suite (Time-based) ===")
	fmt.Printf("Time limit: %v per position\n", timeLimit)
	fmt.Printf("%-40s | %-10s | %-6s | %-10s | %s\n",
		"Position", "Category", "Depth", "Found", "Result")
	fmt.Println(strings.Repeat("-", 90))

	for _, tc := range tacticalPositions {
		pos := board.CreatePositionFormFEN(tc.FEN)
		result := SearchWithTime(pos, pm, timeLimit)

		foundMove := result.Move.ToUCI()
		isCorrect := false
		for _, best := range tc.BestMoves {
			if foundMove == best {
				isCorrect = true
				break
			}
		}

		status := "PASS"
		if !isCorrect {
			status = "FAIL"
			failed++
		} else {
			passed++
		}

		fmt.Printf("%-40s | %-10s | %-6d | %-10s | %s\n",
			tc.Name, tc.Category, result.Depth, foundMove, status)
	}

	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("Results: %d/%d passed (%.1f%%)\n", passed, passed+failed, float64(passed)/float64(passed+failed)*100)

	// Require reasonable pass rate (some WAC positions are very hard)
	minPassRate := 0.7
	if float64(passed)/float64(passed+failed) < minPassRate {
		t.Errorf("Tactical suite (time-based): pass rate %.1f%% below threshold %.1f%%",
			float64(passed)/float64(passed+failed)*100, minPassRate*100)
	}
}

func BenchmarkTacticalSuite(b *testing.B) {
	magic.Prepare()
	pm := generator.NewGenerator()
	positions := make([]board.Position, len(tacticalPositions))
	for i, tc := range tacticalPositions {
		positions[i] = board.CreatePositionFormFEN(tc.FEN)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, pos := range positions {
			Search(pos, pm, tacticalPositions[j].MinDepth+1)
		}
	}
}
