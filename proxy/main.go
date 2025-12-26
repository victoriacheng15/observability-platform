package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"proxy/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *sql.DB
var mongoClient *mongo.Client


func readingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"placeholder": "reading placeholder still"})
}

// --- ETL: Sync Mongo to Postgres ---
func syncReadingHandler(w http.ResponseWriter, r *http.Request) {
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	coll := mongoClient.Database(mongoDBName).Collection(mongoCollection)

	// 1. Extract: Find documents where status is "ingested"
	ctx := context.TODO()
	filter := bson.M{"status": "ingested"}
	// Limit to batch size of 10 to avoid timeouts
	opts := options.Find().SetLimit(10)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("❌ ETL_ERROR: Failed to query Mongo: %v", err)
		http.Error(w, "Failed to query Mongo", 500)
		return
	}
	defer cursor.Close(ctx)

	processedCount := 0

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("⚠️ ETL_WARN: Failed to decode document: %v", err)
			continue
		}

		// Get the ObjectID to update later
		objID, ok := doc["_id"].(primitive.ObjectID)
		if !ok {
			log.Printf("⚠️ ETL_WARN: Document missing ObjectID")
			continue
		}

		// 2. Load: Insert into Postgres as JSONB
		// We extract specific fields for indexing, but keep the full doc in payload
		eventType, _ := doc["event_type"].(string)
		jsonData, _ := json.Marshal(doc)

		_, err = db.Exec(
			`INSERT INTO reading_analytics (mongo_id, event_type, payload, created_at) 
			 VALUES ($1, $2, $3, NOW())
			 ON CONFLICT (mongo_id) DO NOTHING`,
			objID.Hex(), eventType, jsonData,
		)

		if err != nil {
			log.Printf("❌ ETL_ERROR: Failed to insert into Postgres (ID: %s): %v", objID.Hex(), err)
			continue
		}

		// 3. Update Source: Mark as processed in Mongo (Commented out for testing)
		update := bson.M{"$set": bson.M{"status": "processed"}}
		_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			log.Printf("⚠️ ETL_WARN: Failed to update Mongo status (ID: %s): %v", objID.Hex(), err)
		} else {
			processedCount++
		}

		processedCount++
	}

	log.Printf("✅ ETL_SUCCESS: Processed batch of %d documents", processedCount)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "success",
		"processed_count": processedCount,
	})
}

func main() {
	// Try loading from the current directory or the parent directory (root)
	godotenv.Load(".env", "../.env")

	db = utils.InitPostgres()
	mongoClient = utils.InitMongo()

	// Ensure table exists with JSONB support
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS reading_analytics (
		id SERIAL PRIMARY KEY,
		mongo_id TEXT UNIQUE NOT NULL,
		event_type TEXT,
		payload JSONB,
		created_at TIMESTAMP DEFAULT NOW()
	)`)
	if err != nil {
		log.Fatal("❌ Failed to create reading_analytics table:", err)
	}

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	// Register HTTP handlers with logging middleware
	http.HandleFunc("/", utils.WithLogging(utils.HomeHandler))
	http.HandleFunc("/api/reading", utils.WithLogging(readingHandler))
	http.HandleFunc("/api/sync/reading", utils.WithLogging(syncReadingHandler))

	log.Printf("🚀 The GO proxy listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
