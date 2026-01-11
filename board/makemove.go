package board

// UndoInfo stores the state needed to unmake a move.
// This allows efficient make/unmake without copying the entire position.
type UndoInfo struct {
	CapturedPiece Piece    // piece that was captured (Empty if no capture)
	CastleSide    uint8    // castling rights before the move
	EnPassant     Bitboard // en passant square before the move
	HalfMoveClock uint8    // half-move clock before the move
}

// MakeMove executes a move on the position and returns undo information.
// This modifies the position in-place for performance.
// Call UnmakeMove with the returned UndoInfo to reverse the move.
func (pos *Position) MakeMove(m Move) UndoInfo {
	// Save state for unmake
	undo := UndoInfo{
		CapturedPiece: m.Captured,
		CastleSide:    pos.CastleSide,
		EnPassant:     pos.EnPassant,
		HalfMoveClock: pos.HalfMoveClock,
	}

	// Get color bitboards
	var ourColor, enemyColor *Bitboard
	if pos.WhiteMove {
		ourColor = &pos.White
		enemyColor = &pos.Black
	} else {
		ourColor = &pos.Black
		enemyColor = &pos.White
	}

	// Remove piece from source square
	pieceBB := pos.GetPiece(m.Piece)
	*pieceBB &^= m.From
	*ourColor &^= m.From

	// Handle capture (remove enemy piece)
	if m.Captured != Empty {
		if m.Flags&FlagEnPassant != 0 {
			// En passant: captured pawn is not on 'To' square
			var capturedPawnSq Bitboard
			if pos.WhiteMove {
				capturedPawnSq = m.To >> 8 // pawn is one rank below
			} else {
				capturedPawnSq = m.To << 8 // pawn is one rank above
			}
			pos.Pawns &^= capturedPawnSq
			*enemyColor &^= capturedPawnSq
		} else {
			// Normal capture
			capturedBB := pos.GetPiece(m.Captured)
			*capturedBB &^= m.To
			*enemyColor &^= m.To
		}
	}

	// Add piece to destination square (handle promotion)
	if m.Promotion != Empty {
		// Promote: add the promoted piece instead of pawn
		promoBB := pos.GetPiece(m.Promotion)
		*promoBB |= m.To
	} else {
		// Normal move: add same piece type
		*pieceBB |= m.To
	}
	*ourColor |= m.To

	// Handle castling (move the rook)
	if m.Flags&FlagCastling != 0 {
		if pos.WhiteMove {
			if m.To == IndexToBitBoard(6) { // Kingside O-O
				pos.Rooks &^= IndexToBitBoard(7) // Remove from h1
				pos.Rooks |= IndexToBitBoard(5)  // Add to f1
				pos.White &^= IndexToBitBoard(7)
				pos.White |= IndexToBitBoard(5)
			} else { // Queenside O-O-O
				pos.Rooks &^= IndexToBitBoard(0) // Remove from a1
				pos.Rooks |= IndexToBitBoard(3)  // Add to d1
				pos.White &^= IndexToBitBoard(0)
				pos.White |= IndexToBitBoard(3)
			}
		} else {
			if m.To == IndexToBitBoard(62) { // Kingside O-O
				pos.Rooks &^= IndexToBitBoard(63) // Remove from h8
				pos.Rooks |= IndexToBitBoard(61)  // Add to f8
				pos.Black &^= IndexToBitBoard(63)
				pos.Black |= IndexToBitBoard(61)
			} else { // Queenside O-O-O
				pos.Rooks &^= IndexToBitBoard(56) // Remove from a8
				pos.Rooks |= IndexToBitBoard(59)  // Add to d8
				pos.Black &^= IndexToBitBoard(56)
				pos.Black |= IndexToBitBoard(59)
			}
		}
	}

	// Update castling rights
	// King move removes both castling rights
	if m.Piece == King {
		if pos.WhiteMove {
			pos.CastleSide &^= CastleWhiteKingSide | CastleWhiteQueenSide
		} else {
			pos.CastleSide &^= CastleBlackKingSide | CastleBlackQueenSide
		}
	}
	// Rook move/capture removes the corresponding right
	if m.From == IndexToBitBoard(0) || m.To == IndexToBitBoard(0) {
		pos.CastleSide &^= CastleWhiteQueenSide
	}
	if m.From == IndexToBitBoard(7) || m.To == IndexToBitBoard(7) {
		pos.CastleSide &^= CastleWhiteKingSide
	}
	if m.From == IndexToBitBoard(56) || m.To == IndexToBitBoard(56) {
		pos.CastleSide &^= CastleBlackQueenSide
	}
	if m.From == IndexToBitBoard(63) || m.To == IndexToBitBoard(63) {
		pos.CastleSide &^= CastleBlackKingSide
	}

	// Update en passant square
	pos.EnPassant = 0
	if m.Piece == Pawn {
		// Check for double pawn push
		fromIdx := bitboardToIndex(m.From)
		toIdx := bitboardToIndex(m.To)
		diff := toIdx - fromIdx
		if diff == 16 { // White double push
			pos.EnPassant = IndexToBitBoard(fromIdx + 8)
		} else if diff == -16 { // Black double push
			pos.EnPassant = IndexToBitBoard(fromIdx - 8)
		}
	}

	// Update half-move clock
	if m.Piece == Pawn || m.Captured != Empty {
		pos.HalfMoveClock = 0
	} else {
		pos.HalfMoveClock++
	}

	// Switch side to move
	pos.WhiteMove = !pos.WhiteMove

	return undo
}

