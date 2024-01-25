package generator

import (
	"chess/board"
)

type steps []struct{ step []int }

func kingMoves() board.SquareMoves {
	knightSteps := steps{
		{[]int{-1, -1}},
		{[]int{-1, 0}},
		{[]int{-1, 1}},
		{[]int{0, -1}},
		{[]int{0, 1}},
		{[]int{1, -1}},
		{[]int{1, 0}},
		{[]int{1, 1}},
	}
	return generateGenericMoves(knightSteps)
}

func knightMoves() board.SquareMoves {
	knightSteps := steps{
		{[]int{2, -1}},
		{[]int{1, -2}},
		{[]int{-2, -1}},
		{[]int{-1, -2}},
		{[]int{-2, 1}},
		{[]int{-1, 2}},
		{[]int{2, 1}},
		{[]int{1, 2}},
	}
	return generateGenericMoves(knightSteps)
}
func generateGenericMoves(steps steps) board.SquareMoves {
	knightMoves := make(board.SquareMoves)

	for pos := 0; pos < 64; pos++ {
		var directions [][]board.Bitboard
		var list []board.Bitboard
		for _, knight := range steps {
			rank, file := rankAndFile(pos)
			file += knight.step[0]
			rank += knight.step[1]
			if rank >= 0 && rank < 8 && file >= 0 && file < 8 {
				var newPos board.Bitboard
				n := (7-rank)*8 + file
				newPos.SetBit(n)
				list = append(list, newPos)
			}
		}
		directions = append(directions, exactSize(list))
		knightMoves[board.IndexToBitBoard(pos)] = exactSize(directions)
	}
	return knightMoves
}
