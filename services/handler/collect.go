package handler

import (
	"github.com/aibotsoft/micro/util"
	"github.com/aibotsoft/trade/pkg/store"
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
var BetTypeMap = map[string]int64{
	"wdw":          1,
	"dc":           2,
	"ah":           3,
	"ahou":         4,
	"tahou,h":      5,
	"tahou,a":      6,
	"score,both":   7,
	"oe":           8,
	"clean,h":      9,
	"clean,a":      10,
	"win_to_nil,h": 11,
	"win_to_nil,a": 12,
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
		//h.log.Infow("away_blank", "eventList", eventList)
		return
	} else if strings.Index(strings.ToLower(away), "test") != -1 {
		h.log.Infow("away_test", "eventList", eventList)
		return
	}

	leagueId := int64(eventList["competition_id"].(float64))

	//if ok {
	//	h.log.Infow("",
	//		"s", sportCode,
	//		"eventId", eventId,
	//		"home", home,
	//		"away", away,
	//		//"competitionId", leagueId,
	//		//"leagueName", leagueName,
	//		//"country", country,
	//		//"starts", starts,
	//		//"eventList", eventList,
	//		"ir_status", irStatus,
	//	)
	//}
	country := eventList["country"].(string)
	starts, ok := eventList["start_ts"].(string)
	if !ok {
		//h.log.Info("starts_not_ok: ", eventList)
		return
	}
	h.store.SaveSport(sportId, sportName)
	h.store.SaveTeam(homeId, home)
	h.store.SaveTeam(awayId, away)
	h.store.SaveLeague(leagueId, leagueName, country, sportId)
	h.store.SaveEvent(eventId, homeId, awayId, leagueId, starts)
	h.store.SaveEventPeriod(eventId, sportCode, true)

	irStatus, ok := eventList["ir_status"].(map[string]interface{})
	if sportName == "Football" && len(irStatus) > 0 {
		eventPeriodId, err := h.store.GetEventPeriodId(eventId, sportCode)
		if err != nil {
			h.log.Error(err)
			return
		}
		var sf = store.ScoreFootball{EventPeriodId: eventPeriodId}
		rc, ok := irStatus["rc"].([]interface{})
		if ok {
			sf.RedHome = util.PtrFloat64(rc[0].(float64))
			sf.RedAway = util.PtrFloat64(rc[1].(float64))
		}

		score, ok := irStatus["score"].([]interface{})
		if ok {
			sf.ScoreHome = util.PtrFloat64(score[0].(float64))
			sf.ScoreAway = util.PtrFloat64(score[1].(float64))
		}

		timeList, ok := irStatus["time"].([]interface{})
		if ok {
			sf.PeriodCode = timeList[0].(string)
			sf.PeriodMin = util.PtrFloat64(timeList[1].(float64))
		}
		h.log.Infow("", "id", eventId, "s", sportCode, "status", irStatus, "rc", rc, "sf", sf)
		h.store.SaveScoreFootball(sf)
	}

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
		return
	}
	var trios []store.Trio
	var duos []store.Duo
	for key, value := range priceList {
		valueList := value.([]interface{})
		switch key {
		case "wdw":
			trio := h.trio(eventPeriodId, BetTypeMap["wdw"], valueList)
			trios = append(trios, trio)
		case "dc":
			trio := h.trio(eventPeriodId, BetTypeMap["dc"], valueList)
			trios = append(trios, trio)
		case "ah":
			d := h.duo(eventPeriodId, BetTypeMap["ah"], valueList)
			duos = append(duos, d...)
		case "ahou":
			d := h.duo(eventPeriodId, BetTypeMap["ahou"], valueList)
			duos = append(duos, d...)
		case "tahou,h":
			d := h.duo(eventPeriodId, BetTypeMap["tahou,h"], valueList)
			duos = append(duos, d...)
		case "tahou,a":
			d := h.duo(eventPeriodId, BetTypeMap["tahou,a"], valueList)
			duos = append(duos, d...)
		case "score,both":
			d := h.duo(eventPeriodId, BetTypeMap["score,both"], valueList)
			duos = append(duos, d...)
		case "oe":
			d := h.duo(eventPeriodId, BetTypeMap["oe"], valueList)
			duos = append(duos, d...)
		case "clean,h":
			d := h.duo(eventPeriodId, BetTypeMap["clean,h"], valueList)
			duos = append(duos, d...)
		case "clean,a":
			d := h.duo(eventPeriodId, BetTypeMap["clean,a"], valueList)
			duos = append(duos, d...)
		case "win_to_nil,h":
			d := h.duo(eventPeriodId, BetTypeMap["win_to_nil,h"], valueList)
			duos = append(duos, d...)
		case "win_to_nil,a":
			d := h.duo(eventPeriodId, BetTypeMap["win_to_nil,a"], valueList)
			duos = append(duos, d...)
		case "cs":
		default:
			//h.log.Infow("", "key: ", key, "valueList", valueList)
		}
	}
	//h.log.Info(len(trios), " ", len(duos))
	for i := range duos {
		//h.log.Infow("", "", duos)
		h.store.SaveDuo(duos[i])
	}
	for i := range trios {
		//h.log.Infow("", "", trios)
		h.store.SaveTrio(trios[i])
	}
}
func (h *Handler) duo(eventPeriodId int64, betTypeId int64, valueList []interface{}) (duos []store.Duo) {
	d := store.Duo{EventPeriodId: eventPeriodId, BetTypeId: betTypeId}
	for i := range valueList {
		priceList := valueList[i].([]interface{})
		switch v := priceList[0].(type) {
		case float64:
			d.Code = int64(v)
		case nil:
			d.Code = 0
		}
		sideList, ok := priceList[1].([]interface{})
		if ok {
			d.IsActive = true
			for i := range sideList {
				switch sideList[i].([]interface{})[0].(string) {
				case "a":
					d.APrice = sideList[i].([]interface{})[1].(float64)
				case "h":
					d.BPrice = sideList[i].([]interface{})[1].(float64)
				case "over":
					d.APrice = sideList[i].([]interface{})[1].(float64)
				case "under":
					d.BPrice = sideList[i].([]interface{})[1].(float64)
				case "no":
					d.APrice = sideList[i].([]interface{})[1].(float64)
				case "yes":
					d.BPrice = sideList[i].([]interface{})[1].(float64)
				case "even":
					d.APrice = sideList[i].([]interface{})[1].(float64)
				case "odd":
					d.BPrice = sideList[i].([]interface{})[1].(float64)
				default:
					h.log.Info(sideList[i].([]interface{})[0].(string))
				}
			}
		}
		duos = append(duos, d)
	}
	return
}

