package handlers

import (
	"github.com/MykhailoKondrat/go-microservices/products/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}
	p.l.Println("handle PUT product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	//prod := &data.Product{}
	//if err := prod.FromJSON(r.Body); err != nil {
	//	http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
	//}
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorProductNotFound {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
}
