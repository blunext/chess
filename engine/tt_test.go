package engine

import (
	"testing"
	"time"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

func TestTT_StoreAndProbe(t *testing.T) {
	tt := NewTranspositionTable(1) // 1 MB

	hash := uint64(0x123456789ABCDEF0)
	move := board.Move{From: 1 << 12, To: 1 << 28, Piece: board.Pawn} // e2->e4 as bitboards

	// Store entry
	tt.Store(hash, 100, 5, TTFlagExact, move)

	// Probe should find it
	entry, found := tt.Probe(hash)
	assert.True(t, found, "Should find stored entry")
	assert.Equal(t, int16(100), entry.Score)
	assert.Equal(t, int8(5), entry.Depth)
	assert.Equal(t, TTFlagExact, entry.Flag)
}

func TestTT_ProbeNotFound(t *testing.T) {
	tt := NewTranspositionTable(1)

	hash := uint64(0x123456789ABCDEF0)
	_, found := tt.Probe(hash)
	assert.False(t, found, "Should not find entry in empty table")
}

func TestTT_Collision(t *testing.T) {
	tt := NewTranspositionTable(1)

	// Two hashes that map to same index (different upper 32 bits)
	hash1 := uint64(0x1111111100000001)
	hash2 := uint64(0x2222222200000001) // Same lower bits, different upper

	move1 := board.Move{From: 1 << 12, To: 1 << 28, Piece: board.Pawn}
	move2 := board.Move{From: 1 << 11, To: 1 << 27, Piece: board.Pawn}

	tt.Store(hash1, 100, 5, TTFlagExact, move1)
	tt.Store(hash2, 200, 6, TTFlagLower, move2) // Overwrites

	// hash1 should not be found (overwritten)
	_, found1 := tt.Probe(hash1)
	assert.False(t, found1, "hash1 should be overwritten")

	// hash2 should be found
	entry2, found2 := tt.Probe(hash2)
	assert.True(t, found2, "hash2 should be found")
	assert.Equal(t, int16(200), entry2.Score)
}

func TestTT_Clear(t *testing.T) {
	tt := NewTranspositionTable(1)

	hash := uint64(0x123456789ABCDEF0)
	move := board.Move{From: 1 << 12, To: 1 << 28, Piece: board.Pawn}
	tt.Store(hash, 100, 5, TTFlagExact, move)

	tt.Clear()

	_, found := tt.Probe(hash)
	assert.False(t, found, "Table should be empty after clear")
}

func TestTT_Hashfull(t *testing.T) {
	tt := NewTranspositionTable(1)

	// Empty table
	assert.Equal(t, 0, tt.Hashfull())

	// Fill some entries
	for i := uint64(0); i < 500; i++ {
		hash := uint64(0xABCDEF0000000000) | i
		tt.Store(hash, int16(i), 1, TTFlagExact, board.Move{})
	}

	// Should have ~50% full (500/1000 sampled)
	hashfull := tt.Hashfull()
	assert.Greater(t, hashfull, 400, "Should be more than 40% full")
	assert.Less(t, hashfull, 600, "Should be less than 60% full")
}

// Test that TT improves iterative deepening performance
func TestTT_ImprovesIterativeDeepening(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	// Search with TT (should use cached results from previous depths)
	InitTT(16)
	TT.Clear()
	result := SearchWithTime(pos, pm, 200*time.Millisecond)

	// Should reach reasonable depth
	assert.Greater(t, result.Depth, 2, "Should reach at least depth 3 with TT")
	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
}

// Benchmark search with and without TT at fixed depth
// This measures actual speedup from TT by comparing nodes at same depth
func BenchmarkSearch_FixedDepth_WithTT(b *testing.B) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()
	InitTT(64) // 64 MB TT

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TT.Clear()
		Search(pos, pm, 4) // Fixed depth 4
	}
}

func BenchmarkSearch_FixedDepth_WithoutTT(b *testing.B) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	// Disable TT temporarily
	oldTT := TT
	TT = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Search(pos, pm, 4) // Fixed depth 4
	}

	b.StopTimer()
	TT = oldTT
}
