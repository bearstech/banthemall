package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsolidate(t *testing.T) {
	s := NewShortTerm()
	l := NewLongTerm()
	c := &Combined{
		"127.0.0.1",
		"05/Jun/2014:21:49:52 +0200",
		"POST",
		"/mt/mt-tb.cgi/6",
		"403",
		"147",
		"http://jechercheunemeuf.info/",
		"PHP/5.2.66",
	}
	s.Add(c)
	l.Add(c)
	c2 := &Combined{
		"127.0.0.1",
		"05/Jun/2014:21:49:53 +0200",
		"POST",
		"/mt/mt-tb.cgi/6",
		"403",
		"147",
		"http://jechercheunemeuf.info/",
		"PHP/5.2.66",
	}
	s.Add(c2)
	l.Add(c2)
	ips := s.IPs()
	assert.Equal(t, 1, len(ips))
	assert.Equal(t, 2, ips[0].hits4)
	assert.Equal(t, 2, l.total)
	shortConsolidation(s, nil, 0, nil)
}
