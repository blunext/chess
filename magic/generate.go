//go:build ignore
// +build ignore

// This program generates Magic Bitboard data for Rooks and Bishops.
// Run with: go run generate.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
)

type Bitboard uint64

type Magic struct {
	Number uint64
	Mask   Bitboard
	Shift  uint8
}

func main() {
	// Note: Go 1.20+ auto-seeds math/rand, no need for rand.Seed()

	fmt.Println("Generating Magic Bitboards...")

	// Generate for Rooks
	fmt.Println("\n=== ROOKS ===")
	rookMagics, rookMoves := generateMagicsForPiece(true)

	// Generate for Bishops (uses smaller table - max 9 bits = 512 entries)
	fmt.Println("\n=== BISHOPS ===")
	bishopMagics, bishopMoves := generateBishopMagics()

	// Encode and save
	fmt.Println("\nSaving to magicData...")
	if err := saveData(rookMagics, rookMoves, bishopMagics, bishopMoves); err != nil {
		fmt.Printf("Error saving: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}

func generateMagicsForPiece(isRook bool) ([64]Magic, [64][4096]Bitboard) {
	var magics [64]Magic
	var moves [64][4096]Bitboard

	for square := 0; square < 64; square++ {
		mask := generateMask(square, isRook)
		bits := popCount(uint64(mask))
		shift := uint8(64 - bits)

		// Generate all blocker combinations and their results
		numCombinations := 1 << bits
		blockerConfigs := make([]Bitboard, numCombinations)
		attackSets := make([]Bitboard, numCombinations)

		for i := 0; i < numCombinations; i++ {
			blockerConfigs[i] = indexToBlockers(i, mask)
			attackSets[i] = calculateAttacks(square, blockerConfigs[i], isRook)
		}

		// Find magic number
		magic := findMagic(square, mask, shift, blockerConfigs, attackSets, isRook)

		magics[square] = Magic{
			Number: magic,
			Mask:   mask,
			Shift:  shift,
		}

		// Fill the moves table
		for i := 0; i < numCombinations; i++ {
			index := (uint64(blockerConfigs[i]) * magic) >> shift
			moves[square][index] = attackSets[i]
		}

		pieceName := "Rook"
		if !isRook {
			pieceName = "Bishop"
		}
		fmt.Printf("%s %s: magic=0x%016x, bits=%d\n",
			pieceName, squareToAlgebraic(square), magic, bits)
	}

	return magics, moves
}

// generateBishopMagics generates magic numbers and lookup tables for bishops.
// Uses 512-entry table (max 9 bits for bishop blockers) instead of 4096.
func generateBishopMagics() ([64]Magic, [64][512]Bitboard) {
	var magics [64]Magic
	var moves [64][512]Bitboard

	for square := 0; square < 64; square++ {
		mask := generateMask(square, false)
		bits := popCount(uint64(mask))
		shift := uint8(64 - bits)

		numCombinations := 1 << bits
		blockerConfigs := make([]Bitboard, numCombinations)
		attackSets := make([]Bitboard, numCombinations)

		for i := 0; i < numCombinations; i++ {
			blockerConfigs[i] = indexToBlockers(i, mask)
			attackSets[i] = calculateAttacks(square, blockerConfigs[i], false)
		}

		magic := findMagic(square, mask, shift, blockerConfigs, attackSets, false)

		magics[square] = Magic{
			Number: magic,
			Mask:   mask,
			Shift:  shift,
		}

		for i := 0; i < numCombinations; i++ {
			index := (uint64(blockerConfigs[i]) * magic) >> shift
			moves[square][index] = attackSets[i]
		}

		fmt.Printf("Bishop %s: magic=0x%016x, bits=%d\n",
			squareToAlgebraic(square), magic, bits)
	}

	return magics, moves
}

// generateMask creates the blocker mask for a square.
// This excludes edge squares since pieces on edges don't affect further movement.
func generateMask(square int, isRook bool) Bitboard {
	var mask Bitboard
	rank := square / 8
	file := square % 8

	if isRook {
		// Vertical (exclude first and last rank)
		for r := 1; r < 7; r++ {
			if r != rank {
				mask |= 1 << (r*8 + file)
			}
		}
		// Horizontal (exclude first and last file)
		for f := 1; f < 7; f++ {
			if f != file {
				mask |= 1 << (rank*8 + f)
			}
		}
	} else {
		// Diagonals (exclude edges)
		directions := [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
		for _, d := range directions {
			r, f := rank+d[0], file+d[1]
			for r > 0 && r < 7 && f > 0 && f < 7 {
				mask |= 1 << (r*8 + f)
				r += d[0]
				f += d[1]
			}
		}
	}

	return mask
}

// indexToBlockers converts an index to a specific blocker configuration.
func indexToBlockers(index int, mask Bitboard) Bitboard {
	var blockers Bitboard
	bits := mask

	for i := 0; bits != 0; i++ {
		// Find lowest set bit
		lsb := bits & -bits
		bits &^= lsb

		if index&(1<<i) != 0 {
			blockers |= lsb
		}
	}

	return blockers
}

// calculateAttacks computes the attack bitboard for a piece on square with blockers.
func calculateAttacks(square int, blockers Bitboard, isRook bool) Bitboard {
	var attacks Bitboard
	rank := square / 8
	file := square % 8

	var directions [][2]int
	if isRook {
		directions = [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	} else {
		directions = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	}

	for _, d := range directions {
		r, f := rank+d[0], file+d[1]
		for r >= 0 && r < 8 && f >= 0 && f < 8 {
			sq := Bitboard(1 << (r*8 + f))
			attacks |= sq
			if blockers&sq != 0 {
				break // Hit a blocker
			}
			r += d[0]
			f += d[1]
		}
	}

	return attacks
}

// findMagic finds a magic number for the given square.
func findMagic(square int, mask Bitboard, shift uint8, blockers, attacks []Bitboard, isRook bool) uint64 {
	tableSize := 1 << (64 - shift)
	used := make([]Bitboard, tableSize)

	for attempts := 0; attempts < 100_000_000; attempts++ {
		// Generate a random candidate with few bits set (sparse magic)
		magic := randomMagicCandidate()

		// Quick rejection: check if magic * mask has enough bits in upper portion
		if popCount(uint64(mask)*magic&0xFF00000000000000) < 6 {
			continue
		}

		// Clear used table
		for i := range used {
			used[i] = 0
		}

		// Test all blocker configurations
		valid := true
		for i := 0; i < len(blockers); i++ {
			index := (uint64(blockers[i]) * magic) >> shift

			if used[index] == 0 {
				used[index] = attacks[i]
			} else if used[index] != attacks[i] {
				// Collision with different attack set
				valid = false
				break
			}
			// Same attack set at same index is OK (constructive collision)
		}

		if valid {
			return magic
		}
	}

	panic(fmt.Sprintf("Failed to find magic for square %d", square))
}

// randomMagicCandidate generates a random 64-bit number with few bits set.
func randomMagicCandidate() uint64 {
	return rand.Uint64() & rand.Uint64() & rand.Uint64()
}

func popCount(x uint64) int {
	count := 0
	for x != 0 {
		count++
		x &= x - 1
	}
	return count
}

func squareToAlgebraic(sq int) string {
	file := sq % 8
	rank := sq / 8
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}

func saveData(rookMagics [64]Magic, rookMoves [64][4096]Bitboard,
	bishopMagics [64]Magic, bishopMoves [64][512]Bitboard) error {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(rookMagics); err != nil {
		return err
	}
	if err := enc.Encode(rookMoves); err != nil {
		return err
	}
	if err := enc.Encode(bishopMagics); err != nil {
		return err
	}
	if err := enc.Encode(bishopMoves); err != nil {
		return err
	}

	return os.WriteFile("magicData", buf.Bytes(), 0644)
}
