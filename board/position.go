package board

import (
	"fmt"
	"log"
	"math/bits"

	"chess/magic"
)

// SquareMoves maps each square (Bitboard representation) to a list of move sequences.
// Each move sequence is a slice of Bitboards representing a path that a piece can take from that square.
type SquareMoves map[Bitboard][][]Bitboard

// PieceMoves associates each chess piece with its possible moves on the board.
// It uses the SquareMoves to represent all legal moves for a piece from any square, considering the rules of movement unique to each piece.
type PieceMoves map[Piece]SquareMoves

const (
	FileA int = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
	FileNB
)

const (
	Rank1 int = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	RankOut
)

type Piece uint8
type Color uint8

const (
	Empty Piece = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	ColorWhite Color = 0
	ColorBlack Color = 1
)
const (
	CastleWhiteKingSide = 1 << iota
	CastleWhiteQueenSide
	CastleBlackKingSide
	CastleBlackQueenSide
)

type Position struct {
	Pawns, Knights, Bishops, Rooks, Queens, Kings Bitboard
	White, Black                                  Bitboard
	WhiteMove                                     bool
	CastleSide                                    uint8
	EnPassant                                     Bitboard
	HalfMoveClock                                 uint8
}

type coloredPiece struct {
	piece Piece
	color Color
}

type coloredBoard [64]coloredPiece

const InitialPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var noPiece = coloredPiece{Empty, 255}

func createPosition(board coloredBoard) Position {
	position := Position{}
	for i, cp := range board {
		switch cp.piece {
		case Pawn:
			position.Pawns.SetBit(i)
		case Knight:
			position.Knights.SetBit(i)
		case Bishop:
			position.Bishops.SetBit(i)
		case Rook:
			position.Rooks.SetBit(i)
		case Queen:
			position.Queens.SetBit(i)
		case King:
			position.Kings.SetBit(i)
		default:
			if cp.piece != noPiece.piece {
				panic("unhandled no piece")
			}
		}
		switch cp.color {
		case ColorWhite:
			position.White.SetBit(i)
		case ColorBlack:
			position.Black.SetBit(i)
		default:
			if cp.color != noPiece.color {
				panic("unhandled no color")
			}
		}
	}
	return position
}

func (position Position) filterColor() Position {
	p := position
	if p.WhiteMove {
		p.Pawns &= p.White
		p.Knights &= p.White
		p.Bishops &= p.White
		p.Rooks &= p.White
		p.Queens &= p.White
		p.Kings &= p.White
		return p
	}
	p.Pawns &= p.Black
	p.Knights &= p.Black
	p.Bishops &= p.Black
	p.Rooks &= p.Black
	p.Queens &= p.Black
	p.Kings &= p.Black
	return p
}

// AllLegalMoves generates all pseudo-legal moves for a given piece type to empty squares.
//
// This is a simplified move generator that:
//   - Only generates moves to EMPTY squares (no captures)
//   - Stops when encountering any piece (friendly or enemy)
//   - Only updates the moved piece's bitboard (e.g., Bishops, Rooks)
//   - Does NOT update White/Black color bitboards (for performance)
//   - Does NOT update turn, castling rights, en passant, or move clocks
//
// The returned positions are partial and intended for move generation/analysis only.
// Full position updates (captures, color bitboards, game state) should be handled separately.
//
// Parameters:
//   - pieceMoves: Pre-generated move patterns for all piece types
//   - pc: The piece type to generate moves for (Pawn, Knight, Bishop, Rook, Queen, King)
//
// Returns:
//   - A slice of Position structs with only the specified piece's bitboard updated
//   - nil if no pieces of the specified type exist for the side to move
func (position Position) AllLegalMoves(pieceMoves PieceMoves, pc Piece) []Position {
	var positions []Position
	color := position.filterColor()           // take only the color to move
	piecesInColorToMove := color.GetPiece(pc) // get the pieces of that color
	if *piecesInColorToMove == 0 {
		return nil
	}
	// get all the pieces on the board flattened to bitboard
	allFlat := position.Bishops | position.Knights | position.Rooks | position.Queens | position.Kings | position.Pawns
	for _, bitBoard := range piecesInColorToMove.ToSlice() {
		directions := pieceMoves[pc][bitBoard]
		for _, direction := range directions {
			for _, move := range direction {
				if allFlat&move == move { // if there is a piece in the way, stop
					break
				}
				pos := position
				piece := pos.GetPiece(pc) // get the piece reference
				*piece &^= bitBoard       // remove the piece from the board
				*piece |= move            // add the piece to the new position
				//fmt.Println(piece.Hex())
				//fmt.Println(piece.Pretty())
				positions = append(positions, pos)
			}
		}
	}
	return positions
}

