package store

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllUsers retrieves all users except the one with the given username.
func GetAllUsers(exceptUsername string) ([]User, error) {
	rows, err := DB.Query("SELECT id, username, created_at FROM users WHERE username != ?", exceptUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
