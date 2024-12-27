package req

import (
	"github.com/go-playground/validator/v10"
)

func Validate(payload any) error {

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return err
	}
	return nil
}
