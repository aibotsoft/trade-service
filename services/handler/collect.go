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

const minPercent = 1.2

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
	//leagueId := int64(sportList[0].(float64))
	sportCode := sportList[1].(string)
	eventId := sportList[2].(string)
	priceList := m[2].(map[string]interface{})
	if len(priceList) == 0 {
		h.store.SaveEventPeriod(eventId, sportCode, false)
		return
	}
	eventPeriodId, err := h.store.GetEventPeriodId(eventId, sportCode)
	if err != nil {
		h.log.Error(err)
	}

	for key, value := range priceList {
		switch key {
		case "wdw":
			h.wdw(eventPeriodId, value)
		case "ah":
			h.ah(eventPeriodId, value)
		case "dc":
			h.dc(eventPeriodId, value)
		case "ahou":
			h.ahou(eventPeriodId, value)
		default:
			//h.log.Info("key: ", key)
		}
	}
}
func (h *Handler) ahou(eventPeriodId int64, value interface{}) {
	//h.log.Infow("ahou", "id", eventPeriodId, "value", value)
	valueList := value.([]interface{})
	for i := range valueList {
		priceList := valueList[i].([]interface{})
		handicap := int64(priceList[0].(float64))
		sideList, ok := priceList[1].([]interface{})
		if !ok {
			//h.log.Infow("ahou_not_ok", "id", eventPeriodId, "handicap", handicap)
			h.store.DeactivateTotal(eventPeriodId, handicap)
			continue
		}
		var under float64 = 1
		over := sideList[0].([]interface{})[1].(float64)
		if over == 0 {
			over = 1
		}
		if len(sideList) > 1 {
			under = sideList[1].([]interface{})[1].(float64)
		}
		margin := util.TruncateFloat(1/(1/over+1/under)*100-100, 3)
		h.store.SaveTotal(eventPeriodId, handicap, over, under, margin, true)
		if margin > minPercent {
			h.log.Infow("ahou", "id", eventPeriodId, "handicap", handicap, "over", over, "under", under, "margin", margin)
		}
	}
}
func (h *Handler) ah(eventPeriodId int64, value interface{}) {
	valueList := value.([]interface{})
	for i := range valueList {
		priceList := valueList[i].([]interface{})
		handicap := int64(priceList[0].(float64))
		sideList, ok := priceList[1].([]interface{})
		if !ok {
			//h.log.Infow("ah_not_ok", "id", eventPeriodId, "handicap", handicap)
			h.store.DeactivateHandicap(eventPeriodId, handicap)
			continue
		}
		var home float64 = 1
		away := sideList[0].([]interface{})[1].(float64)
		if away == 0 {
			away = 1
		}
		if len(sideList) > 1 {
			home = sideList[1].([]interface{})[1].(float64)

		}

		margin := util.TruncateFloat(1/(1/away+1/home)*100-100, 3)
		h.store.SaveHandicap(eventPeriodId, handicap, away, home, margin, true)
		if margin > minPercent {
			h.log.Infow("ah", "id", eventPeriodId, "handicap", handicap, "away", away, "home", home, "margin", margin)
		}
	}
}
func (h *Handler) dc(eventPeriodId int64, value interface{}) {
	valueList := value.([]interface{})
	priceList := valueList[0].([]interface{})
	sideList, ok := priceList[1].([]interface{})
	if !ok {
		//h.log.Infow("dc_not_ok",  "eventPeriodId", eventPeriodId,"v", value)
		h.store.DeactivateDoubleChance(eventPeriodId)
		return
	}
	var homeAway float64 = 1
	var homeDraw float64 = 1

	awayDraw := sideList[0].([]interface{})[1].(float64)
	if awayDraw == 0 {
		awayDraw = 1
	}
	if len(sideList) > 1 {
		homeAway = sideList[1].([]interface{})[1].(float64)
	}
	if len(sideList) > 2 {
		homeDraw = sideList[2].([]interface{})[1].(float64)

	}
	margin := util.TruncateFloat(1/(1/awayDraw+1/homeAway+1/homeDraw)*100-100, 3)
	h.store.SaveDoubleChance(eventPeriodId, awayDraw, homeAway, homeDraw, margin, true)
	if margin > minPercent {
		h.log.Infow("dc",  "ad", awayDraw, "hd", homeDraw, "ha", homeAway, "m", margin)
	}
}

func (h *Handler) wdw(eventPeriodId int64, value interface{}) {
	valueList := value.([]interface{})
	priceList := valueList[0].([]interface{})
	sideList, ok := priceList[1].([]interface{})
	if !ok {
		//h.log.Infow("wdw_not_ok",  "eventPeriodId", eventPeriodId,"v", value)
		h.store.DeactivateWinDrawWin(eventPeriodId)
		return
	}
	var draw float64 = 1
	var home float64 = 1
	away := sideList[0].([]interface{})[1].(float64)
	if away == 0 {
		away = 1
	}
	if len(sideList) > 1 {
		draw = sideList[1].([]interface{})[1].(float64)
	}
	if len(sideList) > 2 {
		home = sideList[2].([]interface{})[1].(float64)
	}
	margin := util.TruncateFloat(1/(1/away+1/draw+1/home)*100-100, 3)
	h.store.SaveWinDrawWin(eventPeriodId, away, home, draw, margin, true)
	if margin > minPercent {
		h.log.Infow("wdw", "id", eventPeriodId, "a", away, "h", home, "d", draw, "m", margin)
	}
}
