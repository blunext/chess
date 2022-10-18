package bitboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFen(t *testing.T) {
	cb := fromFen(InitialPosition)
	assert.Equal(t, coloredPiece{Rook, ColorBlack}, cb[0])
	assert.Equal(t, coloredPiece{Knight, ColorBlack}, cb[1])
	assert.Equal(t, coloredPiece{Bishop, ColorBlack}, cb[2])
	assert.Equal(t, coloredPiece{Queen, ColorBlack}, cb[3])
	assert.Equal(t, coloredPiece{King, ColorBlack}, cb[4])
	assert.Equal(t, coloredPiece{Bishop, ColorBlack}, cb[5])
	assert.Equal(t, coloredPiece{Knight, ColorBlack}, cb[6])
	assert.Equal(t, coloredPiece{Rook, ColorBlack}, cb[7])
	assert.Equal(t, coloredPiece{Pawn, ColorBlack}, cb[8])
	assert.Equal(t, coloredPiece{Pawn, ColorBlack}, cb[12])
	assert.Equal(t, coloredPiece{Pawn, ColorBlack}, cb[15])
	assert.Equal(t, noPiece, cb[16])
	assert.Equal(t, noPiece, cb[16])
	assert.Equal(t, noPiece, cb[26])
	assert.Equal(t, noPiece, cb[27])
	assert.Equal(t, noPiece, cb[39])
	assert.Equal(t, noPiece, cb[47])
	assert.Equal(t, coloredPiece{Pawn, ColorWhite}, cb[48])
	assert.Equal(t, coloredPiece{Pawn, ColorWhite}, cb[55])
	assert.Equal(t, coloredPiece{Bishop, ColorWhite}, cb[56])
	assert.Equal(t, coloredPiece{Knight, ColorWhite}, cb[57])
	assert.Equal(t, coloredPiece{Rook, ColorWhite}, cb[58])
	assert.Equal(t, coloredPiece{Queen, ColorWhite}, cb[59])
	assert.Equal(t, coloredPiece{King, ColorWhite}, cb[60])
	assert.Equal(t, coloredPiece{Bishop, ColorWhite}, cb[61])
	assert.Equal(t, coloredPiece{Knight, ColorWhite}, cb[62])
	assert.Equal(t, coloredPiece{Rook, ColorWhite}, cb[63])

}

func TestPosition(t *testing.T) {
	cb := fromFen(InitialPosition)
	position := createPosition(cb)
	assert.Equal(t, bitboard(0xff00000000ff00), position.Pawns)
	assert.Equal(t, bitboard(0x4200000000000042), position.Knights)
	assert.Equal(t, bitboard(0x2400000000000024), position.Bishops)
	assert.Equal(t, bitboard(0x8100000000000081), position.Rooks)
	assert.Equal(t, bitboard(0x800000000000008), position.Queens)
	assert.Equal(t, bitboard(0x1000000000000010), position.Kings)

	assert.Equal(t, bitboard(0xffff000000000000), position.White)
	assert.Equal(t, bitboard(0xffff), position.Black)
}
