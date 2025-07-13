package api

import (
	"encoding/json"
	"net/http"

	"learning-telegram/internal/store"
)

// GetChatsHandler retrieves all chat partners (users and groups) for the authenticated user.
func GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		http.Error(w, "无法从Token获取用户信息", http.StatusUnauthorized)
		return
	}

	users, err := store.GetAllUsers(username)
	if err != nil {
		http.Error(w, "获取用户列表失败", http.StatusInternalServerError)
		return
	}

	groups, err := store.GetUserGroups(username)
	if err != nil {
		http.Error(w, "获取群组列表失败", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users":  users,
		"groups": groups,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "无法生成响应", http.StatusInternalServerError)
	}
}
