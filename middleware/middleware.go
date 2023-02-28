package middleware

import (
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
			_, err := helper.ParseTokenString(strArr[1])
			if err == nil {
				// c.Params = []gin.Param{{Key: "User", Value: userStr}}
				// ctx := context.WithValue(c.Request.Context(), ginCtxKey, c)
				// _ = c.Request.WithContext(ctx)

				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
