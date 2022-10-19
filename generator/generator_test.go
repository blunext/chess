package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRookDown(t *testing.T) {
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
		t.Run("rook down", func(t *testing.T) {
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

func TestRookUp(t *testing.T) {
	tests := []struct {
		pos      int
		contains possibleMoves
		len      int
	}{
		{5, possibleMoves{}, 0},
		{17, possibleMoves{0x200, 0x2}, 2},
		{63, possibleMoves{0x80000000000000, 0x80000000, 0x80}, 7},
	}
	for _, test := range tests {
		t.Run("rook up", func(t *testing.T) {
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
}

func TestRookRight(t *testing.T) {
	tests := []struct {
		pos      int
		contains possibleMoves
		len      int
	}{
		{7, possibleMoves{}, 0},
		{6, possibleMoves{0x80}, 1},
		{8, possibleMoves{0x200, 0x800, 0x8000}, 7},
		{62, possibleMoves{0x8000000000000000}, 1},
	}
	for _, test := range tests {
		t.Run("rook right", func(t *testing.T) {
			positions := rookRight(test.pos)
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

func TestRookLeft(t *testing.T) {
	tests := []struct {
		pos      int
		contains possibleMoves
		len      int
	}{
		{0, possibleMoves{}, 0},
		{1, possibleMoves{0x1}, 1},
		{7, possibleMoves{0x40, 0x8, 0x1}, 7},
		{10, possibleMoves{0x200, 0x100}, 2},
		{63, possibleMoves{0x4000000000000000, 0x800000000000000, 0x100000000000000}, 7},
	}
	for _, test := range tests {
		t.Run("rook right", func(t *testing.T) {
			positions := rookLeft(test.pos)
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