// GenerateMoves generates all pseudo-legal moves for the side to move.
//
// Returns a slice of Move structs instead of full Position objects,
// making it more memory-efficient for move generation in search.
//
// Uses Magic Bitboards for O(1) sliding piece move generation.
// Supports: Bishop, Rook, Queen (sliding with magic BB), Knight, King (jumping).
func (position Position) GenerateMoves(pieceMoves PieceMoves) []Move {
	// Cache piece masks (computed once, used by all generators)
	var ourPieces, enemyPieces Bitboard
	if position.WhiteMove {
		ourPieces = position.White
		enemyPieces = position.Black
	} else {
		ourPieces = position.Black
		enemyPieces = position.White
	}
	allPieces := position.Bishops | position.Knights | position.Rooks |
		position.Queens | position.Kings | position.Pawns

	// Single allocation for all moves (max ~64 moves typical)
	moves := make([]Move, 0, 64)

	// Sliding pieces using Magic Bitboards
	moves = position.appendSlidingMoves(moves, Bishop, ourPieces, enemyPieces, allPieces)
	moves = position.appendSlidingMoves(moves, Rook, ourPieces, enemyPieces, allPieces)
	moves = position.appendSlidingMoves(moves, Queen, ourPieces, enemyPieces, allPieces)

	// Jumping pieces
	moves = position.appendJumpingMoves(moves, pieceMoves, Knight, ourPieces, enemyPieces)
	moves = position.appendJumpingMoves(moves, pieceMoves, King, ourPieces, enemyPieces)

	// Pawns
	moves = position.appendPawnMoves(moves, ourPieces, enemyPieces, allPieces)

	// Castling
	moves = position.appendCastlingMoves(moves, allPieces)

	return moves
}

// GenerateLegalMoves generates all legal moves for the current position.
// This filters pseudo-legal moves by checking that the king is not in check after each move.
func (position *Position) GenerateLegalMoves(pieceMoves PieceMoves) []Move {
	pseudoLegal := position.GenerateMoves(pieceMoves)
	legal := make([]Move, 0, len(pseudoLegal))

	for _, m := range pseudoLegal {
		// Make the move
		undo := position.MakeMove(m)

		// After MakeMove, WhiteMove has flipped.
		// We want to check if the side that JUST MOVED left their king in check.
		// That means: is the PREVIOUS side's king attacked by the CURRENT side?
		// IsInCheck checks if the CURRENT side's king is in check, which is wrong.
		// We need to check if the OPPONENT (now current side) can attack our king.

		// Find the king of the side that just moved (opposite of current WhiteMove)
		var kingBB Bitboard
		if position.WhiteMove {
			// It's now white's turn, so black just moved - find black king
			kingBB = position.Kings & position.Black
		} else {
			// It's now black's turn, so white just moved - find white king
			kingBB = position.Kings & position.White
		}

		kingSq := bits.TrailingZeros64(uint64(kingBB))
		// Check if current side (enemy of the one who just moved) attacks the king
		inCheck := position.IsSquareAttacked(kingSq, position.WhiteMove)

		// Unmake the move
		position.UnmakeMove(m, undo)

		if !inCheck {
			legal = append(legal, m)
		}
	}

	return legal
}

// rookAttacks returns attack bitboard for a rook at given square with blockers.
func rookAttacks(sq int, blockers Bitboard) Bitboard {
	m := magic.RookMagics[sq]
	idx := (uint64(blockers&Bitboard(m.Mask)) * m.Number) >> m.Shift
	return Bitboard(magic.RookMoves[sq][idx])
}

// bishopAttacks returns attack bitboard for a bishop at given square with blockers.
func bishopAttacks(sq int, blockers Bitboard) Bitboard {
	m := magic.BishopMagics[sq]
	idx := (uint64(blockers&Bitboard(m.Mask)) * m.Number) >> m.Shift
	return Bitboard(magic.BishopMoves[sq][idx])
}

