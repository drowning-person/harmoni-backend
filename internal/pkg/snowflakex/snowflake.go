package snowflakex

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

func NewSnowflakeNode(startTime string, machineId int64) (*sf.Node, error) {
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return nil, err
	}
	sf.Epoch = st.UnixNano() / 1000000

	return sf.NewNode(machineId)
}
