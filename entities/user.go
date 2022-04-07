package entities

type User struct {
	Id       int64  `json:"id"`
	UserId   string `json:"user_id"`
	RoleId   int32  `json:"role_id"`
	Username string `json:"username"`
	Password string `json:"password" `
}
