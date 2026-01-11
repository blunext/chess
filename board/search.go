package board

// SearchResult contains the best move and its evaluation.
type SearchResult struct {
	Move  Move
	Score int
}

// Search finds the best move using minimax algorithm.
func Search(pos Position, pieceMoves PieceMoves, depth int) SearchResult {
	moves := pos.GenerateLegalMoves(pieceMoves)

	if len(moves) == 0 {
		// No legal moves - checkmate or stalemate
		if pos.IsInCheck() {
			// Checkmate
			if pos.WhiteMove {
				return SearchResult{Score: -100000} // White is mated
			}
			return SearchResult{Score: 100000} // Black is mated
		}
		// Stalemate
		return SearchResult{Score: 0}
	}

	var bestMove Move
	var bestScore int

	if pos.WhiteMove {
		bestScore = -1000000
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := minimax(&pos, pieceMoves, depth-1)
			pos.UnmakeMove(move, undo)

			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		}
	} else {
		bestScore = 1000000
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := minimax(&pos, pieceMoves, depth-1)
			pos.UnmakeMove(move, undo)

			if score < bestScore {
				bestScore = score
				bestMove = move
			}
		}
	}

	return SearchResult{Move: bestMove, Score: bestScore}
}

// minimax returns the evaluation score for a position.
func minimax(pos *Position, pieceMoves PieceMoves, depth int) int {
	if depth == 0 {
		return Evaluate(*pos)
	}

	moves := pos.GenerateLegalMoves(pieceMoves)

	if len(moves) == 0 {
		if pos.IsInCheck() {
			// Checkmate
			if pos.WhiteMove {
				return -100000 + (10 - depth) // Prefer faster mates
			}
			return 100000 - (10 - depth)
		}
		// Stalemate
		return 0
	}

	if pos.WhiteMove {
		bestScore := -1000000
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := minimax(pos, pieceMoves, depth-1)
			pos.UnmakeMove(move, undo)

			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	} else {
		bestScore := 1000000
		for _, move := range moves {
			undo := pos.MakeMove(move)
			score := minimax(pos, pieceMoves, depth-1)
			pos.UnmakeMove(move, undo)

			if score < bestScore {
				bestScore = score
			}
		}
		return bestScore
	}
}
