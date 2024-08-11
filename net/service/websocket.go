package service

import (
	"encoding/json"
	"net/http"
	"tiny_talk/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var allowedOrigins = map[string]bool{
	"http://www.websocket-test.com/": true,
	"http://localhost:8081":          true,
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许来自任何源的请求
		return true
		// origin := r.Header.Get("Origin")
		// return allowedOrigins[origin]
	},
}

// TokenResponse 结构体用于封装新的 Token
type TokenResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func HandleWebSocket(c *gin.Context) {
	current_token := c.Query("token")
	logger.Infof("websocket token: %v", current_token)
	id, err := GetUserIdFromRedisByToken(current_token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("websocket failed to get user id error: %v", err)
		return
	}
	logger.Infof("websocket user id: %v", id)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("websocket failed to upgrade error: %v", err)
		return
	}
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			logger.Errorf("websocket failed to read message error: %v", err)
			break
		}
		new_token, err := RefeshTokenExpiration(current_token)
		if err != nil {
			logger.Errorf("websocket failed to refresh token expiration error: %v", err)
			break
		}

		if current_token != new_token {
			logger.Infof("websocket token expired, refresh token: %v", new_token)
			current_token = new_token

			tokenResp := TokenResponse{
				Type:  "new_token",
				Token: new_token,
			}
			tokenJSON, err := json.Marshal(tokenResp)
			if err != nil {
				logger.Infof("Error marshaling token response: %v", err)
				continue
			}

			err = ws.WriteMessage(websocket.TextMessage, tokenJSON)
			if err != nil {
				logger.Errorf("Error writing token message: %v, error: %v", message, err)
				break
			}
		}

		err = ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.Errorf("websocket failed to write message %v error: %v", message, err)
			break
		}
	}
}
