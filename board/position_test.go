package board

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPosition_filterColor(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		expected string
	}{
		{"initial position", InitialPosition, "8/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{"initial position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1", "rnbqkbnr/pppppppp/8/8/8/8/8/8 b KQkq - 0 1"},
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
			assert.Equalf(t, expected, board, "filterColor()")
		})
	}
}
