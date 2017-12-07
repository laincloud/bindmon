package udp

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func Check(line string) bool {
	arguments := strings.Split(line[strings.Index(line, ";check_udp ")+1:], " ")
	var CommandLine = flag.NewFlagSet(arguments[0], flag.ExitOnError)
	host := CommandLine.String("H", "", "")
	port := CommandLine.Int("p", 0, "")
	timeout := CommandLine.Int("t", 2, "")
	CommandLine.Parse(arguments[1:])
	t := time.Duration(*timeout) * time.Second
	conn, err := net.DialTimeout("udp", *host+":"+strconv.Itoa(*port), t)
	if err != nil {
		log.Println(err)
		return false
	} else {
		defer conn.Close()
		return true
	}
}
