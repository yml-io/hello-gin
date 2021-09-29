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
return current user info if user_id not passed
otherwise, get other user info
*/
func GetFollowList() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")
		queryUserId := c.DefaultQuery("uid", strconv.Itoa(uid))

		var follows []types.Follow

		conn := getConn()

		if result := conn.Where("follower = ?", queryUserId).Find(&follows); result.Error == nil && result.RowsAffected != 0 {
			c.JSON(http.StatusOK, follows)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code":  1,
			"error": "Resoure Not Found",
		})
	}
}

func CreateFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		var follow types.Follow
		if err := c.ShouldBind(&follow); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Invalid Update Params",
			})
			return
		}
		conn := getConn()

		// fill other fields
		follow.Followee = uid

		fmt.Printf("%v", follow)

		// select 可以保证只更新这些字段
		// 还有一些tag可以加在 struct 上面做类型校验
		if result := conn.Select("Follower", "Followee").Create(&follow); result.Error != nil {
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

func DeleteFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")
		deleteUid := c.Param("uid")

		conn := getConn()
		if result := conn.Delete(&types.Follow{}, "followee = ? and follower = ?", uid, deleteUid); result.Error != nil {
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
