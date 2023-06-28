// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"go.uber.org/zap"
	"harmoni/internal/conf"
	"harmoni/internal/cron"
	"harmoni/internal/data/mysql"
	"harmoni/internal/data/redis"
	"harmoni/internal/handler"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/snowflakex"
	"harmoni/internal/repository/auth"
	"harmoni/internal/repository/comment"
	"harmoni/internal/repository/email"
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
)

// Injectors from wire.go:

func initApplication(appConf *conf.App, dbconf *conf.DB, rdbconf *conf.Redis, authConf *conf.Auth, emailConf *conf.Email, messageConf *conf.MessageQueue, logConf *conf.Log) (*Application, func(), error) {
	zapLogger, err := logger.NewZapLogger(logConf)
	if err != nil {
		return nil, nil, err
	}
	client, cleanup, err := redis.NewRedis(rdbconf)
	if err != nil {
		return nil, nil, err
	}
	sugaredLogger := sugar(zapLogger)
	authRepo := auth.NewAuthRepo(client, sugaredLogger)
	authUseCase := usecase.NewAuthUseCase(authConf, authRepo, sugaredLogger)
	db, cleanup2, err := mysql.NewDB(dbconf, zapLogger)
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
	postRepo := post.NewPostRepo(db, client, uniqueIDRepo, sugaredLogger)
	commentRepo := comment.NewCommentRepo(db, client, uniqueIDRepo, sugaredLogger)
	likeUsecase, cleanup3, err := usecase.NewLikeUsecase(messageConf, likeRepo, postRepo, commentRepo, userRepo, sugaredLogger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	userUseCase := usecase.NewUserUseCase(userRepo, authUseCase, likeUsecase, sugaredLogger)
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
	userService := service.NewUserService(userUseCase, authUseCase, accountUsecase, sugaredLogger)
	userHandler := handler.NewUserHandler(userService)
	postUseCase := usecase.NewPostUseCase(postRepo, likeUsecase, sugaredLogger)
	postService := service.NewPostService(postUseCase, sugaredLogger)
	postHandler := handler.NewPostHandler(postService)
	tagService := service.NewTagService(tagUseCase, sugaredLogger)
	tagHandler := handler.NewTagHandler(tagService)
	commentUseCase := usecase.NewCommentUseCase(commentRepo, likeUsecase, sugaredLogger)
	commentService := service.NewCommentService(commentUseCase, sugaredLogger)
	commentHandler := handler.NewCommentHandler(commentService)
	likeService := service.NewLikeUsecase(likeUsecase, sugaredLogger)
	likeHandler := handler.NewLikeHandler(likeService, sugaredLogger)
	harmoniAPIRouter := router.NewHarmoniAPIRouter(accountHandler, followHandler, userHandler, postHandler, tagHandler, commentHandler, likeHandler)
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
