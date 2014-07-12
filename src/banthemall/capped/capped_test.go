package capped

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCapped(t *testing.T) {
	c := NewCappedMetrics(3, time.Minute)
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
	assert.Equal(t, 3, values[0].Value)
	assert.Equal(t, 2, values[1].Value)
}

func TestMetric(t *testing.T) {
	n := time.Now()
	c := NewCappedMetrics(3, time.Minute)
	err := c.Metric(n.Add(-1*time.Minute), 37)
	assert.Equal(t, nil, err)
	err = c.Metric(n, 5)
	assert.Equal(t, nil, err)
	values := c.Values()
	assert.Equal(t, 2, len(values))
	assert.True(t, values[1].Time.After(values[0].Time))
	assert.Equal(t, 37, values[0].Value)
	assert.Equal(t, 5, values[1].Value)
	nn := n.Truncate(time.Minute)
	assert.Equal(t, nn.Add(-1*time.Minute), values[0].Time)
	assert.Equal(t, nn, values[1].Time)
}
