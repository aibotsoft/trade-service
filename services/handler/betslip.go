package handler

import (
	"context"
	api "github.com/aibotsoft/gen/blackapi"
	"github.com/pkg/errors"
)

func (h *Handler) GetBetSlip(sport string, eventId string, betType string) (slip api.BetSlipData, err error) {
	auth, err := h.auth.Auth(context.Background())
	if err != nil {
		return
	}
	//h.log.Infow("get_bet_slip", "sport", sport, "event", eventId, "betType", betType)
	resp, err := h.client.BetSlip(auth, sport, eventId, betType, false)
	if err != nil {
		//h.log.Infow(side.Check.StatusInfo, "err", err, "urlEvent", u, "sport", side.SportName, "league", side.LeagueName,
		//	"home", side.Home, "away", side.Away, "market", side.MarketName, "fortedPrice", side.Price)
		return
	}
	if resp.GetStatus() != "ok" {
		h.log.Infow("status_not_ok", "resp", resp)
		return slip, errors.Errorf("status_not_ok")
	}
	return resp.GetData(), nil
}
