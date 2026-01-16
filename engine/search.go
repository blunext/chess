package engine

import (
	"math/rand"
	"sort"
	"time"

	"chess/board"
	"chess/book"
)

const (
	infinity  = 1000000
	mateScore = 100000
)

// Piece values for move ordering (MVV-LVA)
var pieceValues = map[board.Piece]int{
	board.Pawn:   100,
	board.Knight: 320,
	board.Bishop: 330,
	board.Rook:   500,
	board.Queen:  900,
	board.King:   20000,
}

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
	sort.Slice(moves, func(i, j int) bool {
		return moveScore(moves[i]) > moveScore(moves[j])
	})
}

// OpeningBook is the global opening book (nil if not loaded)
var OpeningBook *book.Book

// bookRng is used for random book move selection
var bookRng = rand.New(rand.NewSource(time.Now().UnixNano()))

// SetOpeningBook sets the opening book to use.
func SetOpeningBook(b *book.Book) {
	OpeningBook = b
}

// SearchResult contains the best move and its evaluation.
type SearchResult struct {
	Move     board.Move
	Score    int
	Nodes    int64 // nodes searched (for debugging)
	FromBook bool  // true if move came from opening book
}

// SearchWithBook probes the opening book first, then falls back to search.
func SearchWithBook(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	// Try opening book first
	if OpeningBook != nil {
		polyHash := book.PolyglotHash(pos)
		if bookMove, ok := OpeningBook.ProbeRandom(polyHash, bookRng); ok {
			// Find matching legal move to get full move info (piece type, etc.)
			legalMoves := pos.GenerateLegalMoves(pieceMoves)
			for _, m := range legalMoves {
				if m.From == bookMove.From && m.To == bookMove.To && m.Promotion == bookMove.Promotion {
					return SearchResult{Move: m, FromBook: true}
				}
			}
		}
	}

	// Fall back to normal search
	return Search(pos, pieceMoves, depth)
}

// Search finds the best move using alpha-beta pruning.
func Search(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	moves := pos.GenerateLegalMoves(pieceMoves)
	sortMoves(moves)

	if len(moves) == 0 {
		if pos.IsInCheck() {
			if pos.WhiteMove {
				return SearchResult{Score: -mateScore}
			}
			return SearchResult{Score: mateScore}
		}
		return SearchResult{Score: 0}
	}

	var bestMove board.Move
	var bestScore int
	var nodes int64

	alpha := -infinity
	beta := infinity

	if pos.WhiteMove {
		bestScore = -infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score, n := alphaBeta(&pos, pieceMoves, depth-1, alpha, beta)
			nodes += n
			pos.UnmakeMove(move, undo)

			if score > bestScore {
				bestScore = score
				bestMove = move
			}
			if score > alpha {
				alpha = score
			}
		}
	} else {
		bestScore = infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score, n := alphaBeta(&pos, pieceMoves, depth-1, alpha, beta)
			nodes += n
			pos.UnmakeMove(move, undo)

			if score < bestScore {
				bestScore = score
				bestMove = move
			}
			if score < beta {
				beta = score
			}
		}
	}

	return SearchResult{Move: bestMove, Score: bestScore, Nodes: nodes}
}

// alphaBeta returns the evaluation score using alpha-beta pruning.
// Alpha = best score the maximizer (white) can guarantee
// Beta = best score the minimizer (black) can guarantee
// Returns (score, nodesSearched)
func alphaBeta(pos *board.Position, pieceMoves board.PieceMoves, depth int, alpha, beta int) (int, int64) {
	if depth == 0 {
		return Evaluate(*pos), 1
	}

	moves := pos.GenerateLegalMoves(pieceMoves)
	sortMoves(moves)

	if len(moves) == 0 {
		if pos.IsInCheck() {
			if pos.WhiteMove {
				return -mateScore + (10 - depth), 1 // Prefer faster mates
			}
			return mateScore - (10 - depth), 1
		}
		return 0, 1 // Stalemate
	}

	var nodes int64

	if pos.WhiteMove {
		// Maximizing player
		bestScore := -infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score, n := alphaBeta(pos, pieceMoves, depth-1, alpha, beta)
			nodes += n
			pos.UnmakeMove(move, undo)

			if score > bestScore {
				bestScore = score
			}
			if score > alpha {
				alpha = score
			}
			// Beta cutoff - black already has a better option
			if alpha >= beta {
				break
			}
		}
		return bestScore, nodes
	} else {
		// Minimizing player
		bestScore := infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score, n := alphaBeta(pos, pieceMoves, depth-1, alpha, beta)
			nodes += n
			pos.UnmakeMove(move, undo)

			if score < bestScore {
				bestScore = score
			}
			if score < beta {
				beta = score
			}
			// Alpha cutoff - white already has a better option
			if alpha >= beta {
				break
			}
		}
		return bestScore, nodes
	}
}
