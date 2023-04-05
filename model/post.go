package model

import (
	"harmoni/model/redis"
	"harmoni/pkg/zap"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Status   int32         `gorm:"not null"`
	PostID   int64         `gorm:"uniqueIndex"`
	AuthorID int64         `gorm:"index"`
	TagID    Int64toString `gorm:"type:varchar(128)"`
	TagName  string        `gorm:"type:varchar(512);index:,class:FULLTEXT,option:WITH PARSER ngram"`
	Title    string        `gorm:"type:varchar(128)"`
	Content  string        `gorm:"type:varchar(512)"`
}

type PostDetail struct {
	Status     int32     `json:"status"`
	Like       int64     `json:"like"`
	PostID     int64     `json:"post_id,string"`
	AuthorID   int64     `json:"author_id,string"`
	AuthorName string    `json:"author_name"`
	Tags       []TagInfo `json:"tags"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (p *Post) Format() (*PostDetail, error) {
	tagNames := strings.Split(p.TagName, ",")
	tags := make([]TagInfo, 0, len(p.TagID))
	for i, v := range p.TagID {
		tags = append(tags, TagInfo{
			TagID:   v,
			TagName: tagNames[i],
		})
	}
	pd := PostDetail{
		Status:    p.Status,
		PostID:    p.PostID,
		AuthorID:  p.AuthorID,
		Tags:      tags,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if err := DB.Select("name").Table("users").Where("user_id = ?", p.AuthorID).Scan(&pd.AuthorName).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		zap.Logger.Error(err.Error())
		return nil, err
	}
	var err error
	pd.Like, err = redis.GetPostLikeNumber(p.PostID)
	if err != nil {
		return nil, err
	}
	return &pd, nil
}

func GetPostByID(postID int64) (*PostDetail, error) {
	var post Post
	if err := DB.Where("post_id = ?", postID).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return post.Format()
}

func GetPostsByIDs(postIDs []string) ([]*PostDetail, error) {
	posts := make([]Post, 0, len(postIDs))
	if err := DB.Where("post_id in (?)", postIDs).Order("FIND_IN_SET(post_id,\"" + strings.Join(postIDs, ",") + "\")").Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return PostsToFormats(posts)
}

func PostsToFormats(posts []Post) ([]*PostDetail, error) {
	postFormats := make([]*PostDetail, 0, len(posts))
	for _, v := range posts {
		pd, err := v.Format()
		if err != nil {
			return nil, err
		}
		postFormats = append(postFormats, pd)
	}
	return postFormats, nil
}
