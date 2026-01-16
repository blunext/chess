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

// Piece-Square Tables (PST)
// Indexed by square (0=a1, 63=h8), values in centipawns
// Based on Simplified Evaluation Function from Chess Programming Wiki

// pawnPST encourages central pawns and advancement
var pawnPST = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0, // rank 1 (impossible)
	5, 10, 10, -20, -20, 10, 10, 5, // rank 2
	5, -5, -10, 0, 0, -10, -5, 5, // rank 3
	0, 0, 0, 20, 20, 0, 0, 0, // rank 4
	5, 5, 10, 25, 25, 10, 5, 5, // rank 5
	10, 10, 20, 30, 30, 20, 10, 10, // rank 6
	50, 50, 50, 50, 50, 50, 50, 50, // rank 7
	0, 0, 0, 0, 0, 0, 0, 0, // rank 8 (promotion)
}

// knightPST encourages central knights, penalizes edges
var knightPST = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50, // rank 1
	-40, -20, 0, 5, 5, 0, -20, -40, // rank 2
	-30, 5, 10, 15, 15, 10, 5, -30, // rank 3
	-30, 0, 15, 20, 20, 15, 0, -30, // rank 4
	-30, 5, 15, 20, 20, 15, 5, -30, // rank 5
	-30, 0, 10, 15, 15, 10, 0, -30, // rank 6
	-40, -20, 0, 0, 0, 0, -20, -40, // rank 7
	-50, -40, -30, -30, -30, -30, -40, -50, // rank 8
}

// bishopPST encourages central diagonals, penalizes corners
var bishopPST = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20, // rank 1
	-10, 5, 0, 0, 0, 0, 5, -10, // rank 2
	-10, 10, 10, 10, 10, 10, 10, -10, // rank 3
	-10, 0, 10, 10, 10, 10, 0, -10, // rank 4
	-10, 5, 5, 10, 10, 5, 5, -10, // rank 5
	-10, 0, 5, 10, 10, 5, 0, -10, // rank 6
	-10, 0, 0, 0, 0, 0, 0, -10, // rank 7
	-20, -10, -10, -10, -10, -10, -10, -20, // rank 8
}

// rookPST encourages 7th rank and central files
var rookPST = [64]int{
	0, 0, 0, 5, 5, 0, 0, 0, // rank 1
	-5, 0, 0, 0, 0, 0, 0, -5, // rank 2
	-5, 0, 0, 0, 0, 0, 0, -5, // rank 3
	-5, 0, 0, 0, 0, 0, 0, -5, // rank 4
	-5, 0, 0, 0, 0, 0, 0, -5, // rank 5
	-5, 0, 0, 0, 0, 0, 0, -5, // rank 6
	5, 10, 10, 10, 10, 10, 10, 5, // rank 7
	0, 0, 0, 0, 0, 0, 0, 0, // rank 8
}

// queenPST slightly encourages central squares
var queenPST = [64]int{
	-20, -10, -10, -5, -5, -10, -10, -20, // rank 1
	-10, 0, 5, 0, 0, 0, 0, -10, // rank 2
	-10, 5, 5, 5, 5, 5, 0, -10, // rank 3
	0, 0, 5, 5, 5, 5, 0, -5, // rank 4
	-5, 0, 5, 5, 5, 5, 0, -5, // rank 5
	-10, 0, 5, 5, 5, 5, 0, -10, // rank 6
	-10, 0, 0, 0, 0, 0, 0, -10, // rank 7
	-20, -10, -10, -5, -5, -10, -10, -20, // rank 8
}

// kingMiddlegamePST encourages castling, penalizes center
var kingMiddlegamePST = [64]int{
	20, 30, 10, 0, 0, 10, 30, 20, // rank 1 (castled positions)
	20, 20, 0, 0, 0, 0, 20, 20, // rank 2
	-10, -20, -20, -20, -20, -20, -20, -10, // rank 3
	-20, -30, -30, -40, -40, -30, -30, -20, // rank 4
	-30, -40, -40, -50, -50, -40, -40, -30, // rank 5
	-30, -40, -40, -50, -50, -40, -40, -30, // rank 6
	-30, -40, -40, -50, -50, -40, -40, -30, // rank 7
	-30, -40, -40, -50, -50, -40, -40, -30, // rank 8
}

// Evaluate returns the position evaluation in centipawns.
// Positive = white is better, negative = black is better.
func Evaluate(pos board.Position) int {
	white := materialCount(pos, pos.White)
	black := materialCount(pos, pos.Black)

	whitePST := pstScore(pos, pos.White)
	blackPST := pstScore(pos, pos.Black)

	return (white + whitePST) - (black + blackPST)
}

// pstScore calculates piece-square table bonus for a color
func pstScore(pos board.Position, color board.Bitboard) int {
	score := 0
	isWhite := color == pos.White

	score += pstForPieces(pos.Pawns&color, pawnPST, isWhite)
	score += pstForPieces(pos.Knights&color, knightPST, isWhite)
	score += pstForPieces(pos.Bishops&color, bishopPST, isWhite)
	score += pstForPieces(pos.Rooks&color, rookPST, isWhite)
	score += pstForPieces(pos.Queens&color, queenPST, isWhite)
	score += pstForPieces(pos.Kings&color, kingMiddlegamePST, isWhite)

	return score
}

// pstForPieces sums PST values for all pieces in a bitboard
func pstForPieces(pieces board.Bitboard, table [64]int, isWhite bool) int {
	score := 0
	for pieces != 0 {
		sq := bitScanForward(pieces)
		if isWhite {
			score += table[sq]
		} else {
			// Mirror vertically for black (rank 1 <-> rank 8)
			score += table[sq^56]
		}
		pieces &= pieces - 1 // clear LSB
	}
	return score
}

// bitScanForward returns index of least significant set bit
func bitScanForward(b board.Bitboard) int {
	// De Bruijn multiplication
	const debruijn64 = 0x03f79d71b4cb0a89
	var index = [64]int{
		0, 1, 48, 2, 57, 49, 28, 3,
		61, 58, 50, 42, 38, 29, 17, 4,
		62, 55, 59, 36, 53, 51, 43, 22,
		45, 39, 33, 30, 24, 18, 12, 5,
		63, 47, 56, 27, 60, 41, 37, 16,
		54, 35, 52, 21, 44, 32, 23, 11,
		46, 26, 40, 15, 34, 20, 31, 10,
		25, 14, 19, 9, 13, 8, 7, 6,
	}
	return index[((uint64(b)&-uint64(b))*debruijn64)>>58]
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
