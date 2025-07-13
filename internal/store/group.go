package store

import (
	"database/sql"
	"time"
)

type Group struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatorID int       `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateGroup creates a new group and adds the creator as the first member.
func CreateGroup(name string, creatorUsername string) (int64, error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Rollback on any error

	var creatorID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = ?", creatorUsername).Scan(&creatorID)
	if err != nil {
		return 0, err // Creator user not found
	}

	res, err := tx.Exec("INSERT INTO groups (name, creator_id) VALUES (?, ?)", name, creatorID)
	if err != nil {
		return 0, err
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, creatorID)
	if err != nil {
		return 0, err
	}

	return groupID, tx.Commit()
}

// AddGroupMember adds a user to a group.
func AddGroupMember(groupID int64, username string) error {
	var userID int
	err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows // User not found
		}
		return err
	}

	_, err = DB.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	return err
}

// GetGroupMembers retrieves all member usernames for a given group.
func GetGroupMembers(groupID int64) ([]string, error) {
	rows, err := DB.Query("SELECT u.username FROM users u JOIN group_members gm ON u.id = gm.user_id WHERE gm.group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		members = append(members, username)
	}
	return members, nil
}

// IsUserInGroup checks if a user is a member of a group.
func IsUserInGroup(username string, groupID int64) (bool, error) {
	var userID int
	err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return false, err
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserGroups retrieves all groups a user is a member of.
func GetUserGroups(username string) ([]Group, error) {
	var userID int
	err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT g.id, g.name, g.creator_id, g.created_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatorID, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
