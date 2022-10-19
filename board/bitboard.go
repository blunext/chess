package board

import (
	"fmt"
)

type Bitboard uint64

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
