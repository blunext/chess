package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// CLI flags
	engine1 := flag.String("engine1", "", "Path to first engine (required)")
	engine2 := flag.String("engine2", "", "Path to second engine (required)")
	games := flag.Int("games", 100, "Number of games to play")
	tc := flag.String("tc", "5+0.05", "Time control: base+increment in seconds")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent games (not implemented yet)")
	useSPRT := flag.Bool("sprt", false, "Use SPRT for early stopping")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	// Validate required flags
	if *engine1 == "" || *engine2 == "" {
		fmt.Println("Usage: tournament -engine1 <path> -engine2 <path> [options]")
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse time control
	baseTime, increment, err := parseTimeControl(*tc)
	if err != nil {
		fmt.Printf("Invalid time control format: %s\n", err)
		os.Exit(1)
	}

	// Print tournament info
	fmt.Printf("Tournament: %s vs %s\n", *engine1, *engine2)
	fmt.Printf("Time control: %.1f+%.2fs\n", baseTime.Seconds(), increment.Seconds())
	fmt.Printf("Games: %d\n", *games)
	if *useSPRT {
		fmt.Println("SPRT: enabled [-5, 0]")
	}
	fmt.Println(strings.Repeat("-", 50))

	// Ignore concurrency for now
	_ = concurrency

	// Run tournament
	result, err := RunTournament(Config{
		Engine1Path: *engine1,
		Engine2Path: *engine2,
		Games:       *games,
		BaseTime:    baseTime,
		Increment:   increment,
		UseSPRT:     *useSPRT,
		Verbose:     *verbose,
	})

	if err != nil {
		fmt.Printf("Tournament error: %v\n", err)
		os.Exit(1)
	}

	// Print results
	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("RESULTS")
	fmt.Println(strings.Repeat("=", 50))
	printResults(result)
}

// Config holds tournament configuration
type Config struct {
	Engine1Path string
	Engine2Path string
	Games       int
	BaseTime    time.Duration
	Increment   time.Duration
	UseSPRT     bool
	Verbose     bool
}

// parseTimeControl parses "5+0.05" format into base time and increment
func parseTimeControl(tc string) (time.Duration, time.Duration, error) {
	var base, inc float64
	n, err := fmt.Sscanf(tc, "%f+%f", &base, &inc)
	if err != nil || n != 2 {
		return 0, 0, fmt.Errorf("expected format: base+increment (e.g., 5+0.05)")
	}
	return time.Duration(base * float64(time.Second)),
		time.Duration(inc * float64(time.Second)),
		nil
}

func printResults(r TournamentResult) {
	total := r.Wins + r.Draws + r.Losses
	score := float64(r.Wins) + 0.5*float64(r.Draws)
	pct := 100.0 * score / float64(total)

	fmt.Printf("Results: +%d =%d -%d (%.1f%%)\n", r.Wins, r.Draws, r.Losses, pct)
	fmt.Printf("Elo difference: %+.0f ±%.0f (95%% CI)\n", r.EloDiff, r.EloError)
	fmt.Printf("LOS: %.1f%%\n", r.LOS*100)

	if r.SPRTResult != "" {
		fmt.Printf("\nSPRT [-5, 0]: LLR = %.2f\n", r.LLR)
		fmt.Printf("Conclusion: %s\n", r.SPRTResult)
	}
}
