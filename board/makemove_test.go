package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMove_SimplePush(t *testing.T) {
	pos := CreatePositionFormFEN("8/8/8/8/8/8/4P3/8 w - - 0 1")
	original := pos // copy for comparison

	m := Move{
		From:  IndexToBitBoard(12), // e2
		To:    IndexToBitBoard(20), // e3
		Piece: Pawn,
	}

	undo := pos.MakeMove(m)

	// Verify pawn moved
	assert.False(t, pos.Pawns.IsBitSet(12), "Pawn should not be on e2")
	assert.True(t, pos.Pawns.IsBitSet(20), "Pawn should be on e3")
	assert.False(t, pos.White.IsBitSet(12), "White should not be on e2")
	assert.True(t, pos.White.IsBitSet(20), "White should be on e3")

	// Verify side changed
	assert.False(t, pos.WhiteMove, "Should be black's move now")

	// Unmake and verify restoration
	pos.UnmakeMove(m, undo)
	assert.Equal(t, original.Pawns, pos.Pawns, "Pawns should be restored")
	assert.Equal(t, original.White, pos.White, "White should be restored")
	assert.Equal(t, original.WhiteMove, pos.WhiteMove, "Side should be restored")
}

func TestMakeMove_Capture(t *testing.T) {
	// White rook on e4, black pawn on e5
	pos := CreatePositionFormFEN("8/8/8/4p3/4R3/8/8/8 w - - 0 1")
	original := pos

	m := Move{
		From:     IndexToBitBoard(28), // e4
		To:       IndexToBitBoard(36), // e5
		Piece:    Rook,
		Captured: Pawn,
	}

	undo := pos.MakeMove(m)

	// Verify rook moved
	assert.False(t, pos.Rooks.IsBitSet(28), "Rook should not be on e4")
	assert.True(t, pos.Rooks.IsBitSet(36), "Rook should be on e5")

	// Verify pawn captured
	assert.False(t, pos.Pawns.IsBitSet(36), "Pawn should be captured")
	assert.False(t, pos.Black.IsBitSet(36), "Black piece should be removed")

	// Verify half-move clock reset
	assert.Equal(t, uint8(0), pos.HalfMoveClock, "Half-move clock should reset on capture")

	// Unmake
	pos.UnmakeMove(m, undo)
	assert.Equal(t, original.Rooks, pos.Rooks, "Rooks should be restored")
	assert.Equal(t, original.Pawns, pos.Pawns, "Pawns should be restored")
	assert.Equal(t, original.Black, pos.Black, "Black should be restored")
}

func TestMakeMove_DoublePawnPush(t *testing.T) {
	pos := CreatePositionFormFEN("8/8/8/8/8/8/4P3/8 w - - 5 1")

	m := Move{
		From:  IndexToBitBoard(12), // e2
		To:    IndexToBitBoard(28), // e4
		Piece: Pawn,
	}

	undo := pos.MakeMove(m)

	// Verify en passant square set
	assert.Equal(t, IndexToBitBoard(20), pos.EnPassant, "En passant should be e3")

	// Verify half-move clock reset
	assert.Equal(t, uint8(0), pos.HalfMoveClock, "Half-move clock should reset on pawn move")

	pos.UnmakeMove(m, undo)
	assert.Equal(t, uint8(5), pos.HalfMoveClock, "Half-move clock should be restored")
	assert.Equal(t, Bitboard(0), pos.EnPassant, "En passant should be cleared")
}

func TestMakeMove_EnPassant(t *testing.T) {
	// White pawn on e5, black pawn on d5 (just pushed), en passant on d6
	pos := CreatePositionFormFEN("8/8/8/3pP3/8/8/8/8 w - d6 0 1")
	original := pos

	m := Move{
		From:     IndexToBitBoard(36), // e5
		To:       IndexToBitBoard(43), // d6
		Piece:    Pawn,
		Captured: Pawn,
		Flags:    FlagEnPassant,
	}

	undo := pos.MakeMove(m)

	// Verify white pawn moved to d6
	assert.True(t, pos.Pawns.IsBitSet(43), "White pawn should be on d6")
	assert.True(t, pos.White.IsBitSet(43), "White should be on d6")

	// Verify black pawn on d5 captured (not d6!)
	assert.False(t, pos.Pawns.IsBitSet(35), "Black pawn on d5 should be captured")
	assert.False(t, pos.Black.IsBitSet(35), "Black should not be on d5")

	// Unmake
	pos.UnmakeMove(m, undo)
	assert.Equal(t, original.Pawns, pos.Pawns, "Pawns should be restored")
	assert.Equal(t, original.White, pos.White, "White should be restored")
	assert.Equal(t, original.Black, pos.Black, "Black should be restored")
}

