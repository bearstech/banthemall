package capped

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCapped(t *testing.T) {
	c := NewCappedMetrics(3, time.Now(), time.Minute)
	var i float64
	for i = 0; i < 4; i++ {
		c.Add(i)
	}
	assert.Equal(t, []float64{3, 1, 2}, c.values)
	assert.Equal(t, 1, c.zero)
	c.AddNothing()
	assert.True(t, math.IsNaN(c.values[1]))
	values := c.Values()
	assert.Equal(t, 2, len(values))
	assert.Equal(t, 3, values[0].value)
	assert.Equal(t, 2, values[1].value)
}
