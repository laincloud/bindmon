package http

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Check(line string) bool {
	arguments := strings.Split(line[strings.Index(line, ";check_http")+1:], " ")
	var CommandLine = flag.NewFlagSet(arguments[0], flag.ExitOnError)
	host := CommandLine.String("H", "", "")
	port := CommandLine.Int("p", 80, "")
	url := CommandLine.String("u", "/", "")
	timeout := CommandLine.Int("t", 10, "")
	CommandLine.Parse(arguments[1:])
	t, err := time.ParseDuration(strconv.Itoa(*timeout) + "s")
	if err != nil {
		log.Println(err)
		return false
	}
	client := http.Client{
		Timeout: t,
	}
	_,err = client.Get("http://" + *host + ":" + strconv.Itoa(*port) + *url)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
