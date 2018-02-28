package monitor

import (
	"bufio"
	"fmt"
	"github.com/laincloud/bindmon/check"
	"log"
	"os"
	"strings"
	"sync"
)

type Monitor struct {
	src     string
	dst     string
	lines   []string
	health  []bool
	rewrite bool
	wg      sync.WaitGroup
	fall    int
}

func NewMonitor(src string, dst string, fall int) *Monitor {
	return &Monitor{
		src:     src,
		dst:     dst,
		fall:    fall,
		lines:   file2lines(src),
		health:  make([]bool, len(file2lines(src))),
		rewrite: false,
	}
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
	defer m.wg.Done()
	for i := begin + 1; i < end; i++ {
		now := false
		for j := 0; j < m.fall; j++ {
			if check.Check(m.lines[i]) {
				now = true
				break
			}
		}
		if now != m.health[i] {
			m.health[i] = now
			m.rewrite = true
		}
	}
}

func (m *Monitor) write() {
	log.Println("start write " + m.dst)
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
}

func (m *Monitor) Mon(ch chan int) {
	log.Println("start monitor src " + m.src + " dst " + m.dst)
	begin := -1
	for i, line := range m.lines {
		if line == ";begin-monitor" {
			begin = i
		}
		if line == ";end-monitor" {
			if begin != -1 && i > begin+1 {
				m.wg.Add(1)
				go m.check(begin, i)
			}
			begin = -1
		}
	}
	m.wg.Wait()
	if m.rewrite {
		m.write()
		ch <- 1
		m.rewrite = false
	} else {
		ch <- 0
	}
}
