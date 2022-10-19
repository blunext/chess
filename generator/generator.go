package generator

import (
	"fmt"

	"chess/board"
)

type possibleMoves []board.Bitboard

type squareMoves map[board.Bitboard][]possibleMoves

type generatedMoves map[uint8]squareMoves

func NewGenerator() generatedMoves {
	moves := make(generatedMoves)
	rookGen(moves)
	fmt.Println(len(moves))
	return moves
}

func rookGen(generatedMoves generatedMoves) {
	var squareMoves = make(squareMoves)
	for pos := 0; pos < 64; pos++ {
		var directions []possibleMoves
		down := rookDown(pos)
		if len(down) != 0 {
			directions = append(directions, down)
		}
		up := rookUp(pos)
		if len(up) != 0 {
			directions = append(directions, up)
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

func exactSize(list possibleMoves) possibleMoves {
	exactSlise := make(possibleMoves, len(list))
	copy(exactSlise, list)
	return exactSlise
}
