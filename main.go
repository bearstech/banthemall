package main

import (
	"bufio"
	"fmt"
	"github.com/nranchev/go-libGeoIP"
	"io"
	"os"
	"regexp"
)

func main() {
	gi, err := libgeo.Load("GeoIP.dat")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	again := true
	apachelog, err := regexp.Compile("(.*) (.*) (.*) (\\[.*\\])")
	if err != nil {
		panic(err)
	}
	bio := bufio.NewReader(os.Stdin)
	for again {
		line, hasMoreInLine, err := bio.ReadLine()
		if err == io.EOF {
			//again = false
			continue
		}
		if err != nil {
			panic(err)
		}
		if hasMoreInLine {
			fmt.Println("Line too long")
		}
		mat := apachelog.FindSubmatch(line)
		ip := string(mat[1])
		loc := gi.GetLocationByIP(ip)
		fmt.Print(loc.CountryCode)
		fmt.Print(" ")
		fmt.Println(ip)

	}
}
