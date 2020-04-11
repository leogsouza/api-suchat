package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gookit/validate"
	"github.com/leogsouza/api-suchat/internal/service"
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

type logoutOutput struct {
	Success bool
}

func (h handler) logout(w http.ResponseWriter, r *http.Request) {
	err := h.Logout(r.Context())
	if err != nil {
		errMessage := ErrorResponse{http.StatusBadRequest, err.Error()}
		respond(w, errMessage, http.StatusBadRequest)
		return
	}

	respond(w, logoutOutput{true}, http.StatusOK)
}

func (h *handler) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		if !strings.HasPrefix(a, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		token := a[7:]
		uid, err := h.Service.AuthUserEmailID(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, service.KeyAuthUserID, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) authUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.AuthUser(r.Context())

	if err == service.ErrUnauthenticated {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err == service.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		respondError(w, err)
		return
	}

	respond(w, u, http.StatusOK)

}
