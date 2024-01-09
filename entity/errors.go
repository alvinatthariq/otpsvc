package entity

import "fmt"

const (
	// Code SQL Error from https://github.com/go-sql-driver/mysql/blob/master/errors.go
	CodeMySQLDuplicateEntry             = 1062
	CodeMySQLForeignKeyConstraintFailed = 1452
	CodeMySQLTableNotExist              = 1146
)

var (
	ErrorOTPNotFound         error = fmt.Errorf("OTP Not Found")
	ErrorOTPInvalid          error = fmt.Errorf("OTP Invalid")
	ErrorOTPExpired          error = fmt.Errorf("OTP Expired")
	ErrorUserIDRequired      error = fmt.Errorf("User ID Required")
	ErrorOTPAlreadyValidated error = fmt.Errorf("OTP Already Validated")
)
