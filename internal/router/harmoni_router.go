package router

import (
	"harmoni/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type HarmoniAPIRouter struct {
	accountHandler  *handler.AccountHandler
	followHandler   *handler.FollowHandler
	userHandler     *handler.UserHandler
	postHandler     *handler.PostHandler
	tagHandler      *handler.TagHandler
	commentHandler  *handler.CommentHandler
	likeHandler     *handler.LikeHandler
	timelineHandler *handler.TimeLineHandler
}

func NewHarmoniAPIRouter(
	accountHandler *handler.AccountHandler,
	followHandler *handler.FollowHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	tagHandler *handler.TagHandler,
	commentHandler *handler.CommentHandler,
	likeHandler *handler.LikeHandler,
	timelineHandler *handler.TimeLineHandler,
) *HarmoniAPIRouter {
	return &HarmoniAPIRouter{
		accountHandler:  accountHandler,
		followHandler:   followHandler,
		userHandler:     userHandler,
		postHandler:     postHandler,
		tagHandler:      tagHandler,
		commentHandler:  commentHandler,
		likeHandler:     likeHandler,
		timelineHandler: timelineHandler,
	}
}

func (h *HarmoniAPIRouter) RegisterHarmoniAPIRouter(r fiber.Router) {
	// user account
	account := r.Group("/account")
	account.Post("/email/change", h.accountHandler.ChangeEmail)
	account.Post("/password/change", h.accountHandler.ChangePassword)
	account.Post("/password/reset", h.accountHandler.ResetPassword)

	follow := r.Group("/follow")
	follow.Post("", h.followHandler.Follow)

	follow.Get("/followers", h.followHandler.GetFollowers)
	follow.Get("/followings", h.followHandler.GetFollowings)
	follow.Get("/isFollowing", h.followHandler.IsFollowing)
	follow.Get("/arefolloweachother", h.followHandler.AreFollowEachOther)

	like := r.Group("/like")
	like.Post("", h.likeHandler.Like)

	like.Get("/list", h.likeHandler.LikingList)
	like.Get("/isLiking", h.likeHandler.IsLiking)

	// user
	r.Get("/user", h.userHandler.GetUsers)

	r.Post("/logout", h.userHandler.Logout)

	// tag
	r.Post("/tag", h.tagHandler.CreateTag)

	// post
	r.Post("/post", h.postHandler.CreatePost)

	// timeline
	r.Get("/timeline", h.timelineHandler.GetUserTimeLine)
	r.Get("/timeline/home", h.timelineHandler.GetHomeTimeLine)

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
