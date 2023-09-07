package post

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"harmoni/internal/entity"
	"harmoni/internal/entity/paginator"
	tagentity "harmoni/internal/entity/tag"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Int64toString []int64

func (is Int64toString) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.Grow(64)
	b.WriteByte('[')
	for i, v := range is {
		b.WriteByte('"')
		b.Write([]byte(strconv.FormatInt(v, 10)))
		b.WriteByte('"')
		if i+1 != len(is) {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.Bytes(), nil
}

func (is *Int64toString) UnmarshalJSON(data []byte) error {
	str := strings.Split(strings.Trim(string(data), "[]"), ",")
	for _, v := range str {
		num, err := strconv.ParseInt(strings.Trim(v, " \""), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		*is = append(*is, num)
	}
	return nil
}

// Scan 方法实现了 sql.Scanner 接口
func (is *Int64toString) Scan(v interface{}) error {
	if s, ok := v.([]uint8); !ok {
		return fmt.Errorf("断言失败")
	} else {
		str := strings.Split(string(s), ",")
		for _, v := range str {
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			*is = append(*is, num)
		}
	}
	return nil
}

func (is Int64toString) Value() (driver.Value, error) {
	if len(is) == 0 {
		return nil, nil
	}
	s := make([]string, 0, 4)
	for _, v := range is {
		s = append(s, strconv.FormatInt(v, 10))
	}
	return strings.Join(s, ","), nil
}

type PostContentT int

const (
	ContentTypeTitle PostContentT = iota + 1
	ContentTypeText
	ContentTypeImage
	ContentTypeVideo
	ContentTypeAudio
	ContentTypeLink
)

var (
	mediaContentType = []PostContentT{
		ContentTypeImage,
		ContentTypeVideo,
		ContentTypeAudio,
	}
)

type PostContent struct {
	*entity.BaseModelWithNoSoftDelete
	PostID  int64        `json:"post_id"`
	UserID  int64        `json:"user_id"`
	Content string       `json:"content"`
	Type    PostContentT `json:"type"`
	Sort    int64        `json:"sort"`
}

type Post struct {
	gorm.Model
	Status    int32   `gorm:"not null"`
	PostID    int64   `gorm:"uniqueIndex"`
	AuthorID  int64   `gorm:"index"`
	Title     string  `gorm:"type:varchar(128)"`
	Content   string  `gorm:"type:varchar(512)"`
	LikeCount int64   `gorm:"not null"`
	TagIDs    []int64 `gorm:"-"`
}

func (Post) TableName() string {
	return "post"
}

type PostDetail struct {
	Status    int32               `json:"status,omitempty"`
	LikeCount int64               `json:"like_count"`
	PostID    int64               `json:"post_id,string,omitempty"`
	AuthorID  int64               `json:"author_id,string,omitempty"`
	Tags      []tagentity.TagInfo `json:"tags,omitempty"`
	Title     string              `json:"title,omitempty"`
	Content   string              `json:"content,omitempty"`
	CreatedAt time.Time           `json:"created_at,omitempty"`
	UpdatedAt time.Time           `json:"updated_at,omitempty"`
}

type PostBasicInfo struct {
	AuthorID  int64  `json:"author_id,omitempty"`
	PostID    int64  `json:"post_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	LikeCount int64  `json:"like_count"`
}

func ConvertPostToDisplay(post *Post) PostBasicInfo {
	return PostBasicInfo{
		AuthorID:  post.AuthorID,
		PostID:    post.PostID,
		Title:     post.Title,
		Content:   post.Content,
		LikeCount: post.LikeCount,
	}
}

func ConvertPostToDisplayDetail(post *Post, tags []tagentity.Tag) PostDetail {
	tagInfos := make([]tagentity.TagInfo, len(tags))
	for i, tag := range tags {
		tagInfos[i] = tagentity.ConvertTagToDisplay(&tag)
	}
	pd := PostDetail{
		Status:    post.Status,
		PostID:    post.PostID,
		AuthorID:  post.AuthorID,
		Tags:      tagInfos,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		LikeCount: post.LikeCount,
	}

	return pd
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	Delete(ctx context.Context, postID int64) error
	GetBasicInfoByPostID(ctx context.Context, postID int64) (*Post, bool, error)
	GetByUserID(ctx context.Context, userID int64, queryCond *PostQuery) (paginator.Page[Post], error)
	GetByUserIDs(ctx context.Context, userID []int64, queryCond *PostQuery) (paginator.Page[Post], error)
	GetByPostID(ctx context.Context, postID int64) (*Post, bool, error)
	GetPostsByTagID(ctx context.Context, tagID int64) ([]Post, error)
	BatchByIDs(ctx context.Context, postIDs []int64) ([]Post, error)
	BatchBasicInfoByIDs(ctx context.Context, postID []int64) ([]Post, error)
	GetLikeCount(ctx context.Context, postID int64) (int64, bool, error)
	UpdateLikeCount(ctx context.Context, postID int64, count int64) error
	GetPage(ctx context.Context, queryCond *PostQuery) (paginator.Page[Post], error)
}
