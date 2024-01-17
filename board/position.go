package board

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

func (position Position) AllBishops(sliders Sliders) []Position {
	var positions []Position
	bishops := position.filterColor().Bishops
	if bishops == 0 {
		return nil
	}
	allFlat := position.ToFlat()
	for _, bitBoard := range bishops.ToSlice() {
		directions := sliders[Bishop][bitBoard]
		for _, direction := range directions {
			for _, move := range direction {
				if allFlat&move == move {
					break
				}
				pos := position
				pos.Bishops &^= bitBoard
				pos.Bishops |= move
				positions = append(positions, pos)
			}
		}
	}
	return positions
}

//
//func (position *Position) getPiece(piece Piece) *Bitboard {
//	switch piece {
//	case Pawn:
//		return &position.Pawns
//	case Knight:
//		return &position.Knights
//	case Bishop:
//		return &position.Bishops
//	case Rook:
//		return &position.Rooks
//	case Queen:
//		return &position.Queens
//	case King:
//		return &position.Kings
//	default:
//		log.Fatal("unhandled piece")
//		return nil
//	}
//}
