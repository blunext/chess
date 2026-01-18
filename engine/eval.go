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
// Fixed: removed penalties for d2/e2 (-20 was bad!), added neutral/positive values
var pawnPST = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0, // rank 1 (impossible)
	0, 0, 0, 5, 5, 0, 0, 0, // rank 2 - neutral, slight bonus for central pawns
	0, 0, 0, 10, 10, 0, 0, 0, // rank 3 - bonus for advanced central pawns
	5, 5, 10, 25, 25, 10, 5, 5, // rank 4 - good advancement
	10, 10, 15, 30, 30, 15, 10, 10, // rank 5 - strong advancement
	20, 20, 25, 35, 35, 25, 20, 20, // rank 6 - very advanced
	50, 50, 50, 50, 50, 50, 50, 50, // rank 7 - about to promote
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

// King Safety constants
const (
	PawnShieldBonus      = 10 // per intact shield pawn
	PawnShieldAdvanced   = 5  // per shield pawn on 3rd rank
	MissingShieldPenalty = 25 // per missing shield pawn
	SemiOpenFilePenalty  = 25 // semi-open file near king
	OpenFilePenalty      = 40 // fully open file near king
	UncastledKingPenalty = 50 // king stuck in center
	NoQueensDivisor      = 4  // reduce king safety importance in endgame
)

// Pawn Structure constants
const (
	DoubledPawnPenalty  = 20 // penalty per doubled pawn
	IsolatedPawnPenalty = 15 // penalty per isolated pawn
	PassedPawnBonus     = 20 // base bonus for passed pawn
	PassedPawnRankBonus = 10 // additional bonus per rank advanced
)

// Space Bonus constants
const (
	CentralPawnBonus    = 25 // bonus for pawns on d4/e4/d5/e5
	ExtendedCenterBonus = 15 // bonus for pawns on c4/c5/f4/f5
	AdvancedPawnRank4   = 5  // bonus for any pawn on rank 4
	AdvancedPawnRank5   = 10 // bonus for any pawn on rank 5
	AdvancedPawnRank6   = 15 // bonus for any pawn on rank 6
)

// File masks for king safety and pawn structure calculations
var fileMasks [8]board.Bitboard

// Adjacent file masks for isolated pawn detection
var adjacentFileMasks [8]board.Bitboard

func init() {
	// Initialize file masks (columns a-h)
	for f := 0; f < 8; f++ {
		var mask board.Bitboard
		for r := 0; r < 8; r++ {
			mask |= board.Bitboard(1) << (r*8 + f)
		}
		fileMasks[f] = mask
	}

	// Initialize adjacent file masks
	for f := 0; f < 8; f++ {
		var mask board.Bitboard
		if f > 0 {
			mask |= fileMasks[f-1]
		}
		if f < 7 {
			mask |= fileMasks[f+1]
		}
		adjacentFileMasks[f] = mask
	}
}

// Evaluate returns the position evaluation in centipawns.
// Positive = white is better, negative = black is better.
// Uses PeSTO tables with tapered eval as base, plus king safety and pawn structure.
func Evaluate(pos board.Position) int {
	// PeSTO provides material + PST with tapered eval
	pestoScore := EvaluatePeSTO(pos)

	// Additional evaluation terms (not in PeSTO)
	// King Safety (scaled by game phase)
	whiteKingSafety := kingSafety(pos, true)
	blackKingSafety := kingSafety(pos, false)

	// Pawn Structure (doubled, isolated, passed pawns)
	whitePawnStructure := pawnStructure(pos, true)
	blackPawnStructure := pawnStructure(pos, false)

	// Space Bonus (central pawns, advanced pawns)
	whiteSpace := spaceBonus(pos, true)
	blackSpace := spaceBonus(pos, false)

	return pestoScore + (whiteKingSafety - blackKingSafety) + (whitePawnStructure - blackPawnStructure) + (whiteSpace - blackSpace)
}

// pawnStructure evaluates pawn structure for a color
func pawnStructure(pos board.Position, isWhite bool) int {
	var ourPawns, enemyPawns board.Bitboard
	if isWhite {
		ourPawns = pos.Pawns & pos.White
		enemyPawns = pos.Pawns & pos.Black
	} else {
		ourPawns = pos.Pawns & pos.Black
		enemyPawns = pos.Pawns & pos.White
	}

	score := 0

	// Doubled pawns: more than one pawn on the same file
	score += doubledPawns(ourPawns)

	// Isolated pawns: pawns with no friendly pawns on adjacent files
	score += isolatedPawns(ourPawns)

	// Passed pawns: pawns with no enemy pawns blocking or attacking
	score += passedPawns(ourPawns, enemyPawns, isWhite)

	return score
}

