package product

import (
	"advancedGo/pkg/req"
	"advancedGo/pkg/res"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
}

type productHandler struct {
	ProductRepository *ProductRepository
}

func NewHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &productHandler{
		ProductRepository: deps.ProductRepository,
	}

	router.HandleFunc("GET /product/{id}", handler.GetByID())
	router.HandleFunc("GET /product", handler.GetAll())
	router.HandleFunc("POST /product", handler.Create())
	router.HandleFunc("PATCH /product/{id}", handler.Update())
	router.HandleFunc("DELETE /product/{id}", handler.Delete())

}

func (handler *productHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Println("Error to parse idString")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product, err := handler.ProductRepository.GetProductById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, product, 200)
	}
}

func (handler *productHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		products, err := handler.ProductRepository.GetAllProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, products, 200)
	}
}

func (handler *productHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ProductCreateRequest

		err := req.HandleBody(r, &payload)
		if err != nil {
			log.Println("Handle body error:", err)
			return
		}

		product := &Product{
			Name:        payload.Name,
			Description: payload.Description,
		}

		createdProduct, err := handler.ProductRepository.Create(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdProduct, 201)
	}
}

func (handler *productHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ProductCreateRequest

		err := req.HandleBody(r, &payload)
		if err != nil {
			log.Println("Handle body error:", err)
			return
		}

		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Println("Error to parse idString")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product := &Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        payload.Name,
			Description: payload.Description,
		}
		updatedProduct, err := handler.ProductRepository.UpdateProductById(product)
		if err != nil {
			fmt.Println("Error to update product")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, updatedProduct, 201)
	}
}

func (handler *productHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Println("Error to parse idString")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.ProductRepository.DeleteProductById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, 200)
	}
}
