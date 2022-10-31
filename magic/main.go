package magic

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"log"

	"chess/board"
)

//go:embed magicData
var magicData []byte

type Magic struct {
	Number uint64
	Mask   board.Bitboard
	Shift  uint8
}

var RookMagics [64]Magic
var BishopMagics [64]Magic
var RookMoves [64][4096]board.Bitboard
var BishopMoves [64][512]board.Bitboard

func Prepare() {
	dec := gob.NewDecoder(bytes.NewReader(magicData))

	err := dec.Decode(&RookMagics)
	if err != nil {
		log.Fatal("decode error RookMagics:", err)
	}
	err = dec.Decode(&RookMoves)
	if err != nil {
		log.Fatal("decode error RookMoves:", err)
	}
	err = dec.Decode(&BishopMagics)
	if err != nil {
		log.Fatal("decode error BishopMagics:", err)
	}
	err = dec.Decode(&BishopMoves)
	if err != nil {
		log.Fatal("decode error BishopMoves:", err)
	}
}
