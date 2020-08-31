package handler

import (
	"fmt"
	"github.com/aibotsoft/micro/util"
	"time"
)

func (h *Handler) CheckLoop() {
	h.store.DeleteBetSlips()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for  {
		select {
		case <-ticker.C:
			//h.log.Info(t)
			surebetId := util.UnixMsNow()
			sb, err := h.store.GetDemoSurebet(0.2)
			if err != nil {
				continue
			}
			betTypes := []string{sb.HomeBetType, sb.AwayBetType}
			for i := range betTypes {
				betSlipId := h.store.HasBetSlip(sb.PeriodCode, sb.EventId, betTypes[i])
				if betSlipId == "" {
					h.log.Infow("", "", sb, "betSlipId", betSlipId, "sId", surebetId)
					slip, err := h.GetBetSlip(sb.PeriodCode, sb.EventId, betTypes[i])
					if err != nil {
						h.log.Error(err)
					} else {
						h.store.SaveBetSlip(slip)
					}
				} else {
					price, err := h.store.GetPrice(betSlipId)
					if err != nil {
						h.log.Error(err)
					} else {
						h.log.Infow("price", "", price)
					}
				}
				//h.store.SaveSurebet()
			}
		}
	}
}
func (h *Handler) Check() error {
	//auth, err := h.auth.Auth(context.Background())
	//if err != nil {
	//	return err
	//}
	sb, err := h.store.GetDemoSurebet(-10)
	if err != nil {
		return err
	}
	h.log.Infow("", "", sb)
	betType := fmt.Sprintf("for,ah,h,%d", sb.HandicapCode)
	slip, err := h.GetBetSlip(sb.PeriodCode, sb.EventId, betType)
	if err != nil {
		return err
	}
	h.store.SaveBetSlip(slip)
	return nil
}
//func (h *Handler) Check() error {
//	auth, err := h.auth.Auth(context.Background())
//	if err != nil {
//		return err
//	}
//	sb, err := h.store.GetDemoSurebet(-10)
//	if err != nil {
//		return err
//	}
//	h.log.Infow("", "", sb)
//	betType := fmt.Sprintf("for,ah,h,%d", sb.HandicapCode)
//	slip, err := h.GetBetSlip(auth, sb.PeriodCode, sb.EventId, betType)
//	if err != nil {
//		return err
//	}
//	h.store.SaveBetSlip(slip)
//	return nil
//}
