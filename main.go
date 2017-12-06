package main

import (
	"flag"
	"github.com/laincloud/bindmon/monitor"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	log.Println("start bindmon")
	src := flag.String("src", "", "")
	dst := flag.String("dst", "", "")
	pid := flag.String("pid", "", "")
	flag.Parse()
	log.Println("src is "+*src)
	log.Println("dst is "+*dst)
	log.Println("pid is "+*pid)
	files, _ := ioutil.ReadDir(*src)
	for _, fi := range files {
		if fi.IsDir() {
			continue
		} else {
			m := monitor.NewMonitor(*src+"/"+fi.Name(), *dst+"/"+fi.Name(), *pid)
			go m.Mon()
		}
	}
	for {
		time.Sleep(time.Hour)
	}
}
