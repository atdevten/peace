package value_objects

import (
	"fmt"
	"strconv"
)

type QuoteID struct {
	value int
}

func NewQuoteID() *QuoteID {
	return &QuoteID{
		value: 0,
	}
}

func NewQuoteIDFromInt(value int) *QuoteID {
	return &QuoteID{
		value: value,
	}
}

func NewQuoteIDFromString(value string) (*QuoteID, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid quote id: %w", err)
	}

	return NewQuoteIDFromInt(intValue), nil
}

func (q *QuoteID) Value() int {
	return q.value
}

func (q *QuoteID) String() string {
	return strconv.Itoa(q.value)
}
