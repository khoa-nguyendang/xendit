package repositories

const (
	FindUserByUsernameQuery = `SELECT id, 
		user_id, 
		username, 
		password
	FROM users
	WHERE username = ?;`
)
