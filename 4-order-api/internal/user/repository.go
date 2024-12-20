package user

import "advancedGo/pkg/db"

type UserRepository struct {
	Database *db.DB
}

func NewUserRepository(database *db.DB) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) FindByPhone(phone string) (*User, error) {

	var user User
	result := repo.Database.DB.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) Update(user *User) error {
	result := repo.Database.DB.Model(user).Updates(map[string]interface{}{
		"session_id":        user.SessionId,
		"verification_code": user.VerificationCode,
		"is_verified":       user.IsVerified,
	})
	return result.Error
}
