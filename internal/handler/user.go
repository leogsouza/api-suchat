package handler

import (
	"encoding/json"
	"net/http"

	"github.com/leogsouza/api-suchat/internal/model"
	"gopkg.in/dealancer/validate.v2"
)

type registerUserInput struct {
	Name     string `validate:"empty=false & format=alnum_unicode"`
	Email    string `validate:"empty=false & format=email"`
	Lastname string `validate:"empty=false & format=alnum_unicode"`
	Password string `validate:"empty=false & gte=8"`
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var in registerUserInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Validate(&in); err != nil {
		switch err.(type) {
		case validate.ErrorSyntax:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case validate.ErrorValidation:
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user := &model.User{
		Name:     in.Name,
		Email:    in.Email,
		Lastname: in.Lastname,
		Password: in.Password,
	}

	err := h.Register(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
