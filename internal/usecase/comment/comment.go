package comment

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	likeentity "harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/entity/user"
	"harmoni/internal/usecase/comment/events"
	"harmoni/internal/usecase/file"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSetComment = wire.NewSet(
	NewCommentUseCase,
	events.NewCommentEventsHandler,
)

type CommentUseCase struct {
	commentRepo commententity.CommentRepository
	likeRepo    likeentity.LikeRepository
	userRepo    user.UserRepository
	fileUsecase *file.FileUseCase
	logger      *zap.SugaredLogger
}

func NewCommentUseCase(
	commentRepo commententity.CommentRepository,
	likeRepo likeentity.LikeRepository,
	userRepo user.UserRepository,
	fileUsecase *file.FileUseCase,
	logger *zap.SugaredLogger) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		userRepo:    userRepo,
		fileUsecase: fileUsecase,
		logger:      logger,
	}
}

func (u *CommentUseCase) Create(ctx context.Context, comment *commententity.Comment) error {
	comment.EscapeContent()
	return u.commentRepo.Create(ctx, comment)
}

func (u *CommentUseCase) GetPage(ctx context.Context, commentQuery *commententity.CommentQuery) (*paginator.Page[*commententity.Comment], error) {
	comments, err := u.commentRepo.List(ctx, commentQuery)
	if err != nil {
		return nil, err
	}

	err = u.fillUsers(ctx, comments.Data)
	if err != nil {
		return nil, err
	}
	err = u.fillLikeCount(ctx, comments.Data)
	if err != nil {
		return nil, err
	}

	if commentQuery.RootID == 0 && len(comments.Data) > 0 {
		err = u.fillRootCommentsWithNSubComments(ctx, &comments)
		if err != nil {
			return nil, err
		}
	}
	return &comments, err
}

func (u *CommentUseCase) fillRootCommentsWithNSubComments(ctx context.Context,
	commentPage *paginator.Page[*commententity.Comment]) error {
	commentIDs := make([]int64, len(commentPage.Data))
	for i, comment := range commentPage.Data {
		commentIDs[i] = comment.CommentID
	}
	subComments, err := u.commentRepo.ListNSubComments(ctx, commentIDs)
	if err != nil {
		return err
	}
	err = u.fillUsers(ctx, subComments)
	if err != nil {
		return err
	}
	err = u.fillLikeCount(ctx, subComments)
	if err != nil {
		return err
	}
	subCommentMap := commententity.CommentList(subComments).ToRooIDMap()
	for i := range commentPage.Data {
		comment := commentPage.Data[i]
		comment.Children = append(comment.Children, subCommentMap[comment.CommentID]...)
	}

	return nil
}

func (u *CommentUseCase) fillUsers(ctx context.Context, comments []*commententity.Comment) error {
	userIDs := make([]int64, 0, len(comments))
	fileIDs := make([]int64, 0, len(comments))
	visitedUserIDs := make(map[int64]bool)
	visitedFileIDs := make(map[int64]bool)

	for _, comment := range comments {
		// Add author userIDs
		if !visitedUserIDs[comment.Author.UserID] {
			userIDs = append(userIDs, comment.Author.UserID)
			visitedUserIDs[comment.Author.UserID] = true
		}
		for _, toUser := range comment.ToMembers {
			// Add toUser userIDs
			if !visitedUserIDs[toUser.UserID] {
				userIDs = append(userIDs, toUser.UserID)
				visitedUserIDs[toUser.UserID] = true
			}
		}
	}
	users, err := u.userRepo.GetByUserIDs(ctx, userIDs)
	if err != nil {
		return err
	}
	userMap := make(map[int64]*user.User)
	for i, user := range users {
		userMap[user.UserID] = &users[i]
	}

	for _, user := range users {
		// Add author fileIDs
		if !visitedFileIDs[user.Avatar] {
			fileIDs = append(fileIDs, user.Avatar)
			visitedFileIDs[user.Avatar] = true
		}
	}
	files, err := u.fileUsecase.ListFileLinkMap(ctx, fileIDs)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		userTmp := userMap[comment.Author.UserID]
		if userTmp != nil {
			userBasic := userTmp.ToBasicInfo(files[userTmp.Avatar])
			comment.Author.Avatar = userBasic.Avatar
			comment.Author.Name = userBasic.Name
		} else {
			comment.Author.Avatar = ""
			comment.Author.Name = ""
		}

		for _, toUser := range comment.ToMembers {
			toUserTmp := userMap[toUser.UserID]
			if toUserTmp != nil {
				userBasic := toUserTmp.ToBasicInfo(files[toUserTmp.Avatar])
				toUser.Avatar = userBasic.Avatar
				toUser.Name = userBasic.Name
			} else {
				toUser.Avatar = ""
				toUser.Name = ""
			}
		}
	}

	return nil
}
func (u *CommentUseCase) fillLikeCount(ctx context.Context, comments []*commententity.Comment) error {
	commentIDs := make([]int64, len(comments))
	for i, comment := range comments {
		commentIDs[i] = comment.CommentID
	}

	likes, err := u.likeRepo.BatchLikeCountByIDs(ctx, commentIDs, likeentity.LikeComment)
	if err != nil {
		return err
	}
	for i := range comments {
		comments[i].LikeCount = likes[comments[i].CommentID]
	}

	return nil
}
