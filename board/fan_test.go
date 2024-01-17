package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColoredBoard(t *testing.T) {
	tests := []struct {
		piece coloredPiece
		index int
	}{
		{coloredPiece{Rook, ColorWhite}, 0},
		{coloredPiece{Knight, ColorWhite}, 1},
		{coloredPiece{Bishop, ColorWhite}, 2},
		{coloredPiece{Queen, ColorWhite}, 3},
		{coloredPiece{King, ColorWhite}, 4},
		{coloredPiece{Bishop, ColorWhite}, 5},
		{coloredPiece{Knight, ColorWhite}, 6},
		{coloredPiece{Rook, ColorWhite}, 7},
		{coloredPiece{Pawn, ColorWhite}, 8},
		{coloredPiece{Pawn, ColorWhite}, 12},
		{coloredPiece{Pawn, ColorWhite}, 15},
		{noPiece, 16},
		{noPiece, 16},
		{noPiece, 26},
		{noPiece, 27},
		{noPiece, 39},
		{noPiece, 47},
		{coloredPiece{Pawn, ColorBlack}, 48},
		{coloredPiece{Pawn, ColorBlack}, 55},
		{coloredPiece{Rook, ColorBlack}, 56},
		{coloredPiece{Knight, ColorBlack}, 57},
		{coloredPiece{Bishop, ColorBlack}, 58},
		{coloredPiece{Queen, ColorBlack}, 59},
		{coloredPiece{King, ColorBlack}, 60},
		{coloredPiece{Bishop, ColorBlack}, 61},
		{coloredPiece{Knight, ColorBlack}, 62},
		{coloredPiece{Rook, ColorBlack}, 63},
	}

	cb := createColoredBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	for _, test := range tests {
		t.Run("colored piece", func(t *testing.T) {
			assert.Equal(t, test.piece, cb[test.index])
		})
	}
}

func TestPosition(t *testing.T) {
	position := CreatePositionFormFEN(InitialPosition)
	tests := []struct {
		expectedPattern Bitboard
		piecePattern    Bitboard
	}{
		{0xff00000000ff00, position.Pawns},
		{0x4200000000000042, position.Knights},
		{0x2400000000000024, position.Bishops},
		{0x8100000000000081, position.Rooks},
		{0x800000000000008, position.Queens},
		{0x1000000000000010, position.Kings},

		{0xffff000000000000, position.Black},
		{0xffff, position.White},
	}

	for _, test := range tests {
		t.Run("patterns", func(t *testing.T) {
			assert.Equal(t, test.expectedPattern, test.piecePattern)
		})
	}

	assert.Equal(t, uint8(CastleWhiteKingSide|CastleWhiteQueenSide|CastleBlackKingSide|CastleBlackQueenSide), position.CastleSide, "castling failed")
	assert.Equal(t, Bitboard(0), position.EnPassant, "en passant failed")
	assert.Truef(t, position.WhiteMove, "white move")

	position = CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")
	assert.Falsef(t, position.WhiteMove, "black move")

	position = CreatePositionFormFEN("rnbqkbnr/pppp1ppp/4p3/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1")
	tests = []struct {
		expectedPattern Bitboard
		piecePattern    Bitboard
	}{
		{0xef10000008f700, position.Pawns},
		{0x8f7ff, position.White},
		{0xffef100000000000, position.Black},
	}
	for _, test := range tests {
		t.Run("patterns", func(t *testing.T) {
			assert.Equal(t, test.expectedPattern, test.piecePattern)
		})
	}
}

func TestCastling(t *testing.T) {
	tests := []struct {
		inPattern string
		expected  uint8
	}{
		{"KQkq", CastleWhiteKingSide | CastleWhiteQueenSide | CastleBlackKingSide | CastleBlackQueenSide},
		{"Qkq", CastleWhiteQueenSide | CastleBlackKingSide | CastleBlackQueenSide},
		{"Kkq", CastleWhiteKingSide | CastleBlackKingSide | CastleBlackQueenSide},
		{"KQq", CastleWhiteKingSide | CastleWhiteQueenSide | CastleBlackQueenSide},
		{"KQk", CastleWhiteKingSide | CastleWhiteQueenSide | CastleBlackKingSide},
		{"KQ", CastleWhiteKingSide | CastleWhiteQueenSide},
		{"kq", CastleBlackKingSide | CastleBlackQueenSide},
		{"q", CastleBlackQueenSide},
		{"-", 0},
	}

	for _, test := range tests {
		t.Run("castling", func(t *testing.T) {
			castling := castleAbility(test.inPattern)
			assert.Equal(t, castling, test.expected)
		})
	}
}

func TestEnPassant(t *testing.T) {
	tests := []struct {
		inPattern string
		expected  Bitboard
	}{
		{"a6", 0x10000000000},
		{"c6", 0x40000000000},
		{"h6", 0x800000000000},
		{"b3", 0x20000},
		{"g3", 0x400000},
		{"-", 0},
	}

	for _, test := range tests {
		t.Run("en passant", func(t *testing.T) {
			result := enPassant(test.inPattern)
			assert.Equal(t, result, test.expected)
		})
	}
}

func TestBB(t *testing.T) {
	BB()
}
