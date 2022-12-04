package test

import (
	"fiberLearn/model"
	"fiberLearn/pkg/snowflake"
	"fiberLearn/pkg/zap"
	"testing"
)

func TestCreateComment(t *testing.T) {
	snowflake.Init("2022-07-31", 1)
	zap.InitLogger("./log/", "main.log")
	model.InitMysql("utf8mb4")
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
