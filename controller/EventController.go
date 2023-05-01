package controller

import (
	"api/common"
	"api/dto"
	"api/model"
	"api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Follow 关注、取消关注
func Follow(ctx *gin.Context) {
	DB := common.GetDB()

	user, _ := ctx.Get("user")
	userData := user.(model.User)

	//获取参数
	name := ctx.Query("name")
	originatorName := userData.Name
	content := ctx.Query("content")
	Type := "follow"

	//验证参数
	if name == originatorName {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "您不能关注自己")
		return
	}
	var userM model.User
	DB.Where("name = ?", name).First(&userM)
	if userM.ID == 0 {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "此用户不存在")
		return
	}

	var Events model.Event
	DB.Where("originator_name = ? AND target_id = ? AND type = ?", originatorName, name, Type).Find(&Events)
	if Events.ID != 0 {
		DB.Unscoped().Delete(&Events)
		DB.Model(&userData).Update("following_count", userData.FollowingCount-1)
		DB.Model(&userM).Update("followed_count", userM.FollowedCount-1)
		response.Success(ctx, gin.H{"target_id": name}, "取消关注成功")
		return
	}

	Event := model.Event{
		TargetId:       name,
		OriginatorName: originatorName,
		TargetName:     userM.Name,
		Type:           Type,
		Content:        content,
	}
	DB.Create(&Event)
	DB.Model(&userData).Update("following_count", userData.FollowingCount+1)
	DB.Model(&userM).Update("followed_count", userM.FollowedCount+1)
	//返回结果
	response.Success(ctx, gin.H{"target_id": name}, "关注成功")
}

//GetFollowingList 获取关注中列表
func GetFollowingList(ctx *gin.Context) {
	DB := common.GetDB()
	originatorName := ctx.Query("name")

	//验证参数
	var userM model.User
	DB.Where("name = ?", originatorName).First(&userM)
	if userM.ID == 0 {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "此用户不存在")
		return
	}

	var events []model.Event
	DB.Where("originator_name = ? AND type = ?", originatorName, "follow").Find(&events)
	if len(events) == 0 {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "关注中列表为空")
		return
	}
	var userDto []dto.UserDto
	for i := 0; i < len(events); i++ {
		var user model.User
		DB.Where("name = ?", events[i].TargetId).Find(&user)
		userDto = append(userDto, dto.ToUserDto(user))
	}
	response.Success(ctx, gin.H{"following_list": userDto}, "获取关注中列表成功")
}

//GetFollowedList 获取关注者列表
func GetFollowedList(ctx *gin.Context) {
	DB := common.GetDB()
	targetId := ctx.Query("name")

	//验证参数
	var userM model.User
	DB.Where("name = ?", targetId).First(&userM)
	if userM.ID == 0 {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "此用户不存在")
		return
	}

	var events []model.Event
	DB.Where("target_id = ? AND type = ?", targetId, "follow").Find(&events)
	if len(events) == 0 {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "关注者列表为空")
		return
	}
	var userDto []dto.UserDto
	for i := 0; i < len(events); i++ {
		var user model.User
		DB.Where("name = ?", events[i].OriginatorName).Find(&user)
		userDto = append(userDto, dto.ToUserDto(user))
	}
	response.Success(ctx, gin.H{"followed_list": userDto}, "获取关注者列表成功")
}

// Like 点赞、取消点赞
func Like(ctx *gin.Context) {
	DB := common.GetDB()

	user, _ := ctx.Get("user")
	userData := user.(model.User)

	//获取参数
	noteId := ctx.Query("note_id")
	originatorName := userData.Name
	content := ctx.Query("content")
	Type := "like"
	var Events model.Event
	DB.Where("originator_name = ? AND target_id = ? AND type = ?", originatorName, noteId, Type).Find(&Events)
	if Events.ID != 0 {
		DB.Unscoped().Delete(&Events)
		response.Success(ctx, gin.H{"target_id": noteId}, "取消点赞成功")
		return
	}
	var note model.Note
	DB.Where("note_id = ?", noteId).Find(&note)

	Event := model.Event{
		TargetId:       noteId,
		OriginatorName: originatorName,
		TargetName:     note.Name,
		Type:           Type,
		Content:        content,
	}
	DB.Create(&Event)
	//DB.Model(&userData).Update("notes_count", userData.NotesCount+1)
	//返回结果
	response.Success(ctx, gin.H{"target_id": noteId}, "点赞成功")

}

//Reply 回复
func Reply(ctx *gin.Context) {
	DB := common.GetDB()

	user, _ := ctx.Get("user")
	userData := user.(model.User)

	//获取参数
	noteId := ctx.PostForm("note_id")
	originatorName := userData.Name
	content := ctx.PostForm("content")
	Type := "reply"

	if noteId == "" {
		return
	}

	var note model.Note
	DB.Where("note_id = ?", noteId).Find(&note)

	Event := model.Event{
		TargetId:       noteId,
		OriginatorName: originatorName,
		TargetName:     note.Name,
		Type:           Type,
		Content:        content,
	}
	DB.Create(&Event)
	//DB.Model(&userData).Update("notes_count", userData.NotesCount+1)
	//返回结果
	response.Success(ctx, gin.H{"event_id": Event.ID, "target_id": noteId, "content": content}, "评论成功")
}

//GetReply 获取评论内容
func GetReply(ctx *gin.Context) {
	DB := common.GetDB()
	noteId := ctx.Query("note_id")
	if noteId == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "请传入正确参数")
		return
	}
	var events []model.Event
	DB.Where("target_id = ? AND type = ?", noteId, "reply").Order("created_at desc").Find(&events)
	var eventsDto []dto.EventDto
	for i := 0; i < len(events); i++ {
		var user model.User
		DB.Where("name = ?", events[i].OriginatorName).First(&user)
		eventsDto = append(eventsDto, dto.ToEventDto(events[i], dto.ToUserDto(user)))
	}
	response.Success(ctx, gin.H{"events": eventsDto}, "获取评论信息成功")
}

func GetNotifications(ctx *gin.Context) {
	DB := common.DB
	user, _ := ctx.Get("user")
	userN := user.(model.User)
	var events []model.Event
	DB.Where("target_name = ? AND originator_name != ?", userN.Name, userN.Name).Order("created_at desc").Find(&events)
	var eventsDto []dto.EventDto
	for i := 0; i < len(events); i++ {
		var user model.User
		DB.Where("name = ?", events[i].OriginatorName).First(&user)
		eventsDto = append(eventsDto, dto.ToEventDto(events[i], dto.ToUserDto(user)))
	}
	response.Success(ctx, gin.H{"events": eventsDto}, "获取通知成功")
}
