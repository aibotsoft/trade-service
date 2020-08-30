package handler

import (
	"fmt"
	"time"
)

func (h *Handler) SubscribeLoop() {
	//["watch_event",[430,"fb","2020-08-29,27263,27262"]]
	for  {
		events, err := h.store.GetLiveEvents()
		if err != nil {
			h.log.Error(err)
		} else {
			for i := range events {
				//h.log.Info(events[i])
				we := fmt.Sprintf(`["watch_event",[%d,"%s","%s"]]`, events[i].LeagueId, events[i].PeriodCode, events[i].Id)
				h.write <- []byte(we)
				//break
			}
		}
		time.Sleep(time.Hour)
	}
}
