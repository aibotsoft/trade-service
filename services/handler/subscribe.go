package handler

func (h *Handler) SubscribeLoop() {
	h.store.GetLiveEvent()

}
