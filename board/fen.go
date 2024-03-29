package board

import (
	"log"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var fileNumber = map[string]int{
	"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8,
}

var rune2Piece = map[rune]coloredPiece{
	'P': {Pawn, ColorWhite},
	'N': {Knight, ColorWhite},
	'B': {Bishop, ColorWhite},
	'R': {Rook, ColorWhite},
	'Q': {Queen, ColorWhite},
	'K': {King, ColorWhite},
	'p': {Pawn, ColorBlack},
	'n': {Knight, ColorBlack},
	'b': {Bishop, ColorBlack},
	'r': {Rook, ColorBlack},
	'q': {Queen, ColorBlack},
	'k': {King, ColorBlack},
}

func CreatePositionFormFEN(fen string) Position {
	fields := strings.Split(fen, " ")
	if len(fields) != 6 {
		log.Fatal("bad fen")
	}
	coloredBoard := createColoredBoard(fields[0])
	position := createPosition(coloredBoard)

	position.WhiteMove = fields[1] == "w"
	position.CastleSide = castleAbility(fields[2])
	position.EnPassant = enPassant(fields[3])

	halfMoveClock, _ := strconv.Atoi(fields[4])
	position.HalfMoveClock = uint8(halfMoveClock)

	// todo: do we need Fullmove counter?

	return position
}

func enPassant(s string) Bitboard {
	var ep Bitboard
	if s == "-" {
		return ep
	}
	file := fileNumber[s[:1]]
	rank, _ := strconv.Atoi(s[1:])
	ep.SetBit(squareIndex(file-1, rank-1))
	return ep
}

func castleAbility(c string) uint8 {
	var castle int
	for _, ch := range c {
		switch ch {
		case 'K':
			castle |= CastleWhiteKingSide
		case 'Q':
			castle |= CastleWhiteQueenSide
		case 'k':
			castle |= CastleBlackKingSide
		case 'q':
			castle |= CastleBlackQueenSide
		}
	}
	return uint8(castle)
}

func createColoredBoard(piecePlacement string) coloredBoard {
	ranks := strings.Split(piecePlacement, "/")
	slices.Reverse(ranks)
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
				b[bIndex] = rune2Piece[char]
				bIndex++
			}
		}
	}
	return b
}
