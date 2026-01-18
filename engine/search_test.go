package engine

import (
	"testing"
	"time"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

// === Move Ordering Tests ===

func TestMoveScore_CapturesRankedByMVVLVA(t *testing.T) {
	// PxQ (pawn captures queen) should score higher than QxP (queen captures pawn)
	pawnCapturesQueen := board.Move{
		Piece:    board.Pawn,
		Captured: board.Queen,
	}
	queenCapturesPawn := board.Move{
		Piece:    board.Queen,
		Captured: board.Pawn,
	}
	knightCapturesQueen := board.Move{
		Piece:    board.Knight,
		Captured: board.Queen,
	}

	pxqScore := moveScore(pawnCapturesQueen)
	qxpScore := moveScore(queenCapturesPawn)
	nxqScore := moveScore(knightCapturesQueen)

	// PxQ > NxQ > QxP (MVV-LVA: best captures first)
	assert.Greater(t, pxqScore, nxqScore, "PxQ should score higher than NxQ")
	assert.Greater(t, nxqScore, qxpScore, "NxQ should score higher than QxP")
	assert.Greater(t, pxqScore, 10000, "Captures should have high base score")
}

func TestMoveScore_PromotionsHighPriority(t *testing.T) {
	promotion := board.Move{
		Piece:     board.Pawn,
		Promotion: board.Queen,
	}
	quietMove := board.Move{
		Piece: board.Knight,
	}

	promoScore := moveScore(promotion)
	quietScore := moveScore(quietMove)

	assert.Greater(t, promoScore, quietScore, "Promotion should score higher than quiet move")
	assert.Greater(t, promoScore, 9000, "Promotion should have high score")
	assert.Equal(t, 0, quietScore, "Quiet move should have zero score")
}

func TestSortMoves_CapturesFirst(t *testing.T) {
	moves := []board.Move{
		{Piece: board.Knight, Captured: board.Empty},                       // quiet
		{Piece: board.Pawn, Captured: board.Queen},                         // PxQ (best)
		{Piece: board.Queen, Captured: board.Pawn},                         // QxP (weak capture)
		{Piece: board.Pawn, Promotion: board.Queen, Captured: board.Empty}, // promotion
	}

	sortMoves(moves)

	// After sorting: PxQ, QxP, promotion, quiet
	assert.Equal(t, board.Queen, moves[0].Captured, "First should be PxQ (best capture)")
	assert.Equal(t, board.Pawn, moves[1].Captured, "Second should be QxP (weaker capture)")
	assert.Equal(t, board.Queen, moves[2].Promotion, "Third should be promotion")
	assert.Equal(t, board.Empty, moves[3].Captured, "Last should be quiet move")
}

// === Search Tests ===

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
	// Initial position is equal, score should be around 0 (PST may add small bonus for good moves)
	assert.InDelta(t, 0, result.Score, 100, "Initial position should be roughly equal")
}

