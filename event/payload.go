package event

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ouqiang/supervisor-event-listener/utils"
)

type Payload struct {
	IP          string
	ProcessName string // 进程名称
	GroupName   string // 进程组名称
	FromState   string
	Expected    int
	PID         int
}

type Fields map[string]string

func ParsePayload(payload string) (*Payload, error) {
	p := &Payload{}
	fields := parseFields(payload)
	if len(fields) == 0 {
		return nil, ErrParsePayload
	}
	hostname, _ := os.Hostname()
	p.IP = fmt.Sprintf("%s(%s)", utils.GetLocalIp(), hostname)
	p.ProcessName = fields["processname"]
	p.GroupName = fields["groupname"]
	p.FromState = fields["from_state"]
	p.Expected, _ = strconv.Atoi(fields["expected"])
	p.PID, _ = strconv.Atoi(fields["pid"])
	return p, nil
}

func parseFields(data string) Fields {
	fields := make(Fields)
	data = strings.TrimSpace(data)
	if data == "" {
		return fields
	}
	// 格式如下
	// ver:3.0 server:supervisor serial:5
	slice := strings.Split(data, " ")
	if len(slice) == 0 {
		return fields
	}
	for _, item := range slice {
		group := strings.Split(item, ":")
		if len(group) < 2 {
			continue
		}
		key := strings.TrimSpace(group[0])
		value := strings.TrimSpace(group[1])
		fields[key] = value
	}

	return fields
}
