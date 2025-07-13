package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connection successful.")
	createTables()
}

func createTables() {
	usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	groupsTable := `
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		creator_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (creator_id) REFERENCES users (id)
	);`

	groupMembersTable := `
	CREATE TABLE IF NOT EXISTS group_members (
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (group_id, user_id),
		FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);`

	messagesTable := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        receiver_id INTEGER, -- Can be NULL for group messages
        group_id INTEGER,    -- Can be NULL for private messages
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sender_id) REFERENCES users (id),
        FOREIGN KEY (receiver_id) REFERENCES users (id),
        FOREIGN KEY (group_id) REFERENCES groups (id)
    );`

	// Drop old messages table if it exists without group_id to rebuild it.
	// This is a simple approach for development, in production a proper migration tool should be used.
	var tableName string
	err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='messages' AND sql NOT LIKE '%group_id%'").Scan(&tableName)
	if err == nil && tableName == "messages" {
		_, err = DB.Exec("DROP TABLE messages")
		if err != nil {
			log.Fatalf("Could not drop old messages table: %v", err)
		}
		log.Println("Old messages table dropped to be recreated.")
	}

	_, err = DB.Exec(usersTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}
	log.Println("Users table ready.")

	_, err = DB.Exec(groupsTable)
	if err != nil {
		log.Fatalf("Could not create groups table: %v", err)
	}
	log.Println("Groups table ready.")

	_, err = DB.Exec(groupMembersTable)
	if err != nil {
		log.Fatalf("Could not create group_members table: %v", err)
	}
	log.Println("Group members table ready.")

	_, err = DB.Exec(messagesTable)
	if err != nil {
		log.Fatalf("Could not create messages table: %v", err)
	}
	log.Println("Messages table ready.")
}
