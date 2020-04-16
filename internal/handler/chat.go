package handler

import "net/http"

func (h *handler) getChats(w http.ResponseWriter, r *http.Request) {
	chats, err := h.GetChats()
	if err != nil {

		respondHTTPError(w, err, http.StatusBadRequest)
	}

	respond(w, chats, http.StatusOK)
}
