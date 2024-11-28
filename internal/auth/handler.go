package auth

import (
	"advancedGo/configs"
	res "advancedGo/pkg"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type authHandler struct {
	*configs.Config
}

func NewHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &authHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())

}

func (a *authHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := LoginResponse{
			Token: a.Auth.Token,
		}

		res.Json(w, data, 200)

	}
}
func (a *authHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Register successful"))
	}
}
