package uci

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"chess/board"
	"chess/engine"
	"chess/generator"
)

const (
	engineName   = "FromZeroToGM"
	engineAuthor = "FromZeroToGM"
)

// UCI holds the state for the UCI protocol
type UCI struct {
	position   board.Position
	pieceMoves board.PieceMoves
	session    *engine.Session
	logger     *engine.Logger
}

// Start begins the UCI protocol loop
func Start() {
	// Use absolute path or current directory
	// For now, let's use the current working directory but print it out
	cwd, _ := os.Getwd()
	logPath := cwd + "/game.log"

	l, err := engine.NewLogger(logPath)
	if err != nil {
		// Try fallback to just filename if permission issues
		l, _ = engine.NewLogger("game.log")
	}

	uci := &UCI{
		position:   board.CreatePositionFormFEN(board.InitialPosition),
		pieceMoves: generator.NewGenerator(),
		session:    engine.NewSession(engine.DefaultHashMB),
		logger:     l,
	}

	if uci.logger != nil {
		defer uci.logger.Close()
		// Print info to stdout so user might see it in bot logs
		fmt.Printf("info string Logging to %s\n", logPath)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if !uci.handleCommand(line) {
			break
		}
	}
}

// handleCommand processes a single UCI command. Returns false if should quit.
func (uci *UCI) handleCommand(line string) bool {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return true
	}

	cmd := parts[0]

	switch cmd {
	case "uci":
		uci.cmdUCI()
	case "isready":
		uci.cmdIsReady()
	case "setoption":
		uci.cmdSetOption(parts[1:])
	case "ucinewgame":
		uci.cmdNewGame()
	case "position":
		uci.cmdPosition(parts[1:])
	case "go":
		uci.cmdGo(parts[1:])
	case "quit":
		return false
	case "d", "display":
		// Debug: display current position
		fmt.Println(uci.position.Pretty())
	}

	return true
}

// cmdUCI handles the "uci" command
func (uci *UCI) cmdUCI() {
	fmt.Printf("id name %s\n", engineName)
	fmt.Printf("id author %s\n", engineAuthor)
	// Options
	fmt.Printf("option name Hash type spin default %d min 1 max 32768\n", engine.DefaultHashMB)
	fmt.Println("uciok")
}

// cmdIsReady handles the "isready" command
func (uci *UCI) cmdIsReady() {
	fmt.Println("readyok")
}

// cmdSetOption handles "setoption name X value Y"
func (uci *UCI) cmdSetOption(args []string) {
	// setoption name Hash value 128
	if len(args) < 4 {
		return
	}

	// Find "name" and "value" positions
	var name, value string
	for i := 0; i < len(args); i++ {
		if args[i] == "name" && i+1 < len(args) {
			name = args[i+1]
		} else if args[i] == "value" && i+1 < len(args) {
			value = args[i+1]
		}
	}

	switch strings.ToLower(name) {
	case "hash":
		if sizeMB, err := strconv.Atoi(value); err == nil && sizeMB > 0 {
			uci.session.ResizeTT(sizeMB)
			fmt.Printf("info string Hash set to %d MB\n", sizeMB)
		}
	}
}

// cmdNewGame handles the "ucinewgame" command
func (uci *UCI) cmdNewGame() {
	uci.position = board.CreatePositionFormFEN(board.InitialPosition)
	// Clear transposition table for new game
	uci.session.Clear()
}

// cmdPosition handles the "position" command
// position startpos [moves e2e4 e7e5 ...]
// position fen <fen> [moves e2e4 e7e5 ...]
func (uci *UCI) cmdPosition(args []string) {
	if len(args) == 0 {
		return
	}

	var movesIdx int

	if args[0] == "startpos" {
		uci.position = board.CreatePositionFormFEN(board.InitialPosition)
		movesIdx = 1
	} else if args[0] == "fen" {
		// Find where "moves" starts (if present)
		fenEnd := len(args)
		for i, arg := range args {
			if arg == "moves" {
				fenEnd = i
				break
			}
		}
		fenStr := strings.Join(args[1:fenEnd], " ")
		uci.position = board.CreatePositionFormFEN(fenStr)
		movesIdx = fenEnd
	} else {
		return
	}

	// Apply moves if present
	if movesIdx < len(args) && args[movesIdx] == "moves" {
		for _, moveStr := range args[movesIdx+1:] {
			move := uci.parseMove(moveStr)
			if move != (board.Move{}) {
				uci.position.MakeMove(move)
			}
		}
	}
}

