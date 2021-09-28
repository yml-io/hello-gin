package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupCommonRouter(r *gin.Engine) (e *gin.Engine) {
	r = setupBasicRouter(r)
	r = setupStaticResourceRouter(r)
	return r
}

func setupStaticResourceRouter(r *gin.Engine) *gin.Engine {
	// served 静态文件
	r.Static("/static", "./src/static")
	// router.StaticFS("/more_static", http.Dir("my_file_system"))
	// r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	// 或者使用 gin.Context 的成员方法
	// c.File("local/file.go")
	// var fs http.FileSystem = // ...
	// c.FileFromFS("fs/file.go", fs)
	return r
}

func setupBasicRouter(r *gin.Engine) *gin.Engine {
	// TODO 读取配置文件
	getAppVersion := func(c *gin.Context) {
		// gin.H is a shortcut for map[string]interface{}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": "1.0",
		})
	}
	r.GET("/version", getAppVersion)
	return r
}
