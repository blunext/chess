package engine

import (
	"chess/board"
)

// TTFlag indicates what type of bound the score represents.
type TTFlag uint8

const (
	TTFlagNone  TTFlag = 0
	TTFlagExact TTFlag = 1 // Exact score (PV node)
	TTFlagLower TTFlag = 2 // Score is lower bound (failed high / beta cutoff)
	TTFlagUpper TTFlag = 3 // Score is upper bound (failed low / no move improved alpha)
)

// TTEntry represents a single entry in the transposition table.
// Size: 16 bytes (optimized for cache efficiency)
type TTEntry struct {
	Hash     uint32     // Upper 32 bits of Zobrist hash (for verification)
	BestMove board.Move // Best move found (12 bytes, but could be compressed)
	Score    int16      // Evaluation score
	Depth    int8       // Search depth
	Flag     TTFlag     // Type of bound
}

// TranspositionTable is a hash table for storing search results.
// Uses "always replace" strategy for simplicity.
type TranspositionTable struct {
	entries []TTEntry
	size    uint64 // Number of entries
	mask    uint64 // size - 1 (for fast modulo via AND)
}

// DefaultHashMB is the default transposition table size in megabytes.
const DefaultHashMB = 64

// NewTranspositionTable creates a new transposition table with the given size in MB.
func NewTranspositionTable(sizeMB int) *TranspositionTable {
	if sizeMB <= 0 {
		sizeMB = DefaultHashMB
	}

	// Calculate number of entries
	// Each entry is approximately 24 bytes (with padding)
	entrySize := uint64(24)
	numEntries := (uint64(sizeMB) * 1024 * 1024) / entrySize

	// Round down to power of 2 for fast modulo
	size := uint64(1)
	for size*2 <= numEntries {
		size *= 2
	}

	return &TranspositionTable{
		entries: make([]TTEntry, size),
		size:    size,
		mask:    size - 1,
	}
}

// index calculates the table index from a Zobrist hash.
func (tt *TranspositionTable) index(hash uint64) uint64 {
	return hash & tt.mask
}

// verify checks if the stored hash matches (using upper 32 bits).
func (tt *TranspositionTable) verify(hash uint64, entry *TTEntry) bool {
	return entry.Hash == uint32(hash>>32)
}

// Probe looks up a position in the transposition table.
// Returns the entry and true if found, or empty entry and false if not.
func (tt *TranspositionTable) Probe(hash uint64) (TTEntry, bool) {
	idx := tt.index(hash)
	entry := &tt.entries[idx]

	if entry.Flag == TTFlagNone {
		return TTEntry{}, false
	}

	if !tt.verify(hash, entry) {
		return TTEntry{}, false
	}

	return *entry, true
}

// Store saves a position in the transposition table.
// Uses "always replace" strategy - new entries always overwrite old ones.
func (tt *TranspositionTable) Store(hash uint64, score int16, depth int8, flag TTFlag, bestMove board.Move) {
	idx := tt.index(hash)
	tt.entries[idx] = TTEntry{
		Hash:     uint32(hash >> 32),
		Score:    score,
		Depth:    depth,
		Flag:     flag,
		BestMove: bestMove,
	}
}

// Clear resets the transposition table.
func (tt *TranspositionTable) Clear() {
	for i := range tt.entries {
		tt.entries[i] = TTEntry{}
	}
}

// Size returns the number of entries in the table.
func (tt *TranspositionTable) Size() uint64 {
	return tt.size
}

// SizeMB returns the approximate size in megabytes.
func (tt *TranspositionTable) SizeMB() int {
	return int((tt.size * 24) / (1024 * 1024))
}

// Hashfull returns the permille of entries that are used (for UCI info).
func (tt *TranspositionTable) Hashfull() int {
	// Sample first 1000 entries for performance
	sample := uint64(1000)
	if sample > tt.size {
		sample = tt.size
	}

	used := 0
	for i := uint64(0); i < sample; i++ {
		if tt.entries[i].Flag != TTFlagNone {
			used++
		}
	}

	return (used * 1000) / int(sample)
}

// Global transposition table instance
var TT *TranspositionTable

// InitTT initializes the global transposition table.
func InitTT(sizeMB int) {
	TT = NewTranspositionTable(sizeMB)
}

func init() {
	// Initialize with default size
	InitTT(DefaultHashMB)
}
