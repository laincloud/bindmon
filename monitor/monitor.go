package monitor

import (
	"bufio"
	"fmt"
	"github.com/facebookgo/pidfile"
	"github.com/laincloud/bindmon/check"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

type Monitor struct {
	src    string
	dst    string
	pid    string
	lines  []string
	health []bool
}

func NewMonitor(src string, dst string, pid string) *Monitor {
	m := &Monitor{
		src:    src,
		dst:    dst,
		pid:    pid,
		lines:  file2lines(src),
		health: make([]bool, len(file2lines(src))),
	}

	for i := 0; i < len(m.health); i++ {
		m.health[i] = true
	}

	return m
}

func file2lines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if (len(line) > 0 && line[0] != ';') || line == ";begin-monitor" || line == ";end-monitor" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

func (m *Monitor) check(begin int, end int) {
	for {
		reload := false
		for i := begin + 1; i < end; i++ {
			now := check.Check(m.lines[i])
			if m.health[i] != now {
				reload = true
				m.health[i] = now
			}
		}
		if reload {
			m.reload()
		}
		time.Sleep(time.Minute)
	}
}

func (m *Monitor) reload() {
	var file *os.File
	begin := -1
	tmp := 0
	file, err := os.Create(m.dst)
	if err != nil {
		return
	}
	defer file.Close()
	os.Truncate(m.dst, 0)
	for i, line := range m.lines {
		if line == ";begin-monitor" {
			begin = i
			file.WriteString(strings.TrimSpace(line) + "\n")
			continue
		}
		if line == ";end-monitor" {
			if begin != -1 && i > begin+1 && tmp == 0 {
				file.WriteString(strings.TrimSpace(m.lines[begin+1]) + "\n")
			}
			file.WriteString(strings.TrimSpace(line) + "\n")
			begin = -1
			tmp = 0
			continue
		}
		if m.health[i] {
			file.WriteString(strings.TrimSpace(line) + "\n")
			if begin > -1 {
				tmp += 1
			}
		}
	}
	pidfile.SetPidfilePath(m.pid)
	pid, err := pidfile.Read()
	if err != nil {
		log.Println(err)
	}
	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		log.Println(err)
	}
}

func (m *Monitor) Mon() {
	log.Println("start monitor src " + m.src + " dst " + m.dst)
	begin := -1
	for i, line := range m.lines {
		if line == ";begin-monitor" {
			begin = i
		}
		if line == ";end-monitor" {
			if begin != -1 && i > begin+1 {
				go m.check(begin, i)
			}
			begin = -1
		}
	}
}
