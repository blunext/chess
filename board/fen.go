package board

import (
	"log"
	"math/bits"
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

	// Compute Zobrist hash
	position.Hash = position.ComputeHash()

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
				for range n {
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

// ToFEN returns the FEN string for the current position
func (position Position) ToFEN() string {
	var sb strings.Builder

	// 1. Piece placement
	for rank := 7; rank >= 0; rank-- {
		empty := 0
		for file := range 8 {
			sq := rank*8 + file
			bb := Bitboard(1 << sq)

			// Find piece at this square
			piece := Empty
			isWhite := false

			if position.White&bb != 0 {
				isWhite = true
			} else if position.Black&bb != 0 {
				isWhite = false
			} else {
				empty++
				continue
			}

			// If we found a piece, append any accumulated empty squares first
			if empty > 0 {
				sb.WriteString(strconv.Itoa(empty))
				empty = 0
			}

			if position.Pawns&bb != 0 {
				piece = Pawn
			} else if position.Knights&bb != 0 {
				piece = Knight
			} else if position.Bishops&bb != 0 {
				piece = Bishop
			} else if position.Rooks&bb != 0 {
				piece = Rook
			} else if position.Queens&bb != 0 {
				piece = Queen
			} else if position.Kings&bb != 0 {
				piece = King
			}

			// Append piece char
			char := ""
			switch piece {
			case Pawn:
				char = "p"
			case Knight:
				char = "n"
			case Bishop:
				char = "b"
			case Rook:
				char = "r"
			case Queen:
				char = "q"
			case King:
				char = "k"
			}

			if isWhite {
				char = strings.ToUpper(char)
			}
			sb.WriteString(char)
		}
		if empty > 0 {
			sb.WriteString(strconv.Itoa(empty))
		}
		if rank > 0 {
			sb.WriteString("/")
		}
	}

	sb.WriteString(" ")

	// 2. Active color
	if position.WhiteMove {
		sb.WriteString("w")
	} else {
		sb.WriteString("b")
	}

	sb.WriteString(" ")

	// 3. Castling availability
	castling := ""
	if position.CastleSide&CastleWhiteKingSide != 0 {
		castling += "K"
	}
	if position.CastleSide&CastleWhiteQueenSide != 0 {
		castling += "Q"
	}
	if position.CastleSide&CastleBlackKingSide != 0 {
		castling += "k"
	}
	if position.CastleSide&CastleBlackQueenSide != 0 {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	sb.WriteString(castling)

	sb.WriteString(" ")

	// 4. En passant target square
	if position.EnPassant != 0 {
		idx := bits.TrailingZeros64(uint64(position.EnPassant))
		file := idx % 8
		rank := idx / 8
		fileChar := string(rune('a' + file))
		rankChar := strconv.Itoa(rank + 1)
		sb.WriteString(fileChar + rankChar)
	} else {
		sb.WriteString("-")
	}

	sb.WriteString(" ")

	// 5. Halfmove clock
	sb.WriteString(strconv.Itoa(int(position.HalfMoveClock)))

	sb.WriteString(" ")

	// 6. Fullmove number (always 1 for now as we don't track it)
	sb.WriteString("1")

	return sb.String()
}