func TestMakeMove_Promotion(t *testing.T) {
	// White pawn on e7
	pos := CreatePositionFormFEN("8/4P3/8/8/8/8/8/8 w - - 0 1")
	original := pos

	m := Move{
		From:      IndexToBitBoard(52), // e7
		To:        IndexToBitBoard(60), // e8
		Piece:     Pawn,
		Promotion: Queen,
	}

	undo := pos.MakeMove(m)

	// Verify pawn removed from e7
	assert.False(t, pos.Pawns.IsBitSet(52), "Pawn should not be on e7")

	// Verify queen added to e8
	assert.True(t, pos.Queens.IsBitSet(60), "Queen should be on e8")
	assert.True(t, pos.White.IsBitSet(60), "White should be on e8")

	// Pawn should NOT be on e8
	assert.False(t, pos.Pawns.IsBitSet(60), "Pawn should not be on e8 after promotion")

	// Unmake
	pos.UnmakeMove(m, undo)
	assert.Equal(t, original.Pawns, pos.Pawns, "Pawns should be restored")
	assert.Equal(t, original.Queens, pos.Queens, "Queens should be restored")
}

func TestMakeMove_CastlingKingside(t *testing.T) {
	pos := CreatePositionFormFEN("8/8/8/8/8/8/8/4K2R w K - 0 1")
	original := pos

	m := Move{
		From:  IndexToBitBoard(4), // e1
		To:    IndexToBitBoard(6), // g1
		Piece: King,
		Flags: FlagCastling,
	}

	undo := pos.MakeMove(m)

	// Verify king moved
	assert.True(t, pos.Kings.IsBitSet(6), "King should be on g1")
	assert.False(t, pos.Kings.IsBitSet(4), "King should not be on e1")

	// Verify rook moved
	assert.True(t, pos.Rooks.IsBitSet(5), "Rook should be on f1")
	assert.False(t, pos.Rooks.IsBitSet(7), "Rook should not be on h1")

	// Verify castling rights removed
	assert.Equal(t, uint8(0), pos.CastleSide&CastleWhiteKingSide, "White kingside should be gone")

	// Unmake
	pos.UnmakeMove(m, undo)
	assert.Equal(t, original.Kings, pos.Kings, "Kings should be restored")
	assert.Equal(t, original.Rooks, pos.Rooks, "Rooks should be restored")
	assert.Equal(t, original.CastleSide, pos.CastleSide, "Castling rights should be restored")
}

func TestMakeMove_CastlingQueenside(t *testing.T) {
	pos := CreatePositionFormFEN("8/8/8/8/8/8/8/R3K3 w Q - 0 1")

	m := Move{
		From:  IndexToBitBoard(4), // e1
		To:    IndexToBitBoard(2), // c1
		Piece: King,
		Flags: FlagCastling,
	}

	undo := pos.MakeMove(m)

	// Verify king moved
	assert.True(t, pos.Kings.IsBitSet(2), "King should be on c1")

	// Verify rook moved (a1 -> d1)
	assert.True(t, pos.Rooks.IsBitSet(3), "Rook should be on d1")
	assert.False(t, pos.Rooks.IsBitSet(0), "Rook should not be on a1")

	pos.UnmakeMove(m, undo)
	assert.True(t, pos.Rooks.IsBitSet(0), "Rook should be back on a1")
}

