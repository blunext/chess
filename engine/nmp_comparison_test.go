package engine

import (
	"fmt"
	"strings"
	"testing"

	"chess/board"
)

// TestNMPComparison compares search quality with and without Null Move Pruning
func TestNMPComparison(t *testing.T) {
	testCases := []struct {
		name     string
		fen      string
		depth    int
		expected string // expected best move (if known)
	}{
		{
			name:     "Mate in 1 (Qxf7#)",
			fen:      "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 0 1",
			depth:    5,
			expected: "h5f7",
		},
		{
			name:     "Capture hanging queen",
			fen:      "rnb1kbnr/pppppppp/8/8/3q4/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
			depth:    5,
			expected: "e3d4",
		},
		{
			name:     "Knight fork opportunity",
			fen:      "r1bqkb1r/pppp1ppp/2n2n2/4N3/4P3/8/PPPP1PPP/RNBQKB1R w KQkq - 0 1",
			depth:    5,
			expected: "", // Nxc6 or Nxf7 both good
		},
		{
			name:     "Kiwipete (complex)",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    6,
			expected: "", // multiple good moves
		},
		{
			name:     "Initial position",
			fen:      board.InitialPosition,
			depth:    6,
			expected: "", // many reasonable moves
		},
		{
			name:     "Endgame K+P vs K (NMP should be disabled)",
			fen:      "8/8/8/4k3/8/4K3/4P3/8 w - - 0 1",
			depth:    8,
			expected: "", // pushing pawn or king moves
		},
		{
			name:     "Middlegame with tactics",
			fen:      "r1bq1rk1/ppp2ppp/2n2n2/3pp3/1bPP4/2N1PN2/PP2BPPP/R1BQK2R w KQ - 0 1",
			depth:    6,
			expected: "", // complex position
		},
	}

	pm := createTestPieceMoves()

	fmt.Println("\n=== NMP Comparison Test ===")
	fmt.Println("Comparing search results WITH and WITHOUT Null Move Pruning")
	fmt.Printf("%-35s | %-5s | %-10s | %-7s | %-10s | %-7s | %-10s | %s\n",
		"Position", "Depth", "NMP Move", "Score", "No-NMP", "Score", "Nodes", "Verdict")
	fmt.Println(strings.Repeat("-", 120))

	allMatch := true
	totalSpeedup := 0.0
	count := 0

	// Save and disable TT for fair comparison
	oldTT := TT

	for _, tc := range testCases {
		pos := board.CreatePositionFormFEN(tc.fen)

		// Search WITH NMP (fresh TT)
		TT = NewTranspositionTable(16)
		UseNullMovePruning = true
		resultNMP := Search(pos, pm, tc.depth)

		// Search WITHOUT NMP (fresh TT)
		TT = NewTranspositionTable(16)
		UseNullMovePruning = false
		resultNoNMP := Search(pos, pm, tc.depth)

		// Restore NMP
		UseNullMovePruning = true

		// Calculate speedup
		speedup := float64(resultNoNMP.Nodes) / float64(resultNMP.Nodes)
		if resultNMP.Nodes > 0 {
			totalSpeedup += speedup
			count++
		}

		// Check if moves match and compare scores
		moveMatch := resultNMP.Move.ToUCI() == resultNoNMP.Move.ToUCI()
		scoreDiff := resultNMP.Score - resultNoNMP.Score

		verdict := "OK"
		if !moveMatch {
			if scoreDiff >= 0 {
				verdict = "DIFF (NMP same/better)"
			} else if scoreDiff > -50 {
				verdict = "DIFF (minor)"
			} else {
				verdict = "DIFF (NMP worse!)"
				allMatch = false
			}
			// Check if NMP found the expected move (if specified)
			if tc.expected != "" && resultNMP.Move.ToUCI() != tc.expected {
				allMatch = false
				verdict = "NMP WRONG!"
			}
		}

		fmt.Printf("%-35s | %-5d | %-10s | %-7d | %-10s | %-7d | %-10d | %s\n",
			tc.name, tc.depth,
			resultNMP.Move.ToUCI(), resultNMP.Score,
			resultNoNMP.Move.ToUCI(), resultNoNMP.Score,
			resultNMP.Nodes, verdict)

		// For positions with expected moves, verify NMP finds them
		if tc.expected != "" {
			if resultNMP.Move.ToUCI() != tc.expected {
				t.Errorf("%s: NMP found %s, expected %s", tc.name, resultNMP.Move.ToUCI(), tc.expected)
			}
		}
	}

	fmt.Println(strings.Repeat("-", 120))
	if count > 0 {
		fmt.Printf("Average speedup: %.2fx\n", totalSpeedup/float64(count))
	}

	// Restore original state
	TT = oldTT
	UseNullMovePruning = true

	if allMatch {
		fmt.Println("\n✓ NMP finds all expected tactical moves correctly")
	} else {
		fmt.Println("\n✗ NMP missed some expected moves - review needed!")
	}
}
