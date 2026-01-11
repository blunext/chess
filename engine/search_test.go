package engine

import (
	"testing"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

func createTestPieceMoves() board.PieceMoves {
	pm := make(board.PieceMoves)

	// Knight moves (all 8 L-shaped jumps from each square)
	pm[board.Knight] = make(board.SquareMoves)
	knightOffsets := [][2]int{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	for sq := 0; sq < 64; sq++ {
		file := sq & 7
		rank := sq >> 3
		var targets []board.Bitboard
		for _, off := range knightOffsets {
			newFile := file + off[0]
			newRank := rank + off[1]
			if newFile >= 0 && newFile < 8 && newRank >= 0 && newRank < 8 {
				targets = append(targets, board.IndexToBitBoard(newRank*8+newFile))
			}
		}
		pm[board.Knight][board.IndexToBitBoard(sq)] = [][]board.Bitboard{targets}
	}

	// King moves (all 8 directions from each square)
	pm[board.King] = make(board.SquareMoves)
	kingOffsets := [][2]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}
	for sq := 0; sq < 64; sq++ {
		file := sq & 7
		rank := sq >> 3
		var targets []board.Bitboard
		for _, off := range kingOffsets {
			newFile := file + off[0]
			newRank := rank + off[1]
			if newFile >= 0 && newFile < 8 && newRank >= 0 && newRank < 8 {
				targets = append(targets, board.IndexToBitBoard(newRank*8+newFile))
			}
		}
		pm[board.King][board.IndexToBitBoard(sq)] = [][]board.Bitboard{targets}
	}

	return pm
}

func TestSearch_InitialPosition(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	result := Search(pos, pm, 1)

	// Should find some move
	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
	// Initial position is equal, score should be around 0
	assert.Equal(t, 0, result.Score, "Initial position should be equal")
}

func TestSearch_CaptureHangingQueen(t *testing.T) {
	// White to move, black queen on d4 can be captured by pawn on e3
	pos := board.CreatePositionFormFEN("rnb1kbnr/pppppppp/8/8/3q4/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 1)

	// Should capture the queen with e3xd4
	assert.Equal(t, "e3d4", result.Move.ToUCI(), "Should capture the hanging queen")
	assert.Equal(t, QueenValue, result.Score, "Should gain a queen")
}

func TestSearch_AvoidLosingQueen(t *testing.T) {
	// White queen on d4 attacked by black pawn on e5, white to move
	// White should move the queen away
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/8/4p3/3Q4/8/PPPPPPPP/RNB1KBNR w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 2)

	// Should not leave queen on d4 to be captured
	assert.NotEqual(t, board.Move{}, result.Move)
	// After best play, white should not be down a queen
	assert.GreaterOrEqual(t, result.Score, 0, "Should not lose the queen")
}

func TestSearch_MateInOne(t *testing.T) {
	// White to move, Qxf7# is mate (scholar's mate pattern)
	pos := board.CreatePositionFormFEN("r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 2) // Need depth 2 to see mate

	// Should find Qxf7#
	assert.Equal(t, "h5f7", result.Move.ToUCI(), "Should find Qxf7#")
	assert.Greater(t, result.Score, 50000, "Mate score should be very high")
}

func TestSearch_Depth2(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	result := Search(pos, pm, 2)

	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move at depth 2")
}

func TestSearch_BlackToMove(t *testing.T) {
	// Black to move, can capture white queen on d4
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/8/4p3/3Q4/8/PPPPPPPP/RNB1KBNR b KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 1)

	// Should capture the queen with e5xd4
	assert.Equal(t, "e5d4", result.Move.ToUCI(), "Black should capture the queen")
	assert.Equal(t, -QueenValue, result.Score, "Score should reflect black winning queen")
}
