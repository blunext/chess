// Package book provides opening book support using Polyglot format.
package book

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"io"
	"math/rand"
	"sort"

	"chess/board"
)

//go:embed book.bin
var embeddedBook []byte

// Entry represents a single opening book entry.
type Entry struct {
	Key    uint64 // Polyglot Zobrist hash
	Move   uint16 // Encoded move
	Weight uint16 // Move weight/priority
	Learn  uint32 // Learning data (unused)
}

// Book holds opening book entries sorted by hash for binary search.
type Book struct {
	entries []Entry
}

// LoadEmbedded loads the embedded opening book.
func LoadEmbedded() *Book {
	if len(embeddedBook) == 0 {
		return nil
	}
	b, _ := loadFromReader(bytes.NewReader(embeddedBook))
	return b
}

// loadFromReader reads a Polyglot format book from a Reader.
func loadFromReader(r io.Reader) (*Book, error) {
	var entries []Entry
	for {
		var e Entry
		// Polyglot format: big-endian
		err := binary.Read(r, binary.BigEndian, &e.Key)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		binary.Read(r, binary.BigEndian, &e.Move)
		binary.Read(r, binary.BigEndian, &e.Weight)
		binary.Read(r, binary.BigEndian, &e.Learn)
		entries = append(entries, e)
	}

	// Entries should already be sorted by key, but ensure it
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &Book{entries: entries}, nil
}

// Probe looks up a position hash and returns all matching moves.
func (b *Book) Probe(hash uint64) []Entry {
	if b == nil || len(b.entries) == 0 {
		return nil
	}

	// Binary search for first entry with this key
	idx := sort.Search(len(b.entries), func(i int) bool {
		return b.entries[i].Key >= hash
	})

	var matches []Entry
	for idx < len(b.entries) && b.entries[idx].Key == hash {
		matches = append(matches, b.entries[idx])
		idx++
	}
	return matches
}

// ProbeRandom returns a random move from the book, weighted by priority.
func (b *Book) ProbeRandom(hash uint64, rng *rand.Rand) (board.Move, bool) {
	matches := b.Probe(hash)
	if len(matches) == 0 {
		return board.Move{}, false
	}

	// Calculate total weight
	var totalWeight uint32
	for _, e := range matches {
		totalWeight += uint32(e.Weight)
	}

	if totalWeight == 0 {
		// All weights zero, pick first
		return DecodeMove(matches[0].Move), true
	}

	// Random weighted selection
	r := rng.Uint32() % totalWeight
	var cumulative uint32
	for _, e := range matches {
		cumulative += uint32(e.Weight)
		if r < cumulative {
			return DecodeMove(e.Move), true
		}
	}

	return DecodeMove(matches[0].Move), true
}

// DecodeMove converts Polyglot move encoding to board.Move.
// Polyglot encoding:
// - bits 0-5: destination square (0-63)
// - bits 6-11: origin square (0-63)
// - bits 12-14: promotion piece (0=none, 1=knight, 2=bishop, 3=rook, 4=queen)
// Note: castling is encoded as king captures rook.
func DecodeMove(raw uint16) board.Move {
	toSq := int(raw & 0x3F)
	fromSq := int((raw >> 6) & 0x3F)
	promo := int((raw >> 12) & 0x07)

	m := board.Move{
		From: board.IndexToBitBoard(fromSq),
		To:   board.IndexToBitBoard(toSq),
	}

	// Handle promotion
	switch promo {
	case 1:
		m.Promotion = board.Knight
	case 2:
		m.Promotion = board.Bishop
	case 3:
		m.Promotion = board.Rook
	case 4:
		m.Promotion = board.Queen
	}

	return m
}

// Size returns the number of entries in the book.
func (b *Book) Size() int {
	if b == nil {
		return 0
	}
	return len(b.entries)
}
