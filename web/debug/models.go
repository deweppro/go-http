package debug

import "github.com/deweppro/go-http/web/server"

//go:generate easyjson

//easyjson:json
type (
	//Config model
	Config struct {
		Debug server.ConfigItem `yaml:"debug" json:"debug"`
	}
)
