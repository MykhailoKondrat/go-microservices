// Package handlers Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	protos "github.com/MykhailoKondrat/go-microservices/currency/protos"
	"github.com/MykhailoKondrat/go-microservices/products/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type GenericError struct {
	Message string `json:"message"`
}

// A list of prouducts
// swagger:response productsResponse
type productsResponse struct {
	// Respone message
	// in: body
	Body []data.Product
}

// swagger:response noContent
type productNoContent struct {
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete form DB
	// in: path
	// required: true
	ID int `json:"id"`
}

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, c protos.CurrencyClient) *Products {
	return &Products{l: l, cc: c}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle post product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
	p.l.Printf("Prod:%#v", prod)

}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
		}

		// validate the product
		err := prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(w, fmt.Sprintf("Error validating product:%s", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
