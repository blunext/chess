package bitboard

import (
	"log"
	"strconv"
	"strings"
	"unicode"
)

type coloredPiece struct {
	piece uint8
	color uint8
}

var noPiece = coloredPiece{Empty, 255}

func runeToFigure(r rune) coloredPiece {
	switch r {
	case 'P':
		return coloredPiece{Pawn, ColorWhite}
	case 'N':
		return coloredPiece{Knight, ColorWhite}
	case 'B':
		return coloredPiece{Bishop, ColorWhite}
	case 'R':
		return coloredPiece{Rook, ColorWhite}
	case 'Q':
		return coloredPiece{Queen, ColorWhite}
	case 'K':
		return coloredPiece{King, ColorWhite}
	case 'p':
		return coloredPiece{Pawn, ColorBlack}
	case 'n':
		return coloredPiece{Knight, ColorBlack}
	case 'b':
		return coloredPiece{Bishop, ColorBlack}
	case 'r':
		return coloredPiece{Rook, ColorBlack}
	case 'q':
		return coloredPiece{Queen, ColorBlack}
	case 'k':
		return coloredPiece{King, ColorBlack}
	default:
		log.Fatalf("cannot convert rune to coloredPiece: %v", r)
		return noPiece
	}
}

const InitialPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type coloredBoard [64]coloredPiece

func fromFen(fen string) coloredBoard {
	fields := strings.Split(fen, " ")
	if len(fields) != 6 {
		log.Fatal("bad fen")
	}
	ranks := strings.Split(fields[0], "/")
	if len(ranks) != 8 {
		log.Fatal("bad ranks no")
	}
	b := coloredBoard{}
	bIndex := 0
	for _, rank := range ranks {
		// fmt.Println(rank)
		for _, char := range rank {
			switch {
			case unicode.IsDigit(char):
				var n, _ = strconv.Atoi(string(char))
				for i := 0; i < n; i++ {
					b[bIndex] = noPiece
					bIndex++
				}
			case unicode.IsLetter(char):
				b[bIndex] = runeToFigure(char)
				bIndex++
			}
		}
	}
	return b
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
