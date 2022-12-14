package generator

import "chess/board"

func generateKnightMoves() knightMoves {
	knightSteps := []struct {
		step []int
	}{
		{[]int{2, -1}},
		{[]int{1, -2}},

		{[]int{-2, -1}},
		{[]int{-1, -2}},

		{[]int{-2, 1}},
		{[]int{-1, 2}},

		{[]int{2, 1}},
		{[]int{1, 2}},
	}
	knightMoves := make(knightMoves)

	for pos := 0; pos < 64; pos++ {
		var list possibleMoves
		for _, knight := range knightSteps {
			rank, file := rankAndFile(pos)
			file += knight.step[0]
			rank += knight.step[1]
			if rank > 0 && rank < 8 && file > 0 && file < 8 {
				var newPos board.Bitboard
				n := (7-rank)*8 + file
				newPos.SetBit(n)
				list = append(list, newPos)
			}
		}
		if len(list) > 0 {
			knightMoves[board.Bitboard(pos)] = exactSize(list)
		}
	}
	return knightMoves
}
