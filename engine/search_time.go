package engine

import (
	"sync/atomic"
	"time"

	"chess/board"
)

// SearchContext holds state for time-managed search.
type SearchContext struct {
	startTime time.Time
	timeLimit time.Duration
	nodes     int64
	stopped   atomic.Bool
}

// NewSearchContext creates a new search context with the given time limit.
func NewSearchContext(timeLimit time.Duration) *SearchContext {
	return &SearchContext{
		startTime: time.Now(),
		timeLimit: timeLimit,
	}
}

// checkTimeout checks if time has expired (called every N nodes).
func (ctx *SearchContext) checkTimeout() bool {
	if ctx.stopped.Load() {
		return true
	}
	if time.Since(ctx.startTime) >= ctx.timeLimit {
		ctx.stopped.Store(true)
		return true
	}
	return false
}

// Stop signals the search to stop.
func (ctx *SearchContext) Stop() {
	ctx.stopped.Store(true)
}

// Elapsed returns time elapsed since search started.
func (ctx *SearchContext) Elapsed() time.Duration {
	return time.Since(ctx.startTime)
}

// SearchResultTimed contains the best move with depth info.
type SearchResultTimed struct {
	Move     board.Move
	Score    int
	Depth    int
	Nodes    int64
	Time     time.Duration
	FromBook bool
}

// SearchWithTime performs iterative deepening search with time limit.
// Uses the default session for backward compatibility.
func SearchWithTime(pos board.Position, pieceMoves board.PieceMoves, timeLimit time.Duration) SearchResultTimed {
	return getDefaultSession().SearchWithTime(pos, pieceMoves, timeLimit)
}

// Emergency buffer to account for network lag and UCI overhead (in ms).
const emergencyBuffer = 200

// AllocateTime calculates how much time to spend on a move.
// wtime/btime are in milliseconds, returns duration.
// Includes emergency buffer to prevent time losses from lag.
func AllocateTime(wtime, btime, winc, binc int, isWhite bool, movestogo int) time.Duration {
	var myTime, myInc int
	if isWhite {
		myTime = wtime
		myInc = winc
	} else {
		myTime = btime
		myInc = binc
	}

	var allocated int

	// If movestogo is specified, divide time by moves remaining
	if movestogo > 0 {
		// Use most of allotted time + increment
		allocated = myTime/movestogo + myInc*3/4
	} else {
		// Otherwise, assume ~30 moves remaining
		// Allocate time/30 + 3/4 of increment
		allocated = myTime/30 + myInc*3/4

		// Minimum 100ms, maximum 1/3 of remaining time
		if allocated < 100 {
			allocated = 100
		}
		if allocated > myTime/3 {
			allocated = myTime / 3
		}
	}

	// Apply emergency buffer to prevent time losses from network lag
	allocated -= emergencyBuffer
	if allocated < 50 {
		allocated = 50 // Absolute minimum
	}

	return time.Duration(allocated) * time.Millisecond
}
