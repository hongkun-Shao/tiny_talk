package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
	"tiny_talk/infrastructure/crud"
	"tiny_talk/infrastructure/models"
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

type UserConn struct {
	ID      int64
	Channel chan string
	WsConn  *websocket.Conn
}

type WsConnHub struct {
	Connections map[int64]*UserConn
	Locker      sync.RWMutex
	// redis
}

var OnlineUsers = WsConnHub{
	Connections: make(map[int64]*UserConn),
	Locker:      sync.RWMutex{},
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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("websocket failed to upgrade error: %v", err)
		return
	}
	defer conn.Close()

	userConn := &UserConn{
		ID:      id,
		Channel: make(chan string),
		WsConn:  conn,
	}

	OnlineUsers.Locker.Lock()
	OnlineUsers.Connections[id] = userConn
	OnlineUsers.Locker.Unlock()
	logger.Infof("websocket user id: %v", id)

	var wg sync.WaitGroup

	// 启动读写 goroutines
	wg.Add(2)
	go readLoop(userConn)
	go writeLoop(userConn)

	// 等待所有 goroutines 完成
	wg.Wait()

	defer func() {
		OnlineUsers.Locker.Lock()
		delete(OnlineUsers.Connections, id)
		OnlineUsers.Locker.Unlock()
		close(userConn.Channel)
		logger.Infof("websocket user close connection: %v", id)
	}()
}

// 从WebSocket读取消息
func readLoop(userConn *UserConn) {
	for {
		current_time := time.Now()
		_, msg, err := userConn.WsConn.ReadMessage()
		if err != nil {
			logger.Infof("websocket read error: %v", err)
			break
		}
		logger.Infof("websocket read message: %v", string(msg))
		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			logger.Infof("websocket unmarshal error: %v", err)
			break
		}
		destId, _ := strconv.ParseInt(message.ReceiverID, 10, 64)

		logger.Infof("websocket read message: %v", message)
		err = crud.MessageCRUD.Create(&models.MessageBasic{
			UserId:   userConn.ID,
			DestId:   destId,
			Content:  message.Content,
			Type:     message.Type,
			SeqNum:   0,
			Status:   0,
			SendTime: &current_time,
		})
		if err != nil {
			logger.Infof("websocket stroe message error: %v", err)
			break
		}
		logger.Infof("websocket stroe message success: %v", message)
		routeMessage(destId, string(msg))
	}
}

// 将消息路由到目标用户
func routeMessage(targetID int64, message string) {
	OnlineUsers.Locker.RLock()
	targetUser, ok := OnlineUsers.Connections[targetID]
	OnlineUsers.Locker.RUnlock()
	if ok {
		targetUser.Channel <- message
	} else {
		logger.Infof("websocket target user offline: %v", targetID)
		// 用户不在线，处理逻辑
	}
}

// 将消息写入WebSocket连接
func writeLoop(userConn *UserConn) {
	for msg := range userConn.Channel {
		err := userConn.WsConn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			logger.Infof("websocket write error: %v", err)
			break
		}
	}
}
