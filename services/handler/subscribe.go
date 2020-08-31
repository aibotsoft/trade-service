package handler

import (
	"fmt"
	"time"
)

const SubscribeTimeOut = time.Second * 30

func (h *Handler) SubscribeLoop() {
	//["watch_event",[430,"fb","2020-08-29,27263,27262"]]
	for {
		events, err := h.store.GetLiveEvents()
		if err != nil {
			h.log.Error(err)
		} else {
			for i := range events {
				//h.log.Info(events[i])
				we := fmt.Sprintf(`["watch_event",[%d,"%s","%s"]]`, events[i].LeagueId, events[i].PeriodCode, events[i].Id)
				_, b := h.store.Cache.Get(we)
				if b {
					//h.log.Info("from_cache: ", we)
					continue
				}
				h.write <- []byte(we)
				h.store.Cache.SetWithTTL(we, true, 1, time.Hour)
			}
		}
		time.Sleep(SubscribeTimeOut)
	}
}
