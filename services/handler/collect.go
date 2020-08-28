package handler

import (
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
	sport := sportList[0].(string)
	sportName := sportPeriodMap[sport]
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
		h.log.Infow("eventList_not_ok", "m[2]", m[2])
		return
	}
	leagueName := eventList["competition_name"].(string)
	if leagueName == "" {
		h.log.Infow("competition_name_blank", "eventList", eventList)
		return
	} else if strings.Index(strings.ToLower(leagueName), "test") != -1 {
		h.log.Infow("competition_name_test", "eventList", eventList)
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

	country := eventList["country"].(string)
	startTs := eventList["start_ts"]

	h.store.SaveSport(sportId, sportName)
	h.store.SaveTeam(homeId, home)
	h.store.SaveTeam(awayId, away)
	h.store.SaveLeague(leagueId, leagueName, country, sportId)
	//h.store.SaveEvent(leagueId, leagueName, country, sportId)

	h.log.Infow("", "sport", sport,
		"eventId", eventId,
		"home", home,
		"away", away,
		"competitionId", leagueId,
		"leagueName", leagueName,
		"country", country,
		"startTs", startTs,
		"eventList", eventList)
	//sport := ss[0].(string)
	//if !util.StringInList(sport, sportList) {
	//	//h.log.Info("fuck_sport: ", sport)
	//	return
	//}
	//var l store.League
	//var home, away store.Team
	//eventIdStr := ss[1].(string)
	//eventSplit := strings.Split(eventIdStr, ",")
	//if len(eventSplit) != 3 {
	//	h.log.Info("split_event_error: ", eventIdStr)
	//	return
	//}
	//home.Id, _ = strconv.ParseInt(eventSplit[1], 10, 64)
	//away.Id, _ = strconv.ParseInt(eventSplit[2], 10, 64)
	//
	//em, ok := d[2].(map[string]interface{})
	//if !ok {
	//	return
	//}
	//starts, ok := em["start_ts"].(string)
	//if !ok {
	//	//h.log.Infow("", "start_ts_not_ok", em)
	//	return
	//}
	//st, err := time.Parse(time.RFC3339, starts)
	//if err != nil {
	//	h.log.Info("time_error", st)
	//	return
	//}
	//home.Name = em["home"].(string)
	//if home.Name == "" {
	//	return
	//}
	//away.Name = em["away"].(string)
	//if away.Name == "" {
	//	return
	//}
	//l.Name = em["competition_name"].(string)
	//if l.Name == "" {
	//	return
	//}
	//l.Id = int64(em["competitionId"].(float64))
	//if l.Id == 0 {
	//	return
	//}
	//l.Country = em["country"].(string)
	//l.Sport = sport
	//
	//var e = store.Event{
	//	Id:       eventIdStr,
	//	LeagueId: l.Id,
	//	HomeId:   home.Id,
	//	AwayId:   away.Id,
	//	Starts:   st,
	//}
	////h.log.Infow("", "asdf", sport, "", eventIdStr, "", em, "home", home, "away", away, "starts", starts, "e", e, "l", l)
	////h.log.Infow("", "e", e)
	//h.store.SaveLeague(l)
	//h.store.SaveTeam(home)
	//h.store.SaveTeam(away)
	//h.store.SaveEvent(e)
}
