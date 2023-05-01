package controller

import (
	"api/common"
	"api/dto"
	"api/model"
	"api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Publish 发布帖子
func Publish(ctx *gin.Context) {
	DB := common.GetDB()

	user, _ := ctx.Get("user")
	userData := user.(model.User)

	//获取参数
	content := ctx.PostForm("content")
	urlString := ctx.PostForm("urls")
	//数据验证
	if len(content) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "内容不能为空")
		return
	}

	if utf8.RuneCountInString(content) > 3000 {

		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "内容不能超过3000字")
		return
	}
	//随机生成 nodeId
	noteId := time.Now().Unix()
	newNote := model.Note{
		NoteId:  "N" + strconv.FormatInt(noteId, 10),
		Name:    userData.Name,
		Content: content,
	}
	DB.Create(&newNote)

	if urlString != "" {
		//保存图片路径
		urls := strings.Split(urlString, ",")
		saveUrls(urls, newNote.NoteId)
	}
	DB.Model(&userData).Update("notes_count", userData.NotesCount+1)
	//返回结果
	response.Success(ctx, gin.H{"node_id": noteId}, "发布成功")
}

// LocalTimeLine 本地时间线
func LocalTimeLine(ctx *gin.Context) {
	DB := common.GetDB()
	user, _ := ctx.Get("user")
	userN := user.(model.User)
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	offset, _ := strconv.Atoi(ctx.Query("offset"))
	username := ctx.Query("username")
	var notes []model.Note
	if username != "" {
		DB.Where("name = ?", username).Limit(limit).Offset(offset).Order("created_at desc").Find(&notes)
	} else {
		DB.Limit(limit).Offset(offset).Order("created_at desc").Find(&notes)
	}

	if len(notes) == 0 {
		response.Fail(ctx, nil, "无更多数据")
		return
	}
	var notesDto []dto.NoteDto
	//var usersDto []dto.UserDto
	for i := 0; i < len(notes); i++ {
		var user model.User
		DB.Where("name = ?", notes[i].Name).First(&user)
		//usersDto = append(usersDto, dto.ToUserDto(user))
		var event model.Event
		DB.Where("target_id = ? AND type = ? AND originator_name = ?", notes[i].NoteId, "like", userN.Name).Find(&event)

		picturesDto := getUrls(notes[i].NoteId)

		if event.ID != 0 {
			notesDto = append(notesDto, dto.ToNoteDto(notes[i], dto.ToUserDto(user), picturesDto, 1))
		} else {
			notesDto = append(notesDto, dto.ToNoteDto(notes[i], dto.ToUserDto(user), picturesDto, 0))
		}

	}

	response.Success(ctx, gin.H{"notes": notesDto}, "获取帖子信息成功")
}

func GetNoteById(ctx *gin.Context) {
	DB := common.GetDB()
	user, _ := ctx.Get("user")
	userN := user.(model.User)
	noteId := ctx.Query("note_id")
	if noteId == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "参数不能为空")
		return
	}
	var note model.Note
	DB.Where("note_id = ?", noteId).Find(&note)
	if note.ID == 0 {
		response.Fail(ctx, nil, "查无此帖")
		return
	}
	var userM model.User
	DB.Where("name = ?", note.Name).First(&userM)
	var event model.Event
	DB.Where("target_id = ? AND type = ? AND originator_name = ?", note.NoteId, "like", userN.Name).Find(&event)
	var noteDto dto.NoteDto
	picturesDto := getUrls(note.NoteId)
	if event.ID != 0 {
		noteDto = dto.ToNoteDto(note, dto.ToUserDto(userM), picturesDto, 1)
	} else {
		noteDto = dto.ToNoteDto(note, dto.ToUserDto(userM), picturesDto, 0)
	}
	response.Success(ctx, gin.H{"note": noteDto}, "获取帖子信息成功")
}

func Search(ctx *gin.Context) {
	DB := common.GetDB()
	key := ctx.Query("key")
	if key == "" {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "key不能为空")
		return
	}
	var notes []model.Note
	DB.Where("content like ?", "%"+key+"%").Order("created_at desc").Find(&notes)
	var notesDto []dto.NoteDto
	for i := 0; i < len(notes); i++ {
		var user model.User
		DB.Where("name = ?", notes[i].Name).Find(&user)
		notesDto = append(notesDto, dto.ToNoteDto2(notes[i], dto.ToUserDto(user)))
	}
	response.Success(ctx, gin.H{"notes": notesDto}, "搜索成功")
}

func DeleteNote(ctx *gin.Context) {
	DB := common.GetDB()
	NoteId := ctx.Query("note_id")
	if len(NoteId) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "note_id不能为空")
		return
	}
	var Note model.Note
	DB.Where("note_id = ? ", NoteId).Find(&Note)
	if Note.ID == 0 {
		response.Fail(ctx, nil, "查无此帖")
		return
	}
	var User model.User
	DB.Where("name = ?", Note.Name).Find(&User)
	user, _ := ctx.Get("user")
	userM := user.(model.User)
	if user.(model.User).ID != User.ID {
		println(user.(model.User).ID)
		println(User.ID)
		ctx.JSON(http.StatusOK, gin.H{"code": 401, "msg": "权限不足"})
		return
	}

	if Note.ID != 0 {
		//删除帖子
		DB.Unscoped().Delete(&Note)
		//更新帖子数
		DB.Model(&userM).Update("notes_count", userM.NotesCount-1)
		//删除与该帖子相关的事件
		var Events []model.Event
		DB.Where("target_id = ?", NoteId).Find(&Events)
		if len(Events) != 0 {
			println(&Events[0].ID)
			DB.Unscoped().Delete(&Events)
		}
		//删除该帖子的图片
		var Pictures []model.Picture
		DB.Where("note_id = ?", NoteId).Find(&Pictures)
		if len(Pictures) != 0 {
			println(&Pictures[0].ID)
			DB.Unscoped().Delete(&Pictures)
		}

		response.Success(ctx, gin.H{"note_id": NoteId}, "删除帖子成功")
		return
	}
}
