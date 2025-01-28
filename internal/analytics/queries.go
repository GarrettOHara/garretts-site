package analytics

import (
    "database/sql"
)

func GetDeviceStats(db *sql.DB) ([]DeviceStat, error) {
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

func GetBrowserStats(db *sql.DB) ([]BrowserStat, error) {
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

func GetPlatformStats(db *sql.DB) ([]PlatformStat, error) {
    rows, err := db.Query(`
        WITH PlatformCounts AS (
            SELECT 
                CASE 
                    WHEN user_agent LIKE '%Linux%' THEN 'Linux'
                    WHEN user_agent LIKE '%iPhone%' THEN 'iPhone'
                    WHEN user_agent LIKE '%Windows%' THEN 'Windows'
                    WHEN user_agent LIKE '%Mac OS%' THEN 'Mac OS'
		    WHEN user_agent LIKE '%Andriod%' THEN 'Andriod'
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
		    WHEN user_agent LIKE '%Andriod%' THEN 'Andriod'
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

func GetRequestStats(db *sql.DB) ([]TimeSeriesData, error) {
    // Query to get requests over time (grouped by hour for example)
    rows, err := db.Query(`
        SELECT 
    	strftime('%Y-%m-%d %H:00', visited_at) as hour,
    	COUNT(*) as count 
        FROM requests 
        GROUP BY strftime('%Y-%m-%d %H:00', visited_at)
        ORDER BY hour
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var stats []TimeSeriesData
    for rows.Next() {
        var data TimeSeriesData
        err := rows.Scan(&data.Time, &data.Count)
        if err != nil {
            return nil, err
        }
        stats = append(stats, data)
    }
    return stats, nil
}

func GetDistinctIPCount(db *sql.DB) (int, error) {
    var count int
    query := `
        SELECT COUNT(DISTINCT
            SUBSTR(ip_address, 1, INSTR(ip_address, ':') - 1) -- Extract only the IP part before the port
        ) 
        FROM requests
    `

    err := db.QueryRow(query).Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}
