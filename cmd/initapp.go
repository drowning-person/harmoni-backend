package main

import (
	"harmoni/internal/conf"
	datamysql "harmoni/internal/data/mysql"
	dataredis "harmoni/internal/data/redis"
	"harmoni/internal/handler"
	"harmoni/internal/pkg/logger"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/snowflakex"
	"harmoni/internal/pkg/validator"
	authrepo "harmoni/internal/repository/auth"
	commentrepo "harmoni/internal/repository/comment"
	emailrepo "harmoni/internal/repository/email"
	followrepo "harmoni/internal/repository/follow"
	postrepo "harmoni/internal/repository/post"
	tagrepo "harmoni/internal/repository/tag"
	uniquerepo "harmoni/internal/repository/unique"
	userrepo "harmoni/internal/repository/user"
	"harmoni/internal/router"
	"harmoni/internal/server"
	"harmoni/internal/service"
	"harmoni/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

func initApp(appConf *conf.App, dbconf *conf.DB, rdbconf *conf.Redis, authConf *conf.Auth, emailConf *conf.Email, logConf *conf.Log) (*fiber.App, error) {
	if err := validator.InitTrans("zh"); err != nil {
		return nil, err
	}

	logger, err := logger.NewZapLogger(logConf)
	if err != nil {
		return nil, err
	}

	rdb, err := dataredis.NewRedis(rdbconf)
	if err != nil {
		return nil, err
	}

	db, err := datamysql.NewDB(dbconf, logger)
	if err != nil {
		return nil, err
	}

	node, err := snowflakex.NewSnowflakeNode(time.Now().Format("2006-01-02"), 1)
	if err != nil {
		return nil, err
	}

	uniqueIDRepo := uniquerepo.NewUniqueIDRepo(node)

	emailRepo := emailrepo.NewEmailRepo(rdb)
	emailUsecase, err := usecase.NewEmailUsecase(emailConf, emailRepo, logger.Sugar())
	if err != nil {
		return nil, err
	}

	authRepo := authrepo.NewAuthRepo(rdb, logger.Sugar())
	authUsecase := usecase.NewAuthUseCase(authConf, authRepo, logger.Sugar())
	authMiddleware := middleware.NewJwtAuthMiddleware(authConf.Secret, authUsecase)

	userRepo := userrepo.NewUserRepo(db, rdb, uniqueIDRepo, logger.Sugar())
	userUsecase := usecase.NewUserUseCase(userRepo, authUsecase, logger.Sugar())
	accountUsecase := usecase.NewAccountUsecase(authUsecase, userRepo, emailUsecase, userUsecase, logger)
	userService := service.NewUserService(userUsecase, authUsecase, accountUsecase, logger.Sugar())
	userHandler := handler.NewUserHandler(userService)

	accountService := service.NewAccountService(accountUsecase, logger.Sugar())
	accountHandler := handler.NewAccountHandler(accountService, authMiddleware, logger.Sugar())

	tagRepo := tagrepo.NewTagRepo(db, rdb, uniqueIDRepo, logger.Sugar())
	tagUsecase := usecase.NewTagUseCase(tagRepo, logger.Sugar())
	tagServie := service.NewTagService(tagUsecase, logger.Sugar())
	tagHanlder := handler.NewTagHandler(tagServie)

	postRepo := postrepo.NewPostRepo(db, rdb, uniqueIDRepo, logger.Sugar())
	postUsecase := usecase.NewPostUseCase(postRepo, logger.Sugar())
	postServie := service.NewPostService(postUsecase, logger.Sugar())
	postHanlder := handler.NewPostHandler(postServie)

	commentRepo := commentrepo.NewCommentRepo(db, rdb, uniqueIDRepo, logger.Sugar())
	commentUsecase := usecase.NewCommentUseCase(commentRepo, logger.Sugar())
	commentServie := service.NewCommentService(commentUsecase, logger.Sugar())
	commentHanlder := handler.NewCommentHandler(commentServie)

	followRepo := followrepo.NewFollowRepo(db, userRepo, tagRepo, uniqueIDRepo, logger.Sugar())
	followUsecase := usecase.NewFollowUseCase(followRepo, userUsecase, tagUsecase, logger.Sugar())
	followService := service.NewFollowService(followUsecase, userUsecase, logger.Sugar())
	followHandler := handler.NewFollowHandler(followService, logger.Sugar())

	hrouter := router.NewHarmoniAPIRouter(accountHandler, followHandler, userHandler, postHanlder, tagHanlder, commentHanlder)

	app := server.NewHTTPServer(appConf.Debug, logger, hrouter, authMiddleware)
	return app, nil
}
