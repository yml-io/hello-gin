package app

import (
	"fmt"
	"hello-gin/src/types"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO it can reuse the struct in types/person ?
type UserForm struct {
	UserName  string `form:"user_name" json:"user_name" xml:"user_name"`
	NickName  string `form:"nick_name" json:"nick_name" xml:"nick_name"`
	Introduce string `form:"introduce" json:"introduce" xml:"introduce"`
	Sex       string `form:"sex" json:"sex" xml:"sex"`
	Email     string `form:"email" json:"email" xml:"email"`
	Privacy   string `form:"privacy" json:"privacy" xml:"privacy"`
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
return current user info if user_id not passed
otherwise, get other user info
TODO according to blacklist and privacy level
*/
func GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")
		queryUID := c.DefaultQuery("uid", strconv.Itoa(uid))

		conn := getConn()

		println("select uid :" + queryUID)
		var person types.Person
		if result := conn.Where("user_name = ?", queryUID).Find(&person); result.Error == nil && result.RowsAffected != 0 {
			c.JSON(http.StatusOK, person)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code":  1,
			"error": "Resoure Not Found",
		})
	}
}

func getUserUploadPath(c *gin.Context) (string, error) {
	uid := c.GetInt("uid")
	fileName := strconv.Itoa(int(time.Now().Unix()))
	uploadDir := filepath.Join("public", strconv.Itoa(uid))

	_, err := os.Stat(uploadDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0777); err != nil {
			return "", err
		}
	}
	return filepath.Join(uploadDir, fileName), nil
}

// curl -X POST http://localhost:8080/user/upload_avatar \
//   -F "file=@/Users/mark_yu/GolandProjects/hello-gin/src/assets/avatar.png" \
//   -H "Content-Type: multipart/form-data" \
//   -H 'authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJhZG1pbiIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.gOm-gNc5r6xxW7SiWhToyVrRIZ-WHhJRJimzGDQoPxo'
// or
// curl -X POST http://localhost:8080/upload \
//   -F "upload[]=@/Users/appleboy/test1.zip" \
//   -F "upload[]=@/Users/appleboy/test2.zip" \
//   -H "Content-Type: multipart/form-data"
func UploadAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			dst, _ := getUserUploadPath(c)
			log.Println(file.Filename + " => " + dst)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				log.Printf("%v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":  1,
					"error": "Upload Error",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"error": "Upload Successfully",
		})
	}
}

func UpdateUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt("uid")

		var userForm UserForm
		if err := c.ShouldBind(&userForm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Invalid Update Params",
			})
			return
		}

		fmt.Printf("%v", userForm)
		// formMap := c.PostFormMap("user")

		conn := getConn()

		// 只更新非 0 值， 也可以使用 select 更新固定的列
		if result := conn.Model(&types.Person{}).Where("user_name = ?", uid).Updates(userForm); result.Error != nil {
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

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var person types.Person
		if err := c.ShouldBind(&person); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"error": "Invalid Update Params",
			})
			return
		}

		fmt.Printf("%v", person)
		// formMap := c.PostFormMap("user")

		conn := getConn()
		// select 可以保证只更新这些字段
		// 还有一些tag可以加在 struct 上面做类型校验
		if result := conn.Select("UserName", "NickName", "Introduce", "Sex", "Email", "Privacy").Create(&person); result.Error != nil {
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

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		deleteUid := c.Param("uid")

		conn := getConn()
		if result := conn.Delete(&types.Person{}, "user_name = ?", deleteUid); result.Error != nil {
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
