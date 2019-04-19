package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ouqiang/supervisor-event-listener/utils"
)

// Message 消息格式
type Message struct {
	TS      time.Time
	Header  *Header
	Payload *Payload
}

func NewMessage(h *Header, p *Payload) Message {
	return Message{
		TS:      time.Now(),
		Header:  h,
		Payload: p,
	}
}

func (msg *Message) String() string {
	tmpl := `Host: %s
Process: %s
PID: %d
EXITED FROM state: %s
Date: %s`
	return fmt.Sprintf(tmpl,
		msg.Payload.Ip,
		msg.Payload.ProcessName,
		msg.Payload.Pid,
		msg.Payload.FromState,
		msg.TS.Format(time.RFC3339),
	)
}

func (msg *Message) ToJson(indent ...int) string {
	realIndent := 0
	if len(indent) > 0 {
		realIndent = indent[0]
	}
	t := ""
	switch realIndent {
	case 0:
	case 1:
		t = " "
	case 2:
		t = "  "
	case 3:
		t = "   "
	case 4:
		t = "    "
	default:
		t = "    "
	}
	_bytes, _ := json.MarshalIndent(msg, "", t)
	return string(_bytes)
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
