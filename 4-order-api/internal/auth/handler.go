package auth

import (
	"advancedGo/configs"
	"advancedGo/pkg/jwt"
	"advancedGo/pkg/middleware"
	"advancedGo/pkg/req"
	"advancedGo/pkg/res"
	"fmt"

	"log"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type authHandler struct {
	*configs.Config
	*AuthService
}

func NewHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &authHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/register", handler.Register())
	router.Handle("POST /auth/verify", middleware.IsAuth(handler.VerifyCode(), deps.Config))

}

func (a *authHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload RegisterRequest

		if err := req.HandleBody(r, &payload); err != nil {
			log.Printf("Handle body error: %v", err)
			if err := res.Json(w, err.Error(), 401); err != nil {
				log.Printf("Failed to send error response: %v", err)
			}
			return
		}

		response, err := a.AuthService.Register(payload.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.NewJWT(a.Config.Auth.Token).Create(jwt.JWTData{
			Phone: payload.Phone,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		data := RegisterResponse{
			Code:    response.Code,
			Message: response.Message,
			Token:   token,
		}

		if err := res.Json(w, data, 200); err != nil {
			log.Printf("Failed to send error response: %v", err)
			return
		}
	}
}

func (a *authHandler) VerifyCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload VerifyRequest
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if ok {
			fmt.Println(phone)
		}

		if err := req.HandleBody(r, &payload); err != nil {
			log.Printf("Handle body error: %v", err)
			if err := res.Json(w, err.Error(), 401); err != nil {
				log.Printf("Failed to send error response: %v", err)
			}
			return
		}

		err := a.AuthService.VerifyCode(phone, payload.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := res.Json(w, "authorization successed", 200); err != nil {
			log.Printf("Failed to send response: %v", err)
			return
		}
	}
}
