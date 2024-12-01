package verify

type LoginRequest struct {
	Email string `json:"email"`
}

type Verification struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}
