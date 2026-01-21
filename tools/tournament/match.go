package main

import (
	"fmt"
	"strings"
	"time"
)

// GameResult represents the outcome of a single game
type GameResult int

const (
	ResultWhiteWins GameResult = iota
	ResultBlackWins
	ResultDraw
	ResultError
)

// TournamentResult holds cumulative tournament results
type TournamentResult struct {
	Wins       int // Engine1 wins
	Draws      int
	Losses     int // Engine1 losses
	EloDiff    float64
	EloError   float64
	LOS        float64
	LLR        float64
	SPRTResult string
}

// RunTournament runs the full tournament
func RunTournament(cfg Config) (TournamentResult, error) {
	result := TournamentResult{}

	for gameNum := 1; gameNum <= cfg.Games; gameNum++ {
		// Alternate colors each game
		e1White := (gameNum % 2) == 1

		gameResult, err := playGame(cfg, e1White, gameNum, cfg.Verbose)
		if err != nil {
			fmt.Printf("Game %d error: %v\n", gameNum, err)
			continue
		}

		// Update results from engine1's perspective
		switch gameResult {
		case ResultWhiteWins:
			if e1White {
				result.Wins++
			} else {
				result.Losses++
			}
		case ResultBlackWins:
			if e1White {
				result.Losses++
			} else {
				result.Wins++
			}
		case ResultDraw:
			result.Draws++
		}

		// Print progress
		total := result.Wins + result.Draws + result.Losses
		score := float64(result.Wins) + 0.5*float64(result.Draws)
		pct := 100.0 * score / float64(total)
		fmt.Printf("Game %d/%d: +%d =%d -%d (%.1f%%)\n",
			gameNum, cfg.Games, result.Wins, result.Draws, result.Losses, pct)

		// Check SPRT
		if cfg.UseSPRT && total >= 10 {
			llr, conclusion := SPRT(result.Wins, result.Draws, result.Losses, -5, 0)
			result.LLR = llr
			if conclusion != "" {
				result.SPRTResult = conclusion
				fmt.Printf("SPRT stopped: %s\n", conclusion)
				break
			}
		}
	}

	// Calculate final statistics
	result.EloDiff, result.EloError = EloDiff(result.Wins, result.Draws, result.Losses)
	result.LOS = LOS(result.Wins, result.Draws, result.Losses)

	return result, nil
}

// playGame plays a single game between two engines
func playGame(cfg Config, e1White bool, gameNum int, verbose bool) (GameResult, error) {
	// Start engines
	engine1, err := NewEngine(cfg.Engine1Path)
	if err != nil {
		return ResultError, fmt.Errorf("start engine1: %w", err)
	}
	defer engine1.Quit()

	engine2, err := NewEngine(cfg.Engine2Path)
	if err != nil {
		return ResultError, fmt.Errorf("start engine2: %w", err)
	}
	defer engine2.Quit()

	// UCI handshake
	if err := engine1.UCI(); err != nil {
		return ResultError, fmt.Errorf("engine1 UCI: %w", err)
	}
	if err := engine2.UCI(); err != nil {
		return ResultError, fmt.Errorf("engine2 UCI: %w", err)
	}

	// New game
	engine1.NewGame()
	engine2.NewGame()

	if err := engine1.IsReady(); err != nil {
		return ResultError, fmt.Errorf("engine1 isready: %w", err)
	}
	if err := engine2.IsReady(); err != nil {
		return ResultError, fmt.Errorf("engine2 isready: %w", err)
	}

	// Assign colors
	var white, black *Engine
	if e1White {
		white, black = engine1, engine2
	} else {
		white, black = engine2, engine1
	}

	// Game state
	moves := []string{}
	positions := make(map[string]int) // For repetition detection
	halfmoveClock := 0
	whiteToMove := true
	wtime := int(cfg.BaseTime.Milliseconds())
	btime := int(cfg.BaseTime.Milliseconds())
	winc := int(cfg.Increment.Milliseconds())
	binc := int(cfg.Increment.Milliseconds())

	// Play game
	for moveNum := 0; moveNum < 500; moveNum++ { // Max 500 moves
		var currentEngine *Engine
		var currentTime *int

		if whiteToMove {
			currentEngine = white
			currentTime = &wtime
		} else {
			currentEngine = black
			currentTime = &btime
		}

		// Set position
		currentEngine.Position("startpos", moves)

		// Search
		startTime := time.Now()
		move, err := currentEngine.Go(wtime, btime, winc, binc)
		elapsed := time.Since(startTime)

		if err != nil {
			return ResultError, fmt.Errorf("engine error: %w", err)
		}

		// Handle no legal moves (checkmate or stalemate)
		if move == "0000" || move == "(none)" {
			// Current player has no moves - they lost or it's stalemate
			if whiteToMove {
				return ResultBlackWins, nil
			}
			return ResultWhiteWins, nil
		}

		// Update time
		*currentTime -= int(elapsed.Milliseconds())
		if whiteToMove {
			*currentTime += winc
		} else {
			*currentTime += binc
		}

		// Check time loss
		if *currentTime <= 0 {
			if whiteToMove {
				return ResultBlackWins, nil
			}
			return ResultWhiteWins, nil
		}

		moves = append(moves, move)

		// Update halfmove clock (simplified: reset on pawn moves or captures)
		if strings.Contains(move, "x") || len(move) == 4 && (move[1] == '2' || move[1] == '7') {
			halfmoveClock = 0
		} else {
			halfmoveClock++
		}

		// 50 move rule
		if halfmoveClock >= 100 {
			return ResultDraw, nil
		}

		// Detect repetition by checking if last 4 moves repeat 3 times (12-move cycle)
		// Pattern: [A B C D] [A B C D] [A B C D] = draw
		if len(moves) >= 12 {
			last4 := strings.Join(moves[len(moves)-4:], " ")
			prev4 := strings.Join(moves[len(moves)-8:len(moves)-4], " ")
			prev8 := strings.Join(moves[len(moves)-12:len(moves)-8], " ")
			if last4 == prev4 && prev4 == prev8 {
				return ResultDraw, nil
			}
		}
		_ = positions // Keep for future proper FEN-based detection

		if verbose {
			color := "W"
			if !whiteToMove {
				color = "B"
			}
			fmt.Printf("  %s: %s (%.0fms)\n", color, move, elapsed.Seconds()*1000)
		}

		whiteToMove = !whiteToMove
	}

	// Max moves reached
	return ResultDraw, nil
}
