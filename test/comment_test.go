package test

import (
	"harmoni/config"
	"harmoni/model"
	"harmoni/pkg/snowflake"
	"harmoni/pkg/zap"
	"testing"
)

func TestCreateComment(t *testing.T) {
	cfg, err := config.ReadConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}
	snowflake.Init("2022-07-31", 1)
	zap.InitLogger(cfg.Log)
	model.InitMysql(cfg.DB)
	for i := 0; i < 100000; i++ {
		c := model.Comment{
			CommentID: snowflake.GenID(),
			PostID:    4201748478562304,
			AuthorID:  1976392388448256,
			Content:   "sb",
			ParentID:  44250132878725120,
			RootID:    41415298993098752,
		}
		model.DB.Create(&c)
	}
}
