package metrics

import (
	"hash/fnv"
)

/*
Counter tha count distinct stuff
*/
type Counter struct {
	data map[uint32]bool
}

/*
NewCounter return a new Counter.
*/
func NewCounter() *Counter {
	c := new(Counter)
	c.data = make(map[uint32]bool)
	return c
}

/*
Add something to count
*/
func (c Counter) Add(element string) {
	hash := fnv.New32()
	hash.Write([]byte(element))
	h := hash.Sum32()
	c.data[h] = true
}

/*
Size of the set
*/
func (c Counter) Size() (size int) {
	return len(c.data)
}
