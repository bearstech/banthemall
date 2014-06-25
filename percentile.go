package main

import "sort"

type Percentile struct {
	values *[]int
}

func NewPercentile() *Percentile {
	p := new(Percentile)
	p.values = &[]int{}
	return p
}

func (p Percentile) Append(i int) {
	slice := *p.values
	slice = append(slice, i)
	*p.values = slice
}

func (p Percentile) Percentile(per int) int {
	slice := *p.values
	sort.Ints(slice)
	return slice[len(slice)*per/100]
}
