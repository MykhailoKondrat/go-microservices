package handlers

import (
	"github.com/MykhailoKondrat/go-microservices/products/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: productsResponse

// GetProducts returns list of products
func (p *Products) GetProducts(w http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