// Precomputed attack tables for jumping pieces (computed at init)
var knightAttacks [64]Bitboard
var kingAttacks [64]Bitboard

func init() {
	// Precompute knight attacks
	knightOffsets := [][2]int{
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	}
	for sq := 0; sq < 64; sq++ {
		file := sq & 7
		rank := sq >> 3
		for _, off := range knightOffsets {
			newFile := file + off[0]
			newRank := rank + off[1]
			if newFile >= 0 && newFile < 8 && newRank >= 0 && newRank < 8 {
				knightAttacks[sq] |= Bitboard(1 << (newRank*8 + newFile))
			}
		}
	}

	// Precompute king attacks
	kingOffsets := [][2]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}
	for sq := 0; sq < 64; sq++ {
		file := sq & 7
		rank := sq >> 3
		for _, off := range kingOffsets {
			newFile := file + off[0]
			newRank := rank + off[1]
			if newFile >= 0 && newFile < 8 && newRank >= 0 && newRank < 8 {
				kingAttacks[sq] |= Bitboard(1 << (newRank*8 + newFile))
			}
		}
	}
}

// IsSquareAttacked checks if a square is attacked by pieces of the given color.
// Uses magic bitboards for sliding pieces and precomputed tables for jumping pieces.
func (position Position) IsSquareAttacked(sq int, byWhite bool) bool {
	sqBB := Bitboard(1 << sq)
	allPieces := position.Pawns | position.Knights | position.Bishops |
		position.Rooks | position.Queens | position.Kings

	var attackers Bitboard
	if byWhite {
		attackers = position.White
	} else {
		attackers = position.Black
	}

	// Check pawn attacks (reverse direction - where could a pawn attack FROM to hit this square?)
	if byWhite {
		// White pawns attack diagonally up (+7 and +9), so we check diagonally down from target
		whitePawns := position.Pawns & attackers
		// A pawn at sq-7 attacks sq (if sq is not on h-file, to prevent wrap-around)
		// A pawn at sq-9 attacks sq (if sq is not on a-file, to prevent wrap-around)
		var pawnAttackers Bitboard
		if sqBB&fileHMask == 0 { // target not on h-file, check sq-7 attacker
			pawnAttackers |= sqBB >> 7
		}
		if sqBB&fileAMask == 0 { // target not on a-file, check sq-9 attacker
			pawnAttackers |= sqBB >> 9
		}
		if whitePawns&pawnAttackers != 0 {
			return true
		}
	} else {
		// Black pawns attack diagonally down (-7 and -9), so we check diagonally up from target
		blackPawns := position.Pawns & attackers
		var pawnAttackers Bitboard
		if sqBB&fileAMask == 0 { // target not on a-file, check sq+7 attacker
			pawnAttackers |= sqBB << 7
		}
		if sqBB&fileHMask == 0 { // target not on h-file, check sq+9 attacker
			pawnAttackers |= sqBB << 9
		}
		if blackPawns&pawnAttackers != 0 {
			return true
		}
	}

	// Check knight attacks
	knights := position.Knights & attackers
	if knightAttacks[sq]&knights != 0 {
		return true
	}

	// Check king attacks
	king := position.Kings & attackers
	if kingAttacks[sq]&king != 0 {
		return true
	}

	// Check bishop/queen attacks (diagonal)
	bishopsQueens := (position.Bishops | position.Queens) & attackers
	if bishopAttacks(sq, allPieces)&bishopsQueens != 0 {
		return true
	}

	// Check rook/queen attacks (straight)
	rooksQueens := (position.Rooks | position.Queens) & attackers
	if rookAttacks(sq, allPieces)&rooksQueens != 0 {
		return true
	}

	return false
}

// IsInCheck returns true if the current side's king is in check.
func (position Position) IsInCheck() bool {
	// Find our king
	var ourKing Bitboard
	if position.WhiteMove {
		ourKing = position.Kings & position.White
	} else {
		ourKing = position.Kings & position.Black
	}

	if ourKing == 0 {
		return false // No king (shouldn't happen in valid position)
	}

	// Get king square index
	kingSq := bits.TrailingZeros64(uint64(ourKing))

	// Check if enemy pieces attack the king square
	return position.IsSquareAttacked(kingSq, !position.WhiteMove)
}

