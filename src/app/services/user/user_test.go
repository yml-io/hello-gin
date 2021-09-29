package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.New()
	router.GET("/view", GetUserInfo())
	router.POST("/update", UpdateUserInfo())
	router.POST("/upload_avatar", UploadAvatar())
	router.PUT("/create", CreateUser())
	router.DELETE("/delete/:uid", DeleteUser())
	return router
}

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
