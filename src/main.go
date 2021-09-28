package main

import (
	application "hello-gin/src/app"
	cm "hello-gin/src/middleware/common"
	"hello-gin/src/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.Default()
	r := gin.New()
	r = cm.SetupCommonMiddleware(r)
	r = router.SetupCommonRouter(r)

	r = application.SetupRouter(r)
	r.Run()
}
