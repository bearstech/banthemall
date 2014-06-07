package main

import (
	"bufio"
	"fmt"
	"github.com/nranchev/go-libGeoIP"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
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

func rbl(ip string) (status string) {
	blocks := strings.Split(ip, ".")
	name := fmt.Sprintf("%s.%s.%s.%s.zen.spamhaus.org", blocks[3], blocks[2], blocks[1], blocks[0])
	r, err := net.LookupHost(name)
	if err != nil {
		//Optimistic answser
		return "-"
	}
	if r[0] == "127.0.0.10" || r[0] == "127.0.0.11" {
		return "PBL"
	}
	if r[0] == "127.0.0.2" {
		return "SBL"
	}
	if r[0] == "127.0.0.3" {
		return "CSS"
	}
	if r[0] == "127.0.0.4" || r[0] == "127.0.0.5" || r[0] == "127.0.0.6" || r[0] == "127.0.0.7" {
		return "XBL"
	}
	return r[0]
}

func consolidate(gi *libgeo.GeoIP, count chan combined) {
	scores := make(map[string]int)
	long_scores := make(map[string]int)
	total := 0
	long := 0
	c := time.Tick(60 * time.Second)
	var cc string
	for {
		select {
		case combi := <-count:
			ip := combi.ip
			scores[ip] += 1
			long_scores[ip] += 1
			total += 1
		case <-c:
			long += 1
			for ip, n := range scores {
				loc := gi.GetLocationByIP(ip)
				status := rbl(ip)
				if loc == nil {
					cc = "??"
				} else {
					cc = loc.CountryCode
				}
				fmt.Printf("%s %s #%d %s\n", cc, ip, n, status)
			}
			fmt.Printf("\t%d hits from %d ip\n", total, len(scores))
			scores = make(map[string]int)
			total = 0
			if long == 10 {
				long = 0
				for ip, n := range long_scores {
					status := rbl(ip)
					fmt.Printf("\tLong: %s #%d %s\n", ip, n, status)
				}
				long_scores = make(map[string]int)
			}
		}
	}
}

func main() {
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
	go consolidate(gi, count)
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
				mat[0][6]}
			count <- ip
		}
	}
}
