package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nranchev/go-libGeoIP"
	"io"
	"os"
	"regexp"
	"sort"
	"time"
)

type combined struct {
	ip          string
	timeStamp   string
	method      string
	url         string
	status      string
	requestSize string
	referer     string
	browser     string
}

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

func consolidate(gi *libgeo.GeoIP, thresold int, count chan combined) {
	scores := make(map[string]map[int]int)
	agents := make(map[string]*Counter)
	urls := make(map[string]*Counter)
	long_scores := make(map[string]int)
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
			scores[ip][s] += 1
			long_scores[ip] += 1
			if _, ok := agents[ip]; !ok {
				agents[ip] = NewCounter()
			}
			agents[ip].Add(combi.browser)
			if _, ok := urls[ip]; !ok {
				urls[ip] = NewCounter()
			}
			urls[ip].Add(combi.url)
			total += 1
		case <-c:
			long += 1
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
			}
			fmt.Printf("\t%d hits from %d ip\n", total, len(scores))
			scores = make(map[string]map[int]int)
			agents = make(map[string]*Counter)
			urls = make(map[string]*Counter)
			total = 0
			if long == 60 {
				long = 0
				long_total := 0
				for ip, n := range long_scores {
					status := Rbl(ip)
					fmt.Printf("\tLong: %15s #%d %s\n", ip, n, status)
					long_total += n
				}
				fmt.Printf("\tLong total: %d\n\n", long_total)
				long_scores = make(map[string]int)
			}
		}
	}
}

func main() {
	var flagThresold int

	flag.IntVar(&flagThresold, "thresold", 0, "Minimum mumber of hits per 10 seconds")

	flag.Parse()

	gi, err := libgeo.Load("GeoIP.dat")
	if err != nil {
		fmt.Printf("Error Libgeo: %s\n", err.Error())
		return
	}

	/*
	   5.39.32.126 - - [05/Jun/2014:21:49:52 +0200] "POST /mt/mt-tb.cgi/6 HTTP/1.1" 403 147 "http://jechercheunemeuf.info/" "PHP/5.2.66"
	*/
	r := "^(?P<RemoteIP>\\S+) \\S+ \\S+ \\[(?P<Timestamp>[^\\]]+)\\] \"(?P<Method>[A-Z]+) (?P<Url>[^\\s]+)[^\"]*\" (?P<StatusCode>\\d+) (?P<RequestSize>\\d+|-) \"(?P<Referer>[^\"]*)\" \"(?P<Browser>[^\"]*)\""
	apachelog, err := regexp.Compile(r)
	if err != nil {
		fmt.Printf("Error Regexp: %s\n", err.Error())
		return
	}

	bio := bufio.NewReader(os.Stdin)
	count := make(chan combined)
	go consolidate(gi, flagThresold, count)
	for {
		line, err := bio.ReadString('\n')
		if err == io.EOF {
			continue
		}
		if err != nil {
			fmt.Printf("Error Regexp: %s\n", err.Error())
			continue
		}
		mat := apachelog.FindAllStringSubmatch(line, -1)
		if len(mat) > 0 && len(mat[0]) > 0 {
			ip := combined{
				mat[0][1],
				mat[0][2],
				mat[0][3],
				mat[0][4],
				mat[0][5],
				mat[0][6],
				mat[0][7],
				mat[0][8]}
			count <- ip
		}
	}
}
