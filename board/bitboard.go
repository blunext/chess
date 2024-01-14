// Package board Layout 2: https://gekomad.github.io/Cinnamon/BitboardCalculator/
//
//	56	57	58	59	60	61	62	63
//	48	49	50	51	52	53	54	55
//	40	41	42	43	44	45	46	47
//	32	33	34	35	36	37	38	39
//	24	25	26	27	28	29	30	31
//	16	17	18	19	20	21	22	23
//	08	09	10	11	12	13	14	15
//	00	01	02	03	04	05	06	07

package board

import (
	"fmt"
)

type Bitboard uint64

func BB() {
	b := Bitboard(0x40201)
	b.SetBit(27)
	s := b.Pretty()
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

func (b *Bitboard) Pretty() string {
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

func Flat(boards []Bitboard) Bitboard {
	var flatten Bitboard
	for _, b := range boards {
		flatten = flatten | b
	}
	return flatten
}

// ToSlice takes a bitboard and returns a slice of bitboards
// where each bitboard has a single bit set
func (b *Bitboard) ToSlice() []Bitboard {
	// todo: consider nil slice
	slice := []Bitboard{}
	for i := 0; i < 64; i++ {
		mask := Bitboard(1 << i)
		if *b&mask != 0 {
			slice = append(slice, mask)
		}
	}
	return slice
}
