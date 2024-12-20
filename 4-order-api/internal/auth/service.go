package auth

import (
	"advancedGo/internal/user"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Register(phone string) (*RegisterResponse, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)

	// Генерируем SessionId
	newSessionID, err := rand.Int(rand.Reader, big.NewInt(100_000_000))
	if err != nil {
		return nil, fmt.Errorf("error generating sessionID: %w", err)
	}

	// Генерируем код подтверждения (4 цифры)
	verificationCode, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return nil, fmt.Errorf("error generating verification code: %w", err)
	}
	code := fmt.Sprintf("%04d", verificationCode)

	if existedUser != nil {
		existedUser.SessionId = newSessionID.String()
		existedUser.VerificationCode = code
		existedUser.IsVerified = false
		err = service.UserRepository.Update(existedUser)
		if err != nil {
			return nil, fmt.Errorf("error updating session")
		}

		return &RegisterResponse{
			Phone:   existedUser.Phone,
			Code:    code,
			Message: "User already exists",
		}, nil
	}

	user := &user.User{
		Phone:            phone,
		SessionId:        newSessionID.String(),
		VerificationCode: code,
		IsVerified:       false,
	}

	_, err = service.UserRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		Phone: user.Phone,
		Code:  code,
	}, nil
}

func (service *AuthService) VerifyCode(phone, code string) error {
	user, err := service.UserRepository.FindByPhone(phone)
	if err != nil {
		return errors.New("user not found")
	}

	if user.VerificationCode != code {
		return errors.New("invalid code")
	}

	user.IsVerified = true
	return service.UserRepository.Update(user)
}