// appendSlidingMoves appends moves for a sliding piece type using Magic Bitboards.
// Takes pre-computed piece masks to avoid redundant calculations.
func (position Position) appendSlidingMoves(moves []Move, pc Piece, ourPieces, enemyPieces, allPieces Bitboard) []Move {
	// Get pieces of this type for side to move
	pieceBB := *position.GetPiece(pc) & ourPieces
	if pieceBB == 0 {
		return moves
	}

	// Bit-scan through each piece
	for bb := pieceBB; bb != 0; {
		fromIdx := bits.TrailingZeros64(uint64(bb))
		fromBB := Bitboard(1 << fromIdx)
		bb &^= fromBB

		// Get attacks using Magic Bitboards
		var attacks Bitboard
		switch pc {
		case Rook:
			attacks = rookAttacks(fromIdx, allPieces)
		case Bishop:
			attacks = bishopAttacks(fromIdx, allPieces)
		case Queen:
			attacks = rookAttacks(fromIdx, allPieces) | bishopAttacks(fromIdx, allPieces)
		}

		// Remove our own pieces from attacks (can't capture own pieces)
		attacks &^= ourPieces

		// Generate moves for each attack square
		for attacks != 0 {
			toIdx := bits.TrailingZeros64(uint64(attacks))
			toBB := Bitboard(1 << toIdx)
			attacks &^= toBB

			// Detect capture
			captured := Empty
			if enemyPieces&toBB != 0 {
				captured = position.pieceAt(toBB)
			}

			moves = append(moves, Move{
				From:     fromBB,
				To:       toBB,
				Piece:    pc,
				Captured: captured,
			})
		}
	}

	return moves
}

// appendJumpingMoves appends moves for jumping pieces (Knight, King).
// Takes pre-computed piece masks to avoid redundant calculations.
func (position Position) appendJumpingMoves(moves []Move, pieceMoves PieceMoves, pc Piece, ourPieces, enemyPieces Bitboard) []Move {
	// Get pieces of this type for side to move
	pieceBB := *position.GetPiece(pc) & ourPieces
	if pieceBB == 0 {
		return moves
	}

	// Bit-scan through each piece
	for bb := pieceBB; bb != 0; {
		fromIdx := bits.TrailingZeros64(uint64(bb))
		fromBB := Bitboard(1 << fromIdx)
		bb &^= fromBB

		// Jumping pieces have one "direction" with all target squares
		directions := pieceMoves[pc][fromBB]
		if len(directions) == 0 {
			continue
		}
		targets := directions[0] // flat list of all jump targets

		for _, toBB := range targets {
			// Can't land on our own pieces
			if ourPieces&toBB == toBB {
				continue
			}

			// Detect capture
			captured := Empty
			if enemyPieces&toBB != 0 {
				captured = position.pieceAt(toBB)
			}

			moves = append(moves, Move{
				From:     fromBB,
				To:       toBB,
				Piece:    pc,
				Captured: captured,
			})
		}
	}

	return moves
}

// Bitboard masks for pawn move generation
const (
	fileAMask Bitboard = 0x0101010101010101 // a-file (prevents left capture wrap)
	fileHMask Bitboard = 0x8080808080808080 // h-file (prevents right capture wrap)
	rank1Mask Bitboard = 0x00000000000000FF // rank 1 (black promotion)
	rank2Mask Bitboard = 0x000000000000FF00 // white pawn start rank
	rank7Mask Bitboard = 0x00FF000000000000 // black pawn start rank
	rank8Mask Bitboard = 0xFF00000000000000 // rank 8 (white promotion)
)

// Castling constants - square indices and blocking masks
const (
	// White castling
	whiteKingStart       = 4 // e1
	whiteRookKingSide    = 7 // h1
	whiteRookQueenSide   = 0 // a1
	whiteKingKingSideTo  = 6 // g1
	whiteKingQueenSideTo = 2 // c1

	// Black castling
	blackKingStart       = 60 // e8
	blackRookKingSide    = 63 // h8
	blackRookQueenSide   = 56 // a8
	blackKingKingSideTo  = 62 // g8
	blackKingQueenSideTo = 58 // c8
)

