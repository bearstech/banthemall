package main

import (
	"regexp"
)

type Combined struct {
	ip          string
	timeStamp   string
	method      string
	url         string
	status      string
	requestSize string
	referer     string
	browser     string
}
type CombinedParser struct {
	re *regexp.Regexp
}

func NewCombinedParser() (*CombinedParser, error) {
	c := new(CombinedParser)

	/*
	   5.39.32.126 - - [05/Jun/2014:21:49:52 +0200] "POST /mt/mt-tb.cgi/6 HTTP/1.1" 403 147 "http://jechercheunemeuf.info/" "PHP/5.2.66"
	*/
	r := "^(?P<RemoteIP>\\S+) \\S+ \\S+ \\[(?P<Timestamp>[^\\]]+)\\] \"(?P<Method>[A-Z]+) (?P<Url>[^\\s]+)[^\"]*\" (?P<StatusCode>\\d+) (?P<RequestSize>\\d+|-) \"(?P<Referer>[^\"]*)\" \"(?P<Browser>[^\"]*)\""
	apachelog, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}
	c.re = apachelog
	return c, nil
}

func (c CombinedParser) Parse(line string) (Combined, bool) {
	mat := c.re.FindAllStringSubmatch(line, -1)
	if len(mat) > 0 && len(mat[0]) > 0 {
		ip := Combined{
			mat[0][1],
			mat[0][2],
			mat[0][3],
			mat[0][4],
			mat[0][5],
			mat[0][6],
			mat[0][7],
			mat[0][8]}
		return ip, true
	} else {
		return Combined{}, false
	}
}
