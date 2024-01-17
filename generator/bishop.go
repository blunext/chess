package generator

import (
	"chess/board"
)

func generateBishopMoves() MovesArray {
	var squareMoves = make(MovesArray)
	for pos := 0; pos < 64; pos++ {
		var directions []Moves
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
		moves = bishopNE(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}

		squareMoves[board.Bitboard(pos)] = exactSize(directions)
	}
	return squareMoves
}

func bishopSE(pos int) Moves {
	var list Moves
	rank, file := rankAndFile(pos)
	for {
		rank--
		file++
		if rank < 0 || file >= 8 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((7-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopSW(pos int) Moves {
	var list Moves
	rank, file := rankAndFile(pos)
	for {
		rank--
		file--
		if rank < 0 || file < 0 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((7-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopNE(pos int) Moves {
	var list Moves
	rank, file := rankAndFile(pos)
	for {
		rank++
		file++
		if rank >= 8 || file >= 8 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((7-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}

func bishopNW(pos int) Moves {
	var list Moves
	rank, file := rankAndFile(pos)
	for {
		rank++
		file--
		if rank >= 8 || file < 0 {
			break
		}
		var newPos board.Bitboard
		newPos.SetBit((7-rank)*8 + file)
		list = append(list, newPos)
	}
	return exactSize(list)
}
