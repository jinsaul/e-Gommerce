// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (makes env vars available via os.Getenv)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}

	// Read environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "egommerce" // default database name
	}

	// Connect to MongoDB
	db, err := ConnectDB(mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Register routes (Go 1.22+ method-based routing)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/products", handleGetProducts(db))
	mux.HandleFunc("POST /api/v1/products", handleCreateProduct(db))

	// Wrap the mux with CORS middleware
	// (needed because Angular runs on :4200 and Go on :8080)
	handler := enableCORS(mux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// enableCORS wraps an http.Handler to add CORS headers.
// In Phase 1 we keep it simple — allow everything from localhost:4200.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
