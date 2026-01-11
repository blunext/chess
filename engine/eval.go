package engine

import "chess/board"

// Piece values in centipawns
const (
	PawnValue   = 100
	KnightValue = 320
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900
)

// Evaluate returns the position evaluation in centipawns.
// Positive = white is better, negative = black is better.
func Evaluate(pos board.Position) int {
	white := materialCount(pos, pos.White)
	black := materialCount(pos, pos.Black)
	return white - black
}

// materialCount calculates total material for pieces on given squares.
func materialCount(pos board.Position, color board.Bitboard) int {
	score := 0
	score += popCount(pos.Pawns&color) * PawnValue
	score += popCount(pos.Knights&color) * KnightValue
	score += popCount(pos.Bishops&color) * BishopValue
	score += popCount(pos.Rooks&color) * RookValue
	score += popCount(pos.Queens&color) * QueenValue
	return score
}

// popCount returns the number of set bits in a bitboard.
func popCount(b board.Bitboard) int {
	count := 0
	for b != 0 {
		b &= b - 1
		count++
	}
	return count
}
