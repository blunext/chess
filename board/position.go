package board

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

const (
	Empty uint8 = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	ColorWhite uint8 = 0
	ColorBlack uint8 = 1
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
	piece uint8
	color uint8
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
		}
		switch cp.color {
		case ColorWhite:
			position.White.SetBit(i)
		case ColorBlack:
			position.Black.SetBit(i)
		}
	}
	return position
}
