package main

import (
    "fmt"
    "strings"
    "net"
)


func Rbl(ip string) (status string) {
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
