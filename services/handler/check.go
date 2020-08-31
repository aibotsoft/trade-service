package handler

import (
	"context"
	"fmt"
)

func (h *Handler) Check() error {
	auth, err := h.auth.Auth(context.Background())
	if err != nil {
		return err
	}
	sb, err := h.store.GetDemoSurebet()
	if err != nil {
		return err
	}
	h.log.Infow("", "", sb)
	betType:=fmt.Sprintf("for,ah,h,%d", sb.HandicapCode)
	slip, err := h.GetBetSlip(auth, sb.PeriodCode, sb.EventId, betType)
	if err != nil {
		return err
	}
	h.log.Infow("", "slip", slip)
	return nil
}

