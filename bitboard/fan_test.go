package bitboard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColoredBoard(t *testing.T) {
	tests := []struct {
		piece coloredPiece
		index int
	}{
		{coloredPiece{Rook, ColorBlack}, 0},
		{coloredPiece{Knight, ColorBlack}, 1},
		{coloredPiece{Bishop, ColorBlack}, 2},
		{coloredPiece{Queen, ColorBlack}, 3},
		{coloredPiece{King, ColorBlack}, 4},
		{coloredPiece{Bishop, ColorBlack}, 5},
		{coloredPiece{Knight, ColorBlack}, 6},
		{coloredPiece{Rook, ColorBlack}, 7},
		{coloredPiece{Pawn, ColorBlack}, 8},
		{coloredPiece{Pawn, ColorBlack}, 12},
		{coloredPiece{Pawn, ColorBlack}, 15},
		{noPiece, 16},
		{noPiece, 16},
		{noPiece, 26},
		{noPiece, 27},
		{noPiece, 39},
		{noPiece, 47},
		{coloredPiece{Pawn, ColorWhite}, 48},
		{coloredPiece{Pawn, ColorWhite}, 55},
		{coloredPiece{Rook, ColorWhite}, 56},
		{coloredPiece{Knight, ColorWhite}, 57},
		{coloredPiece{Bishop, ColorWhite}, 58},
		{coloredPiece{Queen, ColorWhite}, 59},
		{coloredPiece{King, ColorWhite}, 60},
		{coloredPiece{Bishop, ColorWhite}, 61},
		{coloredPiece{Knight, ColorWhite}, 62},
		{coloredPiece{Rook, ColorWhite}, 63},
	}

	cb := createColoredBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	for _, test := range tests {
		t.Run("colored piece", func(t *testing.T) {
			assert.Equal(t, test.piece, cb[test.index])
		})
		fmt.Println()
	}
}

func TestPosition(t *testing.T) {
	cb := fen2ColoredBoard(InitialPosition)
	position := createPosition(cb)

	tests := []struct {
		expectedPattern bitboard
		piecePattern    bitboard
	}{
		{0xff00000000ff00, position.Pawns},
		{0x4200000000000042, position.Knights},
		{0x2400000000000024, position.Bishops},
		{0x8100000000000081, position.Rooks},
		{0x800000000000008, position.Queens},
		{0x1000000000000010, position.Kings},

		{0xffff000000000000, position.White},
		{0xffff, position.Black},
	}

	for _, test := range tests {
		t.Run("positions", func(t *testing.T) {
			assert.Equal(t, test.expectedPattern, test.piecePattern)
		})
		fmt.Println()
	}

}
