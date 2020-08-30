package handler

import (
	"github.com/aibotsoft/micro/util"
	"strconv"
	"strings"
)

var sportPeriodMap = map[string]string{
	"fb":         "Football",
	"fb_ht":      "Football",
	"fb_corn":    "Football",
	"fb_corn_ht": "Football",
	"fb_htft":    "Football",
	"fb_et":      "Football",
	"basket":     "Basketball",
	"basket_ht":  "Basketball",
	"esports":    "eSports",
	"tennis":     "Tennis",
	"cricket":    "Cricket",
	"baseball":   "Baseball",
	"ih":         "Ice Hockey",
	"mma":        "MMA",
	"boxing":     "Boxing",
	"ru":         "Rugby Union",
	"rl":         "Rugby League",
	"af":         "American Football",
}
var sportMap = map[string]int64{
	"Football":          1,
	"Basketball":        2,
	"eSports":           3,
	"Tennis":            4,
	"Cricket":           5,
	"Baseball":          6,
	"Ice Hockey":        7,
	"MMA":               8,
	"Boxing":            9,
	"Rugby Union":       10,
	"Rugby League":      11,
	"American Football": 12,
}

func (h *Handler) processEvent(m []interface{}) {
	sportList := m[1].([]interface{})
	sportCode := sportList[0].(string)
	sportName := sportPeriodMap[sportCode]
	sportId := sportMap[sportName]

	eventId := sportList[1].(string)
	eventSplit := strings.Split(eventId, ",")
	if len(eventSplit) != 3 {
		h.log.Info("split_event_error: ", eventId)
		return
	}
	homeId, _ := strconv.ParseInt(eventSplit[1], 10, 64)
	awayId, _ := strconv.ParseInt(eventSplit[2], 10, 64)

	eventList, ok := m[2].(map[string]interface{})
	if !ok {
		//h.log.Infow("eventList_not_ok", "m[1]", m[1], "m[2]", m[2])
		h.store.SaveEventPeriod(eventId, sportCode, false)
		return
	}
	leagueName := eventList["competition_name"].(string)
	if leagueName == "" {
		h.log.Infow("competition_name_blank", "eventList", eventList)
		return
	} else if strings.Index(strings.ToLower(leagueName), "test") != -1 {
		//h.log.Infow("competition_name_test", "eventList", eventList)
		return
	}

	home := eventList["home"].(string)
	if home == "" {
		h.log.Infow("home_blank", "eventList", eventList)
		return
	} else if strings.Index(strings.ToLower(home), "test") != -1 {
		h.log.Infow("home_test", "eventList", eventList)
		return
	}
	away := eventList["away"].(string)
	if away == "" {
		h.log.Infow("away_blank", "eventList", eventList)
		return
	} else if strings.Index(strings.ToLower(away), "test") != -1 {
		h.log.Infow("away_test", "eventList", eventList)
		return
	}

	leagueId := int64(eventList["competition_id"].(float64))
	//irStatus, ok := eventList["ir_status"].(map[string]interface{})
	if ok {
		//h.log.Infow("",
		//	"s", sportCode,
		//	"eventId", eventId,
		//	"home", home,
		//	"away", away,
		//	//"competitionId", leagueId,
		//	//"leagueName", leagueName,
		//	//"country", country,
		//	//"starts", starts,
		//	//"eventList", eventList,
		//	"ir_status", irStatus,
		//)
	}

	country := eventList["country"].(string)
	starts, ok := eventList["start_ts"].(string)
	if !ok {
		h.log.Info("starts_not_ok: ", eventList)
		return
	}
	h.store.SaveSport(sportId, sportName)
	h.store.SaveTeam(homeId, home)
	h.store.SaveTeam(awayId, away)
	h.store.SaveLeague(leagueId, leagueName, country, sportId)
	h.store.SaveEvent(eventId, homeId, awayId, leagueId, starts)
	h.store.SaveEventPeriod(eventId, sportCode, true)
}
func (h *Handler) processOffersEvent(m []interface{}) {
	//h.log.Infow("offers_event", "msg", m)
	sportList := m[1].([]interface{})
	leagueId := int64(sportList[0].(float64))
	sportCode := sportList[1].(string)
	eventId := sportList[2].(string)
	priceList := m[2].(map[string]interface{})
	if len(priceList) == 0 {
		h.store.SaveEventPeriod(eventId, sportCode, false)
		return
	}
	for key, value := range priceList {
		switch key {
		case "wdw":
			h.wdw(leagueId, sportCode, eventId, value)
		case "ah":
			h.ah(leagueId, sportCode, eventId, value)
		}
	}
}
func (h *Handler) ah(leagueId int64, sportCode string, eventId string, value interface{}) {
	valueList := value.([]interface{})
	for i := range valueList {
		priceList := valueList[i].([]interface{})
		handicap := priceList[0].(float64)
		sideList, ok := priceList[1].([]interface{})
		if !ok {
			h.log.Infow("ah_not_ok", "handicap", handicap, "side", priceList[1], "", sideList)
			continue
		}
		away := sideList[0].([]interface{})[1].(float64)
		if away == 0 {
			away = 1
		}
		home := sideList[1].([]interface{})[1].(float64)
		if home == 0 {
			home = 1
		}
		margin := util.TruncateFloat(1/(1/away+1/home)*100-100, 3)
		h.log.Infow("ah", "handicap", handicap, "away", away, "home", home, "margin", margin)

	}
	//sideList, ok := priceList[1].([]interface{})
}

func (h *Handler) wdw(leagueId int64, sportCode string, eventId string, value interface{}) {
	valueList := value.([]interface{})
	priceList := valueList[0].([]interface{})
	sideList, ok := priceList[1].([]interface{})
	if !ok {
		h.log.Infow("wdw_not_ok", "s", sportCode, "eventId", eventId,"v", value)
		return
	}
	away := sideList[0].([]interface{})[1].(float64)
	if away == 0 {
		away = 1
	}
	draw := sideList[1].([]interface{})[1].(float64)
	if draw == 0 {
		draw = 1
	}
	home := sideList[2].([]interface{})[1].(float64)
	if home == 0 {
		home = 1
	}
	margin := util.TruncateFloat(1/(1/away+1/draw+1/home)*100-100, 3)
	if margin > 0 {
		h.log.Infow("offers_event", "s", sportCode, "eventId", eventId, "a", away, "h", home, "d", draw, "m", margin)
	}
}
