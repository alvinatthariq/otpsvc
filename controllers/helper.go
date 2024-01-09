package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alvinatthariq/otpsvc/entity"
)

func httpRespError(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	switch err {
	case entity.ErrorOTPNotFound,
		entity.ErrorOTPInvalid,
		entity.ErrorOTPExpired,
		entity.ErrorUserIDRequired,
		entity.ErrorOTPAlreadyValidated:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	statusStr := http.StatusText(statusCode)

	jsonErrResp := &entity.HTTPEmptyResp{
		Meta: entity.Meta{
			Path:       r.URL.String(),
			StatusCode: statusCode,
			Status:     statusStr,
			Message:    fmt.Sprintf("%s %s [%d] %s", r.Method, r.URL.RequestURI(), statusCode, statusStr),
			Error:      err.Error(),
			Timestamp:  time.Now().Format(time.RFC3339),
		},
	}

	raw, err := json.Marshal(jsonErrResp)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(raw)
}

func httpRespSuccess(w http.ResponseWriter, r *http.Request, statusCode int, resp interface{}) {
	meta := entity.Meta{
		Path:       r.URL.String(),
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Message:    fmt.Sprintf("%s %s [%d] %s", r.Method, r.URL.RequestURI(), statusCode, http.StatusText(statusCode)),
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	var (
		raw []byte
		err error
	)
	switch data := resp.(type) {
	case nil:
		httpResp := &entity.HTTPEmptyResp{
			Meta: meta,
		}
		raw, err = json.Marshal(httpResp)
		if err != nil {
			statusCode = http.StatusInternalServerError
		}
	case entity.OTP:
		if data.Status == entity.OTPStatusCreated {
			httpResp := &entity.HTTPRequestOTPSuccessResp{
				UserID: data.UserID,
				OTP:    data.OTP,
			}
			raw, err = json.Marshal(httpResp)
			if err != nil {
				statusCode = http.StatusInternalServerError
			}
		} else if data.Status == entity.OTPStatusValidated {
			httpResp := &entity.HTTPValidateOTPSuccessResp{
				UserID:  data.UserID,
				Message: "OTP validated successfully.",
			}

			raw, err = json.Marshal(httpResp)
			if err != nil {
				statusCode = http.StatusInternalServerError
			}
		} else {
			httpResp := &entity.HTTPValidateOTPErrorResp{
				Error:            data.Status,
				ErrorDescription: entity.GetErrorDescriptionByOTPStatus(data.Status),
			}

			raw, err = json.Marshal(httpResp)
			if err != nil {
				statusCode = http.StatusInternalServerError
			}
		}

	default:
		httpRespError(w, r, fmt.Errorf("cannot cast type of %+v", data))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(raw)
}
