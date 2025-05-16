package middleware

import (
	"errors"
	"os"
	"student-api/internal/model"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type User struct {
	UserName string
	Role     string
}

var identityKey = "username"

func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "student zone",
		Key:         []byte(os.Getenv("JWT_SECRET")),
		Timeout:     time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,

		// Authenticator runs on login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login model.Login
			if err := c.ShouldBindJSON(&login); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			// Hardcoded credentials for demo purpose
			if login.Username == "admin" && login.Password == "password" {
				return &User{UserName: login.Username, Role: "admin"}, nil
			}

			if login.Username == "tester" && login.Password == "password" {
				return &User{UserName: login.Username, Role: "tester"}, nil
			}

			return nil, errors.New("invalid credentials")
		},

		// PayloadFunc defines custom claims
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: user.UserName,
					"role":      user.Role,
				}
			}
			return jwt.MapClaims{}
		},

		// IdentityHandler gets identity from token
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
				Role:     claims["role"].(string),
			}
		},

		// Authorizator checks if the user has permission
		Authorizator: func(data interface{}, c *gin.Context) bool {
			user, ok := data.(*User)
			if !ok {
				return false
			}

			// Allow all routes for admin
			if user.Role == "admin" {
				return true
			}

			// Deny DELETE requests for tester
			if user.Role == "tester" && c.Request.Method == "DELETE" {
				return false
			}

			// Allow other operations for tester
			return true
		},

		// Unauthorized handles unauthorized requests
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
