package tcp

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func Check(line string) bool {
	arguments := strings.Split(line[strings.Index(line, ";check_tcp ")+1:], " ")
	var CommandLine = flag.NewFlagSet(arguments[0], flag.ExitOnError)
	var host string
	var port, timeout int
	CommandLine.StringVar(&host, "H", "", "")
	CommandLine.IntVar(&port, "p", 0, "")
	CommandLine.IntVar(&timeout, "t", 3, "")
	CommandLine.Parse(arguments[1:])
	t := time.Duration(timeout) * time.Second
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), t)
	if err != nil {
		log.Println(err)
		return false
	} else {
		defer conn.Close()
		return true
	}
}
