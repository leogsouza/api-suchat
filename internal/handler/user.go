package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gookit/validate"
)

type registerUserInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required|email"`
	Lastname string `json:"lastname" validate:""`
	Password string `json:"password" validate:"required|minLen:8"`
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var in registerUserInput

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

	err := h.Register(in.Name, in.Email, in.Lastname, in.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