// cmdGo handles the "go" command
// go depth 5
// go wtime 300000 btime 300000 winc 0 binc 0
// go infinite
// go movetime 1000
func (uci *UCI) cmdGo(args []string) {
	var depth int
	var wtime, btime, winc, binc, movestogo, movetime int
	useTimeControl := false
	infinite := false

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "depth":
			if i+1 < len(args) {
				if d, err := strconv.Atoi(args[i+1]); err == nil {
					depth = d
				}
				i++
			}
		case "wtime":
			if i+1 < len(args) {
				wtime, _ = strconv.Atoi(args[i+1])
				useTimeControl = true
				i++
			}
		case "btime":
			if i+1 < len(args) {
				btime, _ = strconv.Atoi(args[i+1])
				useTimeControl = true
				i++
			}
		case "winc":
			if i+1 < len(args) {
				winc, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "binc":
			if i+1 < len(args) {
				binc, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "movestogo":
			if i+1 < len(args) {
				movestogo, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "movetime":
			if i+1 < len(args) {
				movetime, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "infinite":
			infinite = true
		}
	}

	var result engine.SearchResultTimed

	if movetime > 0 {
		// Fixed time per move
		timeLimit := time.Duration(movetime) * time.Millisecond
		result = uci.session.SearchWithTime(uci.position, uci.pieceMoves, timeLimit)
	} else if useTimeControl {
		// Time control: allocate time based on remaining time
		timeLimit := engine.AllocateTime(wtime, btime, winc, binc, uci.position.WhiteMove, movestogo)
		result = uci.session.SearchWithTime(uci.position, uci.pieceMoves, timeLimit)
	} else if depth > 0 {
		// Fixed depth search
		fixedResult := uci.session.SearchWithBook(uci.position, uci.pieceMoves, depth)
		result = engine.SearchResultTimed{
			Move:     fixedResult.Move,
			Score:    fixedResult.Score,
			Depth:    depth,
			Nodes:    fixedResult.Nodes,
			FromBook: fixedResult.FromBook,
		}
	} else if infinite {
		// Infinite: use long time limit
		result = uci.session.SearchWithTime(uci.position, uci.pieceMoves, 24*time.Hour)
	} else {
		// Default: use default depth
		fixedResult := uci.session.SearchWithBook(uci.position, uci.pieceMoves, engine.DefaultSearchDepth)
		result = engine.SearchResultTimed{
			Move:     fixedResult.Move,
			Score:    fixedResult.Score,
			Depth:    engine.DefaultSearchDepth,
			Nodes:    fixedResult.Nodes,
			FromBook: fixedResult.FromBook,
		}
	}

	duration := result.Time
	if duration == 0 {
		duration = time.Millisecond // Avoid division by zero
	}

	// Output best move
	if result.Move != (board.Move{}) {
		fmt.Printf("bestmove %s\n", result.Move.ToUCI())

		// Log the valid move
		if uci.logger != nil {
			scoreStr := fmt.Sprintf("%d cp", result.Score)
			// Simple mate detection in score string
			if result.Score > 90000 {
				scoreStr = "Mate +"
			} else if result.Score < -90000 {
				scoreStr = "Mate -"
			}

			uci.logger.Log(engine.LogInfo{
				Timestamp: time.Now(),
				FEN:       uci.position.ToFEN(),
				Move:      result.Move.ToUCI(),
				Source: func() string {
					if result.FromBook {
						return "Book"
					}
					return "Search"
				}(),
				Score:    scoreStr,
				Depth:    result.Depth,
				Nodes:    result.Nodes,
				Duration: duration,
			})
		}
	} else {
		// No legal moves - should not happen in normal games
		fmt.Println("bestmove 0000")
	}
}

// parseMove converts UCI move string to Move struct
func (uci *UCI) parseMove(moveStr string) board.Move {
	moveStr = strings.ToLower(moveStr)

	// Get legal moves and find matching one
	legalMoves := uci.position.GenerateLegalMoves(uci.pieceMoves)
	for _, m := range legalMoves {
		if m.ToUCI() == moveStr {
			return m
		}
	}

	return board.Move{}
}
