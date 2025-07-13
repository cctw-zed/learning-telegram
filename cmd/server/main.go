package main

import (
	"fmt"
	"log"
	"net/http"

	"learning-telegram/internal/api"
	"learning-telegram/internal/store"
	"learning-telegram/internal/websocket"
)

func main() {
	store.InitDB("telegram.db")

	fmt.Println("Starting server on :8080")

	// Auth routes with CORS
	registerHandler := http.HandlerFunc(api.RegisterHandler)
	loginHandler := http.HandlerFunc(api.LoginHandler)
	http.Handle("/api/register", registerHandler)
	http.Handle("/api/login", loginHandler)

	// Group routes (protected) with CORS
	createGroupHandler := api.AuthMiddleware(http.HandlerFunc(api.CreateGroupHandler))
	inviteToGroupHandler := api.AuthMiddleware(http.HandlerFunc(api.InviteToGroupHandler))
	http.Handle("/api/groups/create", createGroupHandler)
	http.Handle("/api/groups/invite", inviteToGroupHandler)

	// Status route (protected) with CORS
	statusHandler := api.AuthMiddleware(http.HandlerFunc(api.UserStatusHandler))
	http.Handle("/api/status/user", statusHandler)

	// Chat list route (protected) with CORS
	chatsHandler := api.AuthMiddleware(http.HandlerFunc(api.GetChatsHandler))
	http.Handle("/api/me/chats", chatsHandler)

	// Websocket route (auth is handled inside the handler)
	http.HandleFunc("/ws", websocket.HandleConnections)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