func TestMakeMove_RookCaptureRemovesCastlingRights(t *testing.T) {
	// White rook on a1, black rook can capture it
	pos := CreatePositionFormFEN("8/8/8/8/8/8/8/R3K2r b Qq - 0 1")

	m := Move{
		From:     IndexToBitBoard(7), // h1
		To:       IndexToBitBoard(0), // a1
		Piece:    Rook,
		Captured: Rook,
	}

	pos.MakeMove(m)

	// White queenside castling should be removed (rook captured)
	assert.Equal(t, uint8(0), pos.CastleSide&CastleWhiteQueenSide, "White queenside should be gone")
}

func TestMakeMove_PreservesPosition(t *testing.T) {
	// Complex position - make several moves and unmake them all
	pos := CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3")
	original := pos

	moves := []Move{
		{From: IndexToBitBoard(5), To: IndexToBitBoard(12), Piece: Bishop}, // Bf1-e2 (not really legal but tests mechanics)
	}

	var undos []UndoInfo
	for _, m := range moves {
		undos = append(undos, pos.MakeMove(m))
	}

	// Unmake in reverse order
	for i := len(moves) - 1; i >= 0; i-- {
		pos.UnmakeMove(moves[i], undos[i])
	}

	// Position should be exactly as before
	assert.Equal(t, original, pos, "Position should be fully restored")
}

// === GenerateLegalMoves Tests ===

func TestGenerateLegalMoves_FiltersIllegalMoves(t *testing.T) {
	// White king on e1, black rook on e8 - white is in check
	// Only legal moves are those that block or move the king out of check
	pos := CreatePositionFormFEN("4r3/8/8/8/8/8/8/4K3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			// e1 king can move to d1, f1, d2, e2, f2
			{IndexToBitBoard(3), IndexToBitBoard(5), IndexToBitBoard(11), IndexToBitBoard(12), IndexToBitBoard(13)},
		},
	}

	legal := pos.GenerateLegalMoves(pm)

	// King cannot stay on e-file (rook check) - d1, f1, d2, f2 are legal, e2 is not
	for _, m := range legal {
		if m.To == IndexToBitBoard(12) { // e2
			t.Error("e2 should be illegal - still on e-file check")
		}
	}

	// Should have 4 legal moves (d1, f1, d2, f2)
	assert.Len(t, legal, 4, "Should have 4 legal moves")
}

func TestGenerateLegalMoves_PinnedPiece(t *testing.T) {
	// White king on e1, white bishop on e4, black rook on e8
	// Bishop is pinned and cannot move (except along the pin line)
	pos := CreatePositionFormFEN("4r3/8/8/8/4B3/8/8/4K3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			{IndexToBitBoard(3), IndexToBitBoard(5), IndexToBitBoard(11), IndexToBitBoard(13)},
		},
	}

	legal := pos.GenerateLegalMoves(pm)

	// Bishop moves that leave the e-file should be illegal
	for _, m := range legal {
		if m.Piece == Bishop {
			toIdx := bitboardToIndex(m.To)
			// Bishop can only move along e-file (vertically) - but it's a bishop, so no valid moves
			// Actually bishop moves diagonally, so ALL bishop moves should be illegal (breaks pin)
			t.Errorf("Bishop move to %d should be illegal - pinned", toIdx)
		}
	}
}

func TestGenerateLegalMoves_KingCannotCaptureDefendedPiece(t *testing.T) {
	// White king on e1, black pawn on d2 defended by black bishop on a5
	pos := CreatePositionFormFEN("8/8/8/b7/8/8/3p4/4K3 w - - 0 1")

	pm := make(PieceMoves)
	pm[Knight] = SquareMoves{}
	pm[King] = SquareMoves{
		IndexToBitBoard(4): [][]Bitboard{
			{IndexToBitBoard(3), IndexToBitBoard(5), IndexToBitBoard(11), IndexToBitBoard(12), IndexToBitBoard(13)},
		},
	}

	legal := pos.GenerateLegalMoves(pm)

	// King cannot capture pawn on d2 (defended by bishop)
	for _, m := range legal {
		if m.To == IndexToBitBoard(11) { // d2
			t.Error("King cannot capture defended pawn on d2")
		}
	}
}
