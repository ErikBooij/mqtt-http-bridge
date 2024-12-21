package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"mqtt-http-bridge/src/processor"
	"net/http"
	"slices"
	"strconv"
	"sync"
)

const defaultHistorySize = 50

type mqttSocketServer struct {
	conns          map[*websocket.Conn]struct{}
	history        map[string][]historyEntry
	historySize    int
	messageChannel messageChan
	mu             sync.RWMutex
	sequence       int
	upgrader       websocket.Upgrader
}

type historyEntry struct {
	message  processor.MQTTMessage
	sequence int
}

type messageChan <-chan processor.MQTTMessage

func newMqttSocketServer(messageChannel messageChan) *mqttSocketServer {
	return &mqttSocketServer{
		conns:          make(map[*websocket.Conn]struct{}),
		history:        make(map[string][]historyEntry),
		historySize:    defaultHistorySize,
		messageChannel: messageChannel,
		sequence:       0,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *mqttSocketServer) run() {
	go func() {
		for {
			m := <-s.messageChannel

			entry := s.addMessageToHistory(m)

			b := s.historyEntryToSocketMessage(entry)

			if b == nil {
				continue
			}

			s.mu.RLock()
			for conn := range s.conns {
				if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
					_ = conn.Close()
					delete(s.conns, conn)
				}
			}
			s.mu.RUnlock()
		}
	}()
}

func (s *mqttSocketServer) mqttLog(c echo.Context) error {
	if !c.IsWebSocket() {
		return ErrorResponse(c, http.StatusBadRequest, "Endpoint only supports Websocket")
	}

	conn, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	lastSeenSequence := 0

	if lss := c.QueryParam("last_seen_sequence"); lss != "" {
		if l, err := strconv.Atoi(lss); err == nil {
			lastSeenSequence = l
		}
	}

	s.publishHistory(conn, lastSeenSequence)

	s.mu.Lock()
	s.conns[conn] = struct{}{}
	s.mu.Unlock()

	_, _, _ = conn.ReadMessage()

	s.mu.Lock()
	_ = conn.Close()
	delete(s.conns, conn)
	s.mu.Unlock()

	return nil
}

func (s *mqttSocketServer) addMessageToHistory(m processor.MQTTMessage) historyEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence++

	entry := historyEntry{
		message:  m,
		sequence: s.sequence,
	}

	s.history[m.Server] = append(s.history[m.Server], entry)

	if len(s.history[m.Server]) > s.historySize {
		s.history[m.Server] = s.history[m.Server][1:]
	}

	return entry
}

func (s *mqttSocketServer) historyEntryToSocketMessage(e historyEntry) []byte {
	type socketMessage struct {
		Server   string `json:"server"`
		Topic    string `json:"topic"`
		Payload  string `json:"payload"`
		User     string `json:"user"`
		Sequence int    `json:"sequence"`
	}

	b, err := json.Marshal(socketMessage{
		Server:   e.message.Server,
		Topic:    e.message.Topic,
		Payload:  e.message.Payload,
		User:     e.message.User,
		Sequence: e.sequence,
	})

	if err != nil {
		return nil
	}

	return b
}

func (s *mqttSocketServer) publishHistory(conn *websocket.Conn, lastSeenSequence int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages := make([]historyEntry, 0, len(s.history)*s.historySize)

	for _, historyMessages := range s.history {
		messages = append(messages, historyMessages...)
	}

	slices.SortFunc(messages, func(a, b historyEntry) int {
		return a.sequence - b.sequence
	})

	for _, m := range messages {
		if m.sequence <= lastSeenSequence {
			continue
		}

		b := s.historyEntryToSocketMessage(m)

		if b == nil {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
			continue
		}
	}
}
