package value_objects

import (
	"fmt"
)

const (
	HappyLevelMin = 1
	HappyLevelMax = 10
)

type HappyLevel struct {
	value int
}

func NewHappyLevel(value int) (*HappyLevel, error) {
	if value < HappyLevelMin || value > HappyLevelMax {
		return nil, fmt.Errorf("happy level must be between %d and %d", HappyLevelMin, HappyLevelMax)
	}

	return &HappyLevel{value: value}, nil
}

func (h *HappyLevel) Value() int {
	return h.value
}

func (h *HappyLevel) IsEmpty() bool {
	return h.value == 0
}
