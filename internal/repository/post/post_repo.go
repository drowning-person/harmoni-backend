package post

import (
	"context"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/entity/unique"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"html"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	tagRepo      tagentity.TagRepository
	logger       *zap.SugaredLogger
}

func NewPostRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, tagRepo tagentity.TagRepository, logger *zap.SugaredLogger) postentity.PostRepository {
	return &postRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		tagRepo:      tagRepo,
		logger:       logger,
	}
}

func (r *postRepo) Create(ctx context.Context, post *postentity.Post) (err error) {
	m := make(map[int64]struct{}, 4)
	tagNames := strings.Builder{}
	tagNames.Grow(8)

	for _, id := range post.TagID {
		if _, ok := m[id]; !ok {
			m[id] = struct{}{}
		} else {
			return errorx.BadRequest(reason.TagDuplicateInObeject)
		}
		if tag, err := r.tagRepo.GetByTagID(ctx, id); err != nil {
			return err
		} else {
			tagNames.WriteString(tag.TagName)
			tagNames.WriteByte(',')
		}
	}
	post.TagName = tagNames.String()

	post.PostID, err = r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return err
	}

	post.Title = html.EscapeString(post.Title)
	post.Content = html.EscapeString(post.Content)

	err = r.db.WithContext(ctx).Create(post).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	err = r.addPost(ctx, post.PostID)
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *postRepo) GetByPostID(ctx context.Context, postID int64) (postentity.Post, error) {
	post := postentity.Post{}
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return postentity.Post{}, errorx.NotFound(reason.PostNotFound)
		}
		return postentity.Post{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	post.LikeCount, err = r.getPostLikeNumber(ctx, postID)
	if err != nil {
		return postentity.Post{}, err
	}

	return post, nil
}

func (r *postRepo) BatchByIDs(ctx context.Context, postIDs []int64) ([]postentity.Post, error) {
	posts := make([]postentity.Post, 0, len(postIDs))
	if err := r.db.WithContext(ctx).Where("post_id IN ?", postIDs).Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.NotFound(reason.PostNotFound)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return posts, nil
}

func (r *postRepo) GetPage(ctx context.Context, pageSize, pageNum int64, orderCond string) (paginator.Page[postentity.Post], error) {
	db := r.db.WithContext(ctx)

	var (
		ids []string
		err error
	)

	switch orderCond {
	case "newest":
		db.Order("created_at DESC")
	case "score":
		ids, err = r.getPostIDsByScore(ctx, pageSize, pageNum)
		if err != nil {
			return paginator.Page[postentity.Post]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		db.Where("post_id IN ?", ids)
	default:
		db.Order("created_at DESC")
	}

	postPage := paginator.Page[postentity.Post]{CurrentPage: pageNum, PageSize: pageSize}
	err = postPage.SelectPages(db)
	if err != nil {
		return paginator.Page[postentity.Post]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if orderCond == "score" {
		data := make([]postentity.Post, 0, len(postPage.Data))
		postMap := make(map[int64]postentity.Post, len(postPage.Data))
		for _, post := range postPage.Data {
			postMap[post.PostID] = post
		}

		for _, id := range ids {
			idTmp, _ := strconv.Atoi(id)
			if post, ok := postMap[int64(idTmp)]; ok {
				data = append(data, post)
			}
		}

		postPage.Data = data
	}

	return postPage, nil
}

func (r *postRepo) LikePost(ctx context.Context, postID int64, userID int64, direction int8) error {
	_, err := r.getPostLikeNumber(ctx, postID)
	if err != nil {
		if err == redis.Nil {
			return errorx.NotFound(reason.PostNotFound)
		}
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	liked, exist, err := r.checkLike(ctx, postID, userID)
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	switch direction {
	case 1:
		if exist && liked == 1 {
			return errorx.BadRequest(reason.LikeAlreadyExist)
		}
	case 2:
		if !exist || liked == 0 {
			return errorx.BadRequest(reason.LikeCancelFailToNotLiked)
		}
	}

	err = r.doLike(ctx, postID, userID, direction)
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
