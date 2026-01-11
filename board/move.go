package board

import "fmt"

// Move represents a chess move in a compact form.
// Instead of storing entire board positions, we store only the essential
// information about what piece moved from where to where.
type Move struct {
	From     Bitboard // source square (single bit set)
	To       Bitboard // destination square (single bit set)
	Piece    Piece    // type of piece that moved
	Captured Piece    // captured piece type (Empty if no capture)
}

// String returns a human-readable representation of the move.
// Format: "Piece: from -> to" with optional capture info.
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

	fromSquare := indexToAlgebraic(fromIdx)
	toSquare := indexToAlgebraic(toIdx)

	if m.Captured != Empty {
		return fmt.Sprintf("%s: %s x %s", pieceNames[m.Piece], fromSquare, toSquare)
	}
	return fmt.Sprintf("%s: %s -> %s", pieceNames[m.Piece], fromSquare, toSquare)
}

// bitboardToIndex converts a bitboard with a single bit set to its index (0-63).
func bitboardToIndex(bb Bitboard) int {
	for i := 0; i < 64; i++ {
		if bb&(1<<i) != 0 {
			return i
		}
	}
	return -1
}

// indexToAlgebraic converts a square index to algebraic notation (e.g., 0 -> "a1").
func indexToAlgebraic(idx int) string {
	if idx < 0 || idx > 63 {
		return "??"
	}
	file := idx & 7  // 0-7
	rank := idx >> 3 // 0-7
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}
