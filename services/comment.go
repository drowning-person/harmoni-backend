package services

import (
	"fiberLearn/model"
	"fiberLearn/pkg/errcode"
	"fiberLearn/pkg/snowflake"
	"fiberLearn/pkg/zap"

	zaplog "go.uber.org/zap"
)

type CommentInsertService struct {
	PostID   int64  `json:"post_id,string" validate:"required" label:"帖子ID"`
	ParentID int64  `json:"parent_id,string" label:"父评论ID"`
	RootID   int64  `json:"root_id,string" label:"根评论ID"`
	Content  string `json:"content" validate:"required,gte=10,lte=256" label:"回复内容"`
}

func (c *CommentInsertService) Insert(authorID int64) (*model.CommentDetail, *errcode.Error) {
	comment := &model.Comment{
		AuthorID:  authorID,
		PostID:    c.PostID,
		ParentID:  c.ParentID,
		RootID:    c.RootID,
		Content:   c.Content,
		CommentID: snowflake.GenID(),
	}
	if err := model.DB.Create(&comment).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.CreatePostFailed
	}
	cd, err := model.CommentToCommentDetail(comment)
	if err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.CreateCommentFailed
	} else if cd == nil {
		zap.Logger.Warn(err.Error())
	}
	return cd, nil
}

type GetPostCommentsService struct {
	PostID int64 `query:"post_id,string"`
}

func (c *GetPostCommentsService) Retrieve(pager model.ParamListData) ([]model.CommentDetail, int64, *errcode.Error) {
	var total int64
	if err := model.DB.Table("comments").Where("post_id = ?", c.PostID).Count(&total).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, 0, errcode.GetCommentsFailed
	}
	if pager.Order == model.OrderByHot {

	} else if pager.Order == model.OrderByTime {
		comments := []model.Comment{}
		if err := model.DB.Order("comment_id DESC").Offset(int(pager.PageNum)).Limit(int(pager.PageSize)).Find(&comments).Error; err != nil {
			zap.Logger.Error(err.Error())
			return nil, 0, errcode.GetCommentsFailed
		}
		cds, err := model.CommentsToCommentDetails(comments)
		if err != nil {
			zap.Logger.Error(err.Error())
			return nil, 0, errcode.GetCommentsFailed
		} else if len(cds) == 0 {
			zap.Logger.Warn("该帖子没有回复", zaplog.Int("帖子ID", int(c.PostID)))
			return nil, 0, nil
		}
		return cds, total, nil
	}
	return nil, 0, nil
}
