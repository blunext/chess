package bitboard

type coloredPiece struct {
	piece uint8
	color uint8
}

type coloredBoard [64]coloredPiece

const InitialPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var noPiece = coloredPiece{Empty, 255}

type Position struct {
	Pawns, Knights, Bishops, Rooks, Queens, Kings bitboard
	White, Black                                  bitboard
	WhiteMove                                     bool
}

func createPosition(board coloredBoard) Position {
	position := Position{}
	for i, cp := range board {
		switch cp.piece {
		case Pawn:
			position.Pawns.setBit(i)
		case Knight:
			position.Knights.setBit(i)
		case Bishop:
			position.Bishops.setBit(i)
		case Rook:
			position.Rooks.setBit(i)
		case Queen:
			position.Queens.setBit(i)
		case King:
			position.Kings.setBit(i)
		}
		switch cp.color {
		case ColorWhite:
			position.White.setBit(i)
		case ColorBlack:
			position.Black.setBit(i)
		}
	}
	return position
}
