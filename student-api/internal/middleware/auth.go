package middleware

import (
	"errors"
	"os"
	logging "student-api/internal/logger"
	"student-api/internal/model"
	"student-api/internal/service"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"go.uber.org/zap"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName string
	Role     string
}

const (
	adminRole   = "admin"
	regularRole = "regular"
)

var identityKey = "username"

func AuthMiddleware(teacherService service.TeacherService) (*jwt.GinJWTMiddleware, error) {
	tracer := otel.Tracer("auth-middleware")

	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "student zone",
		Key:         []byte(os.Getenv("JWT_SECRET")),
		Timeout:     5 * time.Minute,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,

		// Authenticator runs on login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			_, span := tracer.Start(c.Request.Context(), "authenticate-user")
			defer span.End()

			logging.Logger.Info("[auth-middleware]: Processing login request")

			var login model.Login
			if err := c.ShouldBindJSON(&login); err != nil {
				span.SetAttributes(attribute.String("error", err.Error()))
				span.SetStatus(codes.Error, "missing login values")
				logging.Logger.Error("[auth-middleware]: Missing login values", zap.Error(err))
				return nil, jwt.ErrMissingLoginValues
			}

			span.SetAttributes(attribute.String("username", login.Username))
			teacher, err := teacherService.GetTeacher(login.Username, c.Request.Context())
			if err != nil {
				span.SetAttributes(attribute.String("error", err.Error()))
				span.SetStatus(codes.Error, "invalid username")
				logging.Logger.Error("[auth-middleware]: Invalid username",
					zap.String("username", login.Username),
					zap.Error(err))
				return nil, errors.New("Invalid username")
			}

			err = bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(login.Password))
			if err != nil {
				span.SetAttributes(attribute.String("error", err.Error()))
				span.SetStatus(codes.Error, "invalid password")
				logging.Logger.Error("[auth-middleware]: Invalid password",
					zap.String("username", login.Username),
					zap.Error(err))
				return nil, errors.New("Invalid password")
			}

			span.SetAttributes(
				attribute.String("username", teacher.Username),
				attribute.String("role", teacher.Role),
			)
			span.SetStatus(codes.Ok, "authentication successful")

			logging.Logger.Info("[auth-middleware]: Login successful",
				zap.String("username", teacher.Username),
				zap.String("role", teacher.Role))

			return &User{
				UserName: teacher.Username,
				Role:     teacher.Role,
			}, nil
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
			_, span := tracer.Start(c.Request.Context(), "authorize-request")
			defer span.End()

			user, ok := data.(*User)
			if !ok {
				span.SetStatus(codes.Error, "invalid user data type")
				logging.Logger.Error("[auth-middleware]: Invalid user data type in Authorizator")
				return false
			}

			span.SetAttributes(
				attribute.String("username", user.UserName),
				attribute.String("role", user.Role),
				attribute.String("method", c.Request.Method),
				attribute.String("path", c.Request.URL.Path),
			)

			// Allow all routes for adminRole
			if user.Role == adminRole {
				span.SetStatus(codes.Ok, "admin access granted")
				logging.Logger.Info("[auth-middleware]: Admin access granted",
					zap.String("username", user.UserName),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path))
				return true
			}

			// Deny DELETE requests for regularRole
			if user.Role == regularRole && c.Request.Method == "DELETE" {
				span.SetStatus(codes.Error, "regular user attempted delete operation")
				logging.Logger.Warn("[auth-middleware]: Regular user attempted delete operation",
					zap.String("username", user.UserName),
					zap.String("path", c.Request.URL.Path))
				return false
			}

			// Allow other operations for regularRole
			span.SetStatus(codes.Ok, "regular access granted")
			logging.Logger.Info("[auth-middleware]: Regular access granted",
				zap.String("username", user.UserName),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path))
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
