package services

import (
	"fiberLearn/model"
	"fiberLearn/model/redis"
	"fiberLearn/pkg/errcode"
	"fiberLearn/pkg/snowflake"
	"fiberLearn/pkg/zap"
	"html"
	"strings"

	"gorm.io/gorm"
)

func GetPosts(param *model.ParamListData) ([]*model.PostDetail, int64, *errcode.Error) {
	var total int64
	if err := model.DB.Table("posts").Count(&total).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, 0, errcode.GetPostsFailed
	}
	if param.Order == model.OrderByHot {
		ids, err := redis.GetPostIdsByScore(param.PageSize, param.PageNum)
		if err != nil {
			return nil, 0, errcode.GetPostsFailed
		}
		posts, err := model.GetPostsByIDs(ids)
		if err != nil {
			return nil, 0, errcode.GetPostsFailed
		}
		return posts, total, nil
	} else if param.Order == model.OrderByTime {
		posts := []model.Post{}
		if err := model.DB.Order("post_id DESC").Offset(int(param.PageNum)).Limit(int(param.PageSize)).Find(&posts).Error; err != nil {
			zap.Logger.Error(err.Error())
			return nil, 0, errcode.GetPostsFailed
		}
		postFormats, err := model.PostsToFormats(posts)
		if err != nil {
			return nil, 0, errcode.GetPostsFailed
		}
		return postFormats, total, nil
	}
	return nil, 0, nil
}

func GetPostDetail(postID int64) (*model.PostDetail, *errcode.Error) {
	postFormated, err := model.GetPostByID(postID)
	if err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.GetPostFailed
	}
	if postFormated == nil {
		return nil, errcode.NotFound.WithDetails("帖子不存在")
	}

	return postFormated, nil
}

type PostInsertService struct {
	TagID   model.Int64toString `json:"tag_id" validate:"lte=4" label:"话题ID"`
	Title   string              `json:"title" validate:"required,gte=3,lte=128" label:"帖子标题"`
	Content string              `json:"content" validate:"required,gte=10,lte=512" label:"帖子内容"`
}

func (p *PostInsertService) Insert(authorID int64) (*model.PostDetail, *errcode.Error) {
	m := make(map[int64]struct{}, 4)
	tagNames := make([]string, 0, 4)
	for _, v := range p.TagID {
		tag := model.Tag{}
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
		} else {
			return nil, errcode.SameTag
		}
		if err := model.DB.Table("tags").Where("tag_id = ?", v).First(&tag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				zap.Logger.Warn(err.Error())
				return nil, errcode.TagHasNotExisted
			}
			zap.Logger.Error(err.Error())
			return nil, errcode.CreatePostFailed
		}
		tagNames = append(tagNames, tag.TagName)
	}

	post := model.Post{
		PostID:   snowflake.GenID(),
		AuthorID: authorID,
		TagID:    p.TagID,
		TagName:  strings.Join(tagNames, ","),
		Title:    html.EscapeString(p.Title),
		Content:  html.EscapeString(p.Content),
	}

	if err := model.DB.Table("posts").Create(&post).Error; err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.CreatePostFailed
	}
	if err := redis.AddPost(post.PostID); err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.CreatePostFailed
	}

	postFormated, err := post.Format()
	if err != nil {
		zap.Logger.Error(err.Error())
		return nil, errcode.GetPostFailed
	} else if postFormated == nil {
		zap.Logger.Warn(err.Error())
	}
	return postFormated, nil
}

type PostLikeService struct {
	PostID int64 `json:"post_id,string" validate:"required"`
	Like   int8  `json:"like" validate:"required,oneof=1 2"`
}

func (p *PostLikeService) LikePost(authorID int64) *errcode.Error {
	like, flag, err := redis.CheckLike(p.PostID, authorID)
	if err != nil {
		zap.Logger.Error(err.Error())
		return errcode.LikePostFailed
	}
	if p.Like == 1 && like == 1 {
		return errcode.HasLikedPost
	} else if p.Like == 2 && (like == -1 || !flag) {
		return errcode.HasNotLikedPost
	}
	if err := redis.DoLike(p.PostID, authorID, p.Like); err != nil {
		zap.Logger.Error(err.Error())
		return errcode.LikePostFailed
	}
	return nil
}
