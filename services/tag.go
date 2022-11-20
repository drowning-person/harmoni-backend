package services

import (
	"fiberLearn/model"
	"fiberLearn/pkg/errcode"
	"fiberLearn/pkg/snowflake"
	"fiberLearn/pkg/zap"
	"html"

	"gorm.io/gorm"
)

func GetTags(param *model.ParamListData) ([]*model.TagInfo, int64, *errcode.Error) {
	tags := []*model.TagInfo{}
	if err := model.DB.Table("tags").Offset(int(param.PageNum)).Limit(int(param.PageSize)).Find(&tags).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, 0, errcode.GetTagsFailed
	}
	var total int64
	if err := model.DB.Table("tags").Count(&total).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, 0, errcode.GetTagsFailed
	}
	return tags, total, nil
}

type TagInsertService struct {
	TagName      string `json:"tag_name" validate:"required,lte=30" label:"话题名"`
	Introduction string `json:"introduction" validate:"required,lte=256" label:"话题简介"`
}

func (t *TagInsertService) Insert() (*model.TagDetail, *errcode.Error) {
	if exist, err := model.IsTagExist(t.TagName); err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.CreateTagFailed
	} else if exist {
		return nil, errcode.TagHasExisted
	}
	tag := model.TagDetail{
		TagID:        snowflake.GenID(),
		TagName:      html.EscapeString(t.TagName),
		Introduction: html.EscapeString(t.Introduction),
	}
	if err := model.DB.Table("tags").Create(&tag).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.CreateTagFailed
	}
	return &tag, nil
}

func GetTagDetail(tagID int) (*model.TagDetail, *errcode.Error) {
	var tag model.TagDetail
	if err := model.DB.Table("tags").Where("tag_id = ?", tagID).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.Logger.Warn(err.Error())
			return nil, errcode.NotFound.WithDetails("话题不存在")
		}
		zap.Logger.Error(err.Error())
		return nil, errcode.GetTagFailed
	}
	return &tag, nil
}
