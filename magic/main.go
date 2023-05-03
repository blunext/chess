package magic

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"fmt"

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

func Prepare() error {
	dec := gob.NewDecoder(bytes.NewReader(magicData))

	err := dec.Decode(&RookMagics)
	if err != nil {
		return fmt.Errorf("decode error RookMagics: %w", err)
	}
	err = dec.Decode(&RookMoves)
	if err != nil {
		return fmt.Errorf("decode error RookMoves: %w", err)
	}
	err = dec.Decode(&BishopMagics)
	if err != nil {
		return fmt.Errorf("decode error BishopMagics: %w", err)
	}
	err = dec.Decode(&BishopMoves)
	if err != nil {
		return fmt.Errorf("decode error BishopMoves: %w", err)
	}
	return nil
}
