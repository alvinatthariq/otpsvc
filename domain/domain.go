package domain

import (
	"context"

	"github.com/alvinatthariq/otpsvc/entity"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type domain struct {
	gorm        *gorm.DB
	redisClient *redis.Client
}

type DomainItf interface {
	RequestOTP(ctx context.Context, v entity.RequestOTP) (otp entity.OTP, err error)
	ValidateOTP(ctx context.Context, v entity.OTP) (otp entity.OTP, err error)
}

func Init(gorm *gorm.DB, redisClient *redis.Client) DomainItf {
	return &domain{
		gorm:        gorm,
		redisClient: redisClient,
	}
}
