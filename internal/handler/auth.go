package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gookit/validate"
	"github.com/leogsouza/api-suchat/internal/service"
)

type loginInput struct {
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required|minLen:8"`
}

type loginErrorResponse struct {
	LoginSuccess bool   `json:"loginSuccess"`
	Message      string `json:"message"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var in loginInput

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respond(w, loginErrorResponse{false, "Wrong JSON"}, http.StatusOK)
		return
	}

	v := validate.Struct(in)

	if !v.Validate() {
		respond(w, loginErrorResponse{false, "Unprocessable entity"}, http.StatusOK)
		return
	}

	out, err := h.Login(in.Email, in.Password)

	if err != nil {

		respond(w, loginErrorResponse{false, "BadRequest"}, http.StatusOK)
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

type authResponse struct {
	IsAuth bool `json:"isAuth"`
	Error  bool `json:"error"`
}

func (h *handler) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		log.Printf("authorization %v", a)
		if !strings.HasPrefix(a, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		token := a[7:]
		log.Printf("token", token)
		uid, err := h.Service.AuthUserEmailID(token)

		log.Printf("uid", uid)
		if err != nil {
			respond(w, authResponse{}, http.StatusOK)
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
		respond(w, authResponse{}, http.StatusOK)
		return
	}

	if err == service.ErrUserNotFound {
		respond(w, authResponse{}, http.StatusOK)
		return
	}

	if err != nil {
		respond(w, authResponse{}, http.StatusOK)
		return
	}

	respond(w, u, http.StatusOK)

}
