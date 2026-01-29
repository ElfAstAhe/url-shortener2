package utils

import "fmt"

type PoolOvercrowdedError struct {
	capacity int
}

var ErrPoolOvercrowded *PoolOvercrowdedError

func NewPoolOvercrowded(capacity int) *PoolOvercrowdedError {
	return &PoolOvercrowdedError{
		capacity: capacity,
	}
}

func (p *PoolOvercrowdedError) Error() string {
	if p.capacity > 0 {
		return fmt.Sprintf("pool overcrowded, capacity [%d]", p.capacity)
	}
	return "pool overcrowded"
}
