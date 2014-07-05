package main

import (
	"banthemall/metrics"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type stat struct {
	action rune
	key    string
	value  int
}

/*
Carbon client
*/
type Carbon struct {
	/*Address of the carbond server*/
	address string
	msg     chan stat
	/*Consolidation frequency*/
	freq time.Duration
}

/*
NewCarbon return a fresh Carbon client.
*/
func NewCarbon(address string, freq time.Duration) *Carbon {
	c := new(Carbon)
	c.address = address
	c.msg = make(chan stat)
	c.freq = freq
	if address != "" {
		go c.loop()
	}
	return c
}

func (c Carbon) loop() {
	stat := make(map[string]int)
	lstat := make(map[string]*metrics.Percentile)
	t := time.Tick(c.freq)
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Trouble with hostname", err)
		hostname = "john-doe"
	}
	for {
		select {
		case s := <-c.msg:
			if s.action == 'm' {
				if old, ok := stat[s.key]; ok {
					if s.value > old {
						stat[s.key] = s.value
					}
				} else {
					stat[s.key] = s.value
				}
			}
			if s.action == 's' {
				if old, ok := stat[s.key]; ok {
					stat[s.key] = old + s.value
				} else {
					stat[s.key] = s.value
				}
			}
			if s.action == 'l' {
				if _, ok := lstat[s.key]; !ok {
					lstat[s.key] = metrics.NewPercentile()
				}
				lstat[s.key].Append(s.value)
			}
		case now := <-t:
			if len(stat) == 0 {
				continue
			}
			conn, err := net.Dial("tcp", c.address)
			if err != nil {
				log.Println(err)
				continue
			}
			for k, v := range stat {
				fmt.Fprintf(conn, "servers.%s.%s %d %d\n", hostname, k, v, now.Unix())
			}
			for k, v := range lstat {
				fmt.Fprintf(conn, "servers.%s.%s.50 %d %d\n", hostname, k, v.Percentile(50), now.Unix())
				fmt.Fprintf(conn, "servers.%s.%s.95 %d %d\n", hostname, k, v.Percentile(95), now.Unix())
			}
			stat = make(map[string]int)
			lstat = make(map[string]*metrics.Percentile)
			if err = conn.Close(); err != nil {
				log.Println(err)
			}
		}
	}
	panic("aaaaahhhh")
}

/*
Max counter
*/
func (c Carbon) Max(key string, value int) {
	if c.address != "" {
		c.msg <- stat{'m', key, value}
	}
}

/*
Sum counter
*/
func (c Carbon) Sum(key string, value int) {
	if c.address != "" {
		c.msg <- stat{'s', key, value}
	}
}

func (c Carbon) List(key string, value int) {
	if c.address != "" {
		c.msg <- stat{'l', key, value}
	}
}
