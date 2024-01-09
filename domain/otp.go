package domain

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/alvinatthariq/otpsvc/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *domain) RequestOTP(ctx context.Context, v entity.RequestOTP) (otp entity.OTP, err error) {
	// validate user id required
	if v.UserID == "" {
		return otp, entity.ErrorUserIDRequired
	}

	otp = entity.OTP{
		UserID:    v.UserID,
		OTP:       fmt.Sprintf("%d%d%d%d%d", rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9)),
		ExpiredAt: time.Now().UTC().Add(2 * time.Minute),
		Status:    entity.OTPStatusCreated,
	}

	err = d.gorm.Transaction(func(tx *gorm.DB) error {
		// upsert to otp
		err = d.gorm.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&otp).Error
		if err != nil {
			return err
		}

		otpHistory := otp.ToOTPHistory()

		// create otp history
		err = d.gorm.Create(&otpHistory).Error
		if err != nil {
			return err
		}

		return nil
	})

	return otp, err
}

func (d *domain) ValidateOTP(ctx context.Context, v entity.OTP) (otp entity.OTP, err error) {
	// validate user id required
	if v.UserID == "" {
		return otp, entity.ErrorUserIDRequired
	}

	// get from db
	err = d.gorm.First(&otp, "user_id = ?", v.UserID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		otp.Status = entity.OTPStatusNotFound
		return otp, nil
	}

	if otp.Status == entity.OTPStatusValidated {
		return otp, entity.ErrorOTPAlreadyValidated
	}

	if v.OTP != otp.OTP {
		otp.Status = entity.OTPStatusInvalid
		return otp, nil
	}

	otp.Status = entity.OTPStatusValidated
	if time.Now().UTC().After(otp.ExpiredAt) {
		otp.Status = entity.OTPStatusExpired
	}

	err = d.gorm.Save(otp).Error
	if err != nil {
		return otp, err
	}

	return otp, err
}
