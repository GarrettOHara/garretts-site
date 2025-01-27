package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Global variables for DB and logger
var db *sql.DB
var logger *log.Logger

func main() {
	fmt.Println("Running server...")
	// Set up logging
	logFile, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	logger = log.New(logFile, "", log.LstdFlags)

	// Set up SQLite database
	const createTableQuery = `CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT,
		user_agent TEXT,
		device_type TEXT,
		visited_at DATETIME
	)`
	db, err = sql.Open("sqlite3", "requests.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Register routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/analytics", analyticsHandler)

	// Start the server
	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", nil)
}

// Home route handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
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
}

type Analytics struct {
    DeviceStats    []DeviceStat
    BrowserStats   []BrowserStat
    PlatformStats  []PlatformStat
}

type DeviceStat struct {
    DeviceType string
    Count      int
    Percentage float64
}

type BrowserStat struct {
    Browser    string
    Count      int
    Percentage float64
}

type PlatformStat struct {
    Platform   string
    Count      int
    Percentage float64
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("sqlite3", "requests.db")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Get device statistics
    deviceStats, err := getDeviceStats(db)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Get browser statistics
    browserStats, err := getBrowserStats(db)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Get platform statistics
    platformStats, err := getPlatformStats(db)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    analytics := Analytics{
        DeviceStats:    deviceStats,
        BrowserStats:   browserStats,
        PlatformStats:  platformStats,
    }

    tmpl, err := template.ParseFiles("analytics.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    err = tmpl.Execute(w, analytics)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func getDeviceStats(db *sql.DB) ([]DeviceStat, error) {
    rows, err := db.Query(`
        SELECT device_type, COUNT(*) as count,
        ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM requests), 2) as percentage
        FROM requests
        GROUP BY device_type
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []DeviceStat
    for rows.Next() {
        var stat DeviceStat
        err := rows.Scan(&stat.DeviceType, &stat.Count, &stat.Percentage)
        if err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }
    return stats, nil
}

func getBrowserStats(db *sql.DB) ([]BrowserStat, error) {
    rows, err := db.Query(`
        WITH BrowserCounts AS (
            SELECT 
                CASE 
                    WHEN user_agent LIKE '%Firefox%' THEN 'Firefox'
                    WHEN user_agent LIKE '%Chrome%' THEN 'Chrome'
                    WHEN user_agent LIKE '%Safari%' AND user_agent NOT LIKE '%Chrome%' THEN 'Safari'
                    ELSE 'Other'
                END as browser,
                COUNT(*) as count,
                ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM requests), 2) as percentage
            FROM requests
            GROUP BY 
                CASE 
                    WHEN user_agent LIKE '%Firefox%' THEN 'Firefox'
                    WHEN user_agent LIKE '%Chrome%' THEN 'Chrome'
                    WHEN user_agent LIKE '%Safari%' AND user_agent NOT LIKE '%Chrome%' THEN 'Safari'
                    ELSE 'Other'
                END
        )
        SELECT * FROM BrowserCounts
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []BrowserStat
    for rows.Next() {
        var stat BrowserStat
        err := rows.Scan(&stat.Browser, &stat.Count, &stat.Percentage)
        if err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }
    return stats, nil
}

func getPlatformStats(db *sql.DB) ([]PlatformStat, error) {
    rows, err := db.Query(`
        WITH PlatformCounts AS (
            SELECT 
                CASE 
                    WHEN user_agent LIKE '%Linux%' THEN 'Linux'
                    WHEN user_agent LIKE '%iPhone%' THEN 'iPhone'
                    WHEN user_agent LIKE '%Windows%' THEN 'Windows'
                    WHEN user_agent LIKE '%Mac OS%' THEN 'Mac OS'
                    ELSE 'Other'
                END as platform,
                COUNT(*) as count,
                ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM requests), 2) as percentage
            FROM requests
            GROUP BY 
                CASE 
                    WHEN user_agent LIKE '%Linux%' THEN 'Linux'
                    WHEN user_agent LIKE '%iPhone%' THEN 'iPhone'
                    WHEN user_agent LIKE '%Windows%' THEN 'Windows'
                    WHEN user_agent LIKE '%Mac OS%' THEN 'Mac OS'
                    ELSE 'Other'
                END
        )
        SELECT * FROM PlatformCounts
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []PlatformStat
    for rows.Next() {
        var stat PlatformStat
        err := rows.Scan(&stat.Platform, &stat.Count, &stat.Percentage)
        if err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }
    return stats, nil
}
