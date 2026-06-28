package main

import (
	"net/http"	
	"os"
	"github.com/go-chi/chi/v5"
	"fmt"
)

func main() {
	router := chi.NewRouter()
	
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

}
