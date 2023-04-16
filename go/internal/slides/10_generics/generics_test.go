package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeSumAndDiffFloat(t *testing.T) {
	op := NewOperation(2.0, 1.0)
	sum, diff := ComputeSumAndDiff[float64](op)
	assert.Equal(t, 3.0, sum)
	assert.Equal(t, 1.0, diff)
}
