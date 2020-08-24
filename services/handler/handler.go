package handler

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/trade/pkg/client"
	"go.uber.org/zap"
)

type Handler struct {
	cfg    *config.Config
	log    *zap.SugaredLogger
	client *client.Client
	//store   *store.Store
	//auth    *auth.Auth
	//Conf    *config_client.ConfClient
	//balance balance.Balance
	//wsConn  *websocket.Conn
}
