package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Price is stored in cents to avoid floating-point rounding issues.
// TODO(phase2): separate into model.go + dto.go
type Product struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Price       int           `json:"price" bson:"price"`
	ImageURL    string        `json:"imageUrl" bson:"imageUrl"`
	Category    string        `json:"category" bson:"category"`
	Stock       int           `json:"stock" bson:"stock"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
}

// TODO(phase2): split into handler/service/repository layers
func handleGetProducts(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := db.Collection("products")

		cursor, err := collection.Find(r.Context(), bson.D{})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch products"})
			return
		}
		defer cursor.Close(r.Context())

		var products []Product
		if err := cursor.All(r.Context(), &products); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to decode products"})
			return
		}

		// don't return null
		if products == nil {
			products = []Product{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func handleCreateProduct(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit

		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		if strings.TrimSpace(product.Name) == "" || product.Price <= 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "name is required and price must be a positive number (in cents)"})
			return
		}

		product.CreatedAt = time.Now()

		collection := db.Collection("products")
		result, err := collection.InsertOne(r.Context(), product)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to create product"})
			return
		}

		product.ID = result.InsertedID.(bson.ObjectID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}
