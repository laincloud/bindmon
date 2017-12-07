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
	line = line[strings.Index(line, ";check_udp ")+1:]
	arguments := strings.Split(line, " ")
	var CommandLine = flag.NewFlagSet(arguments[0], flag.ExitOnError)
	host := CommandLine.String("H", "", "")
	port := CommandLine.Int("p", 0, "")
	timeout := CommandLine.Int("t", 3, "")
	CommandLine.Parse(arguments[1:])
	t := time.Duration(*timeout) * time.Second
	conn, err := net.DialTimeout("udp", *host+":"+strconv.Itoa(*port), t)
	if err != nil {
		log.Println(line + " ko")
		return false
	} else {
		defer conn.Close()
		log.Println(line + " ok")
		return true
	}
}
