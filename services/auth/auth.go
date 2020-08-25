package auth

import (
	"context"
	"database/sql"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/trade/pkg/client"
	"github.com/aibotsoft/trade/pkg/store"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

const checkLoginPeriod = time.Minute * 3

type Auth struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	store   *store.Store
	client  *client.Client
	account store.Account
	Token   store.Token
	//bettingStatus bool
	tokenLock sync.Mutex
}

func New(cfg *config.Config, log *zap.SugaredLogger, store *store.Store) *Auth {
	return &Auth{cfg: cfg, log: log, store: store, client: client.New(cfg, log)}
}
func (a *Auth) GetAccount(ctx context.Context) (err error) {
	if a.account.Id == 0 {
		account, err := a.store.LoadAccount(ctx)
		if err != nil {
			return err
		}
		a.account = account
	}
	return
}

func (a *Auth) AuthJob() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := a.AuthRound(ctx)
		cancel()
		if err != nil {
			a.log.Error(err)
			time.Sleep(time.Second * 20)
		} else {
			time.Sleep(checkLoginPeriod)
		}
	}
}
func (a *Auth) GetSession() string {
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()
	return a.Token.Session
}
func (a *Auth) AuthRound(ctx context.Context) (err error) {
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()
	err = a.GetAccount(ctx)
	if err != nil {
		return errors.Wrap(err, "get_account_error")
	}
	a.Token, err = a.store.LoadToken(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			a.log.Info("no_token_in_db_begin_login")
			err = a.Login(ctx)
			if err != nil {
				return
			}
		} else {
			return errors.Wrap(err, "load_token_error")
		}
	}
	err = a.CheckLogin(ctx)
	if err != nil {
		a.log.Info("check_login_error_begin_login")
		err = a.Login(ctx)
		return
	}
	return
}
func (a *Auth) Login(ctx context.Context) (err error) {
	resp, err := a.client.Login(ctx, a.account.Username, a.account.Password)
	if err != nil {
		return errors.Wrap(err, "login_error")
	}
	data := resp.GetData()
	session := data.GetSessionId()
	a.log.Infow("login_ok", "session_id", session)
	a.Token.Session = session
	a.Token.CreatedAt = time.Now()
	a.Token.LastCheckAt = time.Now()
	err = a.store.SaveToken(ctx, a.Token)
	return err
}
func (a *Auth) CheckLogin(ctx context.Context) (err error) {
	resp, err := a.client.CheckLogin(ctx, a.Token.Session)
	if err != nil {
		return
	}
	data := resp.GetData()
	if data.GetUsername() == a.account.Username {
		a.Token.LastCheckAt = time.Now()
		err := a.store.UpdateToken(ctx, a.Token)
		if err != nil {
			a.log.Error(err)
		}
	}
	a.log.Infow("check_login_ok", "resp", resp)
	return
}

//user := os.Getenv("user")
//pass := os.Getenv("pass")
//testSession := os.Getenv("test_session")
//a.log.Info(testSession)
