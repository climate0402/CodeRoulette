package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	PlayerID  string      `json:"player_id,omitempty"`
	MatchID   string      `json:"match_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// handleWebSocket handles WebSocket connections for real-time match communication
func (h *Handlers) handleWebSocket(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established for room: %s", roomID)

	// Handle WebSocket messages
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Process message based on type
		switch msg.Type {
		case "join_room":
			h.handleJoinRoom(conn, roomID, &msg)
		case "code_submission":
			h.handleCodeSubmission(conn, roomID, &msg)
		case "skill_card_use":
			h.handleSkillCardUse(conn, roomID, &msg)
		case "ping":
			h.handlePing(conn, &msg)
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

// handleJoinRoom handles when a player joins a match room
func (h *Handlers) handleJoinRoom(conn *websocket.Conn, roomID string, msg *WebSocketMessage) {
	response := WebSocketMessage{
		Type:      "room_joined",
		Data:      map[string]string{"room_id": roomID, "status": "success"},
		Timestamp: getCurrentTimestamp(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}

	log.Printf("Player %s joined room %s", msg.PlayerID, roomID)
}

// handleCodeSubmission handles code submission via WebSocket
func (h *Handlers) handleCodeSubmission(conn *websocket.Conn, roomID string, msg *WebSocketMessage) {
	// Broadcast submission to other players in the room
	// broadcastMsg := WebSocketMessage{
	// 	Type:      "code_submitted",
	// 	Data:      msg.Data,
	// 	PlayerID:  msg.PlayerID,
	// 	MatchID:   msg.MatchID,
	// 	Timestamp: getCurrentTimestamp(),
	// }

	// In a real implementation, you would broadcast this to all connections in the room
	// For now, we'll just send a confirmation back
	response := WebSocketMessage{
		Type:      "submission_received",
		Data:      map[string]string{"status": "received"},
		Timestamp: getCurrentTimestamp(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}

	log.Printf("Code submission received from player %s in room %s", msg.PlayerID, roomID)
}

// handleSkillCardUse handles skill card usage via WebSocket
func (h *Handlers) handleSkillCardUse(conn *websocket.Conn, roomID string, msg *WebSocketMessage) {
	// Broadcast skill card usage to other players
	// broadcastMsg := WebSocketMessage{
	// 	Type:      "skill_card_used",
	// 	Data:      msg.Data,
	// 	PlayerID:  msg.PlayerID,
	// 	MatchID:   msg.MatchID,
	// 	Timestamp: getCurrentTimestamp(),
	// }

	// In a real implementation, you would broadcast this to all connections in the room
	response := WebSocketMessage{
		Type:      "skill_card_acknowledged",
		Data:      map[string]string{"status": "acknowledged"},
		Timestamp: getCurrentTimestamp(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}

	log.Printf("Skill card used by player %s in room %s", msg.PlayerID, roomID)
}

// handlePing handles ping messages for connection health
func (h *Handlers) handlePing(conn *websocket.Conn, msg *WebSocketMessage) {
	response := WebSocketMessage{
		Type:      "pong",
		Data:      map[string]string{"status": "ok"},
		Timestamp: getCurrentTimestamp(),
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// getCurrentTimestamp returns current timestamp in milliseconds
func getCurrentTimestamp() int64 {
	return 0 // Placeholder - implement actual timestamp
}
