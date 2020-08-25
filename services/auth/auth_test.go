package auth

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/trade/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var a *Auth

func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := store.New(cfg, log, db)
	a = New(cfg, log, sto)
	m.Run()
}

func TestAuth_CheckLogin(t *testing.T) {
	err := a.CheckLogin(context.Background())
	assert.NoError(t, err)
}

func TestAuth_Login(t *testing.T) {
	err := a.Login(context.Background())
	assert.NoError(t, err)
}

func TestAuth_LoginRound(t *testing.T) {
	err := a.AuthRound(context.Background())
	assert.NoError(t, err)

}