func (h *Handler) trio(eventPeriodId int64, betTypeId int64, valueList []interface{}) (t store.Trio) {
	t.EventPeriodId = eventPeriodId
	t.BetTypeId = betTypeId
	priceList := valueList[0].([]interface{})
	sideList, ok := priceList[1].([]interface{})
	if ok {
		t.IsActive = true
		for i := range sideList {
			switch sideList[i].([]interface{})[0].(string) {
			case "h":
				t.APrice = sideList[i].([]interface{})[1].(float64)
			case "d":
				t.BPrice = sideList[i].([]interface{})[1].(float64)
			case "a":
				t.CPrice = sideList[i].([]interface{})[1].(float64)
			case "h,d":
				t.APrice = sideList[i].([]interface{})[1].(float64)
			case "h,a":
				t.BPrice = sideList[i].([]interface{})[1].(float64)
			case "a,d":
				t.CPrice = sideList[i].([]interface{})[1].(float64)
			}
		}
	}
	return
}

//func (h *Handler) dc(eventPeriodId int64, valueList []interface{}) (t store.Trio) {
//	t.EventPeriodId = eventPeriodId
//	t.BetTypeId = BetTypeMap["dc"]
//	priceList := valueList[0].([]interface{})
//	sideList, ok := priceList[1].([]interface{})
//	if ok {
//		t.IsActive = true
//		for i := range sideList {
//			//h.log.Infow("",  "v", sideList[i])
//			switch sideList[i].([]interface{})[0].(string) {
//			case "h,d":
//				t.APrice = sideList[i].([]interface{})[1].(float64)
//			case "h,a":
//				t.BPrice = sideList[i].([]interface{})[1].(float64)
//			case "a,d":
//				t.CPrice = sideList[i].([]interface{})[1].(float64)
//			}
//		}
//	}
//	return
//}

//func (h *Handler) dc(eventPeriodId int64, value interface{}) {
//	valueList := value.([]interface{})
//	priceList := valueList[0].([]interface{})
//	sideList, ok := priceList[1].([]interface{})
//	if !ok {
//		//h.log.Infow("dc_not_ok",  "eventPeriodId", eventPeriodId,"v", value)
//		h.store.DeactivateDoubleChance(eventPeriodId)
//		return
//	}
//	var homeAway float64 = 1
//	var homeDraw float64 = 1
//
//	awayDraw := sideList[0].([]interface{})[1].(float64)
//	if awayDraw == 0 {
//		awayDraw = 1
//	}
//	if len(sideList) > 1 {
//		homeAway = sideList[1].([]interface{})[1].(float64)
//	}
//	if len(sideList) > 2 {
//		homeDraw = sideList[2].([]interface{})[1].(float64)
//
//	}
//	margin := util.TruncateFloat(1/(1/awayDraw+1/homeAway+1/homeDraw)*100-100, 3)
//	h.store.SaveDoubleChance(eventPeriodId, awayDraw, homeAway, homeDraw, margin, true)
//	if margin > minPercent {
//		h.log.Infow("dc",  "ad", awayDraw, "hd", homeDraw, "ha", homeAway, "m", margin)
//	}
//}
//var draw float64 = 1
//var home float64 = 1
//away := sideList[0].([]interface{})[1].(float64)
//if away == 0 {
//	away = 1
//}
//if len(sideList) > 1 {
//	draw = sideList[1].([]interface{})[1].(float64)
//}
//if len(sideList) > 2 {
//	home = sideList[2].([]interface{})[1].(float64)
//}
//margin := util.TruncateFloat(1/(1/away+1/draw+1/home)*100-100, 3)
//h.store.SaveWinDrawWin(eventPeriodId, away, home, draw, margin, true)
//if margin > minPercent {
//	h.log.Infow("wdw", "id", eventPeriodId, "a", away, "h", home, "d", draw, "m", margin)
//}
//func (h *Handler) ah(eventPeriodId int64, valueList []interface{}) (duos []store.Duo) {
//	d := store.Duo{EventPeriodId: eventPeriodId, BetTypeId: BetTypeMap["ah"]}
//	for i := range valueList {
//		priceList := valueList[i].([]interface{})
//		d.Code = int64(priceList[0].(float64))
//		sideList, ok := priceList[1].([]interface{})
//		if ok {
//			d.IsActive = true
//			for i := range sideList {
//				switch sideList[i].([]interface{})[0].(string) {
//				case "h":
//					d.APrice = sideList[i].([]interface{})[1].(float64)
//				case "a":
//					d.BPrice = sideList[i].([]interface{})[1].(float64)
//				}
//			}
//		}
//		duos = append(duos, d)
//	}
//	return
//}
