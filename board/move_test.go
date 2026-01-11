package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlidingMoves_Bishop(t *testing.T) {
	// Bishop at c1, free diagonal
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/2B5 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{
		IndexToBitBoard(2): [][]Bitboard{
			{IndexToBitBoard(11), IndexToBitBoard(20)}, // NE: d2, e3
			{IndexToBitBoard(9), IndexToBitBoard(16)},  // NW: b2, a3
		},
	}
	pm[Rook] = SquareMoves{}
	pm[Queen] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	assert.Len(t, moves, 4)

	// Verify all moves are from c1 (index 2)
	for _, m := range moves {
		assert.Equal(t, IndexToBitBoard(2), m.From)
		assert.Equal(t, Bishop, m.Piece)
		assert.Equal(t, Empty, m.Captured)
	}
}

func TestGenerateSlidingMoves_BlockedByPiece(t *testing.T) {
	// Rook at a1, pawn at a3 - rook should only reach a2
	position := CreatePositionFormFEN("8/8/8/8/8/P7/8/R7 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{}
	pm[Rook] = SquareMoves{
		IndexToBitBoard(0): [][]Bitboard{
			{IndexToBitBoard(8), IndexToBitBoard(16), IndexToBitBoard(24)}, // a2, a3, a4
		},
	}
	pm[Queen] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	assert.Len(t, moves, 1)
	assert.Equal(t, IndexToBitBoard(0), moves[0].From)
	assert.Equal(t, IndexToBitBoard(8), moves[0].To) // only a2
	assert.Equal(t, Rook, moves[0].Piece)
}

func TestGenerateSlidingMoves_AllSlidingPieces(t *testing.T) {
	// Position with bishop, rook, and queen
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/2BRQ3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{
		IndexToBitBoard(2): [][]Bitboard{
			{IndexToBitBoard(11)}, // one move
		},
	}
	pm[Rook] = SquareMoves{
		IndexToBitBoard(3): [][]Bitboard{
			{IndexToBitBoard(11)}, // one move
		},
	}
	pm[Queen] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			{IndexToBitBoard(12), IndexToBitBoard(20)}, // two moves
		},
	}

	moves := position.GenerateMoves(pm)

	// 1 bishop + 1 rook + 2 queen = 4 moves
	assert.Len(t, moves, 4)

	// Count pieces
	pieceCount := map[Piece]int{}
	for _, m := range moves {
		pieceCount[m.Piece]++
	}
	assert.Equal(t, 1, pieceCount[Bishop])
	assert.Equal(t, 1, pieceCount[Rook])
	assert.Equal(t, 2, pieceCount[Queen])
}

func TestMove_String(t *testing.T) {
	m := Move{
		From:     IndexToBitBoard(2),  // c1
		To:       IndexToBitBoard(11), // d2
		Piece:    Bishop,
		Captured: Empty,
	}

	assert.Equal(t, "Bishop: c1 -> d2", m.String())
}

func TestMove_String_Capture(t *testing.T) {
	m := Move{
		From:     IndexToBitBoard(2),  // c1
		To:       IndexToBitBoard(11), // d2
		Piece:    Bishop,
		Captured: Pawn,
	}

	assert.Equal(t, "Bishop: c1 x d2", m.String())
}

func TestGenerateJumpingMoves_Knight(t *testing.T) {
	// Knight at b1
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/1N6 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{}
	pm[Rook] = SquareMoves{}
	pm[Queen] = SquareMoves{}
	pm[Knight] = SquareMoves{
		IndexToBitBoard(1): [][]Bitboard{
			// b1 can jump to: a3(16), c3(18), d2(11)
			{IndexToBitBoard(16), IndexToBitBoard(18), IndexToBitBoard(11)},
		},
	}

	moves := position.GenerateMoves(pm)

	assert.Len(t, moves, 3)
	for _, m := range moves {
		assert.Equal(t, Knight, m.Piece)
		assert.Equal(t, IndexToBitBoard(1), m.From)
	}
}

func TestGenerateJumpingMoves_KnightBlockedByOwnPiece(t *testing.T) {
	// Knight at b1, own pawn at c3 - Knight can't land there
	position := CreatePositionFormFEN("8/8/8/8/8/2P5/8/1N6 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{}
	pm[Rook] = SquareMoves{}
	pm[Queen] = SquareMoves{}
	pm[Knight] = SquareMoves{
		IndexToBitBoard(1): [][]Bitboard{
			// b1 targets: a3(16), c3(18), d2(11)
			{IndexToBitBoard(16), IndexToBitBoard(18), IndexToBitBoard(11)},
		},
	}

	moves := position.GenerateMoves(pm)

	// Should have 2 moves (a3 and d2), not c3 (blocked by own pawn)
	assert.Len(t, moves, 2)

	// Verify c3 is not in the moves
	for _, m := range moves {
		assert.NotEqual(t, IndexToBitBoard(18), m.To, "Knight should not land on own piece at c3")
	}
}

func TestGenerateJumpingMoves_TwoKnights(t *testing.T) {
	// Two knights: b1 and g1
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/1N4N1 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{}
	pm[Rook] = SquareMoves{}
	pm[Queen] = SquareMoves{}
	pm[Knight] = SquareMoves{
		IndexToBitBoard(1): [][]Bitboard{
			{IndexToBitBoard(16), IndexToBitBoard(18)}, // b1: a3, c3
		},
		IndexToBitBoard(6): [][]Bitboard{
			{IndexToBitBoard(21), IndexToBitBoard(23)}, // g1: f3, h3
		},
	}

	moves := position.GenerateMoves(pm)

	// 2 moves from b1 + 2 moves from g1 = 4
	assert.Len(t, moves, 4)

	// Count moves per origin
	fromB1 := 0
	fromG1 := 0
	for _, m := range moves {
		if m.From == IndexToBitBoard(1) {
			fromB1++
		}
		if m.From == IndexToBitBoard(6) {
			fromG1++
		}
	}
	assert.Equal(t, 2, fromB1)
	assert.Equal(t, 2, fromG1)
}

func TestGenerateJumpingMoves_King(t *testing.T) {
	// King at e1
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/4K3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{}
	pm[Rook] = SquareMoves{}
	pm[Queen] = SquareMoves{}
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			// e1 can move to: d1(3), f1(5), d2(11), e2(12), f2(13)
			{IndexToBitBoard(3), IndexToBitBoard(5), IndexToBitBoard(11), IndexToBitBoard(12), IndexToBitBoard(13)},
		},
	}

	moves := position.GenerateMoves(pm)

	assert.Len(t, moves, 5)
	for _, m := range moves {
		assert.Equal(t, King, m.Piece)
		assert.Equal(t, IndexToBitBoard(4), m.From)
	}
}

func TestGenerateJumpingMoves_KingBlockedByOwnPieces(t *testing.T) {
	// King at e1, pawns on d2 and f2
	position := CreatePositionFormFEN("8/8/8/8/8/8/3PP3/4K3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Bishop] = SquareMoves{}
	pm[Rook] = SquareMoves{}
	pm[Queen] = SquareMoves{}
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			// e1 targets: d1(3), f1(5), d2(11), e2(12), f2(13)
			{IndexToBitBoard(3), IndexToBitBoard(5), IndexToBitBoard(11), IndexToBitBoard(12), IndexToBitBoard(13)},
		},
	}

	moves := position.GenerateMoves(pm)

	// Should have 3 moves (d1, f1, e2), not d2 or f2 (blocked by pawns)
	assert.Len(t, moves, 3)
}
