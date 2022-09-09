package handlers

import (
	"github.com/MykhailoKondrat/go-microservices/products/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns nothing
// responses:
//  201: noContent

// DeleteProduct deletes a product from the database
func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}
	p.l.Println("handle delete product")
	err = data.DeleteProduct(id)
	if err == data.ErrorProductNotFound {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
}
