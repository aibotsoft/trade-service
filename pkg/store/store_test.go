package store

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/stretchr/testify/assert"
	"testing"
)



var s *Store
func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	s = New(cfg, log, db)
	m.Run()
	s.Close()
}

func TestStore_GetLiveEvent(t *testing.T) {
	events, err := s.GetLiveEvents()
	if assert.NoError(t, err) {
		t.Log(events)
	}
}
