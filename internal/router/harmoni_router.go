package router

import (
	"harmoni/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type HarmoniAPIRouter struct {
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	tagHandler     *handler.TagHandler
	commentHandler *handler.CommentHandler
}

func NewHarmoniAPIRouter(
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	tagHandler *handler.TagHandler,
	commentHandler *handler.CommentHandler,
) *HarmoniAPIRouter {
	return &HarmoniAPIRouter{
		userHandler:    userHandler,
		postHandler:    postHandler,
		tagHandler:     tagHandler,
		commentHandler: commentHandler,
	}
}

func (h *HarmoniAPIRouter) RegisterHarmoniAPIRouter(r fiber.Router) {
	// user
	r.Get("/user", h.userHandler.GetAllUsers)

	// tag
	r.Post("/tag", h.tagHandler.CreateTag)

	// post
	r.Post("/post", h.postHandler.CreatePost)
	r.Post("/post/like", h.postHandler.LikePost)

	// commnet
	r.Post("/comment", h.commentHandler.CreateComment)
}

func (h *HarmoniAPIRouter) RegisterUnAuthHarmoniAPIRouter(r fiber.Router) {
	// user
	r.Get("/user/:id", h.userHandler.GetUser)

	r.Post("/login", h.userHandler.Login)
	r.Post("/register", h.userHandler.Register)

	// tag
	r.Get("/tag", h.tagHandler.GetTags)
	r.Get("/tag/:id", h.tagHandler.GetTagByID)

	// post
	r.Get("/post", h.postHandler.GetPosts)
	r.Get("/post/:id", h.postHandler.GetPostDetail)

	// commnet
	r.Get("/comment", h.commentHandler.GetComments)
}
