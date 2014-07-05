package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	c := NewCounter()
	c.Add("pim")
	c.Add("pam")
	c.Add("poum")
	c.Add("pim")
	assert.Equal(t, 3, c.Size())
}
