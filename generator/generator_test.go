package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		pos      int
		contains possibleMoves
		len      int
	}{
		{0, possibleMoves{0x100, 0x100000000, 0x100000000000000}, 7},
		{3, possibleMoves{0x800, 0x800000000, 0x800000000000000}, 7},
		{17, possibleMoves{0x2000000, 0x2000000000000, 0x200000000000000}, 5},
		{63, possibleMoves{}, 0},
	}
	for _, test := range tests {
		t.Run("move generator", func(t *testing.T) {
			positions := rookDown(test.pos)
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

func TestGenerator2(t *testing.T) {
	tests := []struct {
		pos      int
		contains possibleMoves
		len      int
	}{
		{63, possibleMoves{0x80000000000000, 0x80000000, 0x80}, 7},
		{17, possibleMoves{0x200, 0x2}, 2},
		{5, possibleMoves{}, 0},
	}
	for _, test := range tests {
		t.Run("move generator", func(t *testing.T) {
			positions := rookUp(test.pos)
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
	NewGenerator()
}
