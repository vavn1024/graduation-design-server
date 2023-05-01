package controller

import (
	"api/common"
	"api/dto"
	"api/model"
	"api/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
	"unicode/utf8"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()

	//使用结构体获取参数
	var requestUser = model.User{}
	_ = ctx.Bind(&requestUser)
	name := requestUser.Name
	password := requestUser.Password
	fmt.Println(requestUser.Name)
	//数据验证
	if len(name) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名不能为空")
		return
	}
	if len(password) < 7 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于7位")
		return
	}
	//判断用户名是否存在
	if isNameExist(DB, name) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名已存在")
		return
	}
	//创建用户
	hasePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Password: string(hasePassword),
	}
	DB.Create(&newUser)
	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")

}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//使用结构体获取参数
	var requestUser = model.User{}
	_ = ctx.Bind(&requestUser)
	name := requestUser.Name
	password := requestUser.Password
	//数据验证
	if len(name) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名不能为空")
		return
	}
	if len(password) < 7 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于7位")
		return
	}
	//判断用户名是否存在
	var user model.User
	DB.Where("name = ?", name).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// Show 查询用户信息
func Show(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	username := ctx.Query("username")
	originatorName := user.(model.User).Name
	if user.(model.User).Name != username && username != "" {
		DB := common.GetDB()
		var user model.User
		DB.Where("name = ?", username).Find(&user)
		if user.ID != 0 {
			var event model.Event
			DB.Where("target_id = ? AND type = ? AND originator_name = ?", username, "follow", originatorName).Find(&event)
			if event.ID != 0 {
				response.Success(ctx, gin.H{"user": dto.ToUserDto(user, 1)}, "获取用户信息成功")
			} else {
				response.Success(ctx, gin.H{"user": dto.ToUserDto(user, 0)}, "获取用户信息成功")
			}
			return
		}
		response.Response(ctx, http.StatusOK, 400, nil, "用户不存在")
		return
	}
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "获取用户信息成功")
}

//Update 更新用户信息
func Update(ctx *gin.Context) {
	DB := common.GetDB()
	//使用结构体获取参数
	var requestUser = model.User{}
	_ = ctx.Bind(&requestUser)
	nickname := requestUser.Nickname
	avatar := requestUser.Avatar
	banner := requestUser.Banner
	describe := requestUser.Describe

	//数据验证
	if utf8.RuneCountInString(nickname) > 10 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "昵称长度最多为10位")
		return
	}
	if utf8.RuneCountInString(describe) > 25 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "个人简介不能超过25字")
		return
	}
	user, _ := ctx.Get("user")
	userN := user.(model.User)
	DB.Model(&userN).Updates(model.User{Nickname: nickname, Avatar: avatar, Banner: banner, Describe: describe})
	response.Success(ctx, gin.H{"user": dto.ToUserDto(userN)}, "更新用户信息成功")
}

//ChangePassword 修改用户密码
func ChangePassword(ctx *gin.Context) {
	DB := common.GetDB()
	user, _ := ctx.Get("user")
	userN := user.(model.User)
	//获取参数
	password := ctx.PostForm("password")
	newPassword := ctx.PostForm("newPassword")

	//数据验证
	if len(newPassword) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	//判断用户名是否存在
	DB.Where("name = ?", userN.Name).First(&userN)
	if userN.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(userN.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	//创建新密码
	hasePassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	DB.Model(&userN).Updates(model.User{Password: string(hasePassword)})
	response.Success(ctx, gin.H{"user": dto.ToUserDto(userN)}, "更新用户信息成功")
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	if err != nil {
		c.JSON(http.StatusOK, fmt.Sprintf("'%s' uploaded", filename))
		return
	}
	filepath := path.Join("./"+"avatar", filename)
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"uploading": "done", "message": "success", "url": "http://" + c.Request.Host + "/avatar/" + filename})
}

func isNameExist(db *gorm.DB, name string) bool {
	var user model.User
	db.Where("name = ?", name).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
