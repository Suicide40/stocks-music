package streamer

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
)

type Tick struct {
	Type           string `json:"type"`
	EventId        int64  `json:"eventId"`
	Timestamp      int    `json:"timestamp"`
	Timestampms    int64  `json:"timestampms"`
	SocketSequence int    `json:"socket_sequence"`
	Events         []struct {
		Type      string `json:"type"`
		Side      string `json:"side"`
		Price     string `json:"price"`
		Remaining string `json:"remaining"`
		Delta     string `json:"delta"`
		Reason    string `json:"reason"`
	} `json:"events"`
}

type Gemini struct {
	dialer *websocket.Dialer
	c      *websocket.Conn
}

func NewGemini(dialer *websocket.Dialer) *Gemini {
	return &Gemini{dialer: dialer}
}

func (s *Gemini) Start(ticker string, output chan<- Tick) error {
	serverAddr := "wss://api.gemini.com/v1/marketdata/" + ticker
	conn, _, err := s.dialer.Dial(serverAddr, http.Header{})
	if err != nil {
		return errors.Wrap(err, "failed to connect")
	}
	s.c = conn

	for {
		_, data, err := s.c.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "failed to ReadMessage")
		}

		var tick Tick
		err = json.Unmarshal(data, &tick)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal tick")
		}

		output <- tick
	}
}

func (s *Gemini) Close() {
	if s.c != nil {
		s.c.Close()
	}
	fmt.Println("gemini closed")
}
