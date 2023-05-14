package snowflakex

import (
	"harmoni/internal/conf"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

func NewSnowflakeNode(conf *conf.App) (*sf.Node, error) {
	var st time.Time
	st, err := time.Parse("2006-01-02", conf.StartTime)
	if err != nil {
		return nil, err
	}
	sf.Epoch = st.UnixNano() / 1000000

	return sf.NewNode(conf.AppID)
}
