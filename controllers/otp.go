package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alvinatthariq/otpsvc/entity"
)

func (c *controller) RequestOTP(w http.ResponseWriter, r *http.Request) {
	// parse request body
	var requestBody entity.RequestOTP
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpRespError(w, r, fmt.Errorf("Error Decode Request Body : %v", err))
		return
	}

	ctx := r.Context()

	otp, err := c.domain.RequestOTP(ctx, requestBody)
	if err != nil {
		httpRespError(w, r, err)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, otp)
}

func (c *controller) ValidateOTP(w http.ResponseWriter, r *http.Request) {
	// parse request body
	var requestBody entity.OTP
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpRespError(w, r, fmt.Errorf("Error Decode Request Body : %v", err))
		return
	}

	ctx := r.Context()

	otp, err := c.domain.ValidateOTP(ctx, requestBody)
	if err != nil {
		httpRespError(w, r, err)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, otp)
}
