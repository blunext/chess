package bench

import (
	"chess/board"
	"chess/generator"
	"chess/magic"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := magic.Prepare(); err != nil {
		panic("Failed to load magic data: " + err.Error())
	}
	os.Exit(m.Run())
}

// BenchmarkGenerateMoves benchmarks move generation from initial position.
func BenchmarkGenerateMoves(b *testing.B) {
	position := board.CreatePositionFormFEN(board.InitialPosition)
	pm := generator.NewGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = position.GenerateMoves(pm)
	}
}

// BenchmarkGenerateMoves_MidGame benchmarks move generation in a typical midgame.
func BenchmarkGenerateMoves_MidGame(b *testing.B) {
	position := board.CreatePositionFormFEN("r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4")
	pm := generator.NewGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = position.GenerateMoves(pm)
	}
}

// BenchmarkGenerateMoves_Complex benchmarks with many sliding pieces active.
func BenchmarkGenerateMoves_Complex(b *testing.B) {
	position := board.CreatePositionFormFEN("r2qr1k1/ppp2ppp/2n1bn2/3p4/3P4/2NBBN2/PPP2PPP/R2QR1K1 w - - 0 10")
	pm := generator.NewGenerator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = position.GenerateMoves(pm)
	}
}
