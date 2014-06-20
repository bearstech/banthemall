package main

import (
	"fmt"
	"net"
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
	address string
	msg     chan stat
	freq    time.Duration
}

func NewCarbon(address string, freq time.Duration) *Carbon {
	c := new(Carbon)
	c.address = address
	c.msg = make(chan stat)
	c.freq = freq
	go c.loop()
	return c
}

func (c Carbon) loop() {
	stat := make(map[string]int)
	t := time.Tick(c.freq)
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
		case now := <-t:
			fmt.Println("Tick", now)
			if len(stat) == 0 {
				continue
			}
			conn, err := net.Dial("tcp", c.address)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for k, v := range stat {
				fmt.Fprintf(conn, "%s %d %d\n", k, v, now.Unix())
			}
			stat = make(map[string]int)
			conn.Close()
		}
	}
	panic("aaaaahhhh")
}

func (c Carbon) Max(key string, value int) {
	c.msg <- stat{'m', key, value}
}

func (c Carbon) Sum(key string, value int) {
	c.msg <- stat{'s', key, value}
}
