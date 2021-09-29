package app

import (
	"fmt"
	"hello-gin/src/types"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO it can reuse the struct in types/posts ?
type PostForm struct {
	Title   string `form:"title" json:"title" xml:"title" binding:"required"`
	Content string `form:"content" json:"content" xml:"content" binding:"required"`
	Auth    uint64 `form:"auth" json:"auth" xml:"auth" binding:"required"`
	// Status string
	// Views  uint64
}

// TODO use connection pool  & code refine
func getConn() *gorm.DB {
	dsn := "host=localhost user=hello password=hello dbname=hello port=15432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db
	// sqlDB, err := db.DB()

	// // SetMaxIdleConns 设置空闲连接池中连接的最大数量
	// sqlDB.SetMaxIdleConns(10)
	// // SetMaxOpenConns 设置打开数据库连接的最大数量。
	// sqlDB.SetMaxOpenConns(100)

	// // SetConnMaxLifetime 设置了连接可复用的最大时间。
	// sqlDB.SetConnMaxLifetime(time.Hour)
}

/*
TODO according to blacklist and privacy level
*/
func GetPostInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// uid := c.GetString("uid")
		// TODO use busniess id instead of internal id
		postId := c.Param("pid")

		conn := getConn()

		var post types.Post
		if result := conn.Where("id = ?", postId).Find(&post); result.Error == nil && result.RowsAffected != 0 {
			c.JSON(http.StatusOK, post)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code":  1,
			"error": "Resoure Not Found",
		})
	}
}

/*
return current user info if user_id not passed
otherwise, get other user info
*/
func GetPostList() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")
		queryUserId := c.DefaultQuery("uid", strconv.Itoa(uid))

		var posts []types.Post

		conn := getConn()

		if result := conn.Where("auth = ?", queryUserId).Find(&posts); result.Error == nil && result.RowsAffected != 0 {
			c.JSON(http.StatusOK, posts)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code":  1,
			"error": "Resoure Not Found",
		})
	}
}

/*
TODO 使用 post 的唯一代表值而不是内部 id
*/
func UpdatePostInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// uid := c.GetString("uid")

		updatePostId := c.Param("pid")

		var postForm PostForm
		if err := c.ShouldBind(&postForm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Invalid Update Params",
			})
			return
		}

		fmt.Printf("%v", postForm)

		conn := getConn()

		// 只更新非 0 值， 也可以使用 select 更新固定的列
		if result := conn.Model(&types.Post{}).Where("id = ?", updatePostId).Updates(postForm); result.Error != nil {
			log.Printf("%v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Update Error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"error": "Update Successfully",
		})
	}
}

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		var post types.Post
		if err := c.ShouldBind(&post); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Invalid Update Params",
			})
			return
		}
		conn := getConn()

		// fill other fields
		post.Auth = uid

		fmt.Printf("%v", post)

		// select 可以保证只更新这些字段
		// 还有一些tag可以加在 struct 上面做类型校验
		if result := conn.Select("Title", "Content", "Auth").Create(&post); result.Error != nil {
			log.Printf("%v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Create Error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"error": "Create Successfully",
		})

	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		deletePid := c.Param("pid")

		conn := getConn()
		if result := conn.Delete(&types.Post{}, "id = ?", deletePid); result.Error != nil {
			// log.Printf("%v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Delete Error",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":   0,
				"data":   "Delete Successfully",
				"number": result.RowsAffected,
			})
		}

	}
}
