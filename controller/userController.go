package controller

import (
	"fmt"
	"gin-use-demo/common"
	"gin-use-demo/dto"
	"gin-use-demo/model"
	"gin-use-demo/response"
	"gin-use-demo/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// UserController 才能调用controller下的方法 且注册service参数
type UserController struct {
	UserService service.IUserService
}

// NewUserController 注册service层的数据内容,能够调用service中的接口方法
func NewUserController() UserController {
	userController := UserController{UserService: service.NewUserService()}
	userController.UserService.(service.UserService).DB.AutoMigrate(model.User{})
	return userController
}

// Login 登录
func (u UserController) Login(ctx *gin.Context) {
	var loginUser dto.UserDto
	ctx.Bind(&loginUser)

	// 获取参数
	name := loginUser.Name
	password := loginUser.Password
	number := loginUser.Number

	// 数据验证
	if len(number) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 调用方法
	user, err := u.UserService.Login(name, number)

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 校验密码
	if user.Password != password {
		response.Fail(ctx, "密码不正确,请重新输入", nil)
		return
	}

	//  校验通过,发放token
	token, err := common.ReleaseToken(*user)
	if err != nil {
		response.Response(ctx, 500, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// Register 注册
func (u UserController) Register(ctx *gin.Context) {
	var user model.User
	ctx.Bind(&user)

	// 注册的用户设置为普通角色
	if user.Role == "" {
		user.Role = strconv.Itoa(1)
	}
	fmt.Println(user)

	// 数据验证
	if len(user.Number) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(user.Password) < 6 && len(user.Password) > 12 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位且大于12位")
		return
	}

	// 判断手机号或者用户名是否存在
	userInfo, err := u.UserService.AddUser(&user)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, gin.H{"user": userInfo}, "注册成功！")

}

// UserInfo 获取用户个人信息
func (u UserController) UserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user == nil {
		response.Fail(ctx, "用户不存在", nil)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserInfoDto(user.(model.User))}})
}

// UserInquire 查询参数
type UserInquire struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

// UserList 获取用户信息列表
func (u UserController) UserList(ctx *gin.Context) {
	// 查询参数
	// 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "3"))
	name := ctx.DefaultQuery("name", "")
	number := ctx.DefaultQuery("number", "")

	userList, err := u.UserService.GetUserPageList(pageNum, pageSize, name, number)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, gin.H{
		"userList": userList,
	}, "查询成功")

}

// DeleteUser 删除一个用户
func (u UserController) DeleteUser(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")

	if !ok {
		response.Fail(ctx, "删除失败!", gin.H{"error": "无效的id"})
		return
	}
	err := u.UserService.DeleteUser(id)
	if err != nil {
		response.Fail(ctx, err.Error(), gin.H{
			"error": "删除失败",
		})
		return
	}
	response.Success(ctx, nil, "删除成功!")

}

// GetUserInfoById 查询一个用户
func (u UserController) GetUserInfoById(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")

	if !ok {
		response.Fail(ctx, "查询失败!", gin.H{"error": "无效的id"})
		return
	}
	userInfo, err := u.UserService.GetUserInfos(id)
	if err != nil {
		response.Fail(ctx, err.Error(), gin.H{
			"error": "查询失败",
		})
		return
	}
	response.Success(ctx, gin.H{
		"userInfo": userInfo,
	}, "查询成功!")

}

func (u UserController) AddUser(ctx *gin.Context) {
	var user model.User
	ctx.ShouldBind(&user)

	// 注册的用户设置为普通角色
	if user.Role == "" {
		user.Role = strconv.Itoa(1)
	}
	fmt.Println(user)

	// 数据验证
	if len(user.Number) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(user.Password) < 6 && len(user.Password) > 12 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位且大于12位")
		return
	}

	// 判断手机号或者用户名是否存在
	userInfo, err := u.UserService.AddUser(&user)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, gin.H{"user": userInfo}, "注册成功！")

}

func (u UserController) UpdateUser(ctx *gin.Context) {
	// 查询数据库内的用户数据
	id, ok := ctx.Params.Get("id")
	if !ok {
		response.Fail(ctx, "删除失败!", gin.H{"error": "无效的id"})
		return
	}
	userInfo, err := u.UserService.GetUserInfos(id)
	if err != nil {
		response.Fail(ctx, err.Error(), gin.H{
			"error": "查询失败",
		})
		return
	}
	var updatedUserInfo *dto.UpdateUser
	if err := ctx.ShouldBindJSON(&updatedUserInfo); err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	if updatedUserInfo.Password == "" {
		updatedUserInfo.Password = userInfo.Password
	}
	if updatedUserInfo.Address == "" {
		updatedUserInfo.Address = userInfo.Address
	}
	if updatedUserInfo.Tag == "" {
		updatedUserInfo.Tag = userInfo.Tag
	}
	if updatedUserInfo.Role == "" {
		updatedUserInfo.Role = userInfo.Role
	}

	err = u.UserService.UpdateUser(id, updatedUserInfo)
	if err != nil {
		response.Fail(ctx, err.Error(), gin.H{
			"error": "更新失败",
		})
		return
	}

	response.Success(ctx, nil, "更新成功")
}
