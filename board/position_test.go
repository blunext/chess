package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_filterColor(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected string
	}{
		{"initial position", InitialPosition, "8/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{"black move", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1", "rnbqkbnr/pppppppp/8/8/8/8/8/8 b KQkq - 0 1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			board := CreatePositionFormFEN(tc.board).filterColor()
			if board.WhiteMove { // we don't clear the color position in filterColor()
				board.Black = 0
			} else {
				board.White = 0
			}
			expected := CreatePositionFormFEN(tc.expected)
			// Clear hash for comparison since filterColor doesn't update hash
			board.Hash = 0
			expected.Hash = 0
			assert.Equalf(t, expected, board, "filterColor()")
		})
	}
}
