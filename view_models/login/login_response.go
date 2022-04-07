package viewmodels

type LogInResponse struct {
	Code               int32  `json:"code,omitempty"`
	Message            string `json:"message,omitempty"`
	Error              string `json:"error,omitempty"`
	Token              string `json:"token,omitempty"`
	RefreshToken       string `json:"refresh_token,omitempty"`
	TokenExpiry        int64  `json:"token_expiry,omitempty"`
	RefreshTokenExpiry int64  `json:"refresh_token_expiry,omitempty"`
}
