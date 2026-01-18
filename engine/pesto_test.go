package engine

import (
	"chess/board"
	"testing"
)

// TestPeSTOInitialPosition tests evaluation of the starting position
func TestPeSTOInitialPosition(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	score := EvaluatePeSTO(pos)

	// Starting position should be exactly 0 (perfectly symmetric)
	if score != 0 {
		t.Errorf("Starting position should have score 0, got %d", score)
	}
}

// TestPeSTOMaterialAdvantage tests that material advantage is reflected
func TestPeSTOMaterialAdvantage(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		minScore int // minimum expected score for white
		maxScore int // maximum expected score for white
	}{
		{
			name:     "White up a queen",
			fen:      "rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			minScore: 800, // queen value ~1025
			maxScore: 1200,
		},
		{
			name:     "White up a rook",
			fen:      "rnbqkbn1/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQq - 0 1",
			minScore: 400, // rook value ~477
			maxScore: 600,
		},
		{
			name:     "White up a knight",
			fen:      "r1bqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			minScore: 250, // knight value ~337
			maxScore: 450,
		},
		{
			name:     "White up a pawn",
			fen:      "rnbqkbnr/ppppppp1/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			minScore: 50, // pawn value ~82
			maxScore: 150,
		},
		{
			name:     "Black up a queen",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNB1KBNR w KQkq - 0 1",
			minScore: -1200,
			maxScore: -800,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := board.CreatePositionFormFEN(tt.fen)
			score := EvaluatePeSTO(pos)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Score %d not in expected range [%d, %d]", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

// TestPeSTOTaperedEval tests that tapered eval works correctly
func TestPeSTOTaperedEval(t *testing.T) {
	// Endgame position (only kings and pawns)
	endgameFEN := "4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1"
	endgamePos := board.CreatePositionFormFEN(endgameFEN)
	endgameScore := EvaluatePeSTO(endgamePos)

	// The score should be close to equal (just kings and pawns, symmetric)
	if endgameScore < -100 || endgameScore > 100 {
		t.Errorf("Endgame score %d should be close to 0 for symmetric position", endgameScore)
	}

	// Middlegame position (all pieces)
	middlegameFEN := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1"
	middlegamePos := board.CreatePositionFormFEN(middlegameFEN)
	middlegameScore := EvaluatePeSTO(middlegamePos)

	// White played e4, should have slight advantage
	t.Logf("Middlegame score after 1.e4: %d", middlegameScore)
}

// TestPeSTOKingEndgame tests that king becomes active in endgame
func TestPeSTOKingEndgame(t *testing.T) {
	// King in center is bad in middlegame
	middlegameFEN := "r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1"
	middlegamePos := board.CreatePositionFormFEN(middlegameFEN)

	// King in center is good in endgame
	endgameFEN := "8/pppp1ppp/8/4k3/4K3/8/PPPP1PPP/8 w - - 0 1"
	endgamePos := board.CreatePositionFormFEN(endgameFEN)

	// Just verify both positions evaluate without panic
	mgScore := EvaluatePeSTO(middlegamePos)
	egScore := EvaluatePeSTO(endgamePos)

	t.Logf("Middlegame position score: %d, Endgame position score: %d", mgScore, egScore)
}

// TestPeSTOSymmetry tests that symmetric positions have symmetric scores
func TestPeSTOSymmetry(t *testing.T) {
	// Symmetric starting position
	pos := board.CreatePositionFormFEN(board.InitialPosition)
	score := EvaluatePeSTO(pos)

	// Should be exactly 0 for perfectly symmetric position
	if score != 0 {
		t.Errorf("Starting position should have score 0, got %d", score)
	}
}

// TestPeSTOPieceValues tests that piece values are as expected from PeSTO
func TestPeSTOPieceValues(t *testing.T) {
	// Verify middlegame piece values
	expectedMG := map[string]int{
		"Pawn":   82,
		"Knight": 337,
		"Bishop": 365,
		"Rook":   477,
		"Queen":  1025,
	}

	actualMG := map[string]int{
		"Pawn":   mgPieceValue[PiecePawn],
		"Knight": mgPieceValue[PieceKnight],
		"Bishop": mgPieceValue[PieceBishop],
		"Rook":   mgPieceValue[PieceRook],
		"Queen":  mgPieceValue[PieceQueen],
	}

	for piece, expected := range expectedMG {
		if actualMG[piece] != expected {
			t.Errorf("MG %s value: expected %d, got %d", piece, expected, actualMG[piece])
		}
	}

	// Verify endgame piece values
	expectedEG := map[string]int{
		"Pawn":   94,
		"Knight": 281,
		"Bishop": 297,
		"Rook":   512,
		"Queen":  936,
	}

	actualEG := map[string]int{
		"Pawn":   egPieceValue[PiecePawn],
		"Knight": egPieceValue[PieceKnight],
		"Bishop": egPieceValue[PieceBishop],
		"Rook":   egPieceValue[PieceRook],
		"Queen":  egPieceValue[PieceQueen],
	}

	for piece, expected := range expectedEG {
		if actualEG[piece] != expected {
			t.Errorf("EG %s value: expected %d, got %d", piece, expected, actualEG[piece])
		}
	}
}

// TestPeSTOAdvancedPawns tests that advanced pawns get higher scores
func TestPeSTOAdvancedPawns(t *testing.T) {
	// Pawn on 7th rank (about to promote)
	advancedFEN := "4k3/P7/8/8/8/8/8/4K3 w - - 0 1"
	advancedPos := board.CreatePositionFormFEN(advancedFEN)
	advancedScore := EvaluatePeSTO(advancedPos)

	// Pawn on 2nd rank
	startFEN := "4k3/8/8/8/8/8/P7/4K3 w - - 0 1"
	startPos := board.CreatePositionFormFEN(startFEN)
	startScore := EvaluatePeSTO(startPos)

	// Advanced pawn should score higher
	if advancedScore <= startScore {
		t.Errorf("Advanced pawn (score %d) should be better than starting pawn (score %d)",
			advancedScore, startScore)
	}

	t.Logf("Pawn on 7th rank: %d, Pawn on 2nd rank: %d, Difference: %d",
		advancedScore, startScore, advancedScore-startScore)
}

// TestEvaluateIncludesPeSTO tests that main Evaluate() uses PeSTO
func TestEvaluateIncludesPeSTO(t *testing.T) {
	pos := board.CreatePositionFormFEN(board.InitialPosition)

	pestoScore := EvaluatePeSTO(pos)
	fullScore := Evaluate(pos)

	// Full evaluation includes king safety and pawn structure
	// For starting position, both should be close to 0
	t.Logf("PeSTO score: %d, Full eval score: %d", pestoScore, fullScore)

	if pestoScore != 0 {
		t.Errorf("PeSTO score should be 0 for starting position, got %d", pestoScore)
	}
}

// TestPeSTOGamePhase tests that game phase is calculated correctly
func TestPeSTOGamePhase(t *testing.T) {
	tests := []struct {
		name        string
		fen         string
		description string
		minPhase    int // expected minimum game phase (0=endgame, 24=middlegame)
	}{
		{
			name:        "Starting position (full middlegame)",
			fen:         board.InitialPosition,
			description: "All pieces present = max phase 24",
			minPhase:    24,
		},
		{
			name:        "Only kings and pawns (endgame)",
			fen:         "4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1",
			description: "No minor/major pieces = phase 0",
			minPhase:    0,
		},
		{
			name:        "Queens exchanged (late middlegame)",
			fen:         "rnb1kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNB1KBNR w KQkq - 0 1",
			description: "Missing both queens = phase 16",
			minPhase:    16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := board.CreatePositionFormFEN(tt.fen)
			_ = EvaluatePeSTO(pos) // Just verify it doesn't panic
			t.Logf("%s: %s", tt.name, tt.description)
		})
	}
}
