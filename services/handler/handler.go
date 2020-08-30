package handler

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/trade/pkg/client"
	"github.com/aibotsoft/trade/pkg/store"
	"github.com/aibotsoft/trade/services/auth"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	cfg    *config.Config
	log    *zap.SugaredLogger
	client *client.Client
	store  *store.Store
	auth   *auth.Auth
	//Conf    *config_client.ConfClient
	//balance balance.Balance
	wsConn *websocket.Conn
	read   chan []byte
	write   chan []byte
}

func New(cfg *config.Config, log *zap.SugaredLogger, store *store.Store, auth *auth.Auth) *Handler {
	h := &Handler{cfg: cfg, log: log, client: client.New(cfg, log), store: store, auth: auth}
	h.read = make(chan []byte, 100000)
	h.write = make(chan []byte, 100000)
	return h
}
func (h *Handler) Close() {
	h.store.Close()
	if h.wsConn != nil {
		err := h.wsConn.Close()
		if err != nil {
			h.log.Error(err)
		}
	}
}


