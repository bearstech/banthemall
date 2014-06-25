package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombined(t *testing.T) {
	p, err := NewCombinedParser()
	assert.Nil(t, err)
	c, ok := p.Parse(`5.39.32.126 - - [05/Jun/2014:21:49:52 +0200] "POST /mt/mt-tb.cgi/6 HTTP/1.1" 403 147 "http://jechercheunemeuf.info/" "PHP/5.2.66"`)
	assert.True(t, ok)
	assert.Equal(t, "5.39.32.126", c.ip)
	assert.Equal(t, "05/Jun/2014:21:49:52 +0200", c.timeStamp)
	assert.Equal(t, "POST", c.method)
	assert.Equal(t, "/mt/mt-tb.cgi/6", c.url)
	assert.Equal(t, "403", c.status)
	assert.Equal(t, "147", c.requestSize)
	assert.Equal(t, "http://jechercheunemeuf.info/", c.referer)
	assert.Equal(t, "PHP/5.2.66", c.browser)
}
