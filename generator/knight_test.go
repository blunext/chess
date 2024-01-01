package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"chess/board"
)

func TestKnight(t *testing.T) {
	tests := []struct {
		pos   board.Bitboard
		moves []board.Bitboard
	}{
		{0, []board.Bitboard{0x20000, 0x400}},
		{1, []board.Bitboard{0x10000, 0x40000, 0x800}},
		{37, []board.Bitboard{0x400000, 0x80000000, 0x800000000000, 0x40000000000000, 0x10000000000000, 0x80000000000, 0x8000000, 0x100000}},
		{49, []board.Bitboard{0x800000000000000, 0x80000000000, 0x400000000, 0x100000000}},
		{62, []board.Bitboard{0x800000000000, 0x200000000000, 0x10000000000000}},
		{63, []board.Bitboard{0x20000000000000, 0x400000000000}},
	}

	for _, test := range tests {
		t.Run("knights", func(t *testing.T) {
			squares := generateKnightMoves()
			moves := squares[test.pos]
			assert.Equal(t, len(test.moves), len(moves))
			for _, i := range test.moves {
				assert.Contains(t, moves, i)
			}
		})
	}
}
