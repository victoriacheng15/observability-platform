package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"system-metrics/collectors"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v4/host"
)

func main() {
	// Load .env (current or parent)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	// 1. Initial Detection
	hInfo, err := host.Info()
	if err != nil {
		log.Fatalf("❌ Error getting host info: %v", err)
	}
	osName := fmt.Sprintf("%s %s", hInfo.Platform, hInfo.PlatformVersion)

	hostName, _ := os.Hostname()
	if hostName == "" {
		hostName = "homelab"
	}

	// 2. Database Connection
	connStr := getConnStr()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(ctx)

	// 3. Ensure Schema
	ensureSchema(ctx, conn)

	// 4. Collect and Store Once
	collectAndStore(ctx, conn, hostName, osName)
}

func collectAndStore(ctx context.Context, conn *pgx.Conn, hostName string, osName string) {
	now := time.Now().UTC().Truncate(time.Second)

	// Collect
	cpu, _ := collectors.GetCPUStats()
	mem, _ := collectors.GetMemoryStats()
	disk, _ := collectors.GetDiskStats()
	net, _ := collectors.GetNetworkStats()

	// Store to DB
	metrics := []struct {
		mType   string
		payload interface{}
	}{
		{"cpu", cpu},
		{"memory", mem},
		{"disk", disk},
		{"network", net},
	}

	var insertErrors []string
	for _, m := range metrics {
		if m.payload == nil {
			continue
		}
		payloadJSON, _ := json.Marshal(m.payload)
		_, err := conn.Exec(ctx,
			"INSERT INTO system_metrics (time, host, os, metric_type, payload) VALUES ($1, $2, $3, $4, $5)",
			now, hostName, osName, m.mType, payloadJSON,
		)
		if err != nil {
			errMsg := fmt.Sprintf("❌ Failed to insert %s metric: %v", m.mType, err)
			log.Println(errMsg)
			insertErrors = append(insertErrors, errMsg)
		}
	}

	// Log success only at the top of the hour and if no errors
	if now.Minute() == 0 {
		if len(insertErrors) == 0 {
			fmt.Printf("[%s] ✅ Metrics stored in database.\n", now.Format("15:04:05"))
		} else {
			fmt.Printf("[%s] ⚠️ Metrics collection completed with %d error(s).\n", now.Format("15:04:05"), len(insertErrors))
		}
	}
}

func ensureSchema(ctx context.Context, conn *pgx.Conn) {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS system_metrics (
			time TIMESTAMPTZ(0) NOT NULL,
			host TEXT NOT NULL,
			os TEXT NOT NULL,
			metric_type TEXT NOT NULL,
			payload JSONB NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("❌ Failed to ensure schema: %v", err)
	}

	// Enable hypertable if TimescaleDB is available
	_, err = conn.Exec(ctx, "SELECT create_hypertable('system_metrics', 'time', if_not_exists => true);")
	if err != nil {
		log.Printf("ℹ️ Hypertable check: %v (ignoring if not using TimescaleDB)", err)
	}
}

func getConnStr() string {
	if connStr := os.Getenv("DATABASE_URL"); connStr != "" {
		return connStr
	}

	host := getEnv("DB_HOST")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER")
	dbname := getEnv("DB_NAME")
	password := os.Getenv("SERVER_DB_PASSWORD")

	if host == "" {
		log.Fatal("❌ DB_HOST is not set")
	}
	if user == "" {
		log.Fatal("❌ DB_USER is not set")
	}
	if dbname == "" {
		log.Fatal("❌ DB_NAME is not set")
	}
	if password == "" {
		log.Fatal("❌ SERVER_DB_PASSWORD is not set")
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func getEnv(key string, fallback ...string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}
