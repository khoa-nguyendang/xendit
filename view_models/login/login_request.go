package viewmodels

type LogInRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
