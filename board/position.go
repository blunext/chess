package board

import (
	"fmt"
	"log"
	"math/bits"
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

// GenerateSlidingMoves generates all pseudo-legal moves for sliding pieces
// (Bishop, Rook, Queen) for the side to move.
//
// Returns a slice of Move structs instead of full Position objects,
// making it more memory-efficient for move generation in search.
//
// Currently generates moves to empty squares only (no captures).
func (position Position) GenerateSlidingMoves(pieceMoves PieceMoves) []Move {
	var moves []Move

	// Generate moves for all sliding piece types
	slidingPieces := []Piece{Bishop, Rook, Queen}

	for _, pc := range slidingPieces {
		pieceMoves := position.generateMovesForPiece(pieceMoves, pc)
		moves = append(moves, pieceMoves...)
	}

	return moves
}

// generateMovesForPiece generates pseudo-legal moves for a single piece type.
// Helper function used by GenerateSlidingMoves.
func (position Position) generateMovesForPiece(pieceMoves PieceMoves, pc Piece) []Move {
	var moves []Move

	// We only need pieces of one type for the side to move.
	var ourPieces Bitboard
	if position.WhiteMove {
		ourPieces = position.White
	} else {
		ourPieces = position.Black
	}

	// Mask: only pieces of type 'pc' that belong to side to move
	pieceBB := *position.GetPiece(pc) & ourPieces
	if pieceBB == 0 {
		return nil
	}

	// Get all pieces on the board
	allFlat := position.Bishops | position.Knights | position.Rooks |
		position.Queens | position.Kings | position.Pawns

	// Bit-scanning: iterate over set bits without allocating a slice.
	// Uses CPU instruction TZCNT/BSF for O(1) bit position lookup.
	// Example: bb = 0b00100100 (pieces on c1 and f1)
	//   Iter 1: TrailingZeros64 → 2 (c1), then clear bit → 0b00100000
	//   Iter 2: TrailingZeros64 → 5 (f1), then clear bit → 0b00000000
	//   Loop ends when bb == 0
	for bb := pieceBB; bb != 0; {
		// Find index of the lowest set bit (0-63)
		fromIdx := bits.TrailingZeros64(uint64(bb))
		// Convert index back to single-bit bitboard for map lookup
		fromBB := Bitboard(1 << fromIdx)
		// Clear this bit so next iteration finds the next piece
		// &^ is AND-NOT operator: bb = bb AND (NOT fromBB)
		bb &^= fromBB

		directions := pieceMoves[pc][fromBB]
		for _, direction := range directions {
			for _, toBB := range direction {
				if allFlat&toBB == toBB { // piece in the way, stop
					break
				}
				moves = append(moves, Move{
					From:     fromBB,
					To:       toBB,
					Piece:    pc,
					Captured: Empty,
				})
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
