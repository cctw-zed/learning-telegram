package websocket

import (
	"log"
	"net/http"

	"learning-telegram/internal/auth"
	"learning-telegram/internal/store"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接，方便测试
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// 支持URL参数中的token
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		// 如果URL参数中没有，再检查header
		tokenString = r.Header.Get("X-Token")
	}
	if tokenString == "" {
		http.Error(w, "未授权：缺少Token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "未授权：无效的Token", http.StatusUnauthorized)
		return
	}
	username := claims.Username

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer ws.Close()

	hub.Register(username, ws)
	defer hub.Unregister(username, ws)

	log.Printf("用户 %s 已连接", username)

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("%s 断开连接: %v", username, err)
			break
		}
		typeVal, _ := msg["type"].(string)
		switch typeVal {
		case "send_message", "private":
			to, _ := msg["to"].(string)
			content, _ := msg["content"].(string)
			if to == "" || content == "" {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "to和content不能为空"})
				continue
			}
			// 存储消息
			err := store.InsertPrivateMessage(username, to, content)
			if err != nil {
				log.Printf("消息存储失败 (from: %s, to: %s): %v", username, to, err)
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "消息存储失败，请检查目标用户是否存在"})
				continue
			}
			// 推送给目标用户所有在线端
			push := map[string]interface{}{
				"type":    "new_message",
				"from":    username,
				"content": content,
				"ts":      store.NowStr(),
			}
			hub.SendToUser(to, push)
			// 也回显给自己（多端同步）
			hub.SendToUser(username, push)
		case "send_group_message", "group":
			groupIDFloat, _ := msg["group_id"].(float64)
			groupID := int64(groupIDFloat)
			content, _ := msg["content"].(string)
			if groupID == 0 || content == "" {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "group_id和content不能为空"})
				continue
			}

			// 1. 存储群消息
			err := store.InsertGroupMessage(username, groupID, content)
			if err != nil {
				log.Printf("群消息存储失败 (user: %s, group: %d): %v", username, groupID, err)
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "群消息存储失败"})
				continue
			}

			// 2. 获取所有群成员
			members, err := store.GetGroupMembers(groupID)
			if err != nil {
				log.Printf("获取群成员失败 (group: %d): %v", groupID, err)
				continue
			}

			// 3. 向所有在线的群成员推送消息
			push := map[string]interface{}{
				"type":     "new_group_message",
				"group_id": groupID,
				"from":     username,
				"content":  content,
				"ts":       store.NowStr(),
			}
			for _, member := range members {
				hub.SendToUser(member, push)
			}
		case "history":
			with, _ := msg["with"].(string)
			if with == "" {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "with不能为空"})
				continue
			}
			msgs, err := store.GetPrivateHistory(username, with)
			if err != nil {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "查询历史失败"})
				continue
			}
			ws.WriteJSON(map[string]interface{}{
				"type":     "history",
				"with":     with,
				"messages": msgs,
			})
		case "history_group":
			groupIDFloat, _ := msg["group_id"].(float64)
			groupID := int64(groupIDFloat)
			if groupID == 0 {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "group_id不能为空"})
				continue
			}
			// 1. 验证用户是否在群组中
			isMember, err := store.IsUserInGroup(username, groupID)
			if err != nil || !isMember {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "无权限访问该群组历史"})
				continue
			}
			// 2. 获取群组历史消息
			msgs, err := store.GetGroupHistory(groupID)
			if err != nil {
				ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "查询群组历史失败"})
				continue
			}
			ws.WriteJSON(map[string]interface{}{
				"type":     "history_group",
				"group_id": groupID,
				"messages": msgs,
			})
		case "typing":
			to, _ := msg["to"].(string)
			groupIDFloat, _ := msg["group_id"].(float64)
			groupID := int64(groupIDFloat)

			if to != "" { // Private chat typing
				push := map[string]interface{}{
					"type": "user_typing",
					"from": username,
				}
				hub.SendToUser(to, push)
			} else if groupID != 0 { // Group chat typing
				isMember, err := store.IsUserInGroup(username, groupID)
				if err != nil || !isMember {
					// Don't send error back, just ignore silently.
					continue
				}
				members, err := store.GetGroupMembers(groupID)
				if err != nil {
					continue
				}
				push := map[string]interface{}{
					"type":     "user_typing",
					"from":     username,
					"group_id": groupID,
				}
				// Broadcast to all members except the sender
				for _, member := range members {
					if member != username {
						hub.SendToUser(member, push)
					}
				}
			}
		default:
			ws.WriteJSON(map[string]interface{}{"type": "error", "msg": "未知消息类型"})
		}
	}
}
