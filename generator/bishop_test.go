package generator

import (
	"chess/board"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBishopSE(t *testing.T) {
	tests := []struct {
		pos      int
		contains []board.Bitboard
		len      int
	}{
		{15, []board.Bitboard{}, 0},
		{10, []board.Bitboard{0x80000, 0x2000000000, 0x80000000000000}, 5},
		{48, []board.Bitboard{0x200000000000000}, 1},
		{2, []board.Bitboard{0x800, 0x100000, 0x800000000000}, 5},
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

func TestBishopSW(t *testing.T) {
	tests := []struct {
		pos      int
		contains []board.Bitboard
		len      int
	}{
		{8, []board.Bitboard{}, 0},
		{2, []board.Bitboard{0x200, 0x10000}, 2},
		{07, []board.Bitboard{0x4000, 0x800000000, 0x100000000000000}, 7},
		{39, []board.Bitboard{0x400000000000, 0x1000000000000000}, 3},
	}

	for _, test := range tests {
		t.Run("bishop SW", func(t *testing.T) {
			positions := bishopSW(test.pos)
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

func TestBishopNE(t *testing.T) {
	tests := []struct {
		pos      int
		contains []board.Bitboard
		len      int
	}{
		{5, []board.Bitboard{}, 0},
		{20, []board.Bitboard{0x2000, 0x40}, 2},
		{59, []board.Bitboard{0x10000000000000, 0x200000000000, 0x80000000}, 4},
		{56, []board.Bitboard{0x2000000000000, 0x10000000, 0x80}, 7},
	}

	for _, test := range tests {
		t.Run("bishop NE", func(t *testing.T) {
			positions := bishopNE(test.pos)
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

func TestBishopNW(t *testing.T) {
	tests := []struct {
		pos      int
		contains []board.Bitboard
		len      int
	}{
		{5, []board.Bitboard{}, 0},
		{20, []board.Bitboard{0x800, 0x4}, 2},
		{60, []board.Bitboard{0x8000000000000, 0x200000000, 0x1000000}, 4},
		{63, []board.Bitboard{0x40000000000000, 0x1000000000, 0x40000, 0x1}, 7},
	}

	for _, test := range tests {
		t.Run("bishop NW", func(t *testing.T) {
			positions := bishopNW(test.pos)
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

func TestGenerators(t *testing.T) {
	moves := generateRookMoves()
	assert.Equal(t, 64, len(moves))
	count := 0
	for _, possibleMoves := range moves {
		for _, moveList := range possibleMoves {
			count += len(moveList)
		}
	}
	assert.Equal(t, 64*7*2, count)

	generateQueenMoves(generateRookMoves(), generateBishopMoves())
	// todo: figure out how many bishop moves is possible and test it

}
