package usecase

import (
	"context"
	likeentity "harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
	tagentity "harmoni/internal/entity/tag"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"

	"go.uber.org/zap"
)

type PostUseCase struct {
	postRepo    postentity.PostRepository
	likeRepo    likeentity.LikeRepository
	userUsecase *UserUseCase
	tagUsecase  *TagUseCase
	logger      *zap.SugaredLogger
}

func NewPostUseCase(
	postRepo postentity.PostRepository,
	likeRepo likeentity.LikeRepository,
	userUsecase *UserUseCase,
	tagUsecase *TagUseCase,
	logger *zap.SugaredLogger,
) *PostUseCase {
	return &PostUseCase{
		postRepo:    postRepo,
		tagUsecase:  tagUsecase,
		likeRepo:    likeRepo,
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (u *PostUseCase) setTagInfos(ctx context.Context, postInfo *postentity.PostBasicInfo, tags []tagentity.Tag) error {
	var err error

	if tags == nil {
		tags, err = u.tagUsecase.GetTagsByPostID(ctx, postInfo.PostID)
		if err != nil {
			return err
		}
	}

	tagInfos := make([]tagentity.TagInfo, len(tags))
	for i, tag := range tags {
		tagInfos[i] = tag.ToBasicInfo()
	}
	postInfo.Tags = tagInfos
	return nil
}

func (u *PostUseCase) setUser(ctx context.Context, postInfo *postentity.PostBasicInfo, userID int64) error {
	author, exist, err := u.userUsecase.GetBasicByUserID(ctx, userID)
	if err != nil {
		return err
	} else if !exist {
		author = &userentity.UserBasicInfo{
			UserID: -1,
			Name:   "deactivated",
		}
	}
	postInfo.User = author
	return nil
}

func (u *PostUseCase) MergeList(ctx context.Context, userID int64, postInfos []postentity.PostBasicInfo) ([]postentity.PostBasicInfo, error) {
	postIDs := make([]int64, len(postInfos))
	for i, post := range postInfos {
		postIDs[i] = post.PostID
	}
	likes, err := u.likeRepo.BatchLikeCountByIDs(ctx, postIDs, likeentity.LikePost)
	if err != nil {
		return nil, err
	}
	for i, postInfo := range postInfos {
		if userID != 0 {
			isLiked, err := u.likeRepo.IsLiking(ctx, &likeentity.Like{LikeType: likeentity.LikePost, UserID: userID, LikingID: postInfo.PostID})
			if err != nil {
				return nil, err
			}
			postInfo.Liked = isLiked
		}

		err = u.setTagInfos(ctx, &postInfos[i], nil)
		if err != nil {
			return nil, err
		}

		// TODO: batch users by userIDs
		err = u.setUser(ctx, &postInfos[i], postInfo.User.UserID)
		if err != nil {
			return nil, err
		}

		postInfo.LikeCount = likes[postInfo.PostID]
	}

	return postInfos, nil
}

func (u *PostUseCase) Merge(ctx context.Context, userID int64, postInfo *postentity.PostInfo, tags []tagentity.Tag) error {
	var err error

	likecount, existCount, err := u.likeRepo.LikeCount(ctx, &likeentity.Like{LikeType: likeentity.LikePost, LikingID: postInfo.PostID})
	if err != nil {
		return err
	} else if existCount {
		postInfo.LikeCount = likecount
	}

	if userID != 0 {
		isLiked, err := u.likeRepo.IsLiking(ctx, &likeentity.Like{LikeType: likeentity.LikePost, UserID: userID, LikingID: postInfo.PostID})
		if err != nil {
			return err
		}
		postInfo.Liked = isLiked
	}

	return u.setUser(ctx, &postInfo.PostBasicInfo, postInfo.User.UserID)
}

func (u *PostUseCase) Create(ctx context.Context, post *postentity.Post) (*postentity.PostInfo, error) {
	var (
		tags []tagentity.Tag
		err  error
	)
	if len(post.TagIDs) != 0 {
		tags, err = u.tagUsecase.GetByTagIDs(ctx, post.TagIDs)
		if err != nil {
			return nil, err
		}

		tagIDMap := map[int64]bool{}
		for _, tag := range tags {
			tagIDMap[tag.TagID] = true
		}

		for _, tagID := range post.TagIDs {
			if !tagIDMap[tagID] {
				return nil, errorx.BadRequest(reason.TagNotFound)
			}
		}
	}

	err = u.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	postInfo := post.ToInfo()
	err = u.Merge(ctx, 0, &postInfo, tags)
	if err != nil {
		return nil, err
	}
	return &postInfo, err
}

func (u *PostUseCase) GetByPostID(ctx context.Context, userID int64, postID int64) (*postentity.PostInfo, bool, error) {
	post, exist, err := u.postRepo.GetByPostID(ctx, postID)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, false, nil
	}

	postInfo := post.ToInfo()
	err = u.Merge(ctx, userID, &postInfo, nil)
	if err != nil {
		return nil, false, err
	}
	return &postInfo, true, nil
}

func (u *PostUseCase) GetBasicInfoByPostID(ctx context.Context, postID int64) (*postentity.Post, bool, error) {
	return u.postRepo.GetBasicInfoByPostID(ctx, postID)
}

// func (u *PostUseCase) BatchByIDs(ctx context.Context, postIDs []int64) ([]Post, error)
func (u *PostUseCase) GetPage(ctx context.Context, queryCond *postentity.PostQuery) (*paginator.Page[postentity.PostBasicInfo], error) {
	posts, err := u.postRepo.GetPage(ctx, queryCond)
	if err != nil {
		return nil, err
	}

	postInfos := make([]postentity.PostBasicInfo, len(posts.Data))
	for i, v := range posts.Data {
		postInfos[i] = v.ToBasic()
	}

	postInfos, err = u.MergeList(ctx, queryCond.UserID, postInfos)
	if err != nil {
		return nil, err
	}

	return &paginator.Page[postentity.PostBasicInfo]{
		CurrentPage: posts.CurrentPage,
		PageSize:    posts.PageSize,
		Pages:       posts.Pages,
		Total:       posts.Total,
		Data:        postInfos,
	}, err
}

func (u *PostUseCase) GetLikeCount(ctx context.Context, postID int64) (int64, bool, error) {
	return u.postRepo.GetLikeCount(ctx, postID)
}

func (u *PostUseCase) UpdateLikeCount(ctx context.Context, postID int64, count int64) error {
	return u.postRepo.UpdateLikeCount(ctx, postID, count)
}

func (u *PostUseCase) BatchByIDs(ctx context.Context, postIDs []int64) ([]postentity.PostBasicInfo, error) {
	posts, err := u.postRepo.BatchByIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}
	postInfos := make([]postentity.PostBasicInfo, len(posts))
	for i, post := range posts {
		postInfos[i] = post.ToBasic()
	}

	return postInfos, nil
}
