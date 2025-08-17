package utils

import (
	"fmt"
)

const (
	UnitKg   = "kg"
	UnitGram = "g"

	UnitLiter      = "l"
	UnitMilliliter = "ml"

	UnitPiece = "pcs"
)

var SliceUnits = []string{UnitKg, UnitGram, UnitLiter, UnitMilliliter, UnitPiece}

func ValidateProductUnit(unit string) error {
	for _, i := range SliceUnits {
		if i == unit {
			return nil
		}
	}

	return fmt.Errorf("unit '%s' is invalid", unit)
}
