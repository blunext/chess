package generator

import (
	"chess/board"
)

func generateBishop(generatedMoves generatedMoves) {
	var squareMoves = make(squareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions []possibleMoves
		moves := bishopSE(pos)
		if len(moves) != 0 {
			directions = append(directions, moves)
		}
		// moves = rookUp(pos)
		// if len(moves) != 0 {
		// 	directions = append(directions, moves)
		// }
		// moves = rookRight(pos)
		// if len(moves) != 0 {
		// 	directions = append(directions, moves)
		// }
		// moves = rookLeft(pos)
		// if len(moves) != 0 {
		// 	directions = append(directions, moves)
		// }

		squareMoves[board.Bitboard(pos)] = directions
	}

	generatedMoves[board.Bishop] = squareMoves
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
		a := (8-rank)*8 + file
		newPos.SetBit(a)
		list = append(list, newPos)
	}
	return exactSize(list)
}
