package services

import (
	"fiberLearn/model"
	"fiberLearn/pkg/snowflake"
	"time"
)

type TagInsertService struct {
	TagName      string `json:"tag_name"`
	Introduction string `json:"introduction"`
}

func (t *TagInsertService) Insert() (model.Tag, error) {
	tag := model.Tag{
		TagID:        snowflake.GenID(),
		TagName:      t.TagName,
		Introduction: t.Introduction,
	}
	if err := model.DB.Create(&tag).Error; err != nil {
		return model.Tag{}, err
	}
	return tag, nil
}

type TagGetService struct {
	TagID        int64  `json:"tag_id"`
	TagName      string `json:"tag_name"`
	Introduction string `json:"introduction"`
}

func GetTags() ([]*TagGetService, error) {
	tags := []*TagGetService{}
	tags1 := []*model.Tag{}
	if err := model.DB.Find(&tags1).Error; err != nil {
		return nil, err
	}
	for _, v := range tags1 {
		tag := TagGetService{
			TagID:        v.TagID,
			TagName:      v.TagName,
			Introduction: v.Introduction,
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}

type TagDetailService struct {
	TagID        int64     `json:"tag_id"`
	TagName      string    `json:"tag_name"`
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (t *TagDetailService) Get() error {
	tag := model.Tag{
		TagID: t.TagID,
	}
	if err := model.DB.Where("tag_id = ?", tag.TagID).First(&tag).Error; err != nil {
		return err
	}
	t.CreatedAt = tag.CreatedAt
	t.UpdatedAt = tag.UpdatedAt
	t.TagName = tag.TagName
	t.Introduction = tag.Introduction
	return nil
}
