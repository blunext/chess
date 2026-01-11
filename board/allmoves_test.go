package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock PieceMoves for testing - simplified move generator
func createMockPieceMoves() PieceMoves {
	pm := make(PieceMoves)

	// Simple knight moves from a1 (0)
	pm[Knight] = SquareMoves{
		IndexToBitBoard(0): [][]Bitboard{
			{IndexToBitBoard(17), IndexToBitBoard(10)}, // Knight moves from a1
		},
	}

	// Simple rook moves (just up direction for simplicity)
	pm[Rook] = SquareMoves{
		IndexToBitBoard(0): [][]Bitboard{
			{IndexToBitBoard(8), IndexToBitBoard(16), IndexToBitBoard(24)}, // a1 -> a2, a3, a4
		},
	}

	// Simple bishop moves (diagonal)
	pm[Bishop] = SquareMoves{
		IndexToBitBoard(0): [][]Bitboard{
			{IndexToBitBoard(9), IndexToBitBoard(18), IndexToBitBoard(27)}, // a1 -> b2, c3, d4
		},
		IndexToBitBoard(2): [][]Bitboard{
			{IndexToBitBoard(11), IndexToBitBoard(20), IndexToBitBoard(29)}, // c1 -> d2, e3, f4
		},
	}

	return pm
}

func TestAllLegalMoves_NoPieces(t *testing.T) {
	// Position with no bishops - should return nil
	position := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	pm := createMockPieceMoves()

	// Filter to white pieces only (bishops at c1, f1)
	result := position.AllLegalMoves(pm, Queen) // No white queens at start

	// Should return nil when no pieces of that type exist
	assert.Nil(t, result)
}

func TestAllLegalMoves_Knight(t *testing.T) {
	// Knight at b1
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/1N6 w - - 0 1")
	pm := createMockPieceMoves()

	// Override with actual knight at b1 (1)
	pm[Knight] = SquareMoves{
		IndexToBitBoard(1): [][]Bitboard{
			{IndexToBitBoard(18), IndexToBitBoard(11), IndexToBitBoard(16)}, // b1 -> c3, d2, a3
		},
	}

	result := position.AllLegalMoves(pm, Knight)

	// Should have 3 moves
	assert.NotNil(t, result)
	assert.Len(t, result, 3)

	// Verify moves are to empty squares
	for _, pos := range result {
		assert.NotEqual(t, Bitboard(0), pos.Knights)
	}
}

func TestAllLegalMoves_StopsAtPiece(t *testing.T) {
	// Rook at a1, pawn at a3 - rook should only reach a2
	position := CreatePositionFormFEN("8/8/8/8/8/P7/8/R7 w - - 0 1")

	pm := make(PieceMoves)
	pm[Rook] = SquareMoves{
		IndexToBitBoard(0): [][]Bitboard{
			{IndexToBitBoard(8), IndexToBitBoard(16), IndexToBitBoard(24), IndexToBitBoard(32)}, // a2, a3, a4, a5
		},
	}

	result := position.AllLegalMoves(pm, Rook)

	assert.NotNil(t, result)
	// Should stop at a2 (index 8) before hitting pawn at a3 (index 16)
	assert.Len(t, result, 1)

	// Verify rook moved to a2
	assert.Equal(t, IndexToBitBoard(8), result[0].Rooks)
}

func TestAllLegalMoves_MultipleDirections(t *testing.T) {
	// Bishop at c1, free diagonal
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/2B5 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{
		IndexToBitBoard(2): [][]Bitboard{
			{IndexToBitBoard(11), IndexToBitBoard(20)}, // NE: d2, e3
			{IndexToBitBoard(9), IndexToBitBoard(16)},  // NW: b2, a3
		},
	}

	result := position.AllLegalMoves(pm, Bishop)

	assert.NotNil(t, result)
	// 2 directions, 2 moves each = 4 total
	assert.Len(t, result, 4)
}

func TestAllLegalMoves_BlockedByOwnPiece(t *testing.T) {
	// Rook at a1, own pawn at a2 - no moves up
	position := CreatePositionFormFEN("8/8/8/8/8/8/P7/R7 w - - 0 1")

	pm := make(PieceMoves)
	pm[Rook] = SquareMoves{
		IndexToBitBoard(0): [][]Bitboard{
			{IndexToBitBoard(8), IndexToBitBoard(16)}, // a2, a3
		},
	}

	result := position.AllLegalMoves(pm, Rook)

	// Should return empty list - blocked immediately by own pawn
	// (not nil, because rook exists, just has no legal moves)
	assert.Empty(t, result)
}

func TestAllLegalMoves_BlackToMove(t *testing.T) {
	// Black rook at a8
	position := CreatePositionFormFEN("r7/8/8/8/8/8/8/8 b - - 0 1")

	pm := make(PieceMoves)
	pm[Rook] = SquareMoves{
		IndexToBitBoard(56): [][]Bitboard{
			{IndexToBitBoard(48), IndexToBitBoard(40)}, // a7, a6
		},
	}

	result := position.AllLegalMoves(pm, Rook)

	assert.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestAllLegalMoves_EmptyBoard(t *testing.T) {
	// Completely empty board, white to move
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/8 w - - 0 1")
	pm := createMockPieceMoves()

	result := position.AllLegalMoves(pm, Bishop)

	// No bishops on board
	assert.Nil(t, result)
}

func TestAllLegalMoves_ComplexPosition(t *testing.T) {
	// Real position with open bishop
	position := CreatePositionFormFEN("rnbqkbnr/pppp1ppp/4p3/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1")

	pm := make(PieceMoves)
	// White bishops at c1 and f1
	// c1 bishop can move to d2, e3, f4, g5, h6
	pm[Bishop] = SquareMoves{
		IndexToBitBoard(2): [][]Bitboard{
			{IndexToBitBoard(11), IndexToBitBoard(20), IndexToBitBoard(29), IndexToBitBoard(38), IndexToBitBoard(47)}, // NE diagonal
		},
		IndexToBitBoard(5): [][]Bitboard{
			{IndexToBitBoard(12), IndexToBitBoard(19)}, // f1 has limited moves due to pawns
		},
	}

	result := position.AllLegalMoves(pm, Bishop)

	assert.NotNil(t, result)
	// c1 bishop: 5 moves, f1 bishop blocked by e2 pawn at first square
	// Should have moves from c1 only
	assert.Greater(t, len(result), 0)
}

func TestIndexToBitBoard(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		expected Bitboard
	}{
		{"a1", 0, 0x1},
		{"b1", 1, 0x2},
		{"a2", 8, 0x100},
		{"h1", 7, 0x80},
		{"a8", 56, 0x100000000000000},
		{"h8", 63, 0x8000000000000000},
		{"d4", 27, 0x8000000},
		{"e5", 36, 0x1000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IndexToBitBoard(tt.index)
			assert.Equal(t, tt.expected, result)
		})
	}
}
