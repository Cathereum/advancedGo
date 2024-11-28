package verify

import (
	"advancedGo/configs"

	"net/http"
)

type VerifyHandlerDeps struct {
	*configs.Config
}

type verifyHandler struct {
	*configs.Config
}

func NewHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &verifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
	router.HandleFunc("POST /send", handler.Send())

}

func (a *verifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mail sended"))

	}
}
func (a *verifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Verify successful"))
	}
}
