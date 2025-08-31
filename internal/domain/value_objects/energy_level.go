package value_objects

import "fmt"

const (
	EnergyLevelMin = 1
	EnergyLevelMax = 10
)

type EnergyLevel struct {
	value int
}

func NewEnergyLevel(value int) (*EnergyLevel, error) {
	if value < EnergyLevelMin || value > EnergyLevelMax {
		return nil, fmt.Errorf("energy level must be between %d and %d", EnergyLevelMin, EnergyLevelMax)
	}

	return &EnergyLevel{value: value}, nil
}

func (e *EnergyLevel) Value() int {
	return e.value
}

func (e *EnergyLevel) IsEmpty() bool {
	return e.value == 0
}
