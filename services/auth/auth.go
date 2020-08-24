package auth

import (
	"github.com/aibotsoft/micro/config"
	"go.uber.org/zap"
)

type Auth struct {
	cfg *config.Config
	log *zap.SugaredLogger
	//store         *store.Store
	//client        *client.Client
	//conf          *config_client.ConfClient
	//account       pb.Account
	//token         token.Token
	//bettingStatus bool
	//tokenLock     sync.Mutex
}
