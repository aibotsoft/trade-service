package handler

import (
	"context"
	"fmt"
	"time"
)

func (h *Handler) CheckLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for  {
		select {
		case <-ticker.C:
			//h.log.Info(t)
			sb, err := h.store.GetDemoSurebet(0)
			if err != nil {
				continue
			}
			h.log.Infow("", "", sb)

		}
	}
}
func (h *Handler) Check() error {
	auth, err := h.auth.Auth(context.Background())
	if err != nil {
		return err
	}
	sb, err := h.store.GetDemoSurebet(-10)
	if err != nil {
		return err
	}
	h.log.Infow("", "", sb)
	betType := fmt.Sprintf("for,ah,h,%d", sb.HandicapCode)
	slip, err := h.GetBetSlip(auth, sb.PeriodCode, sb.EventId, betType)
	if err != nil {
		return err
	}
	h.log.Infow("", "slip", slip)
	//h.store.SaveSlip()
	return nil
}
