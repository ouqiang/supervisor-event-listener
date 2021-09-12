package notify

import "github.com/ouqiang/supervisor-event-listener/event"

type Notifiable interface {
	Send(*event.Message) error
}
