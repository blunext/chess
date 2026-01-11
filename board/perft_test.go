package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Reference perft values from Chess Programming Wiki:
// https://www.chessprogramming.org/Perft_Results

func createTestPieceMoves() PieceMoves {
	// We need to generate proper piece moves for knight and king
	// For now, use a simplified version that covers basic moves
	pm := make(PieceMoves)

	// Knight moves (all 8 L-shaped jumps from each square)
	pm[Knight] = make(SquareMoves)
	knightOffsets := [][2]int{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	for sq := 0; sq < 64; sq++ {
		file := sq & 7
		rank := sq >> 3
		var targets []Bitboard
		for _, off := range knightOffsets {
			newFile := file + off[0]
			newRank := rank + off[1]
			if newFile >= 0 && newFile < 8 && newRank >= 0 && newRank < 8 {
				targets = append(targets, IndexToBitBoard(newRank*8+newFile))
			}
		}
		pm[Knight][IndexToBitBoard(sq)] = [][]Bitboard{targets}
	}

	// King moves (all 8 directions from each square)
	pm[King] = make(SquareMoves)
	kingOffsets := [][2]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}
	for sq := 0; sq < 64; sq++ {
		file := sq & 7
		rank := sq >> 3
		var targets []Bitboard
		for _, off := range kingOffsets {
			newFile := file + off[0]
			newRank := rank + off[1]
			if newFile >= 0 && newFile < 8 && newRank >= 0 && newRank < 8 {
				targets = append(targets, IndexToBitBoard(newRank*8+newFile))
			}
		}
		pm[King][IndexToBitBoard(sq)] = [][]Bitboard{targets}
	}

	return pm
}

func TestPerft_InitialPosition_Depth1(t *testing.T) {
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	assert.Equal(t, uint64(20), result, "Initial position depth 1 should have 20 moves")
}

func TestPerft_InitialPosition_Depth2(t *testing.T) {
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 2)
	assert.Equal(t, uint64(400), result, "Initial position depth 2 should have 400 positions")
}

func TestPerft_InitialPosition_Depth3(t *testing.T) {
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 3)
	assert.Equal(t, uint64(8902), result, "Initial position depth 3 should have 8902 positions")
}

func TestPerft_InitialPosition_Depth4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 4 perft in short mode")
	}
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 4)
	assert.Equal(t, uint64(197281), result, "Initial position depth 4 should have 197281 positions")
}

// Kiwipete position - tests many edge cases
// r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -
func TestPerft_Kiwipete_Depth1(t *testing.T) {
	pos := CreatePositionFormFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	assert.Equal(t, uint64(48), result, "Kiwipete depth 1 should have 48 moves")
}

func TestPerft_Kiwipete_Depth2(t *testing.T) {
	pos := CreatePositionFormFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 2)
	assert.Equal(t, uint64(2039), result, "Kiwipete depth 2 should have 2039 positions")
}

func TestPerft_Kiwipete_Depth3(t *testing.T) {
	pos := CreatePositionFormFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 3)
	assert.Equal(t, uint64(97862), result, "Kiwipete depth 3 should have 97862 positions")
}

func TestPerft_Kiwipete_Depth4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 4 perft in short mode")
	}
	pos := CreatePositionFormFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 4)
	assert.Equal(t, uint64(4085603), result, "Kiwipete depth 4 should have 4085603 positions")
}

// Position 3 - tests en passant and promotions
// 8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - -
func TestPerft_Position3_Depth1(t *testing.T) {
	pos := CreatePositionFormFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	assert.Equal(t, uint64(14), result, "Position 3 depth 1 should have 14 moves")
}

// Position for en passant testing
func TestPerft_EnPassant_Depth1(t *testing.T) {
	// After 1.e4 d5 2.e5 f5 - white can capture en passant
	pos := CreatePositionFormFEN("rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	// Should include e5xf6 en passant
	assert.Greater(t, result, uint64(0), "Should have moves including en passant")
}

// Deeper tests for Position 3
func TestPerft_Position3_Depth2(t *testing.T) {
	pos := CreatePositionFormFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 2)
	assert.Equal(t, uint64(191), result, "Position 3 depth 2 should have 191 positions")
}

