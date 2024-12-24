package order

import (
	"advancedGo/configs"
	"advancedGo/internal/product"
	"advancedGo/pkg/di"
	"advancedGo/pkg/middleware"
	"advancedGo/pkg/req"
	"advancedGo/pkg/res"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type OrderHandlerDeps struct {
	*configs.Config
	OrderRepository    *OrderRepository
	IProductRepository di.IProductRepository
	IUserRepository    di.IUserRepository
}

type orderHandler struct {
	*configs.Config
	OrderRepository    *OrderRepository
	IProductRepository di.IProductRepository
	IUserRepository    di.IUserRepository
}

func NewHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &orderHandler{
		Config:             deps.Config,
		OrderRepository:    deps.OrderRepository,
		IProductRepository: deps.IProductRepository,
		IUserRepository:    deps.IUserRepository,
	}

	router.Handle("GET /order/{id}", middleware.IsAuth(handler.GetByID(), deps.Config))
	router.Handle("GET /my-orders", middleware.IsAuth(handler.GetAll(), deps.Config))
	router.Handle("POST /order", middleware.IsAuth(handler.Create(), deps.Config))

}

func (handler *orderHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := handler.IUserRepository.GetUserId(phone)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		orders, err := handler.OrderRepository.GetUserOrders(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, orders, 200)
	}
}

func (handler *orderHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := handler.IUserRepository.GetUserId(phone)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Println("Error to parse idString")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := handler.OrderRepository.GetOrderById(userId, uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, order, 200)
	}
}

func (handler *orderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		phone, ok := r.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := handler.IUserRepository.GetUserId(phone)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var payload OrderCreateRequest

		err = req.HandleBody(r, &payload)
		if err != nil {
			log.Println("Handle body error:", err)
			return
		}

		var products []product.Product
		for _, productId := range payload.ProductIds {
			product, err := handler.IProductRepository.GetProductById(productId)
			if err != nil {
				log.Printf("Failed to fetch product %d: %v", productId, err)
				return
			}
			products = append(products, *product)
		}

		order := &Order{
			UserId:   userId,
			Products: products,
		}

		createdOrder, err := handler.OrderRepository.Create(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, createdOrder, 201)
	}
}
