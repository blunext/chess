//go:build slow

package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Slow perft tests (depth 5-6) - run with: go test -tags slow ./board/
// These tests run in parallel and take significant time to complete.

func TestPerft_InitialPosition_Depth5(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 5)
	assert.Equal(t, uint64(4865609), result, "Initial position depth 5 should have 4865609 positions")
}

func TestPerft_InitialPosition_Depth6(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN(InitialPosition)
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 6)
	assert.Equal(t, uint64(119060324), result, "Initial position depth 6 should have 119060324 positions")
}

func TestPerft_Kiwipete_Depth5(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 5)
	assert.Equal(t, uint64(193690690), result, "Kiwipete depth 5 should have 193690690 positions")
}

func TestPerft_Position3_Depth5(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 5)
	assert.Equal(t, uint64(674624), result, "Position 3 depth 5 should have 674624 positions")
}

func TestPerft_Position3_Depth6(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 6)
	assert.Equal(t, uint64(11030083), result, "Position 3 depth 6 should have 11030083 positions")
}

func TestPerft_Position4_Depth5(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 5)
	assert.Equal(t, uint64(15833292), result, "Position 4 depth 5 should have 15833292 positions")
}

func TestPerft_Position5_Depth5(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 5)
	assert.Equal(t, uint64(89941194), result, "Position 5 depth 5 should have 89941194 positions")
}

func TestPerft_Position6_Depth5(t *testing.T) {
	t.Parallel()
	pos := CreatePositionFormFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10")
	pm := createTestPieceMoves()

	result := pos.Perft(pm, 5)
	assert.Equal(t, uint64(164075551), result, "Position 6 depth 5 should have 164075551 positions")
}
