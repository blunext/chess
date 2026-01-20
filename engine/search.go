package engine

import (
	"slices"

	"chess/board"
	"chess/book"
)

const (
	infinity           = 1000000
	mateScore          = 100000
	nullMoveReduction  = 2   // Depth reduction for null move pruning
	maxSearchDepth     = 64  // Maximum search depth for killer moves array
	deltaPruningMargin = 200 // Safety margin for delta pruning in quiescence
)

// UseNullMovePruning controls whether null move pruning is enabled.
// Currently disabled due to incorrect pruning at higher depths.
// TODO: Implement verification search before re-enabling.
var UseNullMovePruning = false

// isEndgame returns true if position has only kings and pawns (zugzwang risk)
// In such positions, null move pruning can be dangerous
func isEndgame(pos *board.Position) bool {
	return pos.Queens == 0 && pos.Rooks == 0 && pos.Bishops == 0 && pos.Knights == 0
}

// Piece values for move ordering (MVV-LVA)
// Array indexed by board.Piece (Empty=0, Pawn=1, Knight=2, Bishop=3, Rook=4, Queen=5, King=6)
var pieceValues = [7]int{0, 100, 320, 330, 500, 900, 20000}

// moveScore returns a score for move ordering (higher = search first)
// Uses MVV-LVA: Most Valuable Victim - Least Valuable Attacker
func moveScore(m board.Move) int {
	// Captures: score = victim value - attacker value / 10
	if m.Captured != board.Empty {
		victim := pieceValues[m.Captured]
		attacker := pieceValues[m.Piece] / 10
		return 10000 + victim - attacker // +10000 to ensure captures come first
	}

	// Promotions are also high priority
	if m.Promotion != board.Empty {
		return 9000 + pieceValues[m.Promotion]
	}

	// Non-captures: no bonus
	return 0
}

// sortMoves orders moves for better alpha-beta pruning
func sortMoves(moves []board.Move) {
	slices.SortFunc(moves, func(i, j board.Move) int {
		return moveScore(j) - moveScore(i) // descending
	})
}

// OpeningBook is the global opening book (nil if not loaded)
var OpeningBook *book.Book

// SetOpeningBook sets the opening book to use.
func SetOpeningBook(b *book.Book) {
	OpeningBook = b
}

// defaultSession is used by global Search functions for backward compatibility.
var defaultSession *Session

// getDefaultSession returns the default session, creating it if needed.
func getDefaultSession() *Session {
	if defaultSession == nil {
		defaultSession = NewSession(DefaultHashMB)
	}
	return defaultSession
}

// SearchResult contains the best move and its evaluation.
type SearchResult struct {
	Move     board.Move
	Score    int
	Nodes    int64 // nodes searched (for debugging)
	FromBook bool  // true if move came from opening book
}

// SearchWithBook probes the opening book first, then falls back to search.
// Uses the default session for backward compatibility.
func SearchWithBook(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	return getDefaultSession().SearchWithBook(pos, pieceMoves, depth)
}

// Search finds the best move using alpha-beta pruning with fixed depth.
// Uses the default session for backward compatibility.
func Search(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	return getDefaultSession().Search(pos, pieceMoves, depth)
}

// filterCaptures returns only capture moves from a list of moves.
func filterCaptures(moves []board.Move) []board.Move {
	captures := make([]board.Move, 0, len(moves)/4)
	for _, m := range moves {
		if m.Captured != board.Empty {
			captures = append(captures, m)
		}
	}
	return captures
}
