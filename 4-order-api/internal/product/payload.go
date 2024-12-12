package product

import "github.com/lib/pq"

type ProductCreateRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images"`
}
