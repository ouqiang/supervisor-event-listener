package conf

type DryRun struct {
	Text string `toml:"text"`
}

// 邮件服务器
type Mail struct {
	Receivers      []string `toml:"receivers" validate:"required"`
	ServerUser     string   `toml:"server_user" validate:"required"`
	ServerPassword string   `toml:"server_password" validate:"required"`
	ServerHost     string   `toml:"server_host" validate:"required"`
	ServerPort     int      `toml:"server_port" validate:"required"`
}

type WebHook struct {
	URL     string `toml:"url" validate:"required"`
	Timeout int    `toml:"timeout"`
}

type Slack struct {
	URL     string `toml:"url" validate:"required"`
	Channel string `toml:"channel" validate:"required"`
}

type BearyChat struct {
	URL     string `toml:"url" validate:"required"`
	Channel string `toml:"channel"`
	Timeout int    `toml:"timeout"`
}

type Feishu struct {
	URL     string `toml:"url" validate:"required"`
	Timeout int    `toml:"timeout"`
}
