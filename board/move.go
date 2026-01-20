package board

import "fmt"

// Move represents a chess move in a compact form.
// Instead of storing entire board positions, we store only the essential
// information about what piece moved from where to where.
type Move struct {
	From      Bitboard // source square (single bit set)
	To        Bitboard // destination square (single bit set)
	Piece     Piece    // type of piece that moved
	Captured  Piece    // captured piece type (Empty if no capture)
	Promotion Piece    // promotion piece type (Empty if not a promotion)
	Flags     MoveFlag // special move flags (en passant, castling)
}

// MoveFlag represents special move types
type MoveFlag uint8

const (
	FlagNone      MoveFlag = 0
	FlagEnPassant MoveFlag = 1 << iota
	FlagCastling
)

// String returns a human-readable representation of the move.
// Format: "Piece: from -> to" with optional capture/promotion info.
func (m Move) String() string {
	pieceNames := map[Piece]string{
		Pawn:   "Pawn",
		Knight: "Knight",
		Bishop: "Bishop",
		Rook:   "Rook",
		Queen:  "Queen",
		King:   "King",
	}

	fromIdx := bitboardToIndex(m.From)
	toIdx := bitboardToIndex(m.To)

	fromSquare := IndexToAlgebraic(fromIdx)
	toSquare := IndexToAlgebraic(toIdx)

	var result string
	if m.Captured != Empty {
		result = fmt.Sprintf("%s: %s x %s", pieceNames[m.Piece], fromSquare, toSquare)
	} else {
		result = fmt.Sprintf("%s: %s -> %s", pieceNames[m.Piece], fromSquare, toSquare)
	}

	if m.Promotion != Empty {
		result += fmt.Sprintf("=%s", pieceNames[m.Promotion])
	}
	if m.Flags&FlagEnPassant != 0 {
		result += " e.p."
	}
	if m.Flags&FlagCastling != 0 {
		result += " (castling)"
	}

	return result
}

// ToUCI returns the move in UCI notation (e.g., "e2e4", "e7e8q").
// UCI format: <from><to>[promotion]
// Promotion piece: q=queen, r=rook, b=bishop, n=knight (lowercase)
func (m Move) ToUCI() string {
	fromIdx := bitboardToIndex(m.From)
	toIdx := bitboardToIndex(m.To)

	fromSquare := IndexToAlgebraic(fromIdx)
	toSquare := IndexToAlgebraic(toIdx)

	uci := fromSquare + toSquare

	// Add promotion piece (lowercase)
	if m.Promotion != Empty {
		promoChars := map[Piece]string{
			Queen:  "q",
			Rook:   "r",
			Bishop: "b",
			Knight: "n",
		}
		uci += promoChars[m.Promotion]
	}

	return uci
}

// bitboardToIndex converts a bitboard with a single bit set to its index (0-63).
func bitboardToIndex(bb Bitboard) int {
	for i := range 64 {
		if bb&(1<<i) != 0 {
			return i
		}
	}
	return -1
}

// IndexToAlgebraic converts a square index to algebraic notation (e.g., 0 -> "a1").
func IndexToAlgebraic(idx int) string {
	if idx < 0 || idx > 63 {
		return "??"
	}
	file := idx & 7  // 0-7
	rank := idx >> 3 // 0-7
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}
