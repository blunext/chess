package generator

import (
	"testing"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

var (
	testRookMoves   = generateRookMoves()
	testBishopMoves = generateBishopMoves()
	testQueenMoves  = generateQueenMoves(testRookMoves, testBishopMoves)
)

func TestQueenCombinesRookAndBishop(t *testing.T) {
	// Test that queen has all rook moves + all bishop moves for every position
	for pos := 0; pos < 64; pos++ {
		bb := board.IndexToBitBoard(pos)

		rMoves := testRookMoves[bb]
		bMoves := testBishopMoves[bb]
		qMoves := testQueenMoves[bb]

		// Queen should have rook directions + bishop directions
		expectedDirs := len(rMoves) + len(bMoves)
		assert.Equal(t, expectedDirs, len(qMoves), "pos %d should combine rook and bishop directions", pos)

		// Count moves - queen should have rook + bishop move count
		rCount := countMoves(rMoves)
		bCount := countMoves(bMoves)
		qCount := countMoves(qMoves)

		assert.Equal(t, rCount+bCount, qCount, "pos %d queen moves should equal rook + bishop", pos)

		// Verify each move is a valid single-bit bitboard
		for _, dir := range qMoves {
			for _, move := range dir {
				assert.Equal(t, 1, popCount(move), "each move should have exactly one bit set")
			}
		}
	}
}

func TestQueenSpecialPositions(t *testing.T) {
	tests := []struct {
		name string
		pos  int
		desc string
	}{
		{"corner a1", 0, "should have 3 directions (up, right, NE)"},
		{"corner h1", 7, "should have 3 directions (up, left, NW)"},
		{"corner a8", 56, "should have 3 directions (down, right, SE)"},
		{"corner h8", 63, "should have 3 directions (down, left, SW)"},
		{"center d4", 27, "should have 8 directions (all)"},
		{"center e4", 28, "should have 8 directions (all)"},
		{"center d5", 35, "should have 8 directions (all)"},
		{"center e5", 36, "should have 8 directions (all)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := board.IndexToBitBoard(tt.pos)
			moves := testQueenMoves[bb]

			// Corners have 3 directions, center has 8
			if tt.pos == 0 || tt.pos == 7 || tt.pos == 56 || tt.pos == 63 {
				assert.Equal(t, 3, len(moves), tt.desc)
				assert.Equal(t, 21, countMoves(moves), "corners should have 21 total moves")
			} else if tt.pos >= 27 && tt.pos <= 36 && (tt.pos%8 >= 3 && tt.pos%8 <= 4) {
				// Center 4 squares (d4, e4, d5, e5)
				assert.Equal(t, 8, len(moves), tt.desc)
				assert.GreaterOrEqual(t, countMoves(moves), 25, "center should have at least 25 moves")
			}
		})
	}
}

// Helper: count total moves across all directions
func countMoves(directions [][]board.Bitboard) int {
	total := 0
	for _, dir := range directions {
		total += len(dir)
	}
	return total
}

// Helper: count number of bits set in a bitboard (population count)
func popCount(bb board.Bitboard) int {
	count := 0
	for bb != 0 {
		count++
		bb &= bb - 1 // Clear lowest set bit
	}
	return count
}
