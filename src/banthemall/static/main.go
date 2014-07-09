package main

import (
	"banthemall/combined"
	"banthemall/consolidate"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Where is the log?")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	parser, err := combined.NewCombinedParser()
	if err != nil {
		log.Panic(err)
	}
	bio := bufio.NewReader(f)
	db := []combined.Combined{}
	var short_ts time.Time
	var short *consolidate.ShortTerm
	for {
		line, err := bio.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		c, ok := parser.Parse(line)
		if ok {
			ts, err := c.Time()
			if err != nil {
				log.Panic(err)
			}
			db = append(db, c)
			if short == nil {
				short_ts = *ts
				short = consolidate.NewShortTerm()
			} else {
				if ts.Sub(short_ts) > 10*time.Second {
					// consolidate
					fmt.Println(short_ts, short.Size(), short.Hits())
					short = consolidate.NewShortTerm()
					short_ts = short_ts.Add(10 * time.Second)
				}
			}
			short.Add(&c)
		}
	}
	fmt.Println(len(db))
}