// spaceBonus evaluates space control through pawn placement
func spaceBonus(pos board.Position, isWhite bool) int {
	var ourPawns board.Bitboard
	if isWhite {
		ourPawns = pos.Pawns & pos.White
	} else {
		ourPawns = pos.Pawns & pos.Black
	}

	score := 0

	// Central pawn bonus: d4/e4/d5/e5
	// Squares: d4=27, e4=28, d5=35, e5=36
	centralSquares := board.Bitboard((1 << 27) | (1 << 28) | (1 << 35) | (1 << 36))
	score += popCount(ourPawns&centralSquares) * CentralPawnBonus

	// Extended center bonus: c4/c5/f4/f5
	// Squares: c4=26, f4=29, c5=34, f5=37
	extendedSquares := board.Bitboard((1 << 26) | (1 << 29) | (1 << 34) | (1 << 37))
	score += popCount(ourPawns&extendedSquares) * ExtendedCenterBonus

	// Advanced pawn bonus (rank 4, 5, 6)
	// For white: ranks 4,5,6 are indices 24-31, 32-39, 40-47
	// For black: ranks 4,5,6 (from black's view) are indices 32-39, 24-31, 16-23
	tempPawns := ourPawns
	for tempPawns != 0 {
		sq := bitScanForward(tempPawns)
		rank := sq >> 3 // 0-7

		var effectiveRank int
		if isWhite {
			effectiveRank = rank // rank 3=rank4, 4=rank5, 5=rank6
		} else {
			effectiveRank = 7 - rank // flip for black
		}

		switch effectiveRank {
		case 3: // rank 4
			score += AdvancedPawnRank4
		case 4: // rank 5
			score += AdvancedPawnRank5
		case 5: // rank 6
			score += AdvancedPawnRank6
		}

		tempPawns &= tempPawns - 1
	}

	return score
}

// doubledPawns returns penalty for doubled pawns
func doubledPawns(pawns board.Bitboard) int {
	penalty := 0
	for f := 0; f < 8; f++ {
		pawnsOnFile := pawns & fileMasks[f]
		count := popCount(pawnsOnFile)
		if count > 1 {
			// Penalty for each extra pawn on the file
			penalty -= (count - 1) * DoubledPawnPenalty
		}
	}
	return penalty
}

// isolatedPawns returns penalty for isolated pawns
func isolatedPawns(pawns board.Bitboard) int {
	penalty := 0
	tempPawns := pawns
	for tempPawns != 0 {
		sq := bitScanForward(tempPawns)
		file := sq & 7

		// Check if there are any friendly pawns on adjacent files
		friendlyOnAdjacent := pawns & adjacentFileMasks[file]
		if friendlyOnAdjacent == 0 {
			penalty -= IsolatedPawnPenalty
		}
		tempPawns &= tempPawns - 1
	}
	return penalty
}

// passedPawns returns bonus for passed pawns
func passedPawns(ourPawns, enemyPawns board.Bitboard, isWhite bool) int {
	bonus := 0
	tempPawns := ourPawns
	for tempPawns != 0 {
		sq := bitScanForward(tempPawns)
		file := sq & 7
		rank := sq >> 3

		if isPassedPawn(sq, file, rank, enemyPawns, isWhite) {
			bonus += PassedPawnBonus
			// Additional bonus based on how advanced the pawn is
			if isWhite {
				bonus += rank * PassedPawnRankBonus // rank 2-7 -> 20-70 extra
			} else {
				bonus += (7 - rank) * PassedPawnRankBonus // rank 6-1 -> 10-60 extra
			}
		}
		tempPawns &= tempPawns - 1
	}
	return bonus
}

// isPassedPawn checks if a pawn has no enemy pawns blocking or attacking its path
func isPassedPawn(sq, file, rank int, enemyPawns board.Bitboard, isWhite bool) bool {
	// Create a mask of squares that enemy pawns would block this pawn
	var blockingMask board.Bitboard

	// Include the file itself and adjacent files
	blockingMask = fileMasks[file]
	if file > 0 {
		blockingMask |= fileMasks[file-1]
	}
	if file < 7 {
		blockingMask |= fileMasks[file+1]
	}

	// Mask to only ranks ahead of the pawn
	if isWhite {
		// For white, block ranks above the pawn (rank+1 to 7)
		for r := rank + 1; r < 8; r++ {
			for f := 0; f < 8; f++ {
				if blockingMask&(board.Bitboard(1)<<(r*8+f)) != 0 {
					// Keep this square in blocking mask
				} else {
					// Not in our files of interest
				}
			}
		}
		// Create rank mask for ranks ahead
		var aheadMask board.Bitboard
		for r := rank + 1; r < 8; r++ {
			for f := 0; f < 8; f++ {
				aheadMask |= board.Bitboard(1) << (r*8 + f)
			}
		}
		blockingMask &= aheadMask
	} else {
		// For black, block ranks below the pawn (0 to rank-1)
		var aheadMask board.Bitboard
		for r := 0; r < rank; r++ {
			for f := 0; f < 8; f++ {
				aheadMask |= board.Bitboard(1) << (r*8 + f)
			}
		}
		blockingMask &= aheadMask
	}

	// If no enemy pawns in the blocking area, it's a passed pawn
	return (enemyPawns & blockingMask) == 0
}

