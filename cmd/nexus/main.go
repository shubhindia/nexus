package main

import (
	"log"
	"net/http"

	"github.com/shubhindia/nexus/internal/gateway"
	"github.com/shubhindia/nexus/internal/provider"
)

func main() {
	p := provider.NewLlamaCPP("http://localhost:8080")
	_ = gateway.New(p)

	log.Println("Nexus starting on :8081")

	if err := http.ListenAndServe(":8081", http.NewServeMux()); err != nil {
		log.Fatal(err)
	}
}
