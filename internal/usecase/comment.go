package usecase

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	"harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/entity/user"
	"harmoni/internal/usecase/file"

	"go.uber.org/zap"
)

type CommentUseCase struct {
	commentRepo commententity.CommentRepository
	userRepo    user.UserRepository
	fileUsecase *file.FileUseCase
	likeUsecase *LikeUsecase
	logger      *zap.SugaredLogger
}

func NewCommentUseCase(
	commentRepo commententity.CommentRepository,
	userRepo user.UserRepository,
	fileUsecase *file.FileUseCase,
	likeUsecase *LikeUsecase,
	logger *zap.SugaredLogger) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		fileUsecase: fileUsecase,
		likeUsecase: likeUsecase,
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

	commentIDs := make([]int64, len(comments.Data))
	userIDs := make([]int64, 0, len(comments.Data))
	userIDMap := map[int64]bool{}
	for i, comment := range comments.Data {
		commentIDs[i] = comment.CommentID
		userID := comment.Author.UserID
		if !userIDMap[userID] {
			userIDs = append(userIDs, userID)
			userIDMap[userID] = true
		}
		for _, toUser := range comment.ToMembers {
			if userIDMap[toUser.UserID] {
				continue
			}
			userIDMap[toUser.UserID] = true
			userIDs = append(userIDs, toUser.UserID)
		}
	}

	users, err := u.userRepo.GetByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	userDetailMap := map[int64]*user.User{}
	userMap := map[int64]user.UserBasicInfo{}
	fileIDMap := map[int64]bool{}
	fileIDs := make([]int64, 0, len(users))
	for i := range users {
		user := users[i]
		userDetailMap[user.UserID] = &users[i]
	}
	for _, userID := range userIDs {
		var fileID int64
		user := userDetailMap[userID]
		if user != nil {
			fileID = user.Avatar
		}
		if !fileIDMap[fileID] {
			fileIDs = append(fileIDs, fileID)
		}
	}
	files, err := u.fileUsecase.ListFileLinkMap(ctx, fileIDs)
	if err != nil {
		return nil, err
	}

	for _, userID := range userIDs {
		userTmp := userDetailMap[userID]
		if userTmp == nil {
			userTmp = &user.User{UserID: userID}
		}
		userMap[userID] = userTmp.ToBasicInfo(files[userTmp.Avatar])
	}

	likes, err := u.likeUsecase.BatchLikeCountByIDs(ctx, commentIDs, like.LikeComment)
	if err != nil {
		return nil, err
	}
	for i := range comments.Data {
		comments.Data[i].LikeCount = likes[comments.Data[i].CommentID]
		userBasic := userMap[comments.Data[i].Author.UserID]
		comments.Data[i].Author = &userBasic

		members := comments.Data[i].ToMembers
		for j, member := range members {
			userBasic := userMap[member.UserID]
			members[j] = &userBasic
		}
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

	subCommentMap := commententity.CommentList(subComments).ToRooIDMap()
	for i := range commentPage.Data {
		comment := commentPage.Data[i]
		comment.Children = append(comment.Children, subCommentMap[comment.CommentID]...)
	}

	return nil
}
