package main

import (
	"context"
	"log"
	"microservice/data"
	"microservice/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	products := handlers.CreateInstance(l, &data.ProductsData{}, data.Product{})
	sm := mux.NewRouter()

	// sm.Handle("/products", products).Methods
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", products.GetProducts)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", products.UpdateProducts)
	putRouter.Use(products.MiddlewareValidateProduct)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", products.AddProduct)
	postRouter.Use(products.MiddlewareValidateProduct)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
