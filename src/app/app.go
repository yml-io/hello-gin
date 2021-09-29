package app

import (
	"hello-gin/src/middleware/auth"

	favoriteR "hello-gin/src/app/services/favorite"
	followR "hello-gin/src/app/services/follow"
	postR "hello-gin/src/app/services/post"
	userR "hello-gin/src/app/services/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	r = setupPublicRouter(r)
	r = setupAuthRouter(r)
	return r
}

// 公共业务路由处理逻辑
func setupPublicRouter(r *gin.Engine) *gin.Engine {
	unAuthorized := r.Group("/", auth.AuthRequired(r))

	userRouter := unAuthorized.Group("/user")
	{
		userRouter.GET("/view", userR.GetUserInfo())
		userRouter.POST("/update", userR.UpdateUserInfo())
		userRouter.POST("/upload_avatar", userR.UploadAvatar())
		// need check admin
		userRouter.PUT("/create", userR.CreateUser())
		userRouter.DELETE("/delete/:uid", userR.DeleteUser())
	}
	postsRouter := unAuthorized.Group("/posts")
	{
		postsRouter.GET("/view/:pid", postR.GetPostInfo())
		postsRouter.GET("/list", postR.GetPostList())
		postsRouter.PUT("/create", postR.CreatePost())
		postsRouter.DELETE("/delete/:pid", postR.DeletePost())
		postsRouter.POST("/update/:pid", postR.UpdatePostInfo())
	}
	followRouter := unAuthorized.Group("/follow")
	{
		followRouter.GET("/list", followR.GetFollowList())
		followRouter.PUT("/create", followR.CreateFollow())
		followRouter.DELETE("/delete/:uid", followR.DeleteFollow())
	}
	favoriteRouter := unAuthorized.Group("/favorite")
	{
		favoriteRouter.GET("/list", favoriteR.GetFavoriteList())
		favoriteRouter.PUT("/create", favoriteR.CreateFavorite())
		favoriteRouter.DELETE("/delete/:pid", favoriteR.DeleteFavorite())
	}

	// // get posts from follow user which latest posts
	// unAuthorized.GET("/recommend_post", GetRecommendPost())

	return r
}

// 受限业务路由处理逻辑
func setupAuthRouter(r *gin.Engine) *gin.Engine {
	// router 对象可以增加其他的 中间件，当和 group 子路由一起使用时，可以实现 auth 验证特权路由的效果
	authorized := r.Group("/admin")
	authorized.Use(auth.AdminRequired(r))

	// priorityUserRouter := authorized.Group("/user")
	// {
	// 	priorityUserRouter.GET("/list", ListAllUser())
	// }

	// priorityPostsRouter := authorized.Group("/posts")
	// {
	// 	priorityPostsRouter.GET("/list", ListAllPosts())
	// }

	return r
}
