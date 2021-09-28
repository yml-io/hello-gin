package app

import (
	"hello-gin/src/middleware/auth"

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
		userRouter.PUT("/create", userR.CreateUser())
	}
	// postsRouter := unAuthorized.Group("/posts")
	// {
	// 	postsRouter.GET("/view/:post_id", GetPostInfo())
	// 	postsRouter.GET("/list/:user_id", GetPostList())
	// 	postsRouter.PUT("/create", CreatePost())
	// 	postsRouter.DELETE("/:post_id", DeletePost())
	// 	postsRouter.POST("/update/:post_id", UpdatePostInfo())
	// }
	// followRouter := unAuthorized.Group("/follow")
	// {
	// 	followRouter.GET("/list", GetFollowList())
	// 	followRouter.PUT("/create/:user_id", CreateFollow())
	// 	followRouter.DELETE("/:user_id", DeleteFollow())
	// }
	// favoriteRouter := unAuthorized.Group("/favorite")
	// {
	// 	favoriteRouter.GET("/list", GetFavoriteList())
	// 	favoriteRouter.PUT("/create/:user_id", CreateFavorite())
	// 	favoriteRouter.DELETE("/:user_id", DeleteFavorite())
	// }

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
