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
func GetFavoriteList() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")
		queryUserId := c.DefaultQuery("uid", strconv.Itoa(uid))

		var favorites []types.Favorite

		conn := getConn()

		if result := conn.Where("uid = ?", queryUserId).Find(&favorites); result.Error == nil && result.RowsAffected != 0 {
			c.JSON(http.StatusOK, favorites)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code":  1,
			"error": "Resoure Not Found",
		})
	}
}

func CreateFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		var favorite types.Favorite
		if err := c.ShouldBind(&favorite); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Invalid Update Params",
			})
			return
		}
		conn := getConn()

		// fill other fields
		favorite.Uid = uid

		fmt.Printf("%v", favorite)

		// select 可以保证只更新这些字段
		// 还有一些tag可以加在 struct 上面做类型校验
		if result := conn.Select("Follower", "Followee").Create(&favorite); result.Error != nil {
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

func DeleteFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")
		deletePid := c.Param("pid")

		conn := getConn()
		if result := conn.Delete(&types.Favorite{}, "pid = ?", uid, deletePid); result.Error != nil {
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
