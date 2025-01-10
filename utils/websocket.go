package utils

import (
	"encoding/json"
	"sync"
	// "time"

	"shuttle/logger"
	"shuttle/repositories"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketServiceInterface interface {
	HandleWebSocketConnection(c *websocket.Conn)
}

type WebSocketService struct {
	userRepository repositories.UserRepositoryInterface
	authRepository repositories.AuthRepositoryInterface
}

func NewWebSocketService(userRepository repositories.UserRepositoryInterface, authRepository repositories.AuthRepositoryInterface) WebSocketServiceInterface {
	return &WebSocketService{
		userRepository: userRepository,
		authRepository: authRepository,
	}
}

var (
	activeConnections = make(map[string]*websocket.Conn) // Save active WebSocket connections
	mutex             = &sync.Mutex{}                    // Ensure atomic operations
)

func AddConnection(ID string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	activeConnections[ID] = conn
}

func RemoveConnection(ID string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(activeConnections, ID)
}

func GetConnection(ID string) (*websocket.Conn, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	conn, exists := activeConnections[ID]
	return conn, exists
}

// Handle WebSocket connection
var (
	shuttleGroups = make(map[string]map[string]*websocket.Conn) // Save active WebSocket connections
	groupMutex    = &sync.Mutex{}                               // Ensure atomic operations
)

func AddToShuttleGroup(shuttleUUID, userUUID string, conn *websocket.Conn) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	if _, exists := shuttleGroups[shuttleUUID]; !exists {
		shuttleGroups[shuttleUUID] = make(map[string]*websocket.Conn)
	}
	shuttleGroups[shuttleUUID][userUUID] = conn
}

func RemoveFromShuttleGroup(shuttleUUID, userUUID string) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	if group, exists := shuttleGroups[shuttleUUID]; exists {
		delete(group, userUUID)
		if len(group) == 0 {
			delete(shuttleGroups, shuttleUUID)
		}
	}
}

func BroadcastToShuttleGroup(shuttleUUID string, message []byte) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	if group, exists := shuttleGroups[shuttleUUID]; exists {
		for _, conn := range group {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.LogError(err, "WebSocket Broadcast Error", nil)
			}
		}
	}
}

func (s *WebSocketService) HandleWebSocketConnection(c *websocket.Conn) {
	userUUID := c.Params("id")
	shuttleUUID := c.Query("shuttle_uuid")

	// Validate user access to shuttle group
	// if !s.userRepository.HasAccessToShuttle(userUUID, shuttleUUID) {
	// 	c.WriteMessage(websocket.TextMessage, []byte("Unauthorized access to shuttle group"))
	// 	c.Close()
	// 	return
	// }

	AddToShuttleGroup(shuttleUUID, userUUID, c)
	defer func() {
		RemoveFromShuttleGroup(shuttleUUID, userUUID)
		logger.LogInfo("WebSocket Connection Removed from Group", map[string]interface{}{"ShuttleUUID": shuttleUUID, "UserUUID": userUUID})
	}()

	logger.LogInfo("WebSocket Connection Added to Group", map[string]interface{}{"ShuttleUUID": shuttleUUID, "UserUUID": userUUID})
	c.WriteMessage(websocket.TextMessage, []byte("Connected to shuttle group"))

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			logger.LogError(err, "WebSocket Error Reading Message", nil)
			break
		}

		var data struct {
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		}

		if err := json.Unmarshal(msg, &data); err != nil || data.Longitude == 0 || data.Latitude == 0 {
			errorResponse := struct {
				Code    int    `json:"code"`
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Code:    400,
				Status:  "Bad Request",
				Message: "Invalid message format. Must contain 'longitude' and 'latitude'.",
			}
			responseMsg, _ := json.Marshal(errorResponse)
			c.WriteMessage(websocket.TextMessage, responseMsg)
			continue
		}

		logger.LogInfo("Broadcasting Message", map[string]interface{}{
			"ShuttleUUID": shuttleUUID,
			"UserUUID":    userUUID,
			"Longitude":   data.Longitude,
			"Latitude":    data.Latitude,
		})
		BroadcastToShuttleGroup(shuttleUUID, msg)

		response := struct {
			Code    int    `json:"code"`
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Code:    200,
			Status:  "OK",
			Message: "Message broadcasted to shuttle group",
		}
		responseMsg, _ := json.Marshal(response)
		c.WriteMessage(mt, responseMsg)
	}
}