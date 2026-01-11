package board

// Perft counts the number of leaf nodes at a given depth.
// This is the standard debugging/verification technique for move generators.
// Results can be compared with known reference values.
func (pos *Position) Perft(pieceMoves PieceMoves, depth int) uint64 {
	if depth == 0 {
		return 1
	}

	moves := pos.GenerateLegalMoves(pieceMoves)

	if depth == 1 {
		return uint64(len(moves))
	}

	var nodes uint64
	for _, m := range moves {
		undo := pos.MakeMove(m)
		nodes += pos.Perft(pieceMoves, depth-1)
		pos.UnmakeMove(m, undo)
	}

	return nodes
}

// Divide runs perft for each move at depth-1 and prints results.
// Useful for debugging when perft results don't match.
func (pos *Position) Divide(pieceMoves PieceMoves, depth int) map[string]uint64 {
	result := make(map[string]uint64)
	moves := pos.GenerateLegalMoves(pieceMoves)

	for _, m := range moves {
		undo := pos.MakeMove(m)
		nodes := pos.Perft(pieceMoves, depth-1)
		pos.UnmakeMove(m, undo)

		result[m.ToUCI()] = nodes
	}

	return result
}
