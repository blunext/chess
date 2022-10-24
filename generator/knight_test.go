package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnight(t *testing.T) {
	tests := []struct {
		pos      uint64
		contains possibleMoves
		len      int
	}{
		{36, possibleMoves{0x400000, 0x80000000, 0x800000000000, 0x40000000000000, 0x10000000000000, 0x80000000000, 0x80000000000, 0x100000}, 8},
	}

	for _, test := range tests {
		t.Run("bishop SW", func(t *testing.T) {
			squares := generateKnightMoves()
			moves := squares[test.pos]
			assert.Equal(t, test.len, len(moves))

			// TODO: do more
		})
	}
}
