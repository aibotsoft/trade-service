package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const host = "black.betinasia.com"

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
	err := h.WebsocketConnect()
	if err != nil {
		h.log.Info(err)
		//time.Sleep(time.Second * 20)
		//continue
	}
	go h.ReadLoop()
	go h.Read()
	h.Write()
	//if err != nil {
	//	h.log.Info(err)
	//	time.Sleep(time.Second * 10)
	//}
}

//var inMsg chan []byte

func (h *Handler) ReadLoop() {
	for {
		_, message, err := h.wsConn.ReadMessage()
		if err != nil {
			h.log.Info(err)
		} else {
			h.read <- message
		}
	}
}
func (h *Handler) Read() {
	var messages [][]interface{}
	for {
		select {
		case msg := <-h.read:
			//h.log.Info("new_msg: ", msg)
			err := json.Unmarshal(msg, &messages)
			if err != nil {
				h.log.Info(err)
				continue
			}
			for _, m := range messages {
				//h.log.Infow("", "msg", m)
				if m[0] == "event" {
					h.processEvent(m)
				} else if m[0] == "offers_event" {
					h.processOffersEvent(m)
				} else if m[0] == "offers_hcap" {
					h.log.Infow("offers_hcap", "msg", m)
				}
			}
		}
	}
}
func (h *Handler) Write() {
	for {
		select {
		case msg := <-h.write:
			//h.log.Infow("got_write", "msg", msg)
			err := h.wsConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				h.log.Error(err)
			}
		}
	}
}
