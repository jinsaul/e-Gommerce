// product.go
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Product represents a product in the store.
// The same struct is used for both MongoDB documents and JSON responses (Phase 1 — no DTOs yet).
type Product struct {
	ID          bson.ObjectID `json:"id"          bson:"_id,omitempty"`
	Name        string        `json:"name"        bson:"name"`
	Description string        `json:"description" bson:"description"`
	Price       float64       `json:"price"       bson:"price"`
	ImageURL    string        `json:"imageUrl"    bson:"imageUrl"`
	Category    string        `json:"category"    bson:"category"`
	Stock       int           `json:"stock"       bson:"stock"`
	CreatedAt   time.Time     `json:"createdAt"   bson:"createdAt"`
}

// handleGetProducts returns all products from the database.
func handleGetProducts(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := db.Collection("products")

		cursor, err := collection.Find(r.Context(), bson.D{})
		if err != nil {
			http.Error(w, `{"error": "failed to fetch products"}`, http.StatusInternalServerError)
			return
		}

		var products []Product
		if err := cursor.All(r.Context(), &products); err != nil {
			http.Error(w, `{"error": "failed to decode products"}`, http.StatusInternalServerError)
			return
		}

		// Return empty array instead of null when no products exist
		if products == nil {
			products = []Product{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// handleCreateProduct inserts a new product into the database.
func handleCreateProduct(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product Product

		// Decode the JSON request body into the Product struct
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Set the creation timestamp
		product.CreatedAt = time.Now()

		collection := db.Collection("products")
		result, err := collection.InsertOne(r.Context(), product)
		if err != nil {
			http.Error(w, `{"error": "failed to create product"}`, http.StatusInternalServerError)
			return
		}

		// Set the ID from the insert result
		product.ID = result.InsertedID.(bson.ObjectID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}
