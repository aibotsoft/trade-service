package ws

import (
	"github.com/aibotsoft/micro/config"
	"go.uber.org/zap"
)

type Ws struct {
	cfg *config.Config
	log *zap.SugaredLogger
	//client  *client.Client
	//store   *store.Store
	//auth    *auth.Auth
	//Conf    *config_client.ConfClient
	//balance balance.Balance
	//wsConn  *websocket.Conn
}
