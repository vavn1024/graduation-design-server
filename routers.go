package main

import (
	"api/controller"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	//r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/notes/publish", middleware.AuthMiddleware(), controller.Publish)
	r.GET("/api/notes/local-timeline", middleware.AuthMiddleware(), controller.LocalTimeLine)
	r.GET("/api/notes/getNoteById", middleware.AuthMiddleware(), controller.GetNoteById)
	r.GET("/api/notes/search", middleware.AuthMiddleware(), controller.Search)
	r.DELETE("/api/notes/delete", middleware.AuthMiddleware(), controller.DeleteNote)
	r.GET("/api/users/show", middleware.AuthMiddleware(), controller.Show)
	r.POST("/api/users/update", middleware.AuthMiddleware(), controller.Update)
	r.POST("/api/users/changePwd", middleware.AuthMiddleware(), controller.ChangePassword)
	r.POST("/api/users/upload", middleware.AuthMiddleware(), controller.UploadImage)
	r.Static("/avatar", "./avatar")
	r.GET("/api/events/like", middleware.AuthMiddleware(), controller.Like)
	r.GET("/api/events/follow", middleware.AuthMiddleware(), controller.Follow)
	r.GET("/api/events/following_list", middleware.AuthMiddleware(), controller.GetFollowingList)
	r.GET("/api/events/followed_list", middleware.AuthMiddleware(), controller.GetFollowedList)
	r.GET("/api/events/notifications", middleware.AuthMiddleware(), controller.GetNotifications)
	r.POST("/api/events/reply", middleware.AuthMiddleware(), controller.Reply)
	r.GET("/api/events/getReply", middleware.AuthMiddleware(), controller.GetReply)
	return r
}
