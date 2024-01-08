package board

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitboardToSlice(t *testing.T) {
	testCases := []struct {
		name     string
		bitboard Bitboard
		expected []Bitboard
	}{
		{
			name:     "Empty bitboard",
			bitboard: 0,
			expected: []Bitboard{},
		},
		{
			name:     "Single piece",
			bitboard: 1 << 0,
			expected: []Bitboard{1},
		},
		{
			name:     "Mixed positions",
			bitboard: (1 << 0) | (1 << 3) | (1 << 15) | (1 << 30) | (1 << 55) | (1 << 60) | (1 << 63),
			expected: []Bitboard{1, 0x8, 0x8000, 0x40000000, 0x80000000000000, 0x1000000000000000, 0x8000000000000000},
		},
		{
			name:     "Full bottom line",
			bitboard: 0xFF,
			expected: []Bitboard{0x1, 0x2, 0x4, 0x8, 0x10, 0x20, 0x40, 0x80},
		},
		{
			name:     "Full upper line",
			bitboard: 0xFF00000000000000,
			expected: []Bitboard{0x100000000000000, 0x200000000000000, 0x400000000000000, 0x800000000000000, 0x1000000000000000, 0x2000000000000000, 0x4000000000000000, 0x8000000000000000},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bitboard.ToSlice()
			assert.Equal(t, tc.expected, result, "Failed test: %s", tc.name)
		})
	}
}
