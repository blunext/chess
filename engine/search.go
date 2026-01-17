package engine

import (
	"math/rand"
	"sort"
	"time"

	"chess/board"
	"chess/book"
)

const (
	infinity          = 1000000
	mateScore         = 100000
	nullMoveReduction = 2 // Depth reduction for null move pruning
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

// Search finds the best move using alpha-beta pruning with fixed depth.
// Internally uses the timed search with a very long timeout.
func Search(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	ctx := NewSearchContext(24 * time.Hour) // Effectively no time limit
	result := searchRootDepth(pos, pieceMoves, depth, ctx)
	return SearchResult{
		Move:  result.Move,
		Score: result.Score,
		Nodes: ctx.nodes,
	}
}

// searchRootDepth searches to a fixed depth (used by Search and iterative deepening).
func searchRootDepth(pos board.Position, pieceMoves board.PieceMoves, depth int, ctx *SearchContext) SearchResult {
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

	alpha := -infinity
	beta := infinity

	if pos.WhiteMove {
		bestScore = -infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := alphaBeta(&pos, pieceMoves, depth-1, alpha, beta, true, ctx)
			pos.UnmakeMove(move, undo)

			if ctx.stopped.Load() {
				break
			}

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
			score := alphaBeta(&pos, pieceMoves, depth-1, alpha, beta, true, ctx)
			pos.UnmakeMove(move, undo)

			if ctx.stopped.Load() {
				break
			}

			if score < bestScore {
				bestScore = score
				bestMove = move
			}
			if score < beta {
				beta = score
			}
		}
	}

	return SearchResult{Move: bestMove, Score: bestScore, Nodes: ctx.nodes}
}

// alphaBeta returns the evaluation score using alpha-beta pruning.
// Alpha = best score the maximizer (white) can guarantee
// Beta = best score the minimizer (black) can guarantee
// nullMoveAllowed prevents consecutive null moves
func alphaBeta(pos *board.Position, pieceMoves board.PieceMoves, depth int, alpha, beta int, nullMoveAllowed bool, ctx *SearchContext) int {
	// Increment node counter and check timeout every 2048 nodes
	ctx.nodes++
	if ctx.nodes&2047 == 0 && ctx.checkTimeout() {
		return 0
	}

	if depth == 0 {
		return quiescence(pos, pieceMoves, alpha, beta, ctx)
	}

	alphaOrig := alpha
	hash := pos.Hash

	// Probe transposition table
	var ttMove board.Move
	if TT != nil {
		if entry, found := TT.Probe(hash); found {
			ttMove = entry.BestMove
			if int(entry.Depth) >= depth {
				score := int(entry.Score)
				switch entry.Flag {
				case TTFlagExact:
					return score
				case TTFlagLower:
					if score > alpha {
						alpha = score
					}
				case TTFlagUpper:
					if score < beta {
						beta = score
					}
				}
				if alpha >= beta {
					return score
				}
			}
		}
	}

	// Null Move Pruning
	// Idea: if we can skip our move and still beat beta, the position is so good
	// that we can prune. Skip if: in check, depth too shallow, endgame (zugzwang), or just did null move
	inCheck := pos.IsInCheck()
	if UseNullMovePruning && nullMoveAllowed && depth >= 4 && !inCheck && !isEndgame(pos) {
		// Save state for null move
		oldHash := pos.Hash
		oldEnPassant := pos.EnPassant

		// Make null move (flip side to move, clear en passant, update hash)
		pos.WhiteMove = !pos.WhiteMove
		pos.EnPassant = 0
		pos.Hash ^= board.HashSide()

		// Search with reduced depth and zero window
		nullScore := -alphaBeta(pos, pieceMoves, depth-1-nullMoveReduction, -beta, -beta+1, false, ctx)

		// Unmake null move
		pos.WhiteMove = !pos.WhiteMove
		pos.EnPassant = oldEnPassant
		pos.Hash = oldHash

		if nullScore >= beta {
			return beta // Null move cutoff
		}
	}

	moves := pos.GenerateLegalMoves(pieceMoves)

	// Put TT move first if available
	if ttMove != (board.Move{}) {
		for i, m := range moves {
			if m.From == ttMove.From && m.To == ttMove.To && m.Promotion == ttMove.Promotion {
				// Swap to front
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
		// Sort remaining moves (skip first which is TT move)
		if len(moves) > 1 {
			sortMoves(moves[1:])
		}
	} else {
		sortMoves(moves)
	}

	if len(moves) == 0 {
		if pos.IsInCheck() {
			if pos.WhiteMove {
				return -mateScore + (10 - depth)
			}
			return mateScore - (10 - depth)
		}
		return 0 // Stalemate
	}

	var bestMove board.Move
	var bestScore int

	if pos.WhiteMove {
		bestScore = -infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := alphaBeta(pos, pieceMoves, depth-1, alpha, beta, true, ctx)
			pos.UnmakeMove(move, undo)

			if ctx.stopped.Load() {
				return 0
			}

			if score > bestScore {
				bestScore = score
				bestMove = move
			}
			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				break
			}
		}
	} else {
		bestScore = infinity
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := alphaBeta(pos, pieceMoves, depth-1, alpha, beta, true, ctx)
			pos.UnmakeMove(move, undo)

			if ctx.stopped.Load() {
				return 0
			}

			if score < bestScore {
				bestScore = score
				bestMove = move
			}
			if score < beta {
				beta = score
			}
			if alpha >= beta {
				break
			}
		}
	}

	// Store in transposition table
	if TT != nil && !ctx.stopped.Load() {
		var flag TTFlag
		if bestScore <= alphaOrig {
			flag = TTFlagUpper
		} else if bestScore >= beta {
			flag = TTFlagLower
		} else {
			flag = TTFlagExact
		}
		TT.Store(hash, int16(bestScore), int8(depth), flag, bestMove)
	}

	return bestScore
}

// quiescence continues search only for captures to avoid horizon effect.
func quiescence(pos *board.Position, pieceMoves board.PieceMoves, alpha, beta int, ctx *SearchContext) int {
	ctx.nodes++
	if ctx.nodes&2047 == 0 && ctx.checkTimeout() {
		return 0
	}

	standPat := Evaluate(*pos)

	if pos.WhiteMove {
		if standPat >= beta {
			return beta
		}
		if standPat > alpha {
			alpha = standPat
		}

		moves := pos.GenerateLegalMoves(pieceMoves)
		captures := filterCaptures(moves)
		sortMoves(captures)

		for _, move := range captures {
			undo := pos.MakeMove(move)
			score := quiescence(pos, pieceMoves, alpha, beta, ctx)
			pos.UnmakeMove(move, undo)

			if ctx.stopped.Load() {
				return 0
			}

			if score > alpha {
				alpha = score
			}
			if alpha >= beta {
				break
			}
		}
		return alpha
	} else {
		if standPat <= alpha {
			return alpha
		}
		if standPat < beta {
			beta = standPat
		}

		moves := pos.GenerateLegalMoves(pieceMoves)
		captures := filterCaptures(moves)
		sortMoves(captures)

		for _, move := range captures {
			undo := pos.MakeMove(move)
			score := quiescence(pos, pieceMoves, alpha, beta, ctx)
			pos.UnmakeMove(move, undo)

			if ctx.stopped.Load() {
				return 0
			}

			if score < beta {
				beta = score
			}
			if alpha >= beta {
				break
			}
		}
		return beta
	}
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
