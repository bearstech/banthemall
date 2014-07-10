package capped

import (
	"math"
	"sync"
	"time"
)

type CappedMetrics struct {
	zero   int
	values []float64
	mutex  sync.Mutex
	start  time.Time
	step   time.Duration
}

func NewCappedMetrics(size int, start time.Time, step time.Duration) *CappedMetrics {
	c := &CappedMetrics{
		zero:   0,
		values: make([]float64, size),
		start:  start,
		step:   step,
	}
	for i := 0; i < size; i++ {
		c.values[i] = math.NaN()
	}
	return c
}

func (c *CappedMetrics) Add(value float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.values[c.zero] = value
	c.zero++
	if c.zero >= len(c.values) {
		c.zero = 0
	}
	c.start.Add(c.step)
}

func (c *CappedMetrics) AddNothing() {
	c.Add(math.NaN())
}

type Values struct {
	ts    time.Time
	value float64
}

func (c *CappedMetrics) Values() []*Values {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	empty := 0
	for _, v := range c.values {
		if math.IsNaN(v) {
			empty++
		}
	}
	if len(c.values) == empty {
		return nil
	}
	values := make([]*Values, len(c.values)-empty)
	ts := c.start
	i := 0
	for _, v := range c.values {
		if !math.IsNaN(v) {
			val := Values{ts, v}
			values[i] = &val
			i++
		}
		ts.Add(c.step)
	}
	return values
}
