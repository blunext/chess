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
