package main

import (
	"banthemall/combined"
	"fmt"
	"sort"
	"time"

	"github.com/nranchev/go-libGeoIP"
)

type IP struct {
	ip      string
	agents  *Counter
	urls    *Counter
	hits123 int
	hits4   int
	hits5   int
}

func NewIP(ip string) *IP {
	return &IP{
		ip:     ip,
		agents: NewCounter(),
		urls:   NewCounter(),
	}
}

func (ip *IP) Add(statusCode string) {
	f := statusCode[0]
	if f == '1' || f == '2' || f == '3' {
		ip.hits123++
	}
	if f == '4' {
		ip.hits4++
	}
	ip.hits5++
}

func (ip *IP) Hits() int {
	return ip.hits123 + ip.hits4 + ip.hits5
}

type ShortTerm struct {
	bagIP map[string]*IP
	total int
	ips   int
}

func NewShortTerm() *ShortTerm {
	return &ShortTerm{
		bagIP: make(map[string]*IP),
	}
}

func (s *ShortTerm) Add(combi *combined.Combined) {
	s.total++
	ip := combi.IP
	if _, ok := s.bagIP[ip]; !ok {
		s.bagIP[ip] = NewIP(ip)
		s.ips++
	}
	s.bagIP[ip].Add(combi.Status)
	s.bagIP[ip].agents.Add(combi.Browser)
	s.bagIP[ip].urls.Add(combi.URL)
}

func (s *ShortTerm) IPs() []*IP {
	size := len(s.bagIP)
	ss := make([]user, size, size)
	i := 0
	for ip, obj := range s.bagIP {
		ss[i] = user{ip, obj.Hits()}
		i++
	}
	sort.Sort(byscore(ss))
	r := make([]*IP, size, size)
	for i, ip := range ss {
		r[i] = s.bagIP[ip.ip]
	}
	return r
}

func (s *ShortTerm) Size() int {
	return len(s.bagIP)
}

func (s *ShortTerm) Consolidate(gi *libgeo.GeoIP, thresold int, carbon *Carbon) {
	for _, i := range s.IPs() {
		ip := i.ip
		var loc *libgeo.Location
		if gi != nil {
			loc = gi.GetLocationByIP(ip)
		} else {
			loc = nil
		}
		status := Rbl(ip)
		if carbon != nil {
			if status == "-" {
				carbon.Sum("banthemall.spamhaus.RAS", 1)
			} else {
				carbon.Sum("banthemall.spamhaus."+status, 1)
			}
		}
		var cc string
		if loc == nil {
			cc = "??"
		} else {
			cc = loc.CountryCode
		}
		if i.Hits() >= thresold {
			fmt.Printf("%s %15s [23]xx: %4d 4xx: %4d 5xx: %4d #%4d #ua: %4d #url: %4d %s\n",
				cc, ip, i.hits123, i.hits4, i.hits5, i.Hits(), i.agents.Size(),
				i.urls.Size(), status)
		}
		if carbon != nil {
			carbon.Max("banthemall.hit-per-ip.max", i.Hits())
			carbon.List("banthemall.hit-per-ip.percentile", i.Hits())
		}
	}
	if carbon != nil {
		carbon.Max("banthemall.distinct-ip.max", s.ips)
	}
	fmt.Printf("\t%d hits from %d ip\n", s.total, s.ips)
}

type LongTerm struct {
	hits  map[string]int
	total int
}

func NewLongTerm() *LongTerm {
	return &LongTerm{
		hits: make(map[string]int),
	}
}

func (l *LongTerm) Add(combi *combined.Combined) {
	l.total++
	ip := combi.IP
	if _, ok := l.hits[ip]; !ok {
		l.hits[ip] = 1
	} else {
		l.hits[ip]++
	}
}

func (l *LongTerm) Size() int {
	return len(l.hits)
}

func (l *LongTerm) Users() []user {
	s := l.Size()
	users := make(byscore, s, s)
	i := 0
	for ip, n := range l.hits {
		users[i] = user{ip, n}
		i++
	}
	sort.Sort(users)
	return users
}

func (l *LongTerm) Consolidate(carbon *Carbon) {
	for _, user := range l.Users() {
		ip := user.ip
		status := Rbl(ip)
		fmt.Printf("\tLong: %15s #%d %s\n", ip, user.score, status)
		if carbon != nil {
			carbon.Max("banthemall.long.hit-per-ip.max", user.score)
		}
	}
	fmt.Printf("\tLong total: %d\n\n", l.total)
}

type user struct {
	ip    string
	score int
}

type byscore []user

func (b byscore) Len() int { return len(b) }

func (b byscore) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b byscore) Less(i, j int) bool { return b[i].score > b[j].score }

/*
Infinite loop feed with a chan.
*/
func consolidate(gi *libgeo.GeoIP, thresold int, carbon *Carbon, count chan combined.Combined) {
	shortTerm := NewShortTerm()
	longTerm := NewLongTerm()
	long := 0
	c := time.Tick(10 * time.Second)
	for {
		select {
		case combi := <-count:
			shortTerm.Add(&combi)
			longTerm.Add(&combi)
		case <-c:
			shortTerm.Consolidate(gi, thresold, carbon)
			shortTerm = NewShortTerm()
			long++
			if long == 60 {
				longTerm.Consolidate(carbon)
				longTerm = NewLongTerm()
				long = 0
			}
		}
	}
}
