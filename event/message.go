package event

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ouqiang/supervisor-event-listener/utils/errlog"
	"github.com/pkg/errors"
)

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

func ReadMessage(reader *bufio.Reader) (Message, error) {
	header, err := readHeader(reader)
	if err != nil {
		errlog.Error("header:%+v err:%+v", header, err)
		return Message{}, err
	}
	payload, err := readPayload(reader, header.Len)
	if err != nil {
		errlog.Error("payload:%+v err:%+v", payload, err)
		return Message{}, err
	}
	return NewMessage(header, payload), nil
}

// 读取header
func readHeader(reader *bufio.Reader) (*Header, error) {
	// 读取Header
	data, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// 解析Header
	header, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// 读取payload
func readPayload(reader *bufio.Reader, payloadLen int) (*Payload, error) {
	// 读取payload
	buf := make([]byte, payloadLen)
	length, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}
	if payloadLen != length {
		err := ErrPayloadLength
		err = errors.Wrapf(err, " payloadLen:%d != length:%d", payloadLen, length)
		return nil, err
	}
	// 解析payload
	payload, err := ParsePayload(string(buf))
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (msg *Message) String() string {
	tmpl := `Proc: %s
Host: %s
PID:  %d    State: %s
Date: %s`
	return fmt.Sprintf(tmpl,
		msg.Payload.ProcessName,
		msg.Payload.IP,
		msg.Payload.PID, msg.Payload.FromState,
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
