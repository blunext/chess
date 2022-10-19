package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBishopSE(t *testing.T) {
	tests := []struct {
		pos      int
		contains possibleMoves
		len      int
	}{
		{15, possibleMoves{}, 0},
		{10, possibleMoves{0x80000, 0x2000000000, 0x80000000000000}, 5},
		{48, possibleMoves{0x200000000000000}, 1},
		{2, possibleMoves{0x800, 0x100000, 0x800000000000}, 5},
	}

	for _, test := range tests {
		t.Run("bishop SE", func(t *testing.T) {
			positions := bishopSE(test.pos)
			assert.Equal(t, test.len, len(positions))
			for i, expected := range test.contains {
				assert.Contains(t, positions, expected)
				if i == 0 {
					// check if the first one is expected
					assert.Equal(t, positions[0], expected)
				}
			}
		})
	}
}
