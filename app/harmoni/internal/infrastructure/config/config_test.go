package config_test

import (
	"harmoni/app/harmoni/internal/infrastructure/config"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigViperDefalut(t *testing.T) {
	v := viper.New()
	config.SetAppDefault(v)
	config.SetAuthDefault(v)
	config.SetDBDefault(v)
	config.SetLogDefault(v)
	config.SetRedisDefault(v)

	v.BindEnv("app.debug", "HARMONI_DEBUG")
	cf := config.Config{}

	err := v.Unmarshal(&cf)
	if err != nil {
		t.Fatal(err)
	}
	v.Unmarshal(cf.App)
	t.Logf("%#v", cf.App)
}
