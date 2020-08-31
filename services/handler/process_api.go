package handler

func (h *Handler) processApi(m []interface{}) {
	//h.log.Infow("api", "msg", m)
	apiMap, ok := m[1].(map[string]interface{})
	if !ok {
		h.log.Infow("api_not_ok", "m", m)
		return
	}
	data, ok := apiMap["data"].([]interface{})
	if !ok {
		h.log.Infow("data_not_ok", "m", m)
		return
	}
	for i := range data {
		d := data[i].([]interface{})
		if d[0] == "pmm" {
			h.processPmm(d)
		}
	}
}
func (h *Handler) processPmm(d []interface{}) {
	pmm, ok := d[1].(map[string]interface{})
	if !ok {
		h.log.Infow("pmm_not_ok", "msg", d)
		return
	}
	priceList := pmm["price_list"].([]interface{})
	if len(priceList) == 0 {
		return
	}
	h.log.Infow("pmm", "", pmm)


}