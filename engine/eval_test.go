package engine

import (
	"testing"

	"chess/board"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_InitialPosition(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Initial position should be equal material")
}

func TestEvaluate_WhiteMissingPawn(t *testing.T) {
	// Initial position but white is missing e2 pawn
	// Material diff = -100, PST may vary slightly due to missing central pawn
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, -PawnValue, eval, 50, "White missing pawn should be around -100")
}

func TestEvaluate_BlackMissingPawn(t *testing.T) {
	// Initial position but black is missing e7 pawn
	pos := board.CreatePositionFormFEN("rnbqkbnr/pppp1ppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, PawnValue, eval, 50, "Black missing pawn should be around +100")
}

func TestEvaluate_WhiteUpKnight(t *testing.T) {
	// White has extra knight (black missing g8 knight)
	pos := board.CreatePositionFormFEN("rnbqkb1r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, KnightValue, eval, 50, "White up a knight should be around +320")
}

func TestEvaluate_WhiteUpRook(t *testing.T) {
	// Black missing a8 rook
	pos := board.CreatePositionFormFEN("1nbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQk - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, RookValue, eval, 50, "White up a rook should be around +500")
}

func TestEvaluate_WhiteUpQueen(t *testing.T) {
	// Black missing queen
	pos := board.CreatePositionFormFEN("rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.InDelta(t, QueenValue, eval, 50, "White up a queen should be around +900")
}

func TestEvaluate_ComplexPosition(t *testing.T) {
	// FEN: 4kb2/pppppp2/8/8/8/8/PPPPP3/RNBQK3
	// White: K, Q, R, B, N, 5P = 900 + 500 + 330 + 320 + 500 = 2550
	// Black: K, B, 6P = 330 + 600 = 930
	// Material diff = 2550 - 930 = 1620, PST will adjust this
	pos := board.CreatePositionFormFEN("4kb2/pppppp2/8/8/8/8/PPPPP3/RNBQK3 w Q - 0 1")
	eval := Evaluate(pos)

	whiteExpected := QueenValue + RookValue + BishopValue + KnightValue + 5*PawnValue
	blackExpected := BishopValue + 6*PawnValue
	expected := whiteExpected - blackExpected

	assert.InDelta(t, expected, eval, 100, "Complex position should be around material diff")
}

func TestEvaluate_OnlyKings(t *testing.T) {
	pos := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Just kings should be 0")
}

func TestPopCount(t *testing.T) {
	assert.Equal(t, 0, popCount(0))
	assert.Equal(t, 1, popCount(1))
	assert.Equal(t, 1, popCount(0x8000000000000000))
	assert.Equal(t, 8, popCount(0xFF))
	assert.Equal(t, 64, popCount(0xFFFFFFFFFFFFFFFF))
}

// PST-specific tests

func TestPST_KnightCenterBonus(t *testing.T) {
	// Knight on e4 (central) vs knight on a1 (corner)
	// e4 = index 28, a1 = index 0
	centerKnight := board.CreatePositionFormFEN("4k3/8/8/8/4N3/8/8/4K3 w - - 0 1")
	cornerKnight := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/N3K3 w - - 0 1")

	centerEval := Evaluate(centerKnight)
	cornerEval := Evaluate(cornerKnight)

	// Central knight should be worth more than corner knight
	assert.Greater(t, centerEval, cornerEval, "Knight on e4 should score higher than knight on a1")
	// Difference should be significant (central bonus vs edge penalty)
	assert.Greater(t, centerEval-cornerEval, 50, "PST difference should be significant")
}

func TestPST_PawnAdvancement(t *testing.T) {
	// Pawn on e6 (advanced) vs pawn on e2 (starting)
	advancedPawn := board.CreatePositionFormFEN("4k3/8/4P3/8/8/8/8/4K3 w - - 0 1")
	startingPawn := board.CreatePositionFormFEN("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1")

	advancedEval := Evaluate(advancedPawn)
	startingEval := Evaluate(startingPawn)

	// Advanced pawn should be worth more
	assert.Greater(t, advancedEval, startingEval, "Advanced pawn should score higher")
}

func TestPST_RookOnSeventhRank(t *testing.T) {
	// Rook on 7th rank vs rook on 1st rank
	seventhRank := board.CreatePositionFormFEN("4k3/4R3/8/8/8/8/8/4K3 w - - 0 1")
	firstRank := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/R3K3 w - - 0 1")

	seventhEval := Evaluate(seventhRank)
	firstEval := Evaluate(firstRank)

	// Rook on 7th should be worth more
	assert.Greater(t, seventhEval, firstEval, "Rook on 7th rank should score higher")
}

func TestPST_KingSafetyMidgame(t *testing.T) {
	// King in castled position (g1) vs king in center (e4)
	castledKing := board.CreatePositionFormFEN("4k3/8/8/8/8/8/8/5RK1 w - - 0 1")
	centerKing := board.CreatePositionFormFEN("8/4k3/8/8/4K3/8/8/5R2 w - - 0 1")

	castledEval := Evaluate(castledKing)
	centerEval := Evaluate(centerKing)

	// Castled king should score better in middlegame PST
	assert.Greater(t, castledEval, centerEval, "Castled king should score higher than central king")
}

func TestPST_SymmetricPosition(t *testing.T) {
	// Initial position should be exactly 0 due to symmetry
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Symmetric initial position should evaluate to 0")
}

func TestPST_BlackMirroring(t *testing.T) {
	// White knight on e4 should have same PST value as black knight on e5
	// (mirrored position - both are centrally placed for their side)
	whiteKnightE4 := board.CreatePositionFormFEN("4k3/8/8/8/4N3/8/8/4K3 w - - 0 1")
	blackKnightE5 := board.CreatePositionFormFEN("4k3/8/8/4n3/8/8/8/4K3 w - - 0 1")

	whiteEval := Evaluate(whiteKnightE4)
	blackEval := Evaluate(blackKnightE5)

	// Should be roughly opposite (white positive, black negative, similar magnitude)
	// e4 for white = e5 for black (both are rank 4 from their perspective)
	assert.InDelta(t, whiteEval, -blackEval, 20, "Mirrored knight positions should have opposite scores")
}

// === King Safety Tests ===

func TestKingSafety_CastledKingWithPawnShield(t *testing.T) {
	// White king castled kingside with full pawn shield (f2,g2,h2)
	pos := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 1")

	whiteKS := kingSafety(pos, true)
	blackKS := kingSafety(pos, false)

	// Both have castled with intact pawn shields - should have similar scores
	assert.InDelta(t, whiteKS, blackKS, 30, "Both castled kings with pawn shields should have similar safety")
}

func TestKingSafety_BrokenPawnShield(t *testing.T) {
	// White king castled but g2 pawn is missing (g-file open)
	brokenShield := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5NP1/PPPP1P1P/RNBQ1RK1 w - - 0 1")

	// Compare with intact shield
	intactShield := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 1")

	brokenKS := kingSafety(brokenShield, true)
	intactKS := kingSafety(intactShield, true)

	// Broken shield should be worse (pawn on g3 instead of g2)
	// Note: g3 is advanced, f2/h2 still intact
	assert.LessOrEqual(t, brokenKS, intactKS, "Broken pawn shield should give lower score")
}

func TestKingSafety_UncastledKingPenalty(t *testing.T) {
	// White king still in center (e1)
	uncastled := board.CreatePositionFormFEN("r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1")

	// Castled position
	castled := board.CreatePositionFormFEN("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 0 1")

	uncastledKS := kingSafety(uncastled, true)
	castledKS := kingSafety(castled, true)

	// Uncastled king should have worse safety
	assert.Less(t, uncastledKS, castledKS, "Uncastled king should have worse safety than castled")
}

func TestKingSafety_OpenFileNearKing(t *testing.T) {
	// King castled but h-file is completely open (no pawns)
	openFile := board.CreatePositionFormFEN("r1bq1rk1/pppp1pp1/2n2n2/2b1p2p/2B1P3/5N2/PPPP1PP1/RNBQ1RK1 w - - 0 1")

	ks := kingSafety(openFile, true)

	// Should have penalty for open file
	// Exact value depends on other factors, but should be negative overall
	assert.Less(t, ks, 0, "Open file near king should result in negative safety score")
}

func TestKingSafety_EndgameScaling(t *testing.T) {
	// Position without queens - king safety should be reduced
	noQueens := board.CreatePositionFormFEN("r1b2rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNB2RK1 w - - 0 1")

	// Position with queens
	withQueens := board.CreatePositionFormFEN("r1bq1rk1/pppp1ppp/2n2n2/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 1")

	noQueensKS := kingSafety(noQueens, true)
	withQueensKS := kingSafety(withQueens, true)

	// Without queens, king safety should matter less
	// The absolute value should be smaller
	absNoQueens := noQueensKS
	if absNoQueens < 0 {
		absNoQueens = -absNoQueens
	}
	absWithQueens := withQueensKS
	if absWithQueens < 0 {
		absWithQueens = -absWithQueens
	}

	// Without enemy queen, king safety penalty/bonus is divided by 4
	assert.LessOrEqual(t, absNoQueens, absWithQueens, "King safety should matter less without queens")
}

func TestKingSafety_PawnShieldFunction(t *testing.T) {
	// Test pawnShield directly
	// King on g1 with perfect shield (f2,g2,h2)
	perfectShield := board.CreatePositionFormFEN("8/8/8/8/8/8/5PPP/6K1 w - - 0 1")
	kingSq := 6 // g1

	score := pawnShield(perfectShield, kingSq, true)

	// Should have bonus for all three pawns
	assert.Greater(t, score, 0, "Perfect pawn shield should give positive score")
}

func TestKingSafety_Queensidecastle(t *testing.T) {
	// White king castled queenside (c1) with pawn shield (a2,b2,c2)
	queenside := board.CreatePositionFormFEN("r3kbnr/ppp2ppp/2nqb3/3pp3/3PP3/2NQB3/PPP2PPP/2KR1BNR w kq - 0 1")

	ks := kingSafety(queenside, true)

	// Should still evaluate pawn shield for queenside castling
	// With intact pawns on a2,b2,c2
	assert.GreaterOrEqual(t, ks, -50, "Queenside castled king with shield should not be heavily penalized")
}
