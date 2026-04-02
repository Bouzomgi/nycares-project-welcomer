package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/routes"
	"github.com/akrylysov/algnhsa"
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

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda mode: Function URL sends APIGatewayV2 payloads; algnhsa auto-detects
		algnhsa.ListenAndServe(r, nil)
	} else {
		// Local/Docker mode: unchanged
		fmt.Println("Mock server running on http://localhost:3001")
		if err := http.ListenAndServe(":3001", r); err != nil {
			panic(err)
		}
	}
}
