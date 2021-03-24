package server

import "time"

//go:generate easyjson

//easyjson:json
type (
	//Config model
	Config struct {
		HTTP ConfigItem `yaml:"http" json:"http"`
	}
	//ConfigItem model
	ConfigItem struct {
		Addr            string        `yaml:"addr" json:"addr"`
		ReadTimeout     time.Duration `yaml:"read_timeout" json:"read_timeout"`
		WriteTimeout    time.Duration `yaml:"write_timeout" json:"write_timeout"`
		IdleTimeout     time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
	}
)
