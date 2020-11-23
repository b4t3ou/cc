package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/b4t3ou/cc/service"
)

const (
	url = "wss://streamer.cryptocompare.com/v2?api_key="
)

type Stream struct {
	apiKey    string
	stream    *websocket.Conn
	interrupt chan os.Signal
	done      chan struct{}
	param     service.JSONParams
}

func NewStream(apiKey string, param service.JSONParams) (*Stream, error) {
	var (
		err error
	)

	s := &Stream{
		apiKey:    apiKey,
		interrupt: make(chan os.Signal, 1),
		done:      make(chan struct{}),
		param:     param,
	}

	s.stream, _, err = websocket.DefaultDialer.Dial(url+s.apiKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Dial error: %+v", err)
	}

	return s, nil
}

func (sc Stream) Init() (err error) {
	signal.Notify(sc.interrupt, os.Interrupt)

	if err = sc.initStream(); err != nil {
		return err
	}

	defer sc.stream.Close()

	for {
		select {
		case <-sc.done:
			return nil
		case <-sc.interrupt:
			log.Println("interrupt")
			err := sc.stream.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				return fmt.Errorf("failed to close: %+v", err)
			}

			select {
			case <-sc.done:
			case <-time.After(time.Second):
			}

			return nil
		}
	}
}

func (sc Stream) Process(callbacks ...func([]byte) error) {
	defer close(sc.done)

	for {
		_, message, err := sc.stream.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		for _, cb := range callbacks {
			if err := cb(message); err != nil {
				log.Printf("failed to process message: %+v", err)
				return
			}
		}
	}
}

func (sc Stream) initStream() error {
	s, _ := json.Marshal(sc.param)

	if err := sc.stream.WriteMessage(websocket.TextMessage, []byte(string(s))); err != nil {
		return fmt.Errorf("failed to init stream: %+v", err)
	}

	return nil
}
