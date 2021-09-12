package conf

import (
	"fmt"
	"sync/atomic"

	validator "github.com/go-playground/validator/v10"
)

var (
	// c = &Config{}
	cfg atomic.Value
)

type Config struct {
	DryRun    *DryRun
	WebHook   *WebHook
	Mail      *Mail
	Slack     *Slack
	BearyChat *BearyChat
	Feishu    *Feishu
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("nil config")
	}

	if c.BearyChat != nil {
		if c.BearyChat.URL == "" {
			return fmt.Errorf("Invalid bearcychat config")
		}
	}
	if c.Mail != nil {
		if len(c.Mail.Receivers) <= 0 {
			return fmt.Errorf("Invalid mail config")
		}
	}
	if c.Slack != nil {
		if c.Slack.URL == "" {
			return fmt.Errorf("Invalid slack config")
		}
	}

	if c.WebHook != nil {
		if c.WebHook.URL == "" {
			return fmt.Errorf("Invalid slack config")
		}
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}
	return nil
}