// Masks for squares that must be empty for castling
var (
	// White kingside: f1, g1 (indices 5, 6)
	whiteKingSideEmpty Bitboard = (1 << 5) | (1 << 6)
	// White queenside: b1, c1, d1 (indices 1, 2, 3)
	whiteQueenSideEmpty Bitboard = (1 << 1) | (1 << 2) | (1 << 3)
	// Black kingside: f8, g8 (indices 61, 62)
	blackKingSideEmpty Bitboard = (1 << 61) | (1 << 62)
	// Black queenside: b8, c8, d8 (indices 57, 58, 59)
	blackQueenSideEmpty Bitboard = (1 << 57) | (1 << 58) | (1 << 59)
)

// appendCastlingMoves generates castling moves if legal.
// Checks:
// - Castle rights (CastleSide flags)
// - Empty squares between king and rook
// - King not in check, and doesn't pass through attacked squares
func (position Position) appendCastlingMoves(moves []Move, allPieces Bitboard) []Move {
	if position.CastleSide == 0 {
		return moves
	}

	// Determine attacker color (enemy)
	enemyIsWhite := !position.WhiteMove

	if position.WhiteMove {
		// White castling - squares must not be attacked by black

		// White kingside: O-O (e1 -> g1)
		// King passes through e1, f1, g1
		if position.CastleSide&CastleWhiteKingSide != 0 {
			if allPieces&whiteKingSideEmpty == 0 {
				// Check e1, f1, g1 are not attacked
				if !position.IsSquareAttacked(4, enemyIsWhite) && // e1
					!position.IsSquareAttacked(5, enemyIsWhite) && // f1
					!position.IsSquareAttacked(6, enemyIsWhite) { // g1
					moves = append(moves, Move{
						From:  IndexToBitBoard(whiteKingStart),
						To:    IndexToBitBoard(whiteKingKingSideTo),
						Piece: King,
						Flags: FlagCastling,
					})
				}
			}
		}

		// White queenside: O-O-O (e1 -> c1)
		// King passes through e1, d1, c1 (b1 can be attacked, only rook passes)
		if position.CastleSide&CastleWhiteQueenSide != 0 {
			if allPieces&whiteQueenSideEmpty == 0 {
				// Check e1, d1, c1 are not attacked
				if !position.IsSquareAttacked(4, enemyIsWhite) && // e1
					!position.IsSquareAttacked(3, enemyIsWhite) && // d1
					!position.IsSquareAttacked(2, enemyIsWhite) { // c1
					moves = append(moves, Move{
						From:  IndexToBitBoard(whiteKingStart),
						To:    IndexToBitBoard(whiteKingQueenSideTo),
						Piece: King,
						Flags: FlagCastling,
					})
				}
			}
		}
	} else {
		// Black castling - squares must not be attacked by white

		// Black kingside: O-O (e8 -> g8)
		if position.CastleSide&CastleBlackKingSide != 0 {
			if allPieces&blackKingSideEmpty == 0 {
				// Check e8, f8, g8 are not attacked
				if !position.IsSquareAttacked(60, enemyIsWhite) && // e8
					!position.IsSquareAttacked(61, enemyIsWhite) && // f8
					!position.IsSquareAttacked(62, enemyIsWhite) { // g8
					moves = append(moves, Move{
						From:  IndexToBitBoard(blackKingStart),
						To:    IndexToBitBoard(blackKingKingSideTo),
						Piece: King,
						Flags: FlagCastling,
					})
				}
			}
		}

		// Black queenside: O-O-O (e8 -> c8)
		if position.CastleSide&CastleBlackQueenSide != 0 {
			if allPieces&blackQueenSideEmpty == 0 {
				// Check e8, d8, c8 are not attacked
				if !position.IsSquareAttacked(60, enemyIsWhite) && // e8
					!position.IsSquareAttacked(59, enemyIsWhite) && // d8
					!position.IsSquareAttacked(58, enemyIsWhite) { // c8
					moves = append(moves, Move{
						From:  IndexToBitBoard(blackKingStart),
						To:    IndexToBitBoard(blackKingQueenSideTo),
						Piece: King,
						Flags: FlagCastling,
					})
				}
			}
		}
	}

	return moves
}

// promotionPieces are the pieces a pawn can promote to
var promotionPieces = [4]Piece{Queen, Rook, Bishop, Knight}

// appendPawnMove adds a pawn move, generating 4 moves if it's a promotion
func appendPawnMove(moves []Move, fromBB, toBB Bitboard, captured Piece, isPromotion bool) []Move {
	if isPromotion {
		for _, promo := range promotionPieces {
			moves = append(moves, Move{
				From:      fromBB,
				To:        toBB,
				Piece:     Pawn,
				Captured:  captured,
				Promotion: promo,
			})
		}
	} else {
		moves = append(moves, Move{
			From:     fromBB,
			To:       toBB,
			Piece:    Pawn,
			Captured: captured,
		})
	}
	return moves
}

