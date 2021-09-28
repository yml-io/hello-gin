package app

import (
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
		uid := c.GetString("uid")
		queryUID := c.DefaultQuery("uid", uid)

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
	uid := c.GetString("uid")
	fileName := strconv.Itoa(int(time.Now().Unix()))
	uploadDir := filepath.Join("public", uid)

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
		uid := c.GetString("uid")
		formMap := c.PostFormMap("user")

		conn := getConn()
		var updatePerson types.Person
		// assemble the update object
		if v, ok := formMap["user_name"]; ok {
			updatePerson.UserName = v
		}
		if v, ok := formMap["nick_name"]; ok {
			updatePerson.NickName = v
		}
		if v, ok := formMap["introduce"]; ok {
			updatePerson.Introduce = v
		}
		if v, ok := formMap["sex"]; ok {
			updatePerson.Sex = v
		}
		if v, ok := formMap["email"]; ok {
			updatePerson.Email = v
		}
		if v, ok := formMap["privacy"]; ok {
			updatePerson.Privacy = v
		}

		if result := conn.Model(&types.Person{}).Where("user_name = ?", uid).Updates(updatePerson); result.Error != nil {
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
		formMap := c.PostFormMap("user")

		conn := getConn()
		var insertPerson types.Person
		// assemble the update object
		if v, ok := formMap["user_name"]; ok {
			insertPerson.UserName = v
		}
		if v, ok := formMap["nick_name"]; ok {
			insertPerson.NickName = v
		}
		if v, ok := formMap["introduce"]; ok {
			insertPerson.Introduce = v
		}
		if v, ok := formMap["sex"]; ok {
			insertPerson.Sex = v
		}
		if v, ok := formMap["email"]; ok {
			insertPerson.Email = v
		}
		if v, ok := formMap["privacy"]; ok {
			insertPerson.Privacy = v
		}
		conn.Select("UserName", "NickName", "Introduce", "Sex", "Email", "Privacy").Create(&insertPerson)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.PostForm("uid")

		conn := getConn()
		var deletePerson types.Person
		conn.Where("user_name = ?", uid).First(&deletePerson)
		conn.Delete(&deletePerson)
	}
}
