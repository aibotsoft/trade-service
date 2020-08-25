package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const host = "black.betinasia.com"

//type Ws struct {
//	cfg *config.Config
//	log *zap.SugaredLogger
//	//client  *client.Client
//	//store   *store.Store
//	//auth    *auth.Auth
//	//Conf    *config_client.ConfClient
//	//balance balance.Balance
//	//wsConn  *websocket.Conn
//}

func (h *Handler) WebsocketConnect() error {
	var url = fmt.Sprintf("wss://%v/cpricefeed/?token=%v&lang=en&prices_bookies=", host, h.auth.GetSession())
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	h.wsConn = conn
	return nil
}
func (h *Handler) WebsocketJob() {
	time.Sleep(time.Second)
	for {
		err := h.WebsocketConnect()
		if err != nil {
			h.log.Info(err)
			time.Sleep(time.Second * 20)
			continue
		}

		err = h.ReadLoop()
		if err != nil {
			h.log.Info(err)
			time.Sleep(time.Second * 10)
		}
	}
}
func (h *Handler) ReadLoop() error {
	var messages [][]interface{}
	for {
		_, message, err := h.wsConn.ReadMessage()
		if err != nil {
			return err
		}
		err = json.Unmarshal(message, &messages)
		if err != nil {
			return err
		}
		for _, m := range messages {
			//h.log.Infow("", "msg", m)
			if m[0] == "event" {
				h.processEvent(m)
			}

		}
	}
}
