package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Set up logging
	logFile, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	// Set up SQLite database
	const createTableQuery = `CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT,
		user_agent TEXT,
		device_type TEXT,
		visited_at DATETIME
	)`
	db, err := sql.Open("sqlite3", "requests.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		userAgent := r.UserAgent()
		deviceType := "Desktop"
		if strings.Contains(strings.ToLower(userAgent), "mobile") {
			deviceType = "Mobile"
		}
		visitedAt := time.Now()

		// Log to file
		logger.Printf("IP: %s | User-Agent: %s | Time: %s\n", ipAddress, userAgent, visitedAt.Format(time.RFC3339))

		// Insert into database
		_, err := db.Exec("INSERT INTO requests (ip_address, user_agent, device_type, visited_at) VALUES (?, ?, ?, ?)", ipAddress, userAgent, deviceType, visitedAt)
		if err != nil {
			log.Printf("Failed to insert record: %v", err)
		}

		// Serve HTML
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT device_type, COUNT(*) FROM requests GROUP BY device_type")
		if err != nil {
			http.Error(w, "Failed to fetch analytics", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		fmt.Fprintf(w, "<h1>Analytics</h1>")
		fmt.Fprintf(w, "<table border='1'><tr><th>Device Type</th><th>Visits</th></tr>")
		for rows.Next() {
			var deviceType string
			var count int
			if err := rows.Scan(&deviceType, &count); err != nil {
				http.Error(w, "Failed to scan analytics", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "<tr><td>%s</td><td>%d</td></tr>", deviceType, count)
		}
		fmt.Fprintf(w, "</table>")
	})

	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", nil)
}
