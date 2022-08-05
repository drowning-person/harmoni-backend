package services

import (
	"fiberLearn/model"
	"fiberLearn/pkg/snowflake"
	"html"
)

func GetTags(offset, limit int) ([]*model.TagInfo, int64, error) {
	tags := []*model.TagInfo{}
	if err := model.DB.Table("tags").Offset(offset).Limit(limit).Find(&tags).Error; err != nil {
		return nil, 0, err
	}
	var total int64
	if err := model.DB.Table("tags").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return tags, total, nil
}

type TagInsertService struct {
	TagName      string `json:"tag_name" validate:"required,lte=30" label:"话题名"`
	Introduction string `json:"introduction" validate:"required,lte=256" label:"话题简介"`
}

func (t *TagInsertService) Insert() (*model.TagDetail, error) {
	if exist, err := model.IsTagExist(t.TagName); err != nil {
		return nil, err
	} else if exist {
		return nil, nil
	}
	tag := model.TagDetail{
		TagID:        snowflake.GenID(),
		TagName:      html.EscapeString(t.TagName),
		Introduction: html.EscapeString(t.Introduction),
	}
	if err := model.DB.Table("tags").Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func GetTagDetail(tagID int) (*model.TagDetail, error) {
	var tag model.TagDetail
	if err := model.DB.Table("tags").Where("tag_id = ?", tagID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}
