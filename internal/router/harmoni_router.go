package router

import (
	"harmoni/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type HarmoniAPIRouter struct {
	accountHandler *handler.AccountHandler
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	tagHandler     *handler.TagHandler
	commentHandler *handler.CommentHandler
}

func NewHarmoniAPIRouter(
	accountHandler *handler.AccountHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	tagHandler *handler.TagHandler,
	commentHandler *handler.CommentHandler,
) *HarmoniAPIRouter {
	return &HarmoniAPIRouter{
		accountHandler: accountHandler,
		userHandler:    userHandler,
		postHandler:    postHandler,
		tagHandler:     tagHandler,
		commentHandler: commentHandler,
	}
}

func (h *HarmoniAPIRouter) RegisterHarmoniAPIRouter(r fiber.Router) {
	// user account
	account := r.Group("/account")
	account.Post("/email/change", h.accountHandler.ChangeEmail)
	account.Post("/password/change", h.accountHandler.ChangePassword)
	account.Post("/password/reset", h.accountHandler.ResetPassword)

	// user
	r.Get("/user", h.userHandler.GetUsers)

	r.Post("/logout", h.userHandler.Logout)

	// tag
	r.Post("/tag", h.tagHandler.CreateTag)

	// post
	r.Post("/post", h.postHandler.CreatePost)
	r.Post("/post/like", h.postHandler.LikePost)

	// commnet
	r.Post("/comment", h.commentHandler.CreateComment)
}

func (h *HarmoniAPIRouter) RegisterUnAuthHarmoniAPIRouter(r fiber.Router) {
	// user account
	account := r.Group("/account")
	account.Post("/mail/send", h.accountHandler.MailSend)
	account.Post("/mail/check", h.accountHandler.MailCheck)

	// user
	r.Get("/user/:id", h.userHandler.GetUser)
	r.Post("/user/token/refresh", h.userHandler.RefreshToken)

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
