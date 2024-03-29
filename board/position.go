package board

import (
	"log"
)

type SquareMoves map[Bitboard][]Bitboard
type Generics map[Piece]SquareMoves
type SliderSquareMoves map[Bitboard][][]Bitboard
type Sliders map[Piece]SliderSquareMoves

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

func (position Position) ToFlat() Bitboard {
	return toFlat(position.Bishops, position.Knights, position.Rooks, position.Queens, position.Kings, position.Pawns)
}

func (position Position) AllSliders(sliders Sliders, pc Piece) []Position {
	var positions []Position
	color := position.filterColor()           // take only the color to move
	piecesInColorToMove := color.GetPiece(pc) // get the pieces of that color
	if *piecesInColorToMove == 0 {
		return nil
	}
	allFlat := position.ToFlat() // get all the pieces on the board flattened to bitboard
	for _, bitBoard := range piecesInColorToMove.ToSlice() {
		directions := sliders[pc][bitBoard]
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

func (position Position) AllGenerics(generics Generics, pc Piece) []Position {
	var positions []Position
	color := position.filterColor()           // take only the color to move
	piecesInColorToMove := color.GetPiece(pc) // get the pieces of that color
	if *piecesInColorToMove == 0 {
		return nil
	}
	allFlat := position.ToFlat() // get all the pieces on the board flattened to bitboard
	for _, bitBoard := range piecesInColorToMove.ToSlice() {
		moves := generics[pc][bitBoard]
		for _, move := range moves {
			if allFlat&move == move { // if there is a piece in the way, stop
				break
			}
			pos := position
			piece := pos.GetPiece(pc) // get the piece reference
			*piece &^= bitBoard       // remove the piece from the board
			*piece |= move            // add the piece to the new position
			positions = append(positions, pos)
		}
	}
	return positions
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
