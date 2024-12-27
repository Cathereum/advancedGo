package req

import (
	"encoding/json"
	"net/http"
)

func Decode(req *http.Request, payload any) error {

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return err
	}

	return nil
}
