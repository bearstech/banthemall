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

func consolidate(gi *libgeo.GeoIP, count chan string) {
	scores := make(map[string]int)
	c := time.Tick(30 * time.Second)
	var loc *libgeo.Location
	var status string
	var cc string

	for {
		select {
		case ip := <-count:
			scores[ip] += 1
		case <-c:
			for ip, n := range scores {
				loc = gi.GetLocationByIP(ip)
				status = rbl(ip)
				cc = loc.CountryCode
				fmt.Printf("%s %s #%d %s\n", cc, ip, n, status)
				delete(scores, ip)
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

	apachelog, err := regexp.Compile("(.*) (.*) (.*) (\\[.*\\])")
	if err != nil {
		fmt.Printf("Error Regexp: %s\n", err.Error())
		return
	}

	bio := bufio.NewReader(os.Stdin)
	count := make(chan string)
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
		mat := apachelog.FindStringSubmatch(line)
		ip := string(mat[1])
		count <- ip
	}
}
