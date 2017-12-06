package tcp

import (
	"flag"
	"net"
	"strconv"
	"time"
	"strings"
	"log"
)

func Check(line string) bool {
	arguments := strings.Split(line[strings.Index(line, ";check_tcp")+1:], " ")
	var CommandLine = flag.NewFlagSet(arguments[0], flag.ExitOnError)
	var host string
	var port, timeout int
	CommandLine.StringVar(&host,"H","","")
	CommandLine.IntVar(&port,"p",0,"")
	CommandLine.IntVar(&timeout,"t",10,"")
	CommandLine.Parse(arguments[1:])
	t, err := time.ParseDuration(strconv.Itoa(timeout) + "s")
	if err != nil {
		log.Println(err)
		return false
	}
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), t)
	if err != nil {
		log.Println(err)
		return false
	} else {
		defer conn.Close()
		return true
	}
}
