package combined

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombined(t *testing.T) {
	p, err := NewCombinedParser()
	assert.Nil(t, err)
	c, ok := p.Parse(`5.39.32.126 - - [05/Jun/2014:21:49:52 +0200] "POST /mt/mt-tb.cgi/6 HTTP/1.1" 403 147 "http://jechercheunemeuf.info/" "PHP/5.2.66"`)
	assert.True(t, ok)
	assert.Equal(t, "5.39.32.126", c.IP)
	assert.Equal(t, "05/Jun/2014:21:49:52 +0200", c.TimeStamp)
	assert.Equal(t, "POST", c.Method)
	assert.Equal(t, "/mt/mt-tb.cgi/6", c.URL)
	assert.Equal(t, "403", c.Status)
	assert.Equal(t, "147", c.RequestSize)
	assert.Equal(t, "http://jechercheunemeuf.info/", c.Referer)
	assert.Equal(t, "PHP/5.2.66", c.Browser)
}
