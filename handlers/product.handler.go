package handlers

import (
	"log"
	"microservice/data"
	"net/http"
	"strconv"
	"strings"
)

type Products struct {
	log             *log.Logger
	productsMethods *data.ProductsData
	productModel    *data.Product
}

func CreateInstance(l *log.Logger, products *data.ProductsData, product *data.Product) *Products {
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
	err := p.productModel.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unabled to save this product", http.StatusInternalServerError)
	}

	p.productsMethods.AddProducts(p.productModel)

	rw.WriteHeader(http.StatusAccepted)
}

func (p *Products) UpdateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	err := p.productModel.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unabled to save this product", http.StatusInternalServerError)
	}
	p.productsMethods.UpdateProducts(id, p.productModel)

	rw.WriteHeader(http.StatusAccepted)
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		p.GetProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.AddProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		params := strings.Split(r.URL.Path, "/")

		id, _ := strconv.Atoi(params[2])
		p.UpdateProducts(id, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
	return
}
