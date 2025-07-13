package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"learning-telegram/internal/auth"
	"learning-telegram/internal/store"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "密码加密失败", http.StatusInternalServerError)
		return
	}

	_, err = store.DB.Exec(
		"INSERT INTO users (username, password_hash, created_at) VALUES (?, ?, ?)",
		req.Username, string(hash), time.Now(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "用户名已存在", http.StatusConflict)
			return
		}
		http.Error(w, "注册失败", http.StatusInternalServerError)
		return
	}

	// 注册成功后自动生成token
	tokenString, err := auth.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "生成Token失败", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	var hash string
	row := store.DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", req.Username)
	err := row.Scan(&hash)
	if err == sql.ErrNoRows {
		http.Error(w, "用户不存在", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "登录失败", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		http.Error(w, "密码错误", http.StatusUnauthorized)
		return
	}

	tokenString, err := auth.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "生成Token失败", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
