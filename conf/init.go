package conf

import (
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

func Init(fpath string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	c := Config{}
	if err := toml.Unmarshal(data, &c); err != nil {
		return err
	}
	if err := c.Validate(); err != nil {
		return err
	}
	cfg.Store(&c)
	return nil
}

func Reload(fpath string) error {
	return Init(fpath)
}

func Get() *Config {
	return cfg.Load().(*Config)
}
