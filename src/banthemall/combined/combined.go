package combined

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

type Combined struct {
	IP          string
	TimeStamp   string
	Method      string
	URL         string
	Status      string
	RequestSize string
	Referer     string
	Browser     string
	time        *time.Time
}

func month(m string) (i time.Month, err error) {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug",
		"Sep", "Oct", "Nov", "Dec"}
	for i, month := range months {
		if m == month {
			return time.Month(i + 1), nil
		}
	}
	return -1, errors.New("Not a month : " + m)
}

func (c *Combined) Time() (t *time.Time, err error) {
	if c.time == nil {
		ts := c.TimeStamp
		// 05/Jun/2014:21:49:52 +0200
		d, err := strconv.Atoi(ts[0:2])
		if err != nil {
			return nil, err
		}
		m, err := month(ts[3:6])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(ts[7:11])
		if err != nil {
			return nil, err
		}
		h, err := strconv.Atoi(ts[12:14])
		if err != nil {
			return nil, err
		}
		min, err := strconv.Atoi(ts[15:17])
		if err != nil {
			return nil, err
		}
		s, err := strconv.Atoi(ts[18:20])
		if err != nil {
			return nil, err
		}
		l, err := strconv.Atoi(ts[22:24]) // hour part
		if err != nil {
			return nil, err
		}
		if ts[21] == '-' {
			l = -l
		}
		t := time.Date(y, m, d, h, min, s, 0, time.FixedZone(ts[21:25], l))
		c.time = &t
	}
	return c.time, nil
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
			IP:          mat[0][1],
			TimeStamp:   mat[0][2],
			Method:      mat[0][3],
			URL:         mat[0][4],
			Status:      mat[0][5],
			RequestSize: mat[0][6],
			Referer:     mat[0][7],
			Browser:     mat[0][8],
		}
		return ip, true
	}
	return Combined{}, false
}
