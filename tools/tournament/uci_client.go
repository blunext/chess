package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Engine represents a UCI chess engine process
type Engine struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
	name    string
	mu      sync.Mutex
	timeout time.Duration
}

// NewEngine starts a new engine process
func NewEngine(path string) (*Engine, error) {
	cmd := exec.Command(path)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start engine: %w", err)
	}

	e := &Engine{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  bufio.NewReader(stdout),
		name:    path,
		timeout: 10 * time.Second,
	}

	return e, nil
}

// UCI performs the UCI handshake
func (e *Engine) UCI() error {
	e.send("uci")
	_, err := e.readUntil("uciok")
	return err
}

// IsReady sends isready and waits for readyok
func (e *Engine) IsReady() error {
	e.send("isready")
	_, err := e.readUntil("readyok")
	return err
}

// NewGame sends ucinewgame command
func (e *Engine) NewGame() {
	e.send("ucinewgame")
}

// Position sets the position
func (e *Engine) Position(fen string, moves []string) {
	var cmd string
	if fen == "" || fen == "startpos" {
		cmd = "position startpos"
	} else {
		cmd = fmt.Sprintf("position fen %s", fen)
	}

	if len(moves) > 0 {
		cmd += " moves " + strings.Join(moves, " ")
	}

	e.send(cmd)
}

// Go starts search and returns best move
func (e *Engine) Go(wtime, btime, winc, binc int) (string, error) {
	cmd := fmt.Sprintf("go wtime %d btime %d winc %d binc %d", wtime, btime, winc, binc)
	e.send(cmd)

	// Read until bestmove
	for {
		line, err := e.readLine()
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(line, "bestmove") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
			return "", fmt.Errorf("invalid bestmove: %s", line)
		}
	}
}

// Quit terminates the engine
func (e *Engine) Quit() {
	e.send("quit")
	_ = e.cmd.Wait()
}

// send writes a command to the engine
func (e *Engine) send(cmd string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, _ = fmt.Fprintln(e.stdin, cmd)
}

// readLine reads a single line with timeout
func (e *Engine) readLine() (string, error) {
	type result struct {
		line string
		err  error
	}

	ch := make(chan result, 1)
	go func() {
		line, err := e.stdout.ReadString('\n')
		ch <- result{strings.TrimSpace(line), err}
	}()

	select {
	case r := <-ch:
		return r.line, r.err
	case <-time.After(e.timeout):
		return "", fmt.Errorf("timeout reading from engine")
	}
}

// readUntil reads lines until finding one with the given prefix
func (e *Engine) readUntil(prefix string) (string, error) {
	for {
		line, err := e.readLine()
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(line, prefix) {
			return line, nil
		}
	}
}
