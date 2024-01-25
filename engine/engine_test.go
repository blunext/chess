package engine

import (
	"chess/board"
	"chess/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPosition_AllLegalMoves(t *testing.T) {
	tests := []struct {
		name     string
		board    string
		piece    board.Piece
		expected []board.Bitboard
		all      bool
	}{
		{"initial position", board.InitialPosition, board.Pawn, nil, false},
		{"open bishop c1, white move", "rnbqkbnr/pppp1ppp/4p3/8/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1", board.Bishop,
			[]board.Bitboard{0x2400000000000820, 0x2400000000100020, 0x2400000020000020, 0x2400004000000020, 0x2400800000000020}, true},
		{"open bishop c1, white move", "rnbqkbnr/pppp1p1p/4p3/6p1/8/3P4/PPP1PPPP/RNBQKBNR w KQkq - 0 1", board.Bishop,
			[]board.Bitboard{0x2400000000000820, 0x2400000000100020, 0x2400000020000020}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			position := board.CreatePositionFormFEN(tc.board)
			pieceMoves := generator.NewGenerator()
			positions := position.AllLegalMoves(pieceMoves, board.Bishop)
			if tc.expected == nil {
				assert.Nil(t, positions)
				return
			}
			for _, pos := range positions {
				figure := pos.GetPiece(tc.piece)
				assert.Contains(t, tc.expected, *figure)
			}
			if tc.all {
				assert.Len(t, positions, len(tc.expected))
			}
		})
	}
}