// UnmakeMove reverses a move using the saved undo information.
func (pos *Position) UnmakeMove(m Move, undo UndoInfo) {
	// Switch side back first (so we know whose move it was)
	pos.WhiteMove = !pos.WhiteMove

	// Get color bitboards
	var ourColor, enemyColor *Bitboard
	if pos.WhiteMove {
		ourColor = &pos.White
		enemyColor = &pos.Black
	} else {
		ourColor = &pos.Black
		enemyColor = &pos.White
	}

	// Reverse castling rook move first
	if m.Flags&FlagCastling != 0 {
		if pos.WhiteMove {
			if m.To == IndexToBitBoard(6) { // Kingside O-O
				pos.Rooks |= IndexToBitBoard(7)  // Restore to h1
				pos.Rooks &^= IndexToBitBoard(5) // Remove from f1
				pos.White |= IndexToBitBoard(7)
				pos.White &^= IndexToBitBoard(5)
			} else { // Queenside O-O-O
				pos.Rooks |= IndexToBitBoard(0)  // Restore to a1
				pos.Rooks &^= IndexToBitBoard(3) // Remove from d1
				pos.White |= IndexToBitBoard(0)
				pos.White &^= IndexToBitBoard(3)
			}
		} else {
			if m.To == IndexToBitBoard(62) { // Kingside O-O
				pos.Rooks |= IndexToBitBoard(63)  // Restore to h8
				pos.Rooks &^= IndexToBitBoard(61) // Remove from f8
				pos.Black |= IndexToBitBoard(63)
				pos.Black &^= IndexToBitBoard(61)
			} else { // Queenside O-O-O
				pos.Rooks |= IndexToBitBoard(56)  // Restore to a8
				pos.Rooks &^= IndexToBitBoard(59) // Remove from d8
				pos.Black |= IndexToBitBoard(56)
				pos.Black &^= IndexToBitBoard(59)
			}
		}
	}

	// Remove piece from destination (handle promotion)
	if m.Promotion != Empty {
		promoBB := pos.GetPiece(m.Promotion)
		*promoBB &^= m.To
	} else {
		pieceBB := pos.GetPiece(m.Piece)
		*pieceBB &^= m.To
	}
	*ourColor &^= m.To

	// Restore piece to source square (always original piece type)
	pieceBB := pos.GetPiece(m.Piece)
	*pieceBB |= m.From
	*ourColor |= m.From

	// Restore captured piece
	if undo.CapturedPiece != Empty {
		if m.Flags&FlagEnPassant != 0 {
			// En passant: restore pawn to its actual square
			var capturedPawnSq Bitboard
			if pos.WhiteMove {
				capturedPawnSq = m.To >> 8
			} else {
				capturedPawnSq = m.To << 8
			}
			pos.Pawns |= capturedPawnSq
			*enemyColor |= capturedPawnSq
		} else {
			// Normal capture
			capturedBB := pos.GetPiece(undo.CapturedPiece)
			*capturedBB |= m.To
			*enemyColor |= m.To
		}
	}

	// Restore game state
	pos.CastleSide = undo.CastleSide
	pos.EnPassant = undo.EnPassant
	pos.HalfMoveClock = undo.HalfMoveClock
}
