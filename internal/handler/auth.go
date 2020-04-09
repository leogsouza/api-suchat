package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gookit/validate"
)

type loginInput struct {
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required|minLen:8"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var in loginInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := validate.Struct(in)

	if !v.Validate() {
		respond(w, v.Errors, http.StatusUnprocessableEntity)
		return
	}

	out, err := h.Login(in.Email, in.Password)

	if err != nil {
		errMessage := ErrorResponse{http.StatusBadRequest, err.Error()}
		respond(w, errMessage, http.StatusBadRequest)
		return
	}

	respond(w, out, http.StatusOK)
}
