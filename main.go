package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultHost    = "127.0.0.1" // changed from 0.0.0.0 to localhost-only for personal use
	appName        = "sub2api"
	appVersion     = "dev"
)

func main() {
	var (
		host    string
		port    int
		version bool
	)

	flag.StringVar(&host, "host", getEnvOrDefault("HOST", defaultHost), "Host address to listen on")
	flag.IntVar(&port, "port", getEnvIntOrDefault("PORT", defaultPort), "Port to listen on")
	flag.BoolVar(&version, "version", false, "Print version information and exit")
	flag.Parse()

	if version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/sub", handler.SubHandler)
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/", handler.IndexHandler)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting %s on %s", appName, addr)

	// Added read/write timeouts to avoid hanging connections on my local setup
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnvOrDefault returns the value of the environment variable named by key,
// or the defaultValue if the variable is not set or empty.
func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// getEnvIntOrDefault returns the integer value of the environment variable named
// by key, or the defaultValue if the variable is not set, empty, or not a valid integer.
func getEnvIntOrDefault(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		log.Printf("Warning: invalid integer value for %s, using default %d", key, defaultValue)
	}
	return defaultValue
}
