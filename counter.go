package main

import (
	"hash/fnv"
)

type Counter struct {
	data map[uint32]bool
}

func NewCounter() *Counter {
	c := new(Counter)
	c.data = make(map[uint32]bool)
	return c
}

func (c Counter) Add(element string) {
	hash := fnv.New32()
	hash.Write([]byte(element))
	h := hash.Sum32()
	c.data[h] = true
}

func (c Counter) Size() (size int) {
	return len(c.data)
}
