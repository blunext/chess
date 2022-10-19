package generator

import (
	"chess/board"
)

func generateRooks(generatedMoves generatedMoves) {
	var squareMoves = make(squareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions []possibleMoves
		moves := rookDown(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = rookUp(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = rookRight(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = rookLeft(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}

		squareMoves[board.Bitboard(pos)] = directions
	}

	generatedMoves[board.Rook] = squareMoves
}

func rookDown(pos int) possibleMoves {
	var list possibleMoves
	rank := 8 - pos>>3
	for i := 1; i < rank; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos + i*8)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookUp(pos int) possibleMoves {
	var list possibleMoves
	rank := 8 - pos>>3
	for i := 1; i <= 8-rank; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos - i*8)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookRight(pos int) possibleMoves {
	var list possibleMoves
	file := pos&7 + 1
	for i := 1; i <= 8-file; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos + i)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookLeft(pos int) possibleMoves {
	var list possibleMoves
	file := pos&7 + 1
	for i := 1; i < file; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos - i)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func exactSize(list possibleMoves) possibleMoves {
	exactSlise := make(possibleMoves, len(list))
	copy(exactSlise, list)
	return exactSlise
}
