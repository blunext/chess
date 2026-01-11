package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate_InitialPosition(t *testing.T) {
	pos := CreatePositionFormFEN(InitialPosition)
	eval := Evaluate(pos)
	assert.Equal(t, 0, eval, "Initial position should be equal material")
}

func TestEvaluate_WhiteMissingPawn(t *testing.T) {
	// Initial position but white is missing e2 pawn
	pos := CreatePositionFormFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, -PawnValue, eval, "White missing pawn should be -100")
}

func TestEvaluate_BlackMissingPawn(t *testing.T) {
	// Initial position but black is missing e7 pawn
	pos := CreatePositionFormFEN("rnbqkbnr/pppp1ppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, PawnValue, eval, "Black missing pawn should be +100")
}

func TestEvaluate_WhiteUpKnight(t *testing.T) {
	// White has extra knight (black missing g8 knight)
	pos := CreatePositionFormFEN("rnbqkb1r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, KnightValue, eval, "White up a knight should be +320")
}

func TestEvaluate_WhiteUpRook(t *testing.T) {
	// Black missing a8 rook
	pos := CreatePositionFormFEN("1nbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQk - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, RookValue, eval, "White up a rook should be +500")
}

func TestEvaluate_WhiteUpQueen(t *testing.T) {
	// Black missing queen
	pos := CreatePositionFormFEN("rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := Evaluate(pos)
	assert.Equal(t, QueenValue, eval, "White up a queen should be +900")
}

func TestEvaluate_ComplexPosition(t *testing.T) {
	// FEN: 4kb2/pppppp2/8/8/8/8/PPPPP3/RNBQK3
	// White: K, Q, R, B, N, 5P = 900 + 500 + 330 + 320 + 500 = 2550
	// Black: K, B, 6P = 330 + 600 = 930
	// Diff = 2550 - 930 = 1620
	pos := CreatePositionFormFEN("4kb2/pppppp2/8/8/8/8/PPPPP3/RNBQK3 w Q - 0 1")
	eval := Evaluate(pos)

	whiteExpected := QueenValue + RookValue + BishopValue + KnightValue + 5*PawnValue
	blackExpected := BishopValue + 6*PawnValue
	expected := whiteExpected - blackExpected

	assert.Equal(t, expected, eval, "Complex position material count")
}

func TestEvaluate_OnlyKings(t *testing.T) {
	pos := CreatePositionFormFEN("4k3/8/8/8/8/8/8/4K3 w - - 0 1")
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
