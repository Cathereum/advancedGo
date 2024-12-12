package req

import (
	"fmt"
	"net/http"
)

func HandleBody(req *http.Request, payload any) error {
	if err := Decode(req, payload); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	if err := Validate(payload); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println(payload)

	return nil
}
