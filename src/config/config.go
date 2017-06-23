package config

import (
	"github.com/Nomon/gonfig"
)

var Config *gonfig.Gonfig

func init() {
	Config = gonfig.NewConfig(nil)
	Config.Use("argv", gonfig.NewEnvConfig(""))
	Config.Use("env", gonfig.NewEnvConfig(""))
	Config.Use("local", gonfig.NewJsonConfig("./conf/config.json"))
}
