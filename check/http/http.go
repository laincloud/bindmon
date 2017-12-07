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
	line = line[strings.Index(line, ";check_http ")+1:]
	arguments := strings.Split(line , " ")
	var CommandLine = flag.NewFlagSet(arguments[0], flag.ExitOnError)
	host := CommandLine.String("H", "", "")
	port := CommandLine.Int("p", 80, "")
	url := CommandLine.String("u", "/", "")
	timeout := CommandLine.Int("t", 3, "")
	CommandLine.Parse(arguments[1:])
	t := time.Duration(*timeout) * time.Second
	client := http.Client{
		Timeout: t,
	}
	_, err := client.Get("http://" + *host + ":" + strconv.Itoa(*port) + *url)
	if err != nil {
		log.Println(line + " ko")
		return false
	}
	log.Println(line + " ok")
	return true
}
