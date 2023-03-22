package api

import (
	"github.com/kelseyhightower/envconfig"
	authHelper "github.com/more-than-code/auth-helper"
	"github.com/more-than-code/deploybot-service-api/repository"
)

type Config struct {
	MinuteAt       int    `envconfig:"AT_TTL_MINUTE"`
	HourAt         int    `envconfig:"AT_TTL_HOUR"`
	DayAt          int    `envconfig:"AT_TTL_DAY"`
	MinuteRt       int    `envconfig:"RT_TTL_MINUTE"`
	HourRt         int    `envconfig:"RT_TTL_HOUR"`
	DayRt          int    `envconfig:"RT_TTL_DAY"`
	Secret         []byte `envconfig:"TOKEN_SECRET_KEY"`
	GoogleClientId string `envconfig:"GOOGLE_CLIENT_ID"`
}

type Api struct {
	repo           *repository.Repository
	atHelper       *authHelper.Helper
	rtHelper       *authHelper.Helper
	googleClientId string
}

func NewApi() *Api {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	athelper, _ := authHelper.NewHelper(&authHelper.Config{Secret: cfg.Secret, TtlMinute: cfg.MinuteAt, TtlHour: cfg.HourAt, TtlDay: cfg.DayAt})

	rthelper, _ := authHelper.NewHelper(&authHelper.Config{Secret: cfg.Secret, TtlMinute: cfg.MinuteRt, TtlHour: cfg.HourRt, TtlDay: cfg.DayRt})

	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	return &Api{repo: r, atHelper: athelper, rtHelper: rthelper, googleClientId: cfg.GoogleClientId}
}
