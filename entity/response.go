package entity

type Meta struct {
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
	Error      string `json:"error,omitempty"`
}

type HTTPEmptyResp struct {
	Meta Meta `json:"metadata"`
}

type HTTPRequestOTPSuccessResp struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
}

type HTTPValidateOTPSuccessResp struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type HTTPValidateOTPErrorResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
