package servers

import "time"

const (
	StatusOn  = 1
	StatusOff = 0
)

type Config struct {
	Addr            string        `yaml:"addr"`
	Network         string        `yaml:"network,omitempty"`
	ReadTimeout     time.Duration `yaml:"read_timeout,omitempty"`
	WriteTimeout    time.Duration `yaml:"write_timeout,omitempty"`
	IdleTimeout     time.Duration `yaml:"idle_timeout,omitempty"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout,omitempty"`
}
