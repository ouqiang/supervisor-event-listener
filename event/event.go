package event

import (
	"errors"
	"fmt"
	"github.com/ouqiang/supervisor-event-listener/utils"
	"os"
	"strconv"
	"strings"
)

// Message 消息格式
type Message struct {
	Header  *Header
	Payload *Payload
}

func (msg *Message) String() string {
	return fmt.Sprintf("Host: %s\nProcess: %s\nPID: %d\nEXITED FROM state: %s", msg.Payload.Ip, msg.Payload.ProcessName, msg.Payload.Pid, msg.Payload.FromState)

}

// Header Supervisord触发事件时会先发送Header，根据Header中len字段去读取Payload
type Header struct {
	Ver        string
	Server     string
	Serial     int
	Pool       string
	PoolSerial int
	EventName  string // 事件名称
	Len        int    // Payload长度
}

// Payload
type Payload struct {
	Ip          string
	ProcessName string // 进程名称
	GroupName   string // 进程组名称
	FromState   string
	Expected    int
	Pid         int
}

// Fields
type Fields map[string]string

var (
	ErrParseHeader  = errors.New("解析Header失败")
	ErrParsePayload = errors.New("解析Payload失败")
)

func ParseHeader(header string) (*Header, error) {
	h := &Header{}
	fields := parseFields(header)
	if len(fields) == 0 {
		return h, ErrParseHeader
	}

	h.Ver = fields["ver"]
	h.Server = fields["server"]
	h.Serial, _ = strconv.Atoi(fields["serial"])
	h.Pool = fields["pool"]
	h.PoolSerial, _ = strconv.Atoi(fields["poolserial"])
	h.EventName = fields["eventname"]
	h.Len, _ = strconv.Atoi(fields["len"])

	return h, nil
}

func ParsePayload(payload string) (*Payload, error) {
	p := &Payload{}
	fields := parseFields(payload)
	if len(fields) == 0 {
		return p, ErrParsePayload
	}
	hostname, _ := os.Hostname()
	p.Ip = fmt.Sprintf("%s(%s)", utils.GetLocalIp(), hostname)
	p.ProcessName = fields["processname"]
	p.GroupName = fields["groupname"]
	p.FromState = fields["from_state"]
	p.Expected, _ = strconv.Atoi(fields["expected"])
	p.Pid, _ = strconv.Atoi(fields["pid"])

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
