package main

import (
	"banthemall/combined"
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
	flagCarbon := flag.String("carbon", "", "Send some metrics to carbond")

	flag.Parse()

	carbon := NewCarbon(*flagCarbon, 1*time.Minute)

	//FIXME try official Debian path, local path, and nothing.
	gi, err := libgeo.Load("GeoIP.dat")
	if err != nil {
		fmt.Printf("Error Libgeo: %s\n", err.Error())
		return
	}

	apachelog, err := combined.NewCombinedParser()
	if err != nil {
		fmt.Printf("Error Regexp: %s\n", err.Error())
		return
	}

	bio := bufio.NewReader(os.Stdin)
	count := make(chan combined.Combined)
	go consolidate(gi, *flagThresold, carbon, count)
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
