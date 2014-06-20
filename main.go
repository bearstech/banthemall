package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nranchev/go-libGeoIP"
)

func main() {
	flagThresold := flag.Int("thresold", 0, "Minimum mumber of hits per 10 seconds")

	flag.Parse()

	//FIXME try official Debian path, local path, and nothing.
	gi, err := libgeo.Load("GeoIP.dat")
	if err != nil {
		fmt.Printf("Error Libgeo: %s\n", err.Error())
		return
	}

	apachelog, err := NewCombinedParser()
	if err != nil {
		fmt.Printf("Error Regexp: %s\n", err.Error())
		return
	}

	bio := bufio.NewReader(os.Stdin)
	count := make(chan Combined)
	go consolidate(gi, *flagThresold, count)
	for {
		line, err := bio.ReadString('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if ip, ok := apachelog.Parse(line); ok {
			count <- ip
		}
	}
}
