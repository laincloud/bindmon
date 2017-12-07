package main

import (
	"flag"
	"github.com/laincloud/bindmon/bind"
	"github.com/laincloud/bindmon/monitor"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func main() {
	log.Println("start bindmon")
	src := flag.String("src", "", "")
	dst := flag.String("dst", "", "")
	pid := flag.String("pid", "", "")
	flag.Parse()
	log.Println("src is " + *src)
	log.Println("dst is " + *dst)
	log.Println("pid is " + *pid)
	files, err := ioutil.ReadDir(*src)
	if err != nil {
		log.Panic(err)
	}
	count := 0
	monitors := make([]monitor.Monitor, count)
	for _, fi := range files {
		if fi.IsDir() {
			continue
		} else {
			count++
			m := monitor.NewMonitor(*src+"/"+fi.Name(), *dst+"/"+fi.Name())
			monitors = append(monitors, *m)
		}
	}
	ch := make(chan int, count)
	for {
		var wg sync.WaitGroup
		wg.Add(count)
		for _, m := range monitors {
			go func() {
				m.Mon(ch)
				wg.Done()
			}()
		}
		wg.Wait()
		reload := 0
		for i := 0; i < count; i++ {
			reload += <-ch
		}
		if reload > 0 {
			bind.Reload(*pid)
		}
		time.Sleep(time.Minute)
	}
}
