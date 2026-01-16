package bench

import (
	"fmt"
	"testing"
	"time"

	"chess/board"
	"chess/engine"
	"chess/generator"
)

// TestSearchDepthBenchmark measures search performance at different depths.
// Run with: go test ./bench -run TestSearchDepthBenchmark -v
func TestSearchDepthBenchmark(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	pm := generator.NewGenerator()

	fmt.Println("\n=== Search Depth Benchmark ===")
	fmt.Println("Position: Initial")
	fmt.Printf("%-7s %-10s %-12s %-15s\n", "Depth", "Move", "Nodes", "Time")
	fmt.Println("----------------------------------------------")

	for depth := 1; depth <= 8; depth++ {
		start := time.Now()
		result := engine.Search(pos, pm, depth)
		elapsed := time.Since(start)

		fmt.Printf("%-7d %-10s %-12d %-15v\n",
			depth, result.Move.ToUCI(), result.Nodes, elapsed)

		// Stop if taking too long
		if elapsed > 10*time.Second {
			fmt.Println("Stopping - exceeded 10s threshold")
			break
		}
	}
}

// TestSearchTacticalBenchmark measures search on a tactical position.
func TestSearchTacticalBenchmark(t *testing.T) {
	// Kiwipete position - lots of tactics
	pos := board.CreatePositionFormFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	pm := generator.NewGenerator()

	fmt.Println("\n=== Tactical Position Benchmark ===")
	fmt.Println("Position: Kiwipete")
	fmt.Printf("%-7s %-10s %-12s %-15s\n", "Depth", "Move", "Nodes", "Time")
	fmt.Println("----------------------------------------------")

	for depth := 1; depth <= 6; depth++ {
		start := time.Now()
		result := engine.Search(pos, pm, depth)
		elapsed := time.Since(start)

		fmt.Printf("%-7d %-10s %-12d %-15v\n",
			depth, result.Move.ToUCI(), result.Nodes, elapsed)

		if elapsed > 10*time.Second {
			fmt.Println("Stopping - exceeded 10s threshold")
			break
		}
	}
}
