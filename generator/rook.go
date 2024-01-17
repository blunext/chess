package generator

import (
	"chess/board"
)

func generateRookMoves() MovesArray {
	var squareMoves = make(MovesArray)
	for pos := 0; pos < 64; pos++ {
		var directions []Moves
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

		squareMoves[board.Bitboard(pos)] = exactSize(directions)
	}
	return squareMoves
}

func rookDown(pos int) Moves {
	var list Moves
	rank := rank(pos)
	for i := 1; i <= rank; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos + i*8)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookUp(pos int) Moves {
	var list Moves
	rank := rank(pos)
	for i := 1; i < 8-rank; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos - i*8)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookRight(pos int) Moves {
	var list Moves
	file := file(pos) + 1
	for i := 1; i <= 8-file; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos + i)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func rookLeft(pos int) Moves {
	var list Moves
	file := file(pos) + 1
	for i := 1; i < file; i++ {
		var newPos board.Bitboard
		newPos.SetBit(pos - i)
		list = append(list, newPos)
	}
	return exactSize(list)
}
