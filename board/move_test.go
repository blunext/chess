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

	// Filter only rook moves
	var rookMoves []Move
	for _, m := range moves {
		if m.Piece == Rook {
			rookMoves = append(rookMoves, m)
		}
	}

	// Rook: 1 up (a2) + 7 right (b1-h1) = 8 moves
	assert.Len(t, rookMoves, 8)

	// Check that a2 is reachable
	hasA2 := false
	for _, m := range rookMoves {
		assert.Equal(t, IndexToBitBoard(0), m.From)
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

func TestMove_String_Promotion(t *testing.T) {
	m := Move{
		From:      IndexToBitBoard(52), // e7
		To:        IndexToBitBoard(60), // e8
		Piece:     Pawn,
		Promotion: Queen,
	}

	assert.Equal(t, "Pawn: e7 -> e8=Queen", m.String())
}

func TestMove_String_EnPassant(t *testing.T) {
	m := Move{
		From:     IndexToBitBoard(36), // e5
		To:       IndexToBitBoard(43), // d6
		Piece:    Pawn,
		Captured: Pawn,
		Flags:    FlagEnPassant,
	}

	assert.Equal(t, "Pawn: e5 x d6 e.p.", m.String())
}

func TestMove_String_Castling(t *testing.T) {
	m := Move{
		From:  IndexToBitBoard(4), // e1
		To:    IndexToBitBoard(6), // g1
		Piece: King,
		Flags: FlagCastling,
	}

	assert.Equal(t, "King: e1 -> g1 (castling)", m.String())
}

func TestMove_ToUCI(t *testing.T) {
	tests := []struct {
		name     string
		move     Move
		expected string
	}{
		{
			name: "simple move",
			move: Move{
				From:  IndexToBitBoard(12), // e2
				To:    IndexToBitBoard(28), // e4
				Piece: Pawn,
			},
			expected: "e2e4",
		},
		{
			name: "promotion to queen",
			move: Move{
				From:      IndexToBitBoard(52), // e7
				To:        IndexToBitBoard(60), // e8
				Piece:     Pawn,
				Promotion: Queen,
			},
			expected: "e7e8q",
		},
		{
			name: "promotion to knight",
			move: Move{
				From:      IndexToBitBoard(49), // b7
				To:        IndexToBitBoard(57), // b8
				Piece:     Pawn,
				Promotion: Knight,
			},
			expected: "b7b8n",
		},
		{
			name: "castling kingside",
			move: Move{
				From:  IndexToBitBoard(4), // e1
				To:    IndexToBitBoard(6), // g1
				Piece: King,
				Flags: FlagCastling,
			},
			expected: "e1g1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.move.ToUCI())
		})
	}
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

	// Filter only knight moves
	var knightMoves []Move
	for _, m := range moves {
		if m.Piece == Knight {
			knightMoves = append(knightMoves, m)
		}
	}

	// Should have 2 moves (a3 and d2), not c3 (blocked by own pawn)
	assert.Len(t, knightMoves, 2)

	// Verify c3 is not in the knight moves
	for _, m := range knightMoves {
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

	// Filter only king moves
	var kingMoves []Move
	for _, m := range moves {
		if m.Piece == King {
			kingMoves = append(kingMoves, m)
		}
	}

	// Should have 3 moves (d1, f1, e2), not d2 or f2 (blocked by pawns)
	assert.Len(t, kingMoves, 3)
}

func TestGenerateJumpingMoves_KnightCapture(t *testing.T) {
	// White knight at b1, black pawn at c3
	position := CreatePositionFormFEN("8/8/8/8/8/2p5/8/1N6 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{
		IndexToBitBoard(1): [][]Bitboard{
			{IndexToBitBoard(16), IndexToBitBoard(18), IndexToBitBoard(11)}, // a3, c3, d2
		},
	}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Find capture move to c3
	var captureMove *Move
	for i := range moves {
		if moves[i].To == IndexToBitBoard(18) { // c3
			captureMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, captureMove, "Should have move to c3")
	assert.Equal(t, Knight, captureMove.Piece)
	assert.Equal(t, Pawn, captureMove.Captured, "Should capture the pawn")
}

func TestGenerateJumpingMoves_KingCapture(t *testing.T) {
	// White king at e1, black pawn at d2
	position := CreatePositionFormFEN("8/8/8/8/8/8/3p4/4K3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			{IndexToBitBoard(3), IndexToBitBoard(5), IndexToBitBoard(11), IndexToBitBoard(12), IndexToBitBoard(13)},
		},
	}

	moves := position.GenerateMoves(pm)

	// Find capture move to d2
	var captureMove *Move
	for i := range moves {
		if moves[i].To == IndexToBitBoard(11) { // d2
			captureMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, captureMove, "Should have move to d2")
	assert.Equal(t, King, captureMove.Piece)
	assert.Equal(t, Pawn, captureMove.Captured, "Should capture the pawn")
}

// === Pawn Tests ===

func TestPawnMoves_SinglePush(t *testing.T) {
	// White pawn at e2
	position := CreatePositionFormFEN("8/8/8/8/8/8/4P3/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// e2 can go to e3 and e4 (double push from start)
	assert.Len(t, pawnMoves, 2)
}

func TestPawnMoves_DoublePush(t *testing.T) {
	// White pawn at e2, should be able to push to e4
	position := CreatePositionFormFEN("8/8/8/8/8/8/4P3/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Find double push to e4 (index 28)
	var doublePush *Move
	for i := range moves {
		if moves[i].Piece == Pawn && moves[i].To == IndexToBitBoard(28) {
			doublePush = &moves[i]
			break
		}
	}

	assert.NotNil(t, doublePush, "Should have double push to e4")
	assert.Equal(t, IndexToBitBoard(12), doublePush.From, "Should be from e2")
}

func TestPawnMoves_BlockedPush(t *testing.T) {
	// White pawn at e2, blocked by piece at e3
	position := CreatePositionFormFEN("8/8/8/8/8/4p3/4P3/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves - should have 0 (blocked)
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	assert.Len(t, pawnMoves, 0, "Pawn should be blocked")
}

func TestPawnMoves_Capture(t *testing.T) {
	// White pawn at e4, black pawns at d5 and f5
	position := CreatePositionFormFEN("8/8/8/3p1p2/4P3/8/8/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// e4 can: push to e5, capture d5, capture f5 = 3 moves
	assert.Len(t, pawnMoves, 3)

	// Check captures have Captured set
	captureCount := 0
	for _, m := range pawnMoves {
		if m.Captured == Pawn {
			captureCount++
		}
	}
	assert.Equal(t, 2, captureCount, "Should have 2 captures")
}

func TestPawnMoves_BlackPawn(t *testing.T) {
	// Black pawn at e7
	position := CreatePositionFormFEN("8/4p3/8/8/8/8/8/8 b - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// e7 can go to e6 and e5 (double push from start)
	assert.Len(t, pawnMoves, 2)

	// Verify moves are downward (lower indices)
	for _, m := range pawnMoves {
		assert.Less(t, int(m.To), int(m.From), "Black pawn should move down")
	}
}

func TestPawnMoves_NoWrapAround(t *testing.T) {
	// White pawn at a4 - should not capture on h-file
	position := CreatePositionFormFEN("8/8/8/8/P7/8/8/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// a4 can only push to a5, no captures possible
	assert.Len(t, pawnMoves, 1)
	assert.Equal(t, IndexToBitBoard(32), pawnMoves[0].To, "Should push to a5")
}

// === Pawn Promotion Tests ===

func TestPawnMoves_WhitePromotion(t *testing.T) {
	// White pawn at e7 can promote
	position := CreatePositionFormFEN("8/4P3/8/8/8/8/8/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// e7 -> e8 should generate 4 promotion moves (Q, R, B, N)
	assert.Len(t, pawnMoves, 4, "Should have 4 promotion moves")

	// Verify all promotions
	promotions := make(map[Piece]bool)
	for _, m := range pawnMoves {
		assert.Equal(t, IndexToBitBoard(52), m.From, "Should be from e7")
		assert.Equal(t, IndexToBitBoard(60), m.To, "Should be to e8")
		assert.NotEqual(t, Empty, m.Promotion, "Should have promotion piece")
		promotions[m.Promotion] = true
	}

	assert.True(t, promotions[Queen], "Should have Queen promotion")
	assert.True(t, promotions[Rook], "Should have Rook promotion")
	assert.True(t, promotions[Bishop], "Should have Bishop promotion")
	assert.True(t, promotions[Knight], "Should have Knight promotion")
}

func TestPawnMoves_WhitePromotionCapture(t *testing.T) {
	// White pawn at e7, black rook at f8
	position := CreatePositionFormFEN("5r2/4P3/8/8/8/8/8/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// 4 promotions forward (e7-e8) + 4 capture promotions (e7xf8) = 8
	assert.Len(t, pawnMoves, 8, "Should have 8 promotion moves")

	// Count captures
	captureCount := 0
	for _, m := range pawnMoves {
		if m.Captured == Rook {
			captureCount++
		}
	}
	assert.Equal(t, 4, captureCount, "Should have 4 capture promotions")
}

func TestPawnMoves_BlackPromotion(t *testing.T) {
	// Black pawn at e2 can promote
	position := CreatePositionFormFEN("8/8/8/8/8/8/4p3/8 b - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Filter pawn moves
	var pawnMoves []Move
	for _, m := range moves {
		if m.Piece == Pawn {
			pawnMoves = append(pawnMoves, m)
		}
	}

	// e2 -> e1 should generate 4 promotion moves (Q, R, B, N)
	assert.Len(t, pawnMoves, 4, "Should have 4 promotion moves")

	// Verify all promotions
	promotions := make(map[Piece]bool)
	for _, m := range pawnMoves {
		assert.Equal(t, IndexToBitBoard(12), m.From, "Should be from e2")
		assert.Equal(t, IndexToBitBoard(4), m.To, "Should be to e1")
		promotions[m.Promotion] = true
	}

	assert.True(t, promotions[Queen], "Should have Queen promotion")
	assert.True(t, promotions[Rook], "Should have Rook promotion")
	assert.True(t, promotions[Bishop], "Should have Bishop promotion")
	assert.True(t, promotions[Knight], "Should have Knight promotion")
}

// === En Passant Tests ===

func TestPawnMoves_EnPassantWhite(t *testing.T) {
	// White pawn at e5, black pawn just played d7-d5 (en passant on d6)
	position := CreatePositionFormFEN("8/8/8/3pP3/8/8/8/8 w - d6 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Find en passant move
	var epMove *Move
	for i := range moves {
		if moves[i].Flags&FlagEnPassant != 0 {
			epMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, epMove, "Should have en passant move")
	assert.Equal(t, IndexToBitBoard(36), epMove.From, "Should be from e5")
	assert.Equal(t, IndexToBitBoard(43), epMove.To, "Should be to d6")
	assert.Equal(t, Pawn, epMove.Captured, "Should capture pawn")
	assert.Equal(t, Pawn, epMove.Piece, "Should be pawn move")
}

func TestPawnMoves_EnPassantBlack(t *testing.T) {
	// Black pawn at d4, white pawn just played e2-e4 (en passant on e3)
	position := CreatePositionFormFEN("8/8/8/8/3pP3/8/8/8 b - e3 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Find en passant move
	var epMove *Move
	for i := range moves {
		if moves[i].Flags&FlagEnPassant != 0 {
			epMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, epMove, "Should have en passant move")
	assert.Equal(t, IndexToBitBoard(27), epMove.From, "Should be from d4")
	assert.Equal(t, IndexToBitBoard(20), epMove.To, "Should be to e3")
	assert.Equal(t, Pawn, epMove.Captured, "Should capture pawn")
}

func TestPawnMoves_NoEnPassantWithoutFlag(t *testing.T) {
	// Position without en passant flag - no en passant should be generated
	position := CreatePositionFormFEN("8/8/8/3pP3/8/8/8/8 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{}

	moves := position.GenerateMoves(pm)

	// Check no en passant moves
	for _, m := range moves {
		assert.Equal(t, FlagNone, m.Flags&FlagEnPassant, "Should not have en passant flag")
	}
}

// === Castling Move Generation Tests ===

func TestCastlingMoves_WhiteKingSide(t *testing.T) {
	// White can castle kingside: king on e1, rook on h1, f1 and g1 empty
	position := CreatePositionFormFEN("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{}, // e1 king - no normal moves for this test
	}

	moves := position.GenerateMoves(pm)

	// Find kingside castling move
	var castleMove *Move
	for i := range moves {
		if moves[i].Flags&FlagCastling != 0 && moves[i].To == IndexToBitBoard(6) {
			castleMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, castleMove, "Should have kingside castling move")
	assert.Equal(t, IndexToBitBoard(4), castleMove.From, "Should be from e1")
	assert.Equal(t, IndexToBitBoard(6), castleMove.To, "Should be to g1")
	assert.Equal(t, King, castleMove.Piece, "Should be king move")
}

func TestCastlingMoves_WhiteQueenSide(t *testing.T) {
	// White can castle queenside
	position := CreatePositionFormFEN("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{},
	}

	moves := position.GenerateMoves(pm)

	// Find queenside castling move
	var castleMove *Move
	for i := range moves {
		if moves[i].Flags&FlagCastling != 0 && moves[i].To == IndexToBitBoard(2) {
			castleMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, castleMove, "Should have queenside castling move")
	assert.Equal(t, IndexToBitBoard(4), castleMove.From, "Should be from e1")
	assert.Equal(t, IndexToBitBoard(2), castleMove.To, "Should be to c1")
}

func TestCastlingMoves_BlackKingSide(t *testing.T) {
	// Black can castle kingside
	position := CreatePositionFormFEN("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(60): [][]Bitboard{},
	}

	moves := position.GenerateMoves(pm)

	// Find kingside castling move
	var castleMove *Move
	for i := range moves {
		if moves[i].Flags&FlagCastling != 0 && moves[i].To == IndexToBitBoard(62) {
			castleMove = &moves[i]
			break
		}
	}

	assert.NotNil(t, castleMove, "Should have kingside castling move")
	assert.Equal(t, IndexToBitBoard(60), castleMove.From, "Should be from e8")
	assert.Equal(t, IndexToBitBoard(62), castleMove.To, "Should be to g8")
}

func TestCastlingMoves_BlockedByPiece(t *testing.T) {
	// White kingside blocked by bishop on f1
	position := CreatePositionFormFEN("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3KB1R w KQkq - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{},
	}

	moves := position.GenerateMoves(pm)

	// Should NOT have kingside castling
	for _, m := range moves {
		if m.Flags&FlagCastling != 0 && m.To == IndexToBitBoard(6) {
			t.Error("Should not have kingside castling when blocked")
		}
	}
}

func TestCastlingMoves_NoRights(t *testing.T) {
	// Position without castling rights
	position := CreatePositionFormFEN("r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{},
	}

	moves := position.GenerateMoves(pm)

	// Should have NO castling moves
	for _, m := range moves {
		assert.Equal(t, FlagNone, m.Flags&FlagCastling, "Should not have castling flag")
	}
}
