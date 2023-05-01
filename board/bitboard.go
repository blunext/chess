package board

import (
	"fmt"
)

type Bitboard uint64

func BB() {
	b := Bitboard(0x40201)
	b.SetBit(27)
	s := pretty(b)
	fmt.Println(s)
}

func (b *Bitboard) Print() {
	fmt.Println("")
	for i := 0; i < 64; i++ {
		sq := 0
		if b.IsBitSet(i) {
			sq = 1
		}
		fmt.Print(sq)
		if ((i + 1) % 8) == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func (b *Bitboard) bit(index int) uint64 {
	mask := uint64(1) << index
	return (uint64(*b) & mask) >> index
}

func (b *Bitboard) IsBitSet(index int) bool {
	return b.bit(index) == 1
}

func (b *Bitboard) SetBit(index int) {
	*b |= 1 << index
}

func squareIndex(f, r int) int {
	return (r << 3) + f
}

func pretty(b Bitboard) string {
	s := "+---+---+---+---+---+---+---+---+\n"
	for r := Rank8; r >= Rank1; r-- {
		for f := FileA; f <= FileH; f++ {
			if b.IsBitSet(squareIndex(f, r)) {
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
