package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"chess/board"
	"chess/generator"
)

// DefaultSearchDepth is the default depth for engine search
const DefaultSearchDepth = 8

// Play starts an interactive game in the terminal
func Play() {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := generator.NewGenerator()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Chess Engine Interactive Mode ===")
	fmt.Println("Enter moves in UCI format (e.g., e2e4, e7e8q for promotion)")
	fmt.Println("Commands: 'quit', 'undo', 'fen', 'moves'")
	fmt.Println()

	type historyEntry struct {
		move board.Move
		undo board.UndoInfo
	}
	var history []historyEntry

	for {
		// Display board
		fmt.Println(pos.Pretty())

		// Check game state
		legalMoves := pos.GenerateLegalMoves(pm)
		if len(legalMoves) == 0 {
			if pos.IsInCheck() {
				if pos.WhiteMove {
					fmt.Println("Checkmate! Black wins!")
				} else {
					fmt.Println("Checkmate! White wins!")
				}
			} else {
				fmt.Println("Stalemate! Draw!")
			}
			break
		}

		if pos.IsInCheck() {
			fmt.Println("Check!")
		}

		// Prompt
		side := "White"
		if !pos.WhiteMove {
			side = "Black"
		}
		fmt.Printf("%s to move: ", side)

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}
		input = strings.TrimSpace(input)

		// Handle commands
		switch input {
		case "quit", "q":
			fmt.Println("Goodbye!")
			return
		case "undo", "u":
			if len(history) > 0 {
				last := history[len(history)-1]
				pos.UnmakeMove(last.move, last.undo)
				history = history[:len(history)-1]
				fmt.Println("Move undone.")
			} else {
				fmt.Println("No moves to undo.")
			}
			continue
		case "fen":
			fmt.Printf("FEN: (not implemented)\n")
			continue
		case "moves":
			fmt.Println("Legal moves:")
			for _, m := range legalMoves {
				fmt.Printf("  %s\n", m.ToUCI())
			}
			continue
		case "engine", "e":
			// Let engine play
			fmt.Println("Engine thinking...")
			result := Search(pos, pm, DefaultSearchDepth)
			if result.Move == (board.Move{}) {
				fmt.Println("Engine has no move!")
				continue
			}
			fmt.Printf("Engine plays: %s (score: %d)\n", result.Move.ToUCI(), result.Score)
			undo := pos.MakeMove(result.Move)
			history = append(history, historyEntry{result.Move, undo})
			continue
		}

		// Parse and validate move
		move, ok := parseUCIMove(input, legalMoves)
		if !ok {
			fmt.Printf("Invalid move: %s\n", input)
			fmt.Println("Type 'moves' to see legal moves.")
			continue
		}

		// Make move
		undo := pos.MakeMove(move)
		history = append(history, historyEntry{move, undo})

		// Engine response
		fmt.Println("\nEngine thinking...")
		legalMoves = pos.GenerateLegalMoves(pm)
		if len(legalMoves) == 0 {
			continue // Will be handled at top of loop
		}

		result := Search(pos, pm, DefaultSearchDepth)
		if result.Move == (board.Move{}) {
			continue
		}
		fmt.Printf("Engine plays: %s (score: %d)\n", result.Move.ToUCI(), result.Score)
		undo = pos.MakeMove(result.Move)
		history = append(history, historyEntry{result.Move, undo})
		fmt.Println()
	}
}

// parseUCIMove finds the matching move from legal moves
func parseUCIMove(uci string, legalMoves []board.Move) (board.Move, bool) {
	uci = strings.ToLower(strings.TrimSpace(uci))
	for _, m := range legalMoves {
		if m.ToUCI() == uci {
			return m, true
		}
	}
	return board.Move{}, false
}
