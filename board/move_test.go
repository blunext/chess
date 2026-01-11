package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlidingMoves_Bishop(t *testing.T) {
	// Bishop at c1, free diagonals
	// c1 can go: NE diagonal (d2,e3,f4,g5,h6) + NW diagonal (b2,a3) = 7 squares
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/2B5 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Bishop from c1: NE (d2,e3,f4,g5,h6) = 5 + NW (b2,a3) = 2 => 7 moves
	assert.Len(t, moves, 7)

	// Verify all moves are from c1 (index 2)
	for _, m := range moves {
		assert.Equal(t, IndexToBitBoard(2), m.From)
		assert.Equal(t, Bishop, m.Piece)
		assert.Equal(t, Empty, m.Captured)
	}
}

func TestGenerateSlidingMoves_BlockedByPiece(t *testing.T) {
	// Rook at a1, own pawn at a3 - rook can reach a2 (before pawn), b1-h1 (horizontal)
	// Rook from a1: up (a2 only, blocked by a3) + right (b1,c1,d1,e1,f1,g1,h1)
	position := CreatePositionFormFEN("8/8/8/8/8/P7/8/R7 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Rook: 1 up (a2) + 7 right (b1-h1) = 8 moves
	assert.Len(t, moves, 8)

	// Check that a2 is reachable
	hasA2 := false
	for _, m := range moves {
		assert.Equal(t, IndexToBitBoard(0), m.From)
		assert.Equal(t, Rook, m.Piece)
		if m.To == IndexToBitBoard(8) {
			hasA2 = true
		}
	}
	assert.True(t, hasA2, "Rook should be able to move to a2")
}

func TestGenerateSlidingMoves_AllSlidingPieces(t *testing.T) {
	// Position with bishop at c1, rook at d1, queen at e1
	// They block each other along rank 1
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/2BRQ3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Count pieces - all three sliding pieces should have moves
	pieceCount := map[Piece]int{}
	for _, m := range moves {
		pieceCount[m.Piece]++
	}
	assert.Greater(t, pieceCount[Bishop], 0, "Bishop should have moves")
	assert.Greater(t, pieceCount[Rook], 0, "Rook should have moves")
	assert.Greater(t, pieceCount[Queen], 0, "Queen should have moves")
}

func TestGenerateSlidingMoves_Capture(t *testing.T) {
	// White rook at a1, black pawn at a5
	// Rook should be able to capture the pawn
	position := CreatePositionFormFEN("8/8/8/p7/8/8/8/R7 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Find the capture move (a1 -> a5)
	var captureMove *Move
	for i := range moves {
		if moves[i].To == IndexToBitBoard(32) { // a5
			captureMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, captureMove, "Should have a move to a5")
	assert.Equal(t, Rook, captureMove.Piece)
	assert.Equal(t, Pawn, captureMove.Captured, "Should capture the pawn")
}

func TestGenerateSlidingMoves_NoCapture(t *testing.T) {
	// White rook at a1, empty board - all moves should be non-captures
	position := CreatePositionFormFEN("8/8/8/8/8/8/8/R7 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	for _, m := range moves {
		assert.Equal(t, Empty, m.Captured, "Move to %s should not be a capture", m.To.Pretty())
	}
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