// appendPawnMoves appends all pawn moves (pushes and captures).
// Pawns have unique movement rules:
// - Move forward only (direction depends on color)
// - Can move 2 squares from starting rank
// - Capture diagonally only
// Note: Code is intentionally duplicated for white/black for performance
// (constant shifts are faster than variable shifts).
func (position Position) appendPawnMoves(moves []Move, ourPieces, enemyPieces, allPieces Bitboard) []Move {
	pawns := position.Pawns & ourPieces
	if pawns == 0 {
		return moves
	}

	empty := ^allPieces

	if position.WhiteMove {
		// === WHITE PAWNS (move up: +8) ===

		// Single push: move one square forward to empty square
		singlePush := (pawns << 8) & empty
		for bb := singlePush; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB >> 8
			isPromotion := toBB&rank8Mask != 0

			moves = appendPawnMove(moves, fromBB, toBB, Empty, isPromotion)
		}

		// Double push: from rank 2, both squares must be empty
		// (double push can never be a promotion)
		doublePush := ((pawns & rank2Mask) << 8) & empty // first square empty
		doublePush = (doublePush << 8) & empty           // second square empty
		for bb := doublePush; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB >> 16

			moves = appendPawnMove(moves, fromBB, toBB, Empty, false)
		}

		// Left capture (+7): diagonal left-up, exclude a-file (would wrap to h-file)
		leftCapture := ((pawns &^ fileAMask) << 7) & enemyPieces
		for bb := leftCapture; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB >> 7
			isPromotion := toBB&rank8Mask != 0

			moves = appendPawnMove(moves, fromBB, toBB, position.pieceAt(toBB), isPromotion)
		}

		// Right capture (+9): diagonal right-up, exclude h-file (would wrap to a-file)
		rightCapture := ((pawns &^ fileHMask) << 9) & enemyPieces
		for bb := rightCapture; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB >> 9
			isPromotion := toBB&rank8Mask != 0

			moves = appendPawnMove(moves, fromBB, toBB, position.pieceAt(toBB), isPromotion)
		}

		// En passant for white (capture on rank 6)
		if position.EnPassant != 0 {
			epSquare := position.EnPassant
			// En passant via >>7: pawn at (epSquare-7) captures on epSquare
			// Wrap happens when epSquare is on h-file (epSquare-7 wraps to a-file)
			if epSquare&fileHMask == 0 {
				if fromBB := (epSquare >> 7) & pawns; fromBB != 0 {
					moves = append(moves, Move{
						From:     fromBB,
						To:       epSquare,
						Piece:    Pawn,
						Captured: Pawn,
						Flags:    FlagEnPassant,
					})
				}
			}
			// En passant via >>9: pawn at (epSquare-9) captures on epSquare
			// Wrap happens when epSquare is on a-file (epSquare-9 wraps to h-file)
			if epSquare&fileAMask == 0 {
				if fromBB := (epSquare >> 9) & pawns; fromBB != 0 {
					moves = append(moves, Move{
						From:     fromBB,
						To:       epSquare,
						Piece:    Pawn,
						Captured: Pawn,
						Flags:    FlagEnPassant,
					})
				}
			}
		}
	} else {
		// === BLACK PAWNS (move down: -8) ===

		// Single push: move one square forward to empty square
		singlePush := (pawns >> 8) & empty
		for bb := singlePush; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB << 8
			isPromotion := toBB&rank1Mask != 0

			moves = appendPawnMove(moves, fromBB, toBB, Empty, isPromotion)
		}

		// Double push: from rank 7, both squares must be empty
		// (double push can never be a promotion)
		doublePush := ((pawns & rank7Mask) >> 8) & empty // first square empty
		doublePush = (doublePush >> 8) & empty           // second square empty
		for bb := doublePush; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB << 16

			moves = appendPawnMove(moves, fromBB, toBB, Empty, false)
		}

		// Left capture (-9): diagonal left-down, exclude a-file (would wrap to h-file)
		leftCapture := ((pawns &^ fileAMask) >> 9) & enemyPieces
		for bb := leftCapture; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB << 9
			isPromotion := toBB&rank1Mask != 0

			moves = appendPawnMove(moves, fromBB, toBB, position.pieceAt(toBB), isPromotion)
		}

		// Right capture (-7): diagonal right-down, exclude h-file (would wrap to a-file)
		rightCapture := ((pawns &^ fileHMask) >> 7) & enemyPieces
		for bb := rightCapture; bb != 0; {
			toIdx := bits.TrailingZeros64(uint64(bb))
			toBB := Bitboard(1 << toIdx)
			bb &^= toBB
			fromBB := toBB << 7
			isPromotion := toBB&rank1Mask != 0

			moves = appendPawnMove(moves, fromBB, toBB, position.pieceAt(toBB), isPromotion)
		}

		// En passant for black (capture on rank 3)
		if position.EnPassant != 0 {
			epSquare := position.EnPassant
			// Left en passant: pawn attacks from right (-9 direction)
			// If epSquare is on h-file, no pawn can attack from the left (would wrap from a-file)
			if epSquare&fileHMask == 0 {
				if fromBB := (epSquare << 9) & pawns; fromBB != 0 {
					moves = append(moves, Move{
						From:     fromBB,
						To:       epSquare,
						Piece:    Pawn,
						Captured: Pawn,
						Flags:    FlagEnPassant,
					})
				}
			}
			// Right en passant: pawn attacks from left (-7 direction)
			// If epSquare is on a-file, no pawn can attack from the right (would wrap from h-file)
			if epSquare&fileAMask == 0 {
				if fromBB := (epSquare << 7) & pawns; fromBB != 0 {
					moves = append(moves, Move{
						From:     fromBB,
						To:       epSquare,
						Piece:    Pawn,
						Captured: Pawn,
						Flags:    FlagEnPassant,
					})
				}
			}
		}
	}

	return moves
}

