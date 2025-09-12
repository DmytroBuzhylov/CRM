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

var SliceUnits = []string{UnitGram, UnitMilliliter, UnitPiece}

func ValidateProductUnit(unit string) error {
	for _, i := range SliceUnits {
		if i == unit {
			return nil
		}
	}

	return fmt.Errorf("unit '%s' is invalid", unit)
}

func NormalizeUnit(value float64, unit string) (uint64, string, error) {
	switch unit {
	case UnitKg:
		return uint64(value * 1000), UnitGram, nil
	case UnitGram:
		return uint64(value), UnitGram, nil
	case UnitLiter:
		return uint64(value * 1000), UnitMilliliter, nil
	case UnitMilliliter:
		return uint64(value), UnitMilliliter, nil
	case UnitPiece:
		return uint64(value), UnitPiece, nil
	default:
		return 0, "", fmt.Errorf("unsupported unit: %s", unit)
	}
}
