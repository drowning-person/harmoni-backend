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
	commentrepo "harmoni/internal/repository/comment"
	emailrepo "harmoni/internal/repository/email"
	postrepo "harmoni/internal/repository/post"
	tagrepo "harmoni/internal/repository/tag"
	uniquerepo "harmoni/internal/repository/unique"
	userrepo "harmoni/internal/repository/user"
	"harmoni/internal/router"
	"harmoni/internal/server"
	"harmoni/internal/service"
	"harmoni/internal/usecase"

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

	node, err := snowflakex.NewSnowflakeNode("2022-07-31", 1)
	if err != nil {
		return nil, err
	}

	uniqueIDRepo := uniquerepo.NewUniqueIDRepo(node)

	emailRepo := emailrepo.NewEmailRepo(rdb)
	emailUsecase, err := usecase.NewEmailUsecase(emailConf, emailRepo, logger.Sugar())
	if err != nil {
		return nil, err
	}

	userRepo := userrepo.NewUserRepo(db, rdb, uniqueIDRepo, logger.Sugar())
	authUsecase := usecase.NewAuthUseCase(authConf, userRepo, logger.Sugar())
	userUsecase := usecase.NewUserUseCase(userRepo, authUsecase, emailUsecase, logger.Sugar())
	userService := service.NewUserService(userUsecase, authUsecase, logger.Sugar())
	userHandler := handler.NewUserHandler(userService)

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

	authMiddleware := middleware.NewJwtAuthMiddleware(authConf.Secret, authUsecase)
	hrouter := router.NewHarmoniAPIRouter(userHandler, postHanlder, tagHanlder, commentHanlder)

	app := server.NewHTTPServer(appConf.Debug, logger, hrouter, authMiddleware)
	return app, nil
}
