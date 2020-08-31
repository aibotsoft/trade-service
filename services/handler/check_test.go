package handler

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/trade/pkg/store"
	"github.com/aibotsoft/trade/services/auth"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var h *Handler

func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := store.New(cfg, log, db)
	au := auth.New(cfg, log, sto)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	au.AuthRound(ctx)
	cancel()
	h = New(cfg, log, sto, au)
	go h.WebsocketJob()
	m.Run()
	h.Close()
}

func TestHandler_Check(t *testing.T) {
	err := h.Check()
	assert.NoError(t, err)
	time.Sleep(time.Second*3)
}
