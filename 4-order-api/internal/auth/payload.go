package auth

type RegisterResponse struct {
	Phone   string `json:"phone,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}

type VerifyRequest struct {
	Code string `json:"code" validate:"required"`
}

type RegisterRequest struct {
	Phone string `json:"phone" validate:"required"`
}
