package main

import (
	"net/http"
	"strconv"

	"math/rand"
)

type NumRandomHandler struct{}

func NewNumRandomHAndler(router *http.ServeMux) {

	handler := &NumRandomHandler{}
	router.HandleFunc("/numRandom", handler.getRandom())
}

func (n *NumRandomHandler) getRandom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(rand.Intn(6) + 1)))
	}
}
