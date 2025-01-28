package main

import (
    "fmt"
    "log"
    "os"

    "github.com/garrettohara/garretts-site/internal/db"
    "github.com/garrettohara/garretts-site/internal/server"
)

func main() {
    // Set up logging
    logFile, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer logFile.Close()
    logger := log.New(logFile, "", log.LstdFlags)

    // Initialize database
    database, err := db.Initialize("requests.db")
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.Close()

    // Start server
    srv := server.New(database, logger)
    fmt.Println("Serving on :8080")
    log.Fatal(srv.Start(":8080"))
}
