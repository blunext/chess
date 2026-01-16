package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZobristHash_SamePositionSameHash(t *testing.T) {
	// Same FEN should produce same hash
	pos1 := CreatePositionFormFEN(InitialPosition)
	pos2 := CreatePositionFormFEN(InitialPosition)

	assert.Equal(t, pos1.Hash, pos2.Hash, "Same position should have same hash")
	assert.NotEqual(t, uint64(0), pos1.Hash, "Hash should not be zero")
}

func TestZobristHash_DifferentPositionsDifferentHash(t *testing.T) {
	// Different positions should have different hashes
	pos1 := CreatePositionFormFEN(InitialPosition)
	pos2 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1") // after e4

	assert.NotEqual(t, pos1.Hash, pos2.Hash, "Different positions should have different hashes")
}

func TestZobristHash_SideToMoveMakesHashDifferent(t *testing.T) {
	// Same board, but different side to move = different hash
	pos1 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	pos2 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")

	assert.NotEqual(t, pos1.Hash, pos2.Hash, "Different side to move should have different hashes")
}

func TestZobristHash_MakeUnmakeRestoresHash(t *testing.T) {
	// Hash should be restored after unmake
	pos := CreatePositionFormFEN(InitialPosition)
	originalHash := pos.Hash

	// Generate a legal move and play it
	moves := pos.GenerateLegalMoves(nil)
	if len(moves) == 0 {
		t.Fatal("No legal moves available")
	}

	// Make and unmake the move
	undo := pos.MakeMove(moves[0])
	pos.UnmakeMove(moves[0], undo)

	assert.Equal(t, originalHash, pos.Hash, "Hash should be restored after unmake")
}

func TestZobristHash_EnPassantAffectsHash(t *testing.T) {
	// Positions with different en passant should have different hashes
	pos1 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1") // e3 en passant
	pos2 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1")  // no en passant

	assert.NotEqual(t, pos1.Hash, pos2.Hash, "En passant should affect hash")
}

func TestZobristHash_CastlingRightsAffectHash(t *testing.T) {
	// Positions with different castling rights should have different hashes
	pos1 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	pos2 := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kq - 0 1") // no white queenside

	assert.NotEqual(t, pos1.Hash, pos2.Hash, "Castling rights should affect hash")
}
