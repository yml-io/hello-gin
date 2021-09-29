package middleware

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupCommonMiddleware(r *gin.Engine) *gin.Engine {
	r.Use(Interceptor())
	setupLoggerMiddleware(r)
	setupRecoveryMiddleware(r)
	return r
}

func setupLoggerMiddleware(r *gin.Engine) *gin.Engine {
	// write log to file
	fileName := time.Now().Format("2006010215")
	logPath := filepath.Join("log", fileName)
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logPath, 0777); err != nil {
			fmt.Println("Only write log to stdio, Can not create log file: ", err)
			return r
		}
	}

	gin.DisableConsoleColor()
	f, _ := os.Create(filepath.Join(logPath, "access.log"))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 自定义日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.LoggerWithWriter(io.MultiWriter(f, os.Stdout)))

	return r
}

// Recovery middleware recovers from any panics and writes a 500 if there was one.
func setupRecoveryMiddleware(r *gin.Engine) *gin.Engine {
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	return r
}

// interceptor middleware
// 实现对 request 或者 response 的替换
// 当想要访问之前的路由或者中间件设置的 response body 或 header 的时候
// 需要将 gin.context 替换成自定义的 writer
func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wb *responseBuffer
		w := c.Writer
		wb = NewResponseBuffer(w)
		c.Writer = wb
		// before logic
		c.Next()
		// after logic
		wb.Flush()
	}
}
