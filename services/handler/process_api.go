package handler

import (
	"github.com/aibotsoft/trade/pkg/store"
)

func (h *Handler) processApi(m []interface{}) {
	//h.log.Infow("api", "msg", m)
	apiMap, ok := m[1].(map[string]interface{})
	if !ok {
		h.log.Infow("api_not_ok", "m", m)
		return
	}
	data, ok := apiMap["data"].([]interface{})
	if !ok {
		h.log.Infow("data_not_ok", "m", m)
		return
	}
	for i := range data {
		d := data[i].([]interface{})
		if d[0] == "pmm" {
			prices := h.processPmm(d)
			//h.log.Infow("", "prices", prices)
			for i := range prices {
				if prices[i].Price > 0 {
					h.store.SavePrice(prices[i])
				} else {
					h.store.DeactivatePrice(prices[i])
				}
			}
		}
	}
}

func (h *Handler) processPmm(d []interface{}) (prices []store.Price) {
	pmm := d[1].(map[string]interface{})
	pl := pmm["price_list"].([]interface{})
	if len(pl) == 0 {
		return []store.Price{{
			BetslipId: pmm["betslip_id"].(string),
			BetType:   pmm["bet_type"].(string),
			Bookie:    pmm["bookie"].(string),
			EventId:   pmm["event_id"].(string),
			Sport:     pmm["sport"].(string),
			Status:    "",
			Username:  pmm["username"].(string),
		}}
	}
	for i := range pl {
		priceMap := pl[i].(map[string]interface{})
		bookie := priceMap["bookie"].(map[string]interface{})
		effective := priceMap["effective"].(map[string]interface{})
		prices = append(prices, store.Price{
			BetslipId: pmm["betslip_id"].(string),
			BetType:   pmm["bet_type"].(string),
			Bookie:    pmm["bookie"].(string),
			EventId:   pmm["event_id"].(string),
			Sport:     pmm["sport"].(string),
			Status:    pmm["status"].(map[string]interface{})["code"].(string),
			Username:  pmm["username"].(string),
			Num:       i,
			Price:     effective["price"].(float64),
			Min: bookie["min"].([]interface{})[1].(float64),
			Max: bookie["max"].([]interface{})[1].(float64),
		})
		//h.log.Infow("pmm", "", priceMap, "p", p)
	}
	return
}
