package main

import (
	"apex-challenge/routers"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	fmt.Println("Server Started")
	router := routers.NewRouter()
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	fmt.Print(handler)
}
