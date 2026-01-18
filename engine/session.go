package engine

import (
	"fmt"
	"math/rand"
	"time"

	"chess/board"
	"chess/book"
)

// Session holds per-game state that should be isolated between concurrent games.
// This allows running multiple games in parallel using separate goroutines.
type Session struct {
	TT          *TranspositionTable
	bookRng     *rand.Rand
	debugLogger *Logger // Optional debug logger for detailed search info
}

// NewSession creates a new game session with its own transposition table.
// hashSizeMB specifies the size of the transposition table in megabytes.
func NewSession(hashSizeMB int) *Session {
	return &Session{
		TT:      NewTranspositionTable(hashSizeMB),
		bookRng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Clear resets the session state for a new game.
func (s *Session) Clear() {
	if s.TT != nil {
		s.TT.Clear()
	}
}

// ResizeTT creates a new transposition table with the given size.
func (s *Session) ResizeTT(sizeMB int) {
	s.TT = NewTranspositionTable(sizeMB)
}

// SetDebugLogger sets an optional debug logger for detailed search information.
func (s *Session) SetDebugLogger(logger *Logger) {
	s.debugLogger = logger
}

// Search finds the best move using alpha-beta pruning with fixed depth.
func (s *Session) Search(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	ctx := NewSearchContext(24 * time.Hour) // Effectively no time limit
	result := s.searchRootDepth(pos, pieceMoves, depth, ctx)
	return SearchResult{
		Move:  result.Move,
		Score: result.Score,
		Nodes: ctx.nodes,
	}
}

// SearchWithBook probes the opening book first, then falls back to search.
func (s *Session) SearchWithBook(pos board.Position, pieceMoves board.PieceMoves, depth int) SearchResult {
	// Try opening book first
	if OpeningBook != nil {
		polyHash := book.PolyglotHash(pos)
		if bookMove, ok := OpeningBook.ProbeRandom(polyHash, s.bookRng); ok {
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
	return s.Search(pos, pieceMoves, depth)
}

// SearchWithTime performs iterative deepening search with time limit.
func (s *Session) SearchWithTime(pos board.Position, pieceMoves board.PieceMoves, timeLimit time.Duration) SearchResultTimed {
	// Generate moves early for debug logging
	allMoves := pos.GenerateLegalMoves(pieceMoves)
	sortMoves(allMoves)

	// Debug log: search start with context info
	if s.debugLogger != nil {
		firstMove := ""
		if len(allMoves) > 0 {
			firstMove = allMoves[0].ToUCI()
		}
		ttSize := 0
		if s.TT != nil {
			ttSize = s.TT.SizeMB()
		}
		s.debugLogger.Log(LogInfo{
			Timestamp: time.Now(),
			FEN:       pos.ToFEN(),
			Move:      "START",
			Source:    "Debug",
			Score:     fmt.Sprintf("moves=%d first=%s TT=%dMB", len(allMoves), firstMove, ttSize),
			Depth:     0,
			Nodes:     0,
			Duration:  timeLimit,
		})
	}

	// Try opening book first
	if OpeningBook != nil {
		polyHash := book.PolyglotHash(pos)
		if bookMove, ok := OpeningBook.ProbeRandom(polyHash, s.bookRng); ok {
			legalMoves := pos.GenerateLegalMoves(pieceMoves)
			for _, m := range legalMoves {
				if m.From == bookMove.From && m.To == bookMove.To && m.Promotion == bookMove.Promotion {
					if s.debugLogger != nil {
						s.debugLogger.Log(LogInfo{
							Timestamp: time.Now(),
							FEN:       pos.ToFEN(),
							Move:      "BOOK:" + m.ToUCI(),
							Source:    "Debug",
							Score:     "book",
							Depth:     0,
							Nodes:     0,
							Duration:  0,
						})
					}
					return SearchResultTimed{Move: m, FromBook: true}
				}
			}
		}
	}

	ctx := NewSearchContext(timeLimit)
	var bestResult SearchResultTimed
	var prevMove board.Move // Track previous best move to detect changes

	// Iterative deepening: search depth 1, 2, 3, ... until time runs out
	for depth := 1; depth <= 100; depth++ {
		result := s.searchRootDepth(pos, pieceMoves, depth, ctx)

		// Debug log: iteration result
		if s.debugLogger != nil {
			status := "OK"
			if ctx.stopped.Load() {
				status = "STOPPED"
			}
			// Check if best move changed from previous depth
			moveChanged := depth > 1 && result.Move.ToUCI() != prevMove.ToUCI()
			changeInfo := ""
			if moveChanged {
				changeInfo = fmt.Sprintf(" CHANGED(%s->%s)", prevMove.ToUCI(), result.Move.ToUCI())
			}
			s.debugLogger.Log(LogInfo{
				Timestamp: time.Now(),
				FEN:       pos.ToFEN(),
				Move:      "D" + fmt.Sprint(depth) + ":" + result.Move.ToUCI() + changeInfo,
				Source:    status,
				Score:     fmt.Sprintf("%+d", result.Score),
				Depth:     depth,
				Nodes:     ctx.nodes,
				Duration:  ctx.Elapsed(),
			})
		}

		// If search was stopped mid-way, don't use partial results
		if ctx.stopped.Load() && depth > 1 {
			if s.debugLogger != nil {
				s.debugLogger.Log(LogInfo{
					Timestamp: time.Now(),
					FEN:       pos.ToFEN(),
					Move:      "REJECT",
					Source:    "Debug",
					Score:     fmt.Sprintf("d%d_stopped prev=%s", depth, prevMove.ToUCI()),
					Depth:     depth,
					Nodes:     ctx.nodes,
					Duration:  ctx.Elapsed(),
				})
			}
			break
		}

		// Track previous move before updating
		prevMove = result.Move

		// Update best result
		bestResult = SearchResultTimed{
			Move:  result.Move,
			Score: result.Score,
			Depth: depth,
			Nodes: ctx.nodes,
			Time:  ctx.Elapsed(),
		}

		// Debug log: accepted this depth
		if s.debugLogger != nil && depth > 1 {
			s.debugLogger.Log(LogInfo{
				Timestamp: time.Now(),
				FEN:       pos.ToFEN(),
				Move:      "ACCEPT",
				Source:    "Debug",
				Score:     fmt.Sprintf("d%d=%s", depth, result.Move.ToUCI()),
				Depth:     depth,
				Nodes:     ctx.nodes,
				Duration:  ctx.Elapsed(),
			})
		}

		// If we found a mate, no need to search deeper
		if result.Score > mateScore-100 || result.Score < -mateScore+100 {
			if s.debugLogger != nil {
				s.debugLogger.Log(LogInfo{
					Timestamp: time.Now(),
					FEN:       pos.ToFEN(),
					Move:      "MATE",
					Source:    "Debug",
					Score:     fmt.Sprintf("%+d", result.Score),
					Depth:     depth,
					Nodes:     ctx.nodes,
					Duration:  ctx.Elapsed(),
				})
			}
			break
		}

		// Check if we have time for another iteration
		// Heuristic: next depth takes ~3-4x longer
		if ctx.Elapsed()*4 >= timeLimit {
			if s.debugLogger != nil {
				s.debugLogger.Log(LogInfo{
					Timestamp: time.Now(),
					FEN:       pos.ToFEN(),
					Move:      "TIMECUT",
					Source:    "Debug",
					Score:     fmt.Sprintf("%.1fx", float64(ctx.Elapsed())*4/float64(timeLimit)),
					Depth:     depth,
					Nodes:     ctx.nodes,
					Duration:  ctx.Elapsed(),
				})
			}
			break
		}
	}

	// Debug log: final result
	if s.debugLogger != nil {
		s.debugLogger.Log(LogInfo{
			Timestamp: time.Now(),
			FEN:       pos.ToFEN(),
			Move:      "FINAL:" + bestResult.Move.ToUCI(),
			Source:    "Debug",
			Score:     fmt.Sprintf("%+d", bestResult.Score),
			Depth:     bestResult.Depth,
			Nodes:     bestResult.Nodes,
			Duration:  bestResult.Time,
		})
	}

	return bestResult
}

// searchRootDepth searches to a fixed depth (used by Search and iterative deepening).
func (s *Session) searchRootDepth(pos board.Position, pieceMoves board.PieceMoves, depth int, ctx *SearchContext) SearchResult {
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
			score := s.alphaBeta(&pos, pieceMoves, depth-1, alpha, beta, true, ctx)
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
			score := s.alphaBeta(&pos, pieceMoves, depth-1, alpha, beta, true, ctx)
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
func (s *Session) alphaBeta(pos *board.Position, pieceMoves board.PieceMoves, depth int, alpha, beta int, nullMoveAllowed bool, ctx *SearchContext) int {
	// Increment node counter and check timeout every 2048 nodes
	ctx.nodes++
	if ctx.nodes&2047 == 0 && ctx.checkTimeout() {
		return 0
	}

	if depth == 0 {
		return s.quiescence(pos, pieceMoves, alpha, beta, ctx)
	}

	alphaOrig := alpha
	hash := pos.Hash

	// Check if we're in check (used for extensions and null move)
	inCheck := pos.IsInCheck()

	// Check Extension: Extend search by 1 ply when in check
	// Must be done BEFORE TT probe so depth is correct
	if inCheck {
		depth++
	}

	// Probe transposition table (after check extension)
	var ttMove board.Move
	if s.TT != nil {
		if entry, found := s.TT.Probe(hash); found {
			ttMove = entry.BestMove
			if int(entry.Depth) >= depth {
				score := int(entry.Score)
				switch entry.Flag {
				case TTFlagExact:
					return score
				case TTFlagLower:
					// Lower bound: real score >= stored score
					// Cutoff if stored score >= beta (fail-high)
					if score >= beta {
						return score
					}
				case TTFlagUpper:
					// Upper bound: real score <= stored score
					// Cutoff if stored score <= alpha (fail-low)
					if score <= alpha {
						return score
					}
				}
			}
		}
	}

	// Null Move Pruning
	if UseNullMovePruning && nullMoveAllowed && depth >= 4 && !inCheck && !isEndgame(pos) {
		oldHash := pos.Hash
		oldEnPassant := pos.EnPassant

		pos.WhiteMove = !pos.WhiteMove
		pos.EnPassant = 0
		pos.Hash ^= board.HashSide()

		nullScore := -s.alphaBeta(pos, pieceMoves, depth-1-nullMoveReduction, -beta, -beta+1, false, ctx)

		pos.WhiteMove = !pos.WhiteMove
		pos.EnPassant = oldEnPassant
		pos.Hash = oldHash

		if nullScore >= beta {
			return beta
		}
	}

	moves := pos.GenerateLegalMoves(pieceMoves)

	// Put TT move first if available
	if ttMove != (board.Move{}) {
		for i, m := range moves {
			if m.From == ttMove.From && m.To == ttMove.To && m.Promotion == ttMove.Promotion {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
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
			score := s.alphaBeta(pos, pieceMoves, depth-1, alpha, beta, true, ctx)
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
			score := s.alphaBeta(pos, pieceMoves, depth-1, alpha, beta, true, ctx)
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
	if s.TT != nil && !ctx.stopped.Load() {
		var flag TTFlag
		if bestScore <= alphaOrig {
			flag = TTFlagUpper
		} else if bestScore >= beta {
			flag = TTFlagLower
		} else {
			flag = TTFlagExact
		}
		s.TT.Store(hash, int16(bestScore), int8(depth), flag, bestMove)
	}

	return bestScore
}

// quiescence continues search only for captures to avoid horizon effect.
// Also checks for mate threats to avoid missing tactics like Qa2 bug.
func (s *Session) quiescence(pos *board.Position, pieceMoves board.PieceMoves, alpha, beta int, ctx *SearchContext) int {
	ctx.nodes++
	if ctx.nodes&2047 == 0 && ctx.checkTimeout() {
		return 0
	}

	standPat := Evaluate(*pos)

	// Check if we're in check - must search all evasions, not just captures
	inCheck := pos.IsInCheck()

	// Mate threat detection disabled - 22x performance overhead was too high.
	// See ROADMAP.md for alternative approach (Mate Threat Extensions in main search).
	opponentHasMateThreat := false

	if pos.WhiteMove {
		if standPat >= beta && !opponentHasMateThreat {
			return beta
		}
		if standPat > alpha {
			alpha = standPat
		}

		moves := pos.GenerateLegalMoves(pieceMoves)

		// If in check or under mate threat, search all moves; otherwise just captures
		var searchMoves []board.Move
		if inCheck || opponentHasMateThreat {
			searchMoves = moves
		} else {
			searchMoves = filterCaptures(moves)
		}
		sortMoves(searchMoves)

		for _, move := range searchMoves {
			undo := pos.MakeMove(move)
			score := s.quiescence(pos, pieceMoves, alpha, beta, ctx)
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
		// Black to move
		if standPat <= alpha && !opponentHasMateThreat {
			return alpha
		}
		if standPat < beta {
			beta = standPat
		}

		moves := pos.GenerateLegalMoves(pieceMoves)

		var searchMoves []board.Move
		if inCheck || opponentHasMateThreat {
			searchMoves = moves
		} else {
			searchMoves = filterCaptures(moves)
		}
		sortMoves(searchMoves)

		for _, move := range searchMoves {
			undo := pos.MakeMove(move)
			score := s.quiescence(pos, pieceMoves, alpha, beta, ctx)
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

// hasMateInOne checks if the opponent (not the side to move) can deliver mate in 1.
// This is used in quiescence to detect mate threats before using stand-pat cutoff.
// We simulate giving the opponent a free move and check if they can mate.
func (s *Session) hasMateInOne(pos *board.Position, pieceMoves board.PieceMoves, opponentIsWhite bool) bool {
	// Create a copy and flip the side to simulate opponent having a free move
	tempPos := *pos
	tempPos.WhiteMove = opponentIsWhite
	tempPos.Hash ^= board.HashSide() // Update hash for side change

	moves := tempPos.GenerateLegalMoves(pieceMoves)

	for _, m := range moves {
		undo := tempPos.MakeMove(m)
		// After the move, check if defender has any legal moves
		replies := tempPos.GenerateLegalMoves(pieceMoves)
		isMate := len(replies) == 0 && tempPos.IsInCheck()
		tempPos.UnmakeMove(m, undo)

		if isMate {
			return true
		}
	}

	return false
}
