package handler

import (
	"context"
	api "github.com/aibotsoft/gen/blackapi"
	"github.com/pkg/errors"
)

func (h *Handler) GetBetSlip(ctx context.Context, sport string, eventId string, betType string) (api.BetSlipData, error) {
	h.log.Infow("get_bet_slip", "sport", sport, "event", eventId, "betType", betType)
	resp, err := h.client.BetSlip(ctx, sport, eventId, betType, false)
	if err != nil {
		//h.log.Infow(side.Check.StatusInfo, "err", err, "urlEvent", u, "sport", side.SportName, "league", side.LeagueName,
		//	"home", side.Home, "away", side.Away, "market", side.MarketName, "fortedPrice", side.Price)
		return api.BetSlipData{}, err
	}
	if resp.GetStatus() != "ok" {
		h.log.Infow("status_not_ok", "resp", resp)
		return api.BetSlipData{}, errors.Errorf("status_not_ok")
	}
	data := resp.GetData()
	//h.log.Infow("bet_slip_response", "slipId", data.GetBetslipId(), "bet_type", data.GetBetType(),
	//	"bookies_with_offers_len", len(data.GetBookiesWithOffers()))
	return data, nil
}