func TestPerft_Position3_Depth3(t *testing.T) {
	pos := CreatePositionFormFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 3)
	assert.Equal(t, uint64(2812), result, "Position 3 depth 3 should have 2812 positions")
}

func TestPerft_Position3_Depth4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 4 perft in short mode")
	}
	pos := CreatePositionFormFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 4)
	assert.Equal(t, uint64(43238), result, "Position 3 depth 4 should have 43238 positions")
}

// Position 4 - tests promotions with capture, castling rights
// r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1
func TestPerft_Position4_Depth1(t *testing.T) {
	pos := CreatePositionFormFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	assert.Equal(t, uint64(6), result, "Position 4 depth 1 should have 6 moves")
}

func TestPerft_Position4_Depth2(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 2 perft in short mode")
	}
	pos := CreatePositionFormFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 2)
	assert.Equal(t, uint64(264), result, "Position 4 depth 2 should have 264 positions")
}

func TestPerft_Position4_Depth3(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 3 perft in short mode")
	}
	pos := CreatePositionFormFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 3)
	assert.Equal(t, uint64(9467), result, "Position 4 depth 3 should have 9467 positions")
}

func TestPerft_Position4_Depth4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 4 perft in short mode")
	}
	pos := CreatePositionFormFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 4)
	assert.Equal(t, uint64(422333), result, "Position 4 depth 4 should have 422333 positions")
}

// Position 5 - tests promotion options, discovered check
// rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8
func TestPerft_Position5_Depth1(t *testing.T) {
	pos := CreatePositionFormFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	assert.Equal(t, uint64(44), result, "Position 5 depth 1 should have 44 moves")
}

func TestPerft_Position5_Depth2(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 2 perft in short mode")
	}
	pos := CreatePositionFormFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 2)
	assert.Equal(t, uint64(1486), result, "Position 5 depth 2 should have 1486 positions")
}

func TestPerft_Position5_Depth3(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 3 perft in short mode")
	}
	pos := CreatePositionFormFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 3)
	assert.Equal(t, uint64(62379), result, "Position 5 depth 3 should have 62379 positions")
}

func TestPerft_Position5_Depth4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 4 perft in short mode")
	}
	pos := CreatePositionFormFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 4)
	assert.Equal(t, uint64(2103487), result, "Position 5 depth 4 should have 2103487 positions")
}

// Position 6 - symmetric position, tactical complexity
// r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10
func TestPerft_Position6_Depth1(t *testing.T) {
	pos := CreatePositionFormFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 1)
	assert.Equal(t, uint64(46), result, "Position 6 depth 1 should have 46 moves")
}

func TestPerft_Position6_Depth2(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 2 perft in short mode")
	}
	pos := CreatePositionFormFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 2)
	assert.Equal(t, uint64(2079), result, "Position 6 depth 2 should have 2079 positions")
}

func TestPerft_Position6_Depth3(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 3 perft in short mode")
	}
	pos := CreatePositionFormFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 3)
	assert.Equal(t, uint64(89890), result, "Position 6 depth 3 should have 89890 positions")
}

func TestPerft_Position6_Depth4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping depth 4 perft in short mode")
	}
	pos := CreatePositionFormFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 4)
	assert.Equal(t, uint64(3894594), result, "Position 6 depth 4 should have 3894594 positions")
}

// Divide test for debugging
func TestDivide_InitialPosition(t *testing.T) {
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	divide := pos.Divide(pm, 2)

	// Check that we get the expected structure
	assert.Len(t, divide, 20, "Should have 20 moves in divide")

	// Sum should equal perft(2)
	var total uint64
	for _, nodes := range divide {
		total += nodes
	}
	assert.Equal(t, uint64(400), total, "Divide sum should equal perft(2)")
}
