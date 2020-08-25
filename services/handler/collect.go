package handler

func (h *Handler) processEvent(m []interface{}) {
	sportList := m[1].([]interface{})
	sport := sportList[0].(string)
	eventId := sportList[1].(string)

	eventList, ok := m[2].(map[string]interface{})
	if !ok {
		h.log.Infow("eventList_not_ok", "m[2]", m[2])
		return
	}
	home := eventList["home"]
	away := eventList["away"]
	competitionId := eventList["competition_id"]
	competitionName := eventList["competition_name"]
	country := eventList["country"]
	startTs := eventList["start_ts"]

	h.log.Infow("", "sport", sport,
		"eventId", eventId,
		"home", home,
		"away", away,
		"competitionId", competitionId,
		"competitionName", competitionName,
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
