package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	helper "github.com/more-than-code/auth-helper"
)

type Config struct {
	Secret   []byte `envconfig:"TOKEN_SECRET_KEY"`
	MinuteAt int    `envconfig:"AT_TTL_MINUTE"`
	HourAt   int    `envconfig:"AT_TTL_HOUR"`
	DayAt    int    `envconfig:"AT_TTL_DAY"`
	MinuteRt int    `envconfig:"RT_TTL_MINUTE"`
	HourRt   int    `envconfig:"RT_TTL_HOUR"`
	DayRt    int    `envconfig:"RT_TTL_DAY"`
}

type contextKey struct {
}

var ginCtxKey contextKey

func AuthRequired() gin.HandlerFunc {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	helper, _ := helper.NewHelper(&helper.Config{Secret: cfg.Secret})

	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		strArr := strings.Split(header, "Bearer ")

		if len(strArr) > 1 {
			userStr, err := helper.ParseTokenString(strArr[1])
			if err == nil {
				c.Params = []gin.Param{{Key: "User", Value: userStr}}
				ctx := context.WithValue(c.Request.Context(), ginCtxKey, c)
				_ = c.Request.WithContext(ctx)

				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func CORSEnabled() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
