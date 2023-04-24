package conf_test

import (
	"fmt"
	"harmoni/internal/conf"
	"testing"

	"github.com/gookit/config/v2"
	"github.com/spf13/viper"
)

func TestConfigDefalut(t *testing.T) {
	c := config.New("fly").WithOptions(config.ParseDefault)
	db := conf.Config{}
	// only set name
	c.SetData(map[string]any{
		"db": map[string]any{
			"driver": "${DB_DRIVER | mysql }",
		},
	})
	c.LoadFiles("../../config/config.yaml")
	err := c.Decode(&db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("fuck")
	t.Logf("%#v", db.DB)
}

func TestConfigViperDefalut(t *testing.T) {
	v := viper.New()
	conf.SetAppDefault(v)
	conf.SetAuthDefault(v)
	conf.SetDBDefault(v)
	conf.SetLogDefault(v)
	conf.SetRedisDefault(v)

	v.BindEnv("app.debug", "HARMONI_DEBUG")
	cf := conf.Config{}

	err := v.Unmarshal(&cf)
	if err != nil {
		t.Fatal(err)
	}
	v.Unmarshal(cf.App)
	t.Logf("%#v", cf.App)
}
