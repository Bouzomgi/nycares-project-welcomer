package main

import (
	"fmt"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	routes.RegisterAdminRoutes(r)
	routes.RegisterLoginRoute(r)
	routes.RegisterMessageRoutes(r)
	routes.RegisterScheduleRoute(r)

	// Health check endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	}).Methods("GET")

	fmt.Println("Mock server running on http://localhost:3001")
	http.ListenAndServe(":3001", r)
}
