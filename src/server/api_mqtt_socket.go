package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/processor"
	"net/http"
	"sync"
)

type mqttSocketServer struct {
	conns          map[*websocket.Conn]struct{}
	messageChannel messageChan
	mu             sync.RWMutex
	upgrader       websocket.Upgrader
}

type messageChan <-chan processor.MQTTMessage

func newMqttSocketServer(messageChannel messageChan) *mqttSocketServer {
	return &mqttSocketServer{
		conns:          make(map[*websocket.Conn]struct{}),
		messageChannel: messageChannel,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *mqttSocketServer) mqttLog(c echo.Context) error {
	if !c.IsWebSocket() {
		return ErrorResponse(c, http.StatusBadRequest, "Endpoint only supports Websocket")
	}

	conn, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	s.mu.Lock()
	s.conns[conn] = struct{}{}
	s.mu.Unlock()

	_, _, _ = conn.ReadMessage()

	s.cleanUpConnection(conn)

	return nil
}

func (s *mqttSocketServer) run() {
	type socketMessage struct {
		Server  string `json:"server"`
		Topic   string `json:"topic"`
		Payload string `json:"payload"`
		User    string `json:"user"`
	}

	go func() {
		for {
			m := <-s.messageChannel

			b, err := json.Marshal(socketMessage{
				Server:  m.Server,
				Topic:   m.Topic,
				Payload: m.Payload,
				User:    m.User,
			})

			if err != nil {
				continue
			}

			s.mu.RLock()

			for conn := range s.conns {
				if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
					s.cleanUpConnection(conn)
				}
			}
			s.mu.RUnlock()
		}
	}()
}

func (s *mqttSocketServer) cleanUpConnection(conn *websocket.Conn) {
	s.mu.Lock()
	_ = conn.Close()
	delete(s.conns, conn)
	s.mu.Unlock()
}
