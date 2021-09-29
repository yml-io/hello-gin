package auth

import (
	"errors"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type HeaderWithAuth struct {
	Authorization string `header:"authorization"  binding:"required"`
}

// todo jwt 解码
func AuthRequired(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var header = HeaderWithAuth{}
		if err := c.ShouldBindHeader(&header); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalidate request header for Auth",
			})
			return
		}

		jwtToken := header.Authorization
		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err_msg := "Unexpected signing method: excepted HS256, got " + token.Method.Alg()
				return nil, errors.New(err_msg)
			}

			// TODO should query secret key from store via uid from jwt claim segment
			// for test: all key is "secretkey"
			hmacSecret := []byte("aaa")
			return hmacSecret, nil
		})

		if err != nil {
			goto error_abort
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// TODO uid 应该是一个临时值，而不应该对应数据库的值
			c.Set("uid", int(claims["uid"].(float64)))
			return
		}

	error_abort:
		if err != nil {
			fmt.Printf("%v", err)
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Invalidate AUTHORIZATION header",
		})
	}
}

func AdminRequired(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
