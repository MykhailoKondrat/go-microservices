package main

import (
	"context"
	protos "github.com/MykhailoKondrat/go-microservices/currency/protos"
	"github.com/MykhailoKondrat/go-microservices/products/handlers"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "propduct-api", log.LstdFlags)
	//hh := handlers.NewHello(l)
	//gb := handlers.NewGoodBuy(l)
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cs := protos.NewCurrencyClient(conn)
	ph := handlers.NewProducts(l, cs)

	sm := mux.NewRouter()
	ops := middleware.RedocOpts{
		SpecURL: "./swagger.yaml",
	}
	sh := middleware.Redoc(ops, nil)

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)
	getRouter.Handle("/docs", sh)

	putRouter := sm.Methods("PUT").Subrouter()

	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods("DELETE").Subrouter()

	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProduct)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
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
