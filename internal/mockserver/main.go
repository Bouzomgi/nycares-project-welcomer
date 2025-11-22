package main

import (
	"fmt"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	routes.RegisterLoginRoute(r)
	routes.RegisterMessageRoutes(r)
	routes.RegisterScheduleRoute(r)

	fmt.Println("Mock server running on http://localhost:3001")
	http.ListenAndServe(":3001", r)
}
