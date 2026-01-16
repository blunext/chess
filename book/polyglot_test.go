package book

import (
	"testing"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

func TestPolyglotHash_InitialPosition(t *testing.T) {
	// The initial position should have a known Polyglot hash
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	hash := PolyglotHash(pos)

	// Hash should be non-zero
	assert.NotEqual(t, uint64(0), hash, "Initial position hash should not be zero")
}

func TestPolyglotHash_DifferentPositions(t *testing.T) {
	pos1 := board.CreatePositionFormFEN(board.InitialPosition)
	pos2 := board.CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1") // after e4

	hash1 := PolyglotHash(pos1)
	hash2 := PolyglotHash(pos2)

	assert.NotEqual(t, hash1, hash2, "Different positions should have different hashes")
}

func TestPolyglotHash_SideToMove(t *testing.T) {
	// Same board, different side = different hash
	pos1 := board.CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	pos2 := board.CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")

	hash1 := PolyglotHash(pos1)
	hash2 := PolyglotHash(pos2)

	assert.NotEqual(t, hash1, hash2, "Different side to move should have different hashes")
}

func TestDecodeMove_SimpleMove(t *testing.T) {
	// e2e4: from=12 (e2), to=28 (e4)
	// Encoding: to | (from << 6) = 28 | (12 << 6) = 28 | 768 = 796
	encoded := uint16(28 | (12 << 6))
	move := DecodeMove(encoded)

	assert.Equal(t, board.IndexToBitBoard(12), move.From, "From should be e2")
	assert.Equal(t, board.IndexToBitBoard(28), move.To, "To should be e4")
	assert.Equal(t, board.Empty, move.Promotion, "No promotion")
}

func TestDecodeMove_Promotion(t *testing.T) {
	// e7e8q: from=52 (e7), to=60 (e8), promo=4 (queen)
	// Encoding: to | (from << 6) | (promo << 12)
	encoded := uint16(60 | (52 << 6) | (4 << 12))
	move := DecodeMove(encoded)

	assert.Equal(t, board.IndexToBitBoard(52), move.From, "From should be e7")
	assert.Equal(t, board.IndexToBitBoard(60), move.To, "To should be e8")
	assert.Equal(t, board.Queen, move.Promotion, "Should promote to queen")
}

func TestBook_ProbeEmpty(t *testing.T) {
	book := &Book{}
	matches := book.Probe(0x12345678)

	assert.Empty(t, matches, "Empty book should return no matches")
}

func TestBook_LoadEmbedded(t *testing.T) {
	book := LoadEmbedded()

	assert.NotNil(t, book, "Embedded book should load")
	assert.Greater(t, book.Size(), 0, "Embedded book should have entries")

	// Test initial position lookup
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	hash := PolyglotHash(pos)
	matches := book.Probe(hash)

	assert.NotEmpty(t, matches, "Should find moves for initial position")
}
