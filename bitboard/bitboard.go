package bitboard

import (
	"fmt"
)

type bitboard uint64

func (b *bitboard) Print() {
	fmt.Println("")
	for i := 0; i < 64; i++ {
		sq := 0
		if b.isBitSet(i) {
			sq = 1
		}
		fmt.Print(sq)
		if ((i + 1) % 8) == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func (b *bitboard) bit(index int) uint64 {

	mask := uint64(1) << index
	return (uint64(*b) & mask) >> index
}

func (b *bitboard) isBitSet(index int) bool {
	return b.bit(index) == 1
}

func (b *bitboard) setBit(index int) {
	*b |= 1 << index
}