func (position *Position) GetPiece(piece Piece) *Bitboard {
	switch piece {
	case Pawn:
		return &position.Pawns
	case Knight:
		return &position.Knights
	case Bishop:
		return &position.Bishops
	case Rook:
		return &position.Rooks
	case Queen:
		return &position.Queens
	case King:
		return &position.Kings
	default:
		log.Fatal("unhandled piece")
		return nil
	}
}

// pieceAt returns the piece type at the given square bitboard.
// Returns Empty if no piece is found.
func (position Position) pieceAt(sq Bitboard) Piece {
	if position.Pawns&sq != 0 {
		return Pawn
	}
	if position.Knights&sq != 0 {
		return Knight
	}
	if position.Bishops&sq != 0 {
		return Bishop
	}
	if position.Rooks&sq != 0 {
		return Rook
	}
	if position.Queens&sq != 0 {
		return Queen
	}
	if position.Kings&sq != 0 {
		return King
	}
	return Empty
}

// Pretty returns a compact, human-readable representation of the chess position
// using Unicode chess piece symbols.
func (position *Position) Pretty() string {
	var s string

	// Unicode chess pieces: White (uppercase) and Black (lowercase)
	pieceSymbols := map[Piece]map[Color]string{
		Pawn:   {ColorWhite: "♙", ColorBlack: "♟"},
		Knight: {ColorWhite: "♘", ColorBlack: "♞"},
		Bishop: {ColorWhite: "♗", ColorBlack: "♝"},
		Rook:   {ColorWhite: "♖", ColorBlack: "♜"},
		Queen:  {ColorWhite: "♕", ColorBlack: "♛"},
		King:   {ColorWhite: "♔", ColorBlack: "♚"},
	}

	for r := Rank8; r >= Rank1; r-- {
		s += fmt.Sprintf("%d  ", r+1)
		for f := FileA; f <= FileH; f++ {
			idx := squareIndex(f, r)
			symbol := "·"

			// Check each piece type and color
			for pieceType := Pawn; pieceType <= King; pieceType++ {
				pieceBB := position.GetPiece(pieceType)
				if pieceBB.IsBitSet(idx) {
					if position.White.IsBitSet(idx) {
						symbol = pieceSymbols[pieceType][ColorWhite]
					} else if position.Black.IsBitSet(idx) {
						symbol = pieceSymbols[pieceType][ColorBlack]
					}
					break
				}
			}

			s += symbol + " "
		}
		s += "\n"
	}
	s += "   a b c d e f g h\n"

	return s
}
