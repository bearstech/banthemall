package main

import (
	"banthemall/combined"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsolidate(t *testing.T) {
	s := NewShortTerm()
	l := NewLongTerm()
	c := &combined.Combined{
		IP:          "127.0.0.1",
		TimeStamp:   "05/Jun/2014:21:49:52 +0200",
		Method:      "POST",
		URL:         "/mt/mt-tb.cgi/6",
		Status:      "403",
		RequestSize: "147",
		Referer:     "http://jechercheunemeuf.info/",
		Browser:     "PHP/5.2.66",
	}
	s.Add(c)
	l.Add(c)
	c2 := &combined.Combined{
		IP:          "127.0.0.1",
		TimeStamp:   "05/Jun/2014:21:49:53 +0200",
		Method:      "POST",
		URL:         "/mt/mt-tb.cgi/6",
		Status:      "403",
		RequestSize: "147",
		Referer:     "http://jechercheunemeuf.info/",
		Browser:     "PHP/5.2.66",
	}
	s.Add(c2)
	l.Add(c2)
	ips := s.IPs()
	assert.Equal(t, 1, len(ips))
	assert.Equal(t, 2, ips[0].hits4)
	assert.Equal(t, 2, l.total)
	s.Consolidate(nil, 0, nil)
	l.Consolidate(nil)
}
