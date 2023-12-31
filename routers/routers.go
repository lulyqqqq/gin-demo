package routers

import (
	"gin-use-demo/controller"
	"gin-use-demo/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {

	// 配置跨越中间件以及监测程序是否异常的中间件
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	userController := controller.NewUserController()

	// 登录接口
	r.POST("/login", userController.Login)
	// 注册接口
	r.POST("/register", userController.Register)
	// 获取用户列表信息
	r.GET("/user/list", userController.UserList)

	userRouters := r.Group("/user")
	// 配置token验证
	userRouters.Use(middleware.AuthMiddleware())
	{
		// 获取用户个人信息
		userRouters.GET("/info", userController.UserInfo)
		userRouters.GET("/:id", userController.GetUserInfoById)
		// 增加用户
		userRouters.POST("/add", userController.AddUser)
		// 删除用户
		userRouters.DELETE("/:id", userController.DeleteUser)
		// 修改用户信息
		userRouters.PUT("/:id", userController.UpdateUser)
	}
	return r
}
