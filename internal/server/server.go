package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/garrettohara/garretts-site/internal/analytics"
)

type Server struct {
	db     *sql.DB
	logger *log.Logger
}

func New(db *sql.DB, logger *log.Logger) *Server {
	return &Server{
		db:     db,
		logger: logger,
	}
}

func (s *Server) Start(addr string) error {
	// Serve static files first
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", s.handleHome())
	http.HandleFunc("/analytics", analytics.HandleAnalytics(s.db, s))

	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.CaptureStats(w, r)
		http.ServeFile(w, r, "templates/index.html")
	}
}

func (s *Server) CaptureStats(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	} else {
		/**
		 * X-Forwarded-For (XFF) is a comma-separated list of IPs:
		 *     - First IP: The original client’s IP (usually).
		 *     - Middle IPs: Any intermediate proxies.
		 *     - Last IP: The most recent proxy before reaching your server.
		 **/
		ips := strings.Split(ipAddress, ",")
		ipAddress = strings.TrimSpace(ips[0])
	}

	userAgent := r.UserAgent()
	deviceType := "Desktop"
	if strings.Contains(strings.ToLower(userAgent), "mobile") {
		deviceType = "Mobile"
	}
	visitedAt := time.Now()

	s.logger.Printf("IP: %s | User-Agent: %s | Time: %s\n",
		ipAddress, userAgent, visitedAt.Format(time.RFC3339))
	fmt.Printf("IP: %s | User-Agent: %s | Time: %s\n",
		ipAddress, userAgent, visitedAt.Format(time.RFC3339))

	_, err := s.db.Exec(
		"INSERT INTO requests (ip_address, user_agent, device_type, visited_at) VALUES (?, ?, ?, ?)",
		ipAddress, userAgent, deviceType, visitedAt)
	if err != nil {
		s.logger.Printf("Failed to insert record: %v", err)
	}
}
