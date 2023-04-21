package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculations(t *testing.T) {
	rng.Seed(0)
	avg1, err := calculateByReferences()
	assert.Nil(t, err)

	rng.Seed(0)
	avg2, err := calculatePacked()
	assert.Nil(t, err)
	assert.EqualValues(t, avg1, avg2)

	rng.Seed(0)
	avg3, err := calculatePackedParallel()
	assert.Nil(t, err)
	assert.EqualValues(t, avg1, avg3)
}
