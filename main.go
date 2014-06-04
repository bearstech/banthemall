package main

import (
	"bufio"
	"fmt"
	"github.com/nranchev/go-libGeoIP"
	"io"
	"os"
	"regexp"
	"time"
)

func consolidate(gi *libgeo.GeoIP, count chan string) {
	scores := make(map[string]int)
	c := time.Tick(30 * time.Second)
	for {
		select {
		case ip := <-count:
			scores[ip] += 1
		case <-c:
			for ip, n := range scores {
				loc := gi.GetLocationByIP(ip)
				fmt.Printf("%s %s #%d\n", loc.CountryCode, ip, n)
			}
			scores = make(map[string]int)
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
		line, hasMoreInLine, err := bio.ReadLine()
		if err == io.EOF {
			continue
		}
		if err != nil {
			fmt.Printf("Error Regexp: %s\n", err.Error())
			continue
		}
		if hasMoreInLine {
			fmt.Println("Line too long")
		}
		mat := apachelog.FindSubmatch(line)
		ip := string(mat[1])
		count <- ip
	}
}
