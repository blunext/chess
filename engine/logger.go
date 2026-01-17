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
	Source    string // "Book" or "Search"
	Score     string // e.g. "30cp", "Mate in 5"
	Depth     int
	Nodes     int64
	Duration  time.Duration
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

// Close closes the logger channel and file
func (l *Logger) Close() {
	close(l.queue)
	<-l.done // Wait for writer to finish
	l.file.Close()
}

// writer is the background goroutine that writes to the file
func (l *Logger) writer() {
	for info := range l.queue {
		// Format: [YYYY-MM-DD HH:MM:SS] Move: ... | Score: ... | ...
		//line := fmt.Sprintf("[%s] Move: %-5s | Score: %-8s | Depth: %d | Source: %-6s | Nodes: %-8d | Time: %-6s | FEN: %s\n",
		line := fmt.Sprintf("%s | M: %-5s | Sc: %-8s | %-6s | Ns: %-8d | T: %-8s | FEN: %s\n",
			info.Timestamp.Format("15:04:05"),
			info.Move,
			info.Score,
			//info.Depth,
			info.Source,
			info.Nodes,
			info.Duration.Round(10*time.Millisecond),
			info.FEN,
		)
		l.file.WriteString(line)
	}
	l.done <- true
}
