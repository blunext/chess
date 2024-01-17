package generator

import (
	"chess/board"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenericMoves(t *testing.T) {
	tests := []struct {
		pos   board.Bitboard
		fun   func() SquareMoves
		moves []board.Bitboard
	}{
		{0, knightMoves, []board.Bitboard{0x20000, 0x400}},
		{1, knightMoves, []board.Bitboard{0x10000, 0x40000, 0x800}},
		{37, knightMoves, []board.Bitboard{0x400000, 0x80000000, 0x800000000000, 0x40000000000000, 0x10000000000000, 0x80000000000, 0x8000000, 0x100000}},
		{49, knightMoves, []board.Bitboard{0x800000000000000, 0x80000000000, 0x400000000, 0x100000000}},
		{62, knightMoves, []board.Bitboard{0x800000000000, 0x200000000000, 0x10000000000000}},
		{63, knightMoves, []board.Bitboard{0x20000000000000, 0x400000000000}},
		{0, kingMoves, []board.Bitboard{0x100, 0x200, 0x2}},
		{21, kingMoves, []board.Bitboard{0x20000000, 0x40000000, 0x400000, 0x4000, 0x2000, 0x1000, 0x100000, 0x10000000}},
		{62, kingMoves, []board.Bitboard{0x2000000000000000, 0x20000000000000, 0x40000000000000, 0x80000000000000, 0x8000000000000000}},
	}

	for _, test := range tests {
		t.Run("Moves", func(t *testing.T) {
			squares := test.fun()
			moves := squares[test.pos]
			assert.Equal(t, len(test.moves), len(moves))
			for _, i := range test.moves {
				assert.Contains(t, moves, i)
			}
		})
	}
}
