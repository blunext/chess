package board

import (
	"fmt"
)

type Bitboard uint64

func BB() {
	b := uint64(0x40201)
	b = SetBit(b, 27)
	s := pretty(b)
	fmt.Println(s)
}

func bit(b uint64, index int) uint64 {
	mask := uint64(1) << index
	return (b & mask) >> index
}

func IsBitSet(b uint64, index int) bool {
	return bit(b, index) == 1
}

func SetBit(b uint64, index int) uint64 {
	return b | (1 << index)
}

func OneBit(index int) uint64 {
	return 1 << index
}

func squareIndex(f, r int) int {
	return (r << 3) + f
}

func pretty(b uint64) string {
	s := "+---+---+---+---+---+---+---+---+\n"
	for r := RANK_8; r >= RANK_1; r-- {
		for f := FILE_A; f <= FILE_H; f++ {
			if IsBitSet(b, squareIndex(f, r)) {
				s += "| X "
			} else {
				s += "|   "
			}
		}
		s += fmt.Sprintf("| %d\n+---+---+---+---+---+---+---+---+\n", r+1)
	}
	s += "  a   b   c   d   e   f   g   h\n"
	return s
}