// kingSafety evaluates king safety for the given color
func kingSafety(pos board.Position, isWhite bool) int {
	// Find king position
	var kingBB board.Bitboard
	if isWhite {
		kingBB = pos.Kings & pos.White
	} else {
		kingBB = pos.Kings & pos.Black
	}
	if kingBB == 0 {
		return 0
	}
	kingSq := bitScanForward(kingBB)
	kingFile := kingSq & 7
	kingRank := kingSq >> 3

	score := 0

	// 1. Pawn Shield evaluation
	score += pawnShield(pos, kingSq, isWhite)

	// 2. Open files near king
	score += openFilesNearKing(pos, kingFile)

	// 3. Uncastled king penalty (king on d or e file in middlegame)
	if kingFile == 3 || kingFile == 4 { // d or e file
		if isWhite && kingRank == 0 {
			score -= UncastledKingPenalty
		} else if !isWhite && kingRank == 7 {
			score -= UncastledKingPenalty
		}
	}

	// 4. Scale by game phase (king safety less important without queens)
	var enemyQueens board.Bitboard
	if isWhite {
		enemyQueens = pos.Queens & pos.Black
	} else {
		enemyQueens = pos.Queens & pos.White
	}
	if enemyQueens == 0 {
		score /= NoQueensDivisor
	}

	return score
}

// pawnShield evaluates the pawn shield in front of the king
func pawnShield(pos board.Position, kingSq int, isWhite bool) int {
	kingFile := kingSq & 7
	kingRank := kingSq >> 3
	score := 0

	// Only evaluate pawn shield for castled king positions
	// White: king on rank 0, files f,g,h (kingside) or a,b,c (queenside)
	// Black: king on rank 7, same files
	isCastledPosition := false
	if isWhite && kingRank == 0 && (kingFile >= 5 || kingFile <= 2) {
		isCastledPosition = true
	} else if !isWhite && kingRank == 7 && (kingFile >= 5 || kingFile <= 2) {
		isCastledPosition = true
	}

	if !isCastledPosition {
		return 0
	}

	// Check pawns on the three files around the king
	var ourPawns board.Bitboard
	if isWhite {
		ourPawns = pos.Pawns & pos.White
	} else {
		ourPawns = pos.Pawns & pos.Black
	}

	// Check files: kingFile-1, kingFile, kingFile+1
	for df := -1; df <= 1; df++ {
		f := kingFile + df
		if f < 0 || f > 7 {
			continue
		}

		pawnsOnFile := ourPawns & fileMasks[f]
		if pawnsOnFile == 0 {
			// No pawn on this file - penalty
			score -= MissingShieldPenalty
		} else {
			// Check pawn position
			for pawnsOnFile != 0 {
				pawnSq := bitScanForward(pawnsOnFile)
				pawnRank := pawnSq >> 3

				if isWhite {
					if pawnRank == 1 { // 2nd rank - ideal
						score += PawnShieldBonus
					} else if pawnRank == 2 { // 3rd rank - advanced
						score += PawnShieldAdvanced
					}
				} else {
					if pawnRank == 6 { // 7th rank for black - ideal
						score += PawnShieldBonus
					} else if pawnRank == 5 { // 6th rank for black - advanced
						score += PawnShieldAdvanced
					}
				}
				pawnsOnFile &= pawnsOnFile - 1
			}
		}
	}

	return score
}

// openFilesNearKing penalizes open/semi-open files near the king
func openFilesNearKing(pos board.Position, kingFile int) int {
	score := 0
	allPawns := pos.Pawns

	// Check files around the king (kingFile-1, kingFile, kingFile+1)
	for df := -1; df <= 1; df++ {
		f := kingFile + df
		if f < 0 || f > 7 {
			continue
		}

		whitePawnsOnFile := pos.Pawns & pos.White & fileMasks[f]
		blackPawnsOnFile := pos.Pawns & pos.Black & fileMasks[f]
		anyPawnsOnFile := allPawns & fileMasks[f]

		if anyPawnsOnFile == 0 {
			// Fully open file - big penalty
			score -= OpenFilePenalty
		} else if whitePawnsOnFile == 0 || blackPawnsOnFile == 0 {
			// Semi-open file - smaller penalty
			score -= SemiOpenFilePenalty
		}
	}

	return score
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