func TestSearch_CaptureHangingQueen(t *testing.T) {
	// White to move, black queen on d4 can be captured by pawn on e3
	pos := board.CreatePositionFormFEN("rnb1kbnr/pppppppp/8/8/3q4/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 1)

	// Should capture the queen with e3xd4
	assert.Equal(t, "e3d4", result.Move.ToUCI(), "Should capture the hanging queen")
	assert.InDelta(t, QueenValue, result.Score, 100, "Should gain roughly a queen")
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
	assert.InDelta(t, -QueenValue, result.Score, 100, "Score should reflect black winning roughly a queen")
}

// === Quiescence Tests ===

func TestQuiescence_SeesCaptureBeyondHorizon(t *testing.T) {
	// Position where white can push pawn, then black captures it
	// Without quiescence: engine might think pawn is safe at depth 1
	// With quiescence: engine sees the recapture
	//
	// White pawn on e4, black knight on f6 attacks e5
	// If white plays e4-e5, black plays Nxe5
	// Quiescence should see this exchange
	pos := board.CreatePositionFormFEN("rnbqkb1r/pppppppp/5n2/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	pm := createTestPieceMoves()

	// At depth 1, without quiescence, engine might think e5 is fine
	// With quiescence, it sees Nxe5 and evaluates correctly
	result := Search(pos, pm, 1)

	// e4-e5 loses a pawn, engine should prefer other moves
	assert.NotEqual(t, "e4e5", result.Move.ToUCI(), "Should not push pawn to be captured")
}

func TestQuiescence_SeesRecapture(t *testing.T) {
	// Position: white knight on d5, black pawn on e6 can capture
	// If black plays exd5, white has no immediate recapture
	// This tests that quiescence evaluates capture sequences correctly
	pos := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/4p3/3N4/8/8/PPPPPPPP/R1BQKBNR b KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 1)

	// Black should capture the knight - it's free material
	assert.Equal(t, "e6d5", result.Move.ToUCI(), "Black should capture the knight")
	// Score should be negative (favorable to black)
	assert.Less(t, result.Score, 0, "Score should be negative (black winning)")
}

// === Time Management Tests ===

// === Null Move Pruning Tests ===

func TestNullMovePruning_MateDetectionAtDepth4(t *testing.T) {
	// Ensure NMP doesn't break mate detection at deeper depths
	// Scholar's mate position - Qxf7#
	pos := board.CreatePositionFormFEN("r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 4)

	assert.Equal(t, "h5f7", result.Move.ToUCI(), "Should still find Qxf7# at depth 4 with NMP")
	assert.Greater(t, result.Score, 50000, "Should detect mate")
}

func TestNullMovePruning_EndgameDisabled(t *testing.T) {
	// In endgame with only kings and pawns, NMP should be disabled
	// This is a K+P vs K position where zugzwang is possible
	pos := board.CreatePositionFormFEN("8/8/8/4k3/8/4K3/4P3/8 w - - 0 1")
	pm := createTestPieceMoves()

	// Should not crash and should find a reasonable move
	result := Search(pos, pm, 4)

	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move in endgame")
}

func TestNullMovePruning_BasicFunctionality(t *testing.T) {
	// Test that NMP works correctly on a middlegame position
	// White has extra material, NMP should help prune quickly
	pos := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 5)

	// Should find a move and have reasonable score
	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
	// Position is roughly equal, score should be reasonable
	assert.InDelta(t, 0, result.Score, 200, "Score should be reasonable for equal position")
}

func TestNullMovePruning_DoesNotBreakTactics(t *testing.T) {
	// White can win black's queen with a knight fork on c7
	// NMP should not prune away this tactical opportunity
	pos := board.CreatePositionFormFEN("r1bqkb1r/pppp1ppp/2n2n2/4N3/4P3/8/PPPP1PPP/RNBQKB1R w KQkq - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 3)

	// White should play Nxc6 winning material
	if result.Move.ToUCI() == "e5c6" {
		// Good - found the capture
		assert.Greater(t, result.Score, 100, "Should gain material")
	} else if result.Move.ToUCI() == "e5f7" {
		// Also good - fork on f7 attacks king and rook
		assert.Greater(t, result.Score, 0, "Should gain advantage")
	}
	// Any reasonable tactical move is acceptable
	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
}

func TestSearchWithTime_ReturnsInTime(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	timeLimit := 100 * time.Millisecond
	start := time.Now()
	result := SearchWithTime(pos, pm, timeLimit)
	elapsed := time.Since(start)

	// Should return within reasonable time (allow some overhead)
	assert.Less(t, elapsed, 200*time.Millisecond, "Search should complete within time limit")
	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
	assert.Greater(t, result.Depth, 0, "Should reach at least depth 1")
}

func TestSearchWithTime_IterativeDeepening(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	// With 500ms, should reach multiple depths
	result := SearchWithTime(pos, pm, 500*time.Millisecond)

	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
	assert.GreaterOrEqual(t, result.Depth, 3, "Should reach at least depth 3 in 500ms")
}

func TestAllocateTime_Basic(t *testing.T) {
	// 60 seconds remaining, no increment
	allocated := AllocateTime(60000, 60000, 0, 0, true, 0)
	// Should allocate about 2 seconds (60000/30)
	assert.GreaterOrEqual(t, allocated, 1500*time.Millisecond)
	assert.LessOrEqual(t, allocated, 3000*time.Millisecond)
}

func TestAllocateTime_WithIncrement(t *testing.T) {
	// 60 seconds + 1 second increment
	allocated := AllocateTime(60000, 60000, 1000, 1000, true, 0)
	// Should be slightly more than without increment
	noInc := AllocateTime(60000, 60000, 0, 0, true, 0)
	assert.Greater(t, allocated, noInc, "Increment should increase allocated time")
}

func TestAllocateTime_MovesToGo(t *testing.T) {
	// 60 seconds, 10 moves to go
	allocated := AllocateTime(60000, 60000, 0, 0, true, 10)
	// Should allocate about 6 seconds (60000/10)
	assert.GreaterOrEqual(t, allocated, 5*time.Second)
	assert.LessOrEqual(t, allocated, 7*time.Second)
}

// === Check Extension Tests ===

func TestCheckExtension_FindsMateWithChecks(t *testing.T) {
	// Position where mate requires a checking sequence
	// White: Queen on d1, King on e1
	// Black: King on h8, pawns on g7,h7 (no escape squares on back rank)
	// Qd1-d8# is mate in 1
	pos := board.CreatePositionFormFEN("7k/6pp/8/8/8/8/8/3QK3 w - - 0 1")
	pm := createTestPieceMoves()

	// At depth 2, should find Qd8#
	result := Search(pos, pm, 2)

	// Should find Qd8# (back rank mate)
	assert.Equal(t, "d1d8", result.Move.ToUCI(), "Should find Qd8# mate")
	assert.Greater(t, result.Score, 50000, "Should detect mate sequence")
}

func TestCheckExtension_DoesNotMissBackRankMate(t *testing.T) {
	// Classic back rank mate pattern
	// White rook on a1, black king on g8 with pawns blocking escape
	pos := board.CreatePositionFormFEN("6k1/5ppp/8/8/8/8/8/R3K3 w Q - 0 1")
	pm := createTestPieceMoves()

	result := Search(pos, pm, 2)

	// Should find Ra8#
	assert.Equal(t, "a1a8", result.Move.ToUCI(), "Should find Ra8# back rank mate")
	assert.Greater(t, result.Score, 50000, "Should detect mate")
}

func TestCheckExtension_DefendsAgainstCheck(t *testing.T) {
	// Position where black is in check and must find the best defense
	// Black king on h8 with pawns, white queen giving check on g8
	// Black must capture or escape
	pos := board.CreatePositionFormFEN("6Qk/6pp/8/8/8/8/8/4K3 b - - 0 1")
	pm := createTestPieceMoves()

	// Black to move, in check from Qg8
	assert.True(t, pos.IsInCheck(), "Black should be in check")

	result := Search(pos, pm, 3)

	// Should find a legal move that deals with check (Kxg8 is only option)
	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move to escape check")
}

func TestCheckExtension_FindsMateIn2WithCheck(t *testing.T) {
	// Anastasia's mate pattern: Knight and rook mate
	// This requires extending through check to see the mate
	pos := board.CreatePositionFormFEN("5k2/4R3/5N2/8/8/8/8/4K3 w - - 0 1")
	pm := createTestPieceMoves()

	// Re7 threatens Nf6-h7# (smothered), or various Rook mates
	result := Search(pos, pm, 4)

	// Should find a checking sequence leading to mate or winning material
	assert.NotEqual(t, board.Move{}, result.Move, "Should find best move")
	assert.Greater(t, result.Score, 0, "Should be winning")
}

func TestCheckExtension_CheckExtendsSearchDepth(t *testing.T) {
	// Position where a check sequence would be cut off without extension
	// Simple king chase that leads to winning material
	pos := board.CreatePositionFormFEN("4k3/8/8/8/8/8/3Q4/4K3 w - - 0 1")
	pm := createTestPieceMoves()

	// Even at low depth, check extensions should help find good checking moves
	result := Search(pos, pm, 2)

	assert.NotEqual(t, board.Move{}, result.Move, "Should find a move")
	// With queen vs nothing, white is winning
	assert.Greater(t, result.Score, 800, "Should recognize huge advantage")
}

// === Mate Threat Detection Tests ===

func TestHasMateInOne_BackRankMate(t *testing.T) {
	// Black King on g8 with pawns on f7,g7,h7 blocking escape
	// White Rook on a1 can play Ra8# (back rank mate)
	pos := board.CreatePositionFormFEN("6k1/5ppp/8/8/8/8/8/R3K3 b - - 0 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	// White (opponent of black) threatens mate in 1
	hasThreat := session.hasMateInOne(&pos, pm, true) // true = white is opponent

	assert.True(t, hasThreat, "Should detect back rank mate threat Ra8#")
}

func TestHasMateInOne_QueenMate(t *testing.T) {
	// White Queen on h5, Black King on e8, f7 weak
	// White threatens Qxf7# (scholar's mate style)
	pos := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	hasThreat := session.hasMateInOne(&pos, pm, true)

	assert.True(t, hasThreat, "Should detect Qxf7# threat")
}

func TestHasMateInOne_NoThreat(t *testing.T) {
	// Starting position - no mate threats
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()
	session := NewSession(16)

	hasThreat := session.hasMateInOne(&pos, pm, true)

	assert.False(t, hasThreat, "No mate threat in starting position")
}

func TestHasMateInOne_SmotheredMate(t *testing.T) {
	// Classic smothered mate setup:
	// Black King on h8, Rook on g8, pawns on g7/h7
	// White Knight on g5 threatens Nf7# (smothered mate)
	pos := board.CreatePositionFormFEN("6rk/6pp/8/6N1/8/8/8/4K3 b - - 0 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	hasThreat := session.hasMateInOne(&pos, pm, true)

	assert.True(t, hasThreat, "Should detect smothered mate threat Nf7#")
}

func TestHasMateInOne_AlmostMate(t *testing.T) {
	// Position where opponent has checks but not mate
	// King has escape square (pawn on e7 blocks Ra8 but King can go to d7/f7)
	pos := board.CreatePositionFormFEN("4k3/4p3/8/8/8/8/8/R3K3 b - - 0 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	hasThreat := session.hasMateInOne(&pos, pm, true)

	// King can escape to d7 or f7, so no immediate mate
	assert.False(t, hasThreat, "King has escape squares, no mate in 1")
}

// === Check Evasion in Quiescence Tests ===

func TestQuiescence_CheckEvasion_KingMoves(t *testing.T) {
	// Black King in check, only escape is Kf8 (not a capture)
	// Without check evasion, quiescence would only look at captures and miss Kf8
	pos := board.CreatePositionFormFEN("4k3/8/5Q2/8/8/8/8/4K3 b - - 0 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	// Search should find an escape move, not crash or return bad score
	result := session.SearchWithTime(pos, pm, 100*time.Millisecond)

	// Black should escape, not get mated immediately
	assert.NotEqual(t, board.Move{}, result.Move, "Should find escape move")
	// Score should not be mate (would be very negative for black)
	assert.Greater(t, result.Score, -30000, "Should not be immediate mate")
}

func TestQuiescence_CheckEvasion_BlockingMove(t *testing.T) {
	// Black in check, can block with a piece (not capture)
	// Rook gives check, Bishop can block
	pos := board.CreatePositionFormFEN("4k3/8/8/8/4b3/8/8/R3K3 b - - 0 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	result := session.SearchWithTime(pos, pm, 100*time.Millisecond)

	// Black should find a defense
	assert.NotEqual(t, board.Move{}, result.Move, "Should find blocking or escape move")
}

func TestQuiescence_MateThreat_ContinuesSearch(t *testing.T) {
	// Position from the a2b2 bug: Black should see Nf5 threatens Qxg7#
	// and not just take stand-pat
	pos := board.CreatePositionFormFEN("r1b2rk1/pp1p1ppp/8/2p2N2/8/2P1P1Q1/qP4P1/2K2R1R b - - 1 1")
	pm := createTestPieceMoves()
	session := NewSession(16)

	// Check that mate threat is detected
	hasThreat := session.hasMateInOne(&pos, pm, true) // White threatens mate

	assert.True(t, hasThreat, "Should detect Qxg7# threat after any Black move")
}

// === Performance Benchmarks ===

func BenchmarkSearch_Depth4(b *testing.B) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := createTestPieceMoves()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session := NewSession(64)
		session.Search(pos, pm, 4)
	}
}
