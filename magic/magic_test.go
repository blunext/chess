package magic

import (
	"chess/board"
	"chess/generator"
	"testing"
)

// TestMagicBitboardsCorrectness verifies that pre-computed Magic Bitboards
// produce the same results as our manual move generators.
func TestMagicBitboardsCorrectness(t *testing.T) {
	// Load magic data first
	if err := Prepare(); err != nil {
		t.Fatalf("Failed to load magic data: %v", err)
	}

	// Generate reference moves using our current generators
	refRookMoves := generator.GenerateRookMovesForTesting()
	refBishopMoves := generator.GenerateBishopMovesForTesting()

	// Test a subset of positions (every 8th square for performance)
	testSquares := []int{0, 8, 16, 24, 32, 40, 48, 56, 4, 28, 35, 63}

	for _, square := range testSquares {
		t.Run(board.IndexToAlgebraic(square), func(t *testing.T) {
			// Test Rook
			testRookMagic(t, square, refRookMoves)
			// Test Bishop
			testBishopMagic(t, square, refBishopMoves)
		})
	}
}

func testRookMagic(t *testing.T, square int, refMoves board.SquareMoves) {
	squareBB := board.IndexToBitBoard(square)
	refDirections := refMoves[squareBB]

	// Test with different blocker configurations
	testBlockers := []board.Bitboard{
		0,                  // empty board
		0xFFFFFFFFFFFFFFFF, // full board
		0x0000000000FF0000, // horizontal blockers
		0x0101010101010101, // vertical blockers
	}

	for _, blockers := range testBlockers {
		// Get magic result
		m := RookMagics[square]
		maskedBlockers := blockers & m.Mask
		index := (uint64(maskedBlockers) * m.Number) >> m.Shift
		magicResult := RookMoves[square][index]

		// Get reference result by simulating direction-based generation
		refResult := simulateDirectionMoves(refDirections, blockers)

		if magicResult != refResult {
			t.Errorf("Rook at %s with blockers %016x: magic=%016x, ref=%016x",
				board.IndexToAlgebraic(square), blockers, magicResult, refResult)
		}
	}
}

func testBishopMagic(t *testing.T, square int, refMoves board.SquareMoves) {
	squareBB := board.IndexToBitBoard(square)
	refDirections := refMoves[squareBB]

	testBlockers := []board.Bitboard{
		0,
		0xFFFFFFFFFFFFFFFF,
		0x0055AA0000AA5500, // diagonal pattern
	}

	for _, blockers := range testBlockers {
		m := BishopMagics[square]
		maskedBlockers := blockers & m.Mask
		index := (uint64(maskedBlockers) * m.Number) >> m.Shift
		magicResult := BishopMoves[square][index]

		refResult := simulateDirectionMoves(refDirections, blockers)

		if magicResult != refResult {
			t.Errorf("Bishop at %s with blockers %016x: magic=%016x, ref=%016x",
				board.IndexToAlgebraic(square), blockers, magicResult, refResult)
		}
	}
}

// simulateDirectionMoves mimics the direction-based move generation
// to produce a reference result for comparison.
// NOTE: Magic Bitboards include the blocker square (for captures),
// so we also add the blocker to the result before breaking.
func simulateDirectionMoves(directions [][]board.Bitboard, blockers board.Bitboard) board.Bitboard {
	var result board.Bitboard

	for _, direction := range directions {
		for _, square := range direction {
			result |= square
			if blockers&square == square {
				// Hit a blocker, include it (potential capture) then stop
				break
			}
		}
	}

	return result
}
