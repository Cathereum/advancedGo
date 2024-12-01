package verify

import (
	"advancedGo/configs"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"net/http"
	"net/smtp"

	"github.com/jordan-wright/email"
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
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())

}

func (a *verifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("Request error %v", err)
			return
		}

		hashBytes := make([]byte, 16)
		if _, err := rand.Read(hashBytes); err != nil {
			log.Printf("Error generating hash: %v", err)
			return
		}

		hash := hex.EncodeToString(hashBytes)

		e := email.NewEmail()
		e.From = "yourMail@gmail.com"
		e.To = []string{payload.Email}
		e.Subject = "Email Verification"
		verifyLink := fmt.Sprintf("http://localhost:8081/verify/%s", hash)
		e.HTML = []byte(fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Please click the link below to verify your email:</p>
		<a href="%s">%s</a>
	`, verifyLink, verifyLink))

		err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("",
			"yourMail@gmail.com",
			"yourMail app code",
			"smtp.gmail.com",
		))
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return
		}

		data := Verification{
			Email: payload.Email,
			Hash:  hash,
		}

		file, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Printf("Error creating JSON: %v", err)
			return
		}

		if err := os.WriteFile("verify.json", file, 0644); err != nil {
			log.Printf("Error saving file: %v", err)
			return
		}

		log.Printf("Data saved to verify.json:")
		log.Printf("Email: %v", payload.Email)
		log.Printf("Generated hash: %v", hash)
	}
}
func (a *verifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		data, err := os.ReadFile("verify.json")
		if err != nil {
			log.Printf("Error reading file: %v", err)
			return
		}

		var verification Verification
		if err := json.Unmarshal(data, &verification); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return
		}

		if verification.Hash == hash {
			if err := os.Remove("verify.json"); err != nil {
				log.Printf("Error removing file: %v", err)
			}
			w.Write([]byte("Email successfully verified!"))
			return
		}

		w.Write([]byte("Verification failed: invalid verification code"))
	}
}
