package check

import (
	"github.com/laincloud/bindmon/check/tcp"
	"strings"
	"github.com/laincloud/bindmon/check/udp"
	"github.com/laincloud/bindmon/check/http"
)

func Check(line string) bool {
	if strings.Contains(line, ";check_tcp") {
		return tcp.Check(line)
	} else if strings.Contains(line, ";check_udp") {
		return udp.Check(line)
	} else if strings.Contains(line, ";check_http") {
		return http.Check(line)
	} else {
		return true
	}
}
