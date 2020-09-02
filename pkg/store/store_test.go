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

func TestStore_SaveSurebet(t *testing.T) {
	err := s.SaveSurebet(Side{
		SurebetId:     0,
		BetslipId:     "",
		Price:         0,
		BestPrice:     0,
		WeightedPrice: 0,
		Min:           0,
		Max:           0,
		Volume:        0,
		Bookie:        "",
		PriceCount:    0,
	})
	assert.NoError(t, err)
}