package engine

import (
	"fmt"
	"os"
	"time"
)

// LogInfo contains all the data points to be logged for a move
type LogInfo struct {
	Timestamp time.Time
	FEN       string
	Move      string
	Piece     string // piece that moved (e.g. "Pawn", "Knight", etc.)
	Source    string // "Book" or "Search"
	Score     string // e.g. "30cp", "Mate in 5"
	Depth     int
	Nodes     int64
	Duration  time.Duration
	GoParams  string // go command parameters (e.g. "wtime:180000 btime:178000")
}

// Logger handles threaded logging to a file
type Logger struct {
	file  *os.File
	queue chan LogInfo
	done  chan bool
}

// NewLogger creates a new logger instance
func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		file:  file,
		queue: make(chan LogInfo, 100), // Buffer up to 100 moves
		done:  make(chan bool),
	}

	go l.writer()

	return l, nil
}

// Log sends a log entry to the writer queue
func (l *Logger) Log(info LogInfo) {
	select {
	case l.queue <- info:
		// Queued successfully
	default:
		// Channel full, drop log to avoid blocking engine
		fmt.Println("Warning: Log queue full, dropping entry")
	}
}

// LogGameStart logs the start of a new game with parameters
func (l *Logger) LogGameStart(params string) {
	if l == nil {
		return
	}
	line := fmt.Sprintf("\n=== NEW GAME STARTED === %s | %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		params,
	)
	l.file.WriteString(line)
}

// Close closes the logger channel and file
func (l *Logger) Close() {
	close(l.queue)
	<-l.done // Wait for writer to finish
	l.file.Close()
}

// writer is the background goroutine that writes to the file
func (l *Logger) writer() {
	for info := range l.queue {
		// Get piece prefix (first letter)
		piecePrefix := ""
		if info.Piece != "" {
			piecePrefix = string(info.Piece[0])
		}

		// Get source prefix (B for Book, S for Search)
		sourcePrefix := "S"
		if info.Source == "Book" {
			sourcePrefix = "B"
		}

		// Build line with optional GoParams
		goParams := ""
		if info.GoParams != "" {
			goParams = " | " + info.GoParams
		}

		line := fmt.Sprintf("%s | M/%s: %s%-5s | Sc: %-8s | Ns: %-8d | T: %-8s | FEN: %s%s\n",
			info.Timestamp.Format("01-02 15:04:05"),
			sourcePrefix,
			piecePrefix,
			info.Move,
			info.Score,
			info.Nodes,
			info.Duration.Round(10*time.Millisecond),
			info.FEN,
			goParams,
		)
		l.file.WriteString(line)
	}
	l.done <- true
}
