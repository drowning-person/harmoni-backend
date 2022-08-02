package services

import (
	"fiberLearn/model"
	"fiberLearn/pkg/snowflake"
)

var db = model.DB.Table("tags")

func GetTags(offset, limit int) ([]*model.TagInfo, int64, error) {
	tags := []*model.TagInfo{}
	if err := db.Offset(offset).Limit(limit).Find(&tags).Error; err != nil {
		return nil, 0, err
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return tags, total, nil
}

type TagInsertService struct {
	TagName      string `json:"tag_name"`
	Introduction string `json:"introduction"`
}

func (t *TagInsertService) Insert() (*model.TagInfo, error) {
	tag := model.TagInfo{
		TagID:   snowflake.GenID(),
		TagName: t.TagName,
	}
	if err := db.Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func GetTagDetail(tagID int) (*model.TagDetail, error) {
	var tag model.TagDetail
	if err := db.Where("tag_id = ?", tagID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}
