package auth

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/trade/pkg/client"
	"github.com/aibotsoft/trade/pkg/store"
	"go.uber.org/zap"
	"os"
	"sync"
)

type Auth struct {
	cfg    *config.Config
	log    *zap.SugaredLogger
	store  *store.Store
	client *client.Client
	//account       pb.Account
	//token         token.Token
	//bettingStatus bool
	tokenLock sync.Mutex
}

func New(cfg *config.Config, log *zap.SugaredLogger, store *store.Store) *Auth {
	return &Auth{cfg: cfg, log: log, store: store, client: client.New(cfg, log)}
}

func (a *Auth) LoginRound(ctx context.Context) (err error) {
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()
	err = a.CheckLogin(ctx)
	return
}
func (a *Auth) Login(ctx context.Context) (id string, err error) {
	user := os.Getenv("user")
	pass := os.Getenv("pass")
	resp, err := a.client.Login(ctx, user, pass)
	if err != nil {
		return
	}
	data := resp.GetData()
	a.log.Infow("", "", data.GetSessionId())
	return data.GetSessionId(), nil
}
func (a *Auth) CheckLogin(ctx context.Context) (err error) {
	testSession := os.Getenv("test_session")
	a.log.Info(testSession)
	resp, err := a.client.CheckLogin(ctx, testSession)
	if err != nil {
		return
	}
	a.log.Info(resp)
	return
}
