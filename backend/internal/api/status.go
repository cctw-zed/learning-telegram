package api

import (
	"encoding/json"
	"net/http"

	"learning-telegram/internal/websocket"
)

// UserStatusHandler checks and returns the online status of a user.
func UserStatusHandler(w http.ResponseWriter, r *http.Request) {
	// We expect the username to be a query parameter, e.g., /api/status/user?username=testuser
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "查询参数 'username' 不能为空", http.StatusBadRequest)
		return
	}

	appHub := websocket.GetHub()
	isOnline := appHub.IsUserOnline(username)

	response := map[string]interface{}{
		"username": username,
		"online":   isOnline,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// This is unlikely to happen, but good practice to handle.
		http.Error(w, "无法生成响应", http.StatusInternalServerError)
	}
}
