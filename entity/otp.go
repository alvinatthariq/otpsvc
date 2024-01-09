package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	OTPStatusCreated   = "created"
	OTPStatusValidated = "validated"
	OTPStatusExpired   = "expired"
	OTPStatusInvalid   = "invalid"
	OTPStatusNotFound  = "notfound"
)

func GetErrorDescriptionByOTPStatus(otpStatus string) string {
	var errorDesc string
	switch otpStatus {
	case OTPStatusExpired:
		errorDesc = "OTP has expired"
	case OTPStatusNotFound:
		errorDesc = "OTP Not Found"
	case OTPStatusInvalid:
		errorDesc = "OTP Invalidated"
	}

	return errorDesc
}

type OTP struct {
	UserID    string    `json:"user_id" gorm:"primaryKey"`
	OTP       string    `json:"otp"`
	Status    string    `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

type OTPHistory struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id"`
	OTP       string    `json:"otp"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *OTP) ToOTPHistory() OTPHistory {
	return OTPHistory{
		ID:        uuid.New().String(),
		UserID:    o.UserID,
		OTP:       o.OTP,
		ExpiredAt: o.ExpiredAt,
		CreatedAt: time.Now().UTC(),
	}
}

type RequestOTP struct {
	UserID string `json:"user_id"`
}
