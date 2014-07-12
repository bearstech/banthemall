package capped

import (
	"errors"
	"math"
	"sync"
	"time"
)

type CappedMetrics struct {
	zero   int
	values []float64
	mutex  sync.Mutex
	start  time.Time
	end    time.Time
	step   time.Duration
	size   int
}

func NewCappedMetrics(size int, step time.Duration) *CappedMetrics {
	n := time.Now().Truncate(step)
	start := n.Add(time.Duration(int64(-size) * int64(step)))
	c := &CappedMetrics{
		zero:   0,
		values: make([]float64, size),
		start:  start.Truncate(step),
		step:   step,
	}
	c.end = c.start
	for i := 0; i < size; i++ {
		c.values[i] = math.NaN()
	}
	return c
}

func (c *CappedMetrics) Update() {
	n := time.Now().Truncate(c.step)
	if n.Equal(c.end) {
		return
	}
	diff := n.Sub(c.end)
	for i := int64(0); i < int64(diff); i += int64(c.step) {
		c.AddNothing()
	}
}

func (c *CappedMetrics) Metric(t time.Time, value float64) error {
	n := time.Now()
	t = t.Truncate(c.step)
	if t.After(n) {
		return errors.New("No futur")
	}
	if t.Before(c.start) {
		// Too late
		return errors.New("Too old")
	}
	if t.After(c.end) { // oups, we are late
		c.Update()
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

	diff := t.Sub(c.end)
	poz := int(int64(diff)/int64(c.step)) + c.zero
	l := len(c.values)
	if poz > l {
		poz -= l
	}
	if poz < 0 {
		poz += l
	}
	c.values[poz] = value

	return nil
}

func (c *CappedMetrics) Add(value float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.values[c.zero] = value
	c.zero++
	c.end = c.end.Add(c.step)
	if c.size < len(c.values) {
		c.size++
	} else {
		c.start.Add(c.step)
	}
	if c.zero >= len(c.values) {
		c.zero = 0
	}
}

func (c *CappedMetrics) AddNothing() {
	c.Add(math.NaN())
}

type Values struct {
	Time  time.Time
	Value float64
}

func (c *CappedMetrics) modulo(n int) int {
	l := len(c.values)
	if n >= l {
		return c.modulo(n - l)
	}
	if n < 0 {
		return c.modulo(n + l)
	}
	return n
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
	j := 0
	for i := 0; i < len(c.values); i++ {
		v := c.values[c.modulo(c.zero-len(c.values)+1+i)]
		ts = ts.Add(c.step)
		if !math.IsNaN(v) {
			val := Values{ts, v}
			values[j] = &val
			j++
		}
	}
	return values
}
