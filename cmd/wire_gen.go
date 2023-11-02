// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"go.uber.org/zap"
	"harmoni/internal/conf"
	"harmoni/internal/cron"
	"harmoni/internal/handler"
	"harmoni/internal/infrastructure/data"
	"harmoni/internal/pkg/filesystem"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/snowflakex"
	"harmoni/internal/repository/auth"
	"harmoni/internal/repository/comment"
	"harmoni/internal/repository/email"
	file2 "harmoni/internal/repository/file"
	"harmoni/internal/repository/follow"
	"harmoni/internal/repository/like"
	"harmoni/internal/repository/post"
	"harmoni/internal/repository/tag"
	"harmoni/internal/repository/unique"
	"harmoni/internal/repository/user"
	"harmoni/internal/router"
	"harmoni/internal/server"
	"harmoni/internal/service"
	"harmoni/internal/usecase"
	"harmoni/internal/usecase/file"
)

// Injectors from wire.go:

func initApplication(appConf *conf.App, dbconf *conf.DB, rdbconf *conf.Redis, authConf *conf.Auth, emailConf *conf.Email, messageConf *conf.MessageQueue, fileConf *conf.FileStorage, logConf *conf.Log) (*Application, func(), error) {
	zapLogger, err := logger.NewZapLogger(logConf)
	if err != nil {
		return nil, nil, err
	}
	client, cleanup, err := data.NewRedis(rdbconf)
	if err != nil {
		return nil, nil, err
	}
	sugaredLogger := sugar(zapLogger)
	authRepo := auth.NewAuthRepo(client, sugaredLogger)
	authUseCase := usecase.NewAuthUseCase(authConf, authRepo, sugaredLogger)
	db, cleanup2, err := data.NewDB(dbconf, zapLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	node, err := snowflakex.NewSnowflakeNode(appConf)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	uniqueIDRepo := unique.NewUniqueIDRepo(node)
	userRepo := user.NewUserRepo(db, client, uniqueIDRepo, sugaredLogger)
	emailRepo := email.NewEmailRepo(client)
	emailUsecase, err := usecase.NewEmailUsecase(emailConf, emailRepo, sugaredLogger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	likeRepo := like.NewLikeRepo(db, client, userRepo, sugaredLogger)
	policy := file.NewPolicy(fileConf)
	fileSystem, err := filesystem.NewFileSystem(policy, client)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	fileRepo := file2.NewFileRepository(db, client, uniqueIDRepo, sugaredLogger)
	fileUseCase := file.NewFileUseCase(appConf, client, fileConf, fileSystem, fileRepo, sugaredLogger)
	userUseCase := usecase.NewUserUseCase(likeRepo, userRepo, authUseCase, fileUseCase, sugaredLogger)
	accountUsecase := usecase.NewAccountUsecase(authUseCase, userRepo, emailUsecase, userUseCase, zapLogger)
	accountService := service.NewAccountService(accountUsecase, sugaredLogger)
	jwtAuthMiddleware := middleware.NewJwtAuthMiddleware(authUseCase)
	accountHandler := handler.NewAccountHandler(accountService, jwtAuthMiddleware, sugaredLogger)
	tagRepo := tag.NewTagRepo(db, client, uniqueIDRepo, sugaredLogger)
	followRepo := follow.NewFollowRepo(db, userRepo, tagRepo, uniqueIDRepo, sugaredLogger)
	tagUseCase := usecase.NewTagUseCase(tagRepo, sugaredLogger)
	followUseCase := usecase.NewFollowUseCase(followRepo, userUseCase, tagUseCase, sugaredLogger)
	followService := service.NewFollowService(followUseCase, userUseCase, sugaredLogger)
	followHandler := handler.NewFollowHandler(followService, sugaredLogger)
	fileService := service.NewFileService(fileUseCase, userUseCase, sugaredLogger)
	fileHandler := handler.NewFileHandler(fileService, sugaredLogger)
	userService := service.NewUserService(userUseCase, authUseCase, accountUsecase, sugaredLogger)
	userHandler := handler.NewUserHandler(userService)
	postRepo := post.NewPostRepo(db, client, tagRepo, uniqueIDRepo, sugaredLogger)
	postUseCase := usecase.NewPostUseCase(postRepo, likeRepo, userUseCase, tagUseCase, sugaredLogger)
	postService := service.NewPostService(postUseCase, tagUseCase, sugaredLogger)
	postHandler := handler.NewPostHandler(postService)
	tagService := service.NewTagService(tagUseCase, sugaredLogger)
	tagHandler := handler.NewTagHandler(tagService)
	commentRepo := comment.NewCommentRepo(db, client, uniqueIDRepo, sugaredLogger)
	likeUsecase, cleanup3, err := usecase.NewLikeUsecase(messageConf, likeRepo, postRepo, postUseCase, commentRepo, userRepo, sugaredLogger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	commentUseCase := usecase.NewCommentUseCase(commentRepo, userRepo, fileUseCase, likeUsecase, sugaredLogger)
	commentService := service.NewCommentService(commentUseCase, sugaredLogger)
	commentHandler := handler.NewCommentHandler(commentService)
	likeService := service.NewLikeUsecase(likeUsecase, sugaredLogger)
	likeHandler := handler.NewLikeHandler(likeService, sugaredLogger)
	timeLinePullUsecase := usecase.NewTimeLineUsecase(followRepo, postUseCase, sugaredLogger)
	timeLineService := service.NewTimeLineService(timeLinePullUsecase, sugaredLogger)
	timeLineHandler := handler.NewTimeLineHandler(timeLineService)
	harmoniAPIRouter := router.NewHarmoniAPIRouter(accountHandler, followHandler, fileHandler, userHandler, postHandler, tagHandler, commentHandler, likeHandler, timeLineHandler)
	app := server.NewHTTPServer(appConf, zapLogger, harmoniAPIRouter, jwtAuthMiddleware)
	scheduledTaskManager, cleanup4, err := cron.NewScheduledTaskManager(messageConf, likeUsecase, sugaredLogger)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	application := NewApplication(app, scheduledTaskManager)
	return application, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

func sugar(l *zap.Logger) *zap.SugaredLogger {
	return l.Sugar()
}
