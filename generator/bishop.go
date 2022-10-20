package generator

import (
	"chess/board"
)

func generateBishopMoves() squareMoves {
	var squareMoves = make(squareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions []possibleMoves
		moves := bishopSE(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = bishopSW(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = bishopNW(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		moves = bishopSE(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}

		squareMoves[board.Bitboard(pos)] = directions
	}
	return squareMoves
}

func bishopSE(pos int) possibleMoves {
	var list possibleMoves
	rank := 8 - pos>>3
	file := pos & 7
	for {
		rank -= 1
		file += 1
		if rank <= 0 || file >= 8 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((8-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopSW(pos int) possibleMoves {
	var list possibleMoves
	rank := 8 - pos>>3
	file := pos & 7
	for {
		rank -= 1
		file -= 1
		if rank <= 0 || file < 0 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((8-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopNE(pos int) possibleMoves {
	var list possibleMoves
	rank := 8 - pos>>3
	file := pos & 7
	for {
		rank += 1
		file += 1
		if rank > 8 || file >= 8 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((8-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopNW(pos int) possibleMoves {
	var list possibleMoves
	rank := 8 - pos>>3
	file := pos & 7
	for {
		rank += 1
		file -= 1
		if rank > 8 || file < 0 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((8-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}
