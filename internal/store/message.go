package store

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// InsertPrivateMessage inserts a private message into the database.
func InsertPrivateMessage(sender, receiver, content string) error {
	_, err := DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content, created_at) VALUES ((SELECT id FROM users WHERE username = ?), (SELECT id FROM users WHERE username = ?), ?, ?)",
		sender, receiver, content, time.Now(),
	)
	return err
}

// InsertGroupMessage inserts a group message into the database.
func InsertGroupMessage(sender string, groupID int64, content string) error {
	_, err := DB.Exec(
		"INSERT INTO messages (sender_id, group_id, content, created_at) VALUES ((SELECT id FROM users WHERE username = ?), ?, ?, ?)",
		sender, groupID, content, time.Now(),
	)
	return err
}

// GetPrivateHistory retrieves private chat history between two users.
func GetPrivateHistory(user1, user2 string) ([]Message, error) {
	rows, err := DB.Query(
		`SELECT m.id, u1.username, u2.username, m.content, m.created_at
		 FROM messages m
		 JOIN users u1 ON m.sender_id = u1.id
		 JOIN users u2 ON m.receiver_id = u2.id
		 WHERE (u1.username = ? AND u2.username = ?) OR (u1.username = ? AND u2.username = ?)
		 ORDER BY m.created_at ASC LIMIT 100`,
		user1, user2, user2, user1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Sender, &m.Receiver, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// GetGroupHistory retrieves chat history for a group.
func GetGroupHistory(groupID int64) ([]Message, error) {
	rows, err := DB.Query(
		`SELECT m.id, u.username, m.group_id, m.content, m.created_at
		 FROM messages m
		 JOIN users u ON m.sender_id = u.id
		 WHERE m.group_id = ?
		 ORDER BY m.created_at ASC LIMIT 100`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		// Note: m.Receiver will be empty for group messages.
		if err := rows.Scan(&m.ID, &m.Sender, &m.Receiver, &m.Content, &m.CreatedAt); err != nil {
			// We expect receiver_id to be NULL, so we create a placeholder for scanning.
			var groupID sql.NullInt64
			if scanErr := rows.Scan(&m.ID, &m.Sender, &groupID, &m.Content, &m.CreatedAt); scanErr != nil {
				return nil, scanErr
			}
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// 获取当前时间字符串
func NowStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
