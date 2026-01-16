package board

import "math/rand"

// Zobrist hashing keys for position identification.
// These are pseudo-random 64-bit numbers used to create a unique hash for each position.
// XOR properties enable incremental updates: hash ^= key removes/adds a piece.

var (
	// zobristPiece[color][piece-1][square] - keys for each piece on each square
	// color: 0=white, 1=black; piece: Pawn=1..King=6 (subtract 1 for index)
	zobristPiece [2][6][64]uint64

	// zobristCastling[rights] - keys for castling combinations (0-15)
	zobristCastling [16]uint64

	// zobristEnPassant[file] - keys for en passant file (0-7, a-h)
	zobristEnPassant [8]uint64

	// zobristSide - key for side to move (XOR when black to move)
	zobristSide uint64
)

func init() {
	// Use fixed seed for reproducibility (same hash across runs)
	// Polyglot uses specific keys, but for now we use our own
	rng := rand.New(rand.NewSource(0x12345678DEADBEEF))

	// Generate piece keys
	for color := 0; color < 2; color++ {
		for piece := 0; piece < 6; piece++ {
			for sq := 0; sq < 64; sq++ {
				zobristPiece[color][piece][sq] = rng.Uint64()
			}
		}
	}

	// Generate castling keys
	for i := 0; i < 16; i++ {
		zobristCastling[i] = rng.Uint64()
	}

	// Generate en passant keys
	for i := 0; i < 8; i++ {
		zobristEnPassant[i] = rng.Uint64()
	}

	// Generate side to move key
	zobristSide = rng.Uint64()
}

// pieceIndex converts Piece type to zobrist array index (0-5)
func pieceIndex(p Piece) int {
	return int(p) - 1 // Pawn=1 -> 0, Knight=2 -> 1, etc.
}

// colorIndex converts Color to zobrist array index (0 or 1)
func colorIndex(c Color) int {
	return int(c)
}

// ComputeHash calculates the full Zobrist hash for a position.
// This is used when creating a position from FEN.
func (pos *Position) ComputeHash() uint64 {
	var hash uint64

	// Hash all pieces
	allPieces := pos.Pawns | pos.Knights | pos.Bishops | pos.Rooks | pos.Queens | pos.Kings
	for sq := 0; sq < 64; sq++ {
		sqBB := Bitboard(1 << sq)
		if allPieces&sqBB == 0 {
			continue
		}

		// Determine color
		var color int
		if pos.White&sqBB != 0 {
			color = 0
		} else {
			color = 1
		}

		// Determine piece type
		var pieceIdx int
		switch {
		case pos.Pawns&sqBB != 0:
			pieceIdx = 0
		case pos.Knights&sqBB != 0:
			pieceIdx = 1
		case pos.Bishops&sqBB != 0:
			pieceIdx = 2
		case pos.Rooks&sqBB != 0:
			pieceIdx = 3
		case pos.Queens&sqBB != 0:
			pieceIdx = 4
		case pos.Kings&sqBB != 0:
			pieceIdx = 5
		}

		hash ^= zobristPiece[color][pieceIdx][sq]
	}

	// Hash castling rights
	hash ^= zobristCastling[pos.CastleSide]

	// Hash en passant file (if set)
	if pos.EnPassant != 0 {
		// Get file from en passant square
		epSq := 0
		for i := 0; i < 64; i++ {
			if pos.EnPassant&(1<<i) != 0 {
				epSq = i
				break
			}
		}
		file := epSq & 7
		hash ^= zobristEnPassant[file]
	}

	// Hash side to move
	if !pos.WhiteMove {
		hash ^= zobristSide
	}

	return hash
}

// HashPiece returns the Zobrist key for a piece on a square.
// Used for incremental hash updates in MakeMove/UnmakeMove.
func HashPiece(piece Piece, color Color, sq int) uint64 {
	return zobristPiece[colorIndex(color)][pieceIndex(piece)][sq]
}

// HashCastling returns the Zobrist key for castling rights.
func HashCastling(rights uint8) uint64 {
	return zobristCastling[rights]
}

// HashEnPassant returns the Zobrist key for en passant file.
func HashEnPassant(file int) uint64 {
	return zobristEnPassant[file]
}

// HashSide returns the Zobrist key for side to move.
func HashSide() uint64 {
	return zobristSide
}
