package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/nranchev/go-libGeoIP"
)

func status(statuscode string) (r int) {
	f := statuscode[0]
	if f == '1' || f == '2' || f == '3' {
		return 0
	}
	if f == '4' {
		return 1
	}
	return 2
}

type user struct {
	ip    string
	score int
}

type byscore []user

func (b byscore) Len() int { return len(b) }

func (b byscore) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b byscore) Less(i, j int) bool { return b[i].score > b[j].score }

func consolidate(gi *libgeo.GeoIP, thresold int, carbon *Carbon, count chan Combined) {
	scores := make(map[string]map[int]int)
	agents := make(map[string]*Counter)
	urls := make(map[string]*Counter)
	longScores := make(map[string]int)
	total := 0
	long := 0
	c := time.Tick(10 * time.Second)
	var cc string
	for {
		select {
		case combi := <-count:
			ip := combi.ip
			s := status(combi.status)
			if _, ok := scores[ip]; !ok {
				scores[ip] = make(map[int]int)
			}
			scores[ip][s]++
			longScores[ip]++
			if _, ok := agents[ip]; !ok {
				agents[ip] = NewCounter()
			}
			agents[ip].Add(combi.browser)
			if _, ok := urls[ip]; !ok {
				urls[ip] = NewCounter()
			}
			urls[ip].Add(combi.url)
			total++
		case <-c:
			long++
			ss := []user{}
			for ip, sco := range scores {
				ss = append(ss, user{ip, sco[0] + sco[1] + sco[2]})
			}
			sort.Sort(byscore(ss))
			for _, s := range ss {
				ip := s.ip
				sco := scores[ip]
				loc := gi.GetLocationByIP(ip)
				status := Rbl(ip)
				if status == "-" {
					carbon.Sum("banthemall.spamhaus.RAS", 1)
				} else {
					carbon.Sum("banthemall.spamhaus."+status, 1)
				}
				if loc == nil {
					cc = "??"
				} else {
					cc = loc.CountryCode
				}
				r23 := sco[0]
				r4 := sco[1]
				r5 := sco[2]
				r := r23 + r4 + r5
				if r >= thresold {
					fmt.Printf("%s %15s [23]xx: %4d 4xx: %4d 5xx: %4d #%4d #ua: %4d #url: %4d %s\n",
						cc, ip, r23, r4, r5, r, agents[ip].Size(),
						urls[ip].Size(), status)
				}
				carbon.Max("banthemall.hit-per-ip.max", r)
			}
			carbon.Max("banthemall.distinct-ip.max", len(scores))
			fmt.Printf("\t%d hits from %d ip\n", total, len(scores))
			scores = make(map[string]map[int]int)
			agents = make(map[string]*Counter)
			urls = make(map[string]*Counter)
			total = 0
			if long == 60 {
				long = 0
				longTotal := 0
				ss := []user{}
				for ip, n := range longScores {
					ss = append(ss, user{ip, n})
				}
				sort.Sort(byscore(ss))
				for _, s := range ss {
					ip := s.ip
					n := longScores[ip]
					status := Rbl(ip)
					fmt.Printf("\tLong: %15s #%d %s\n", ip, n, status)
					longTotal += n
					carbon.Max("banthemall.long.hit-per-ip.max", n)
				}
				fmt.Printf("\tLong total: %d\n\n", longTotal)
				longScores = make(map[string]int)
			}
		}
	}
}
