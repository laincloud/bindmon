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
	line = line[strings.Index(line, ";check_tcp ")+1:]
	arguments := strings.Split(line, " ")
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
		log.Println(line + " ko")
		return false
	} else {
		defer conn.Close()
		log.Println(line + " ok")
		return true
	}
}
