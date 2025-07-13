package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"learning-telegram/internal/store"
)

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type InviteToGroupRequest struct {
	GroupID  int64  `json:"group_id"`
	Username string `json:"username"`
}

// CreateGroupHandler handles the creation of a new group.
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	creatorUsername, ok := r.Context().Value("username").(string)
	if !ok {
		http.Error(w, "无法从Token获取用户信息", http.StatusUnauthorized)
		return
	}

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "群组名不能为空", http.StatusBadRequest)
		return
	}

	groupID, err := store.CreateGroup(req.Name, creatorUsername)
	if err != nil {
		http.Error(w, "创建群组失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "群组创建成功",
		"group_id": groupID,
	})
}

// InviteToGroupHandler handles inviting a user to a group.
func InviteToGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Although we are not using the inviter's username here,
	// the middleware has already ensured that the request is authenticated.
	// You could add permission checks here, e.g., if only admins can invite.
	if _, ok := r.Context().Value("username").(string); !ok {
		http.Error(w, "无法从Token获取用户信息", http.StatusUnauthorized)
		return
	}

	var req InviteToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数", http.StatusBadRequest)
		return
	}
	if req.GroupID == 0 || strings.TrimSpace(req.Username) == "" {
		http.Error(w, "group_id 和 username 不能为空", http.StatusBadRequest)
		return
	}

	// In a real app, you should also check if the inviter has permission to add members.
	// For simplicity, we are skipping that check here.

	err := store.AddGroupMember(req.GroupID, req.Username)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "用户已在群组中", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "要邀请的用户不存在", http.StatusNotFound)
			return
		}
		http.Error(w, "邀请失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("邀请成功"))
}
