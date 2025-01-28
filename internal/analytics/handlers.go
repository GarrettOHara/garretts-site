package analytics

import (
    "database/sql"
    "html/template"
    "net/http"
)

type LoggerDB interface {
    CaptureStats(w http.ResponseWriter, r *http.Request)
}

func HandleAnalytics(db *sql.DB, s LoggerDB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
	s.CaptureStats(w, r)

        deviceStats, err := GetDeviceStats(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        browserStats, err := GetBrowserStats(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        platformStats, err := GetPlatformStats(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        requestsOverTime, err := GetRequestStats(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

	distinctIPCount, err := GetDistinctIPCount(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        analytics := Analytics{
            DeviceStats:      deviceStats,
            BrowserStats:     browserStats,
            PlatformStats:    platformStats,
            RequestsOverTime: requestsOverTime,
            DistinctIPCount:  distinctIPCount,
        }

        tmpl, err := template.ParseFiles("templates/analytics.html")
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
}
