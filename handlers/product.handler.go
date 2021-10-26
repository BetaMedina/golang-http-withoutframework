package handlers

import (
	"context"
	"log"
	"microservice/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	log             *log.Logger
	productsMethods *data.ProductsData
	productModel    data.Product
}

func CreateInstance(l *log.Logger, products *data.ProductsData, product data.Product) *Products {
	return &Products{l, products, product}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	lp := p.productsMethods.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unabled to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.productsMethods.AddProducts(prod)

	rw.WriteHeader(http.StatusAccepted)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.productsMethods.UpdateProducts(id, p.productModel)

	rw.WriteHeader(http.StatusAccepted)
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		err := p.productModel.FromJSON(r.Body)
		if err != nil {
			p.log.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, p.productModel)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
