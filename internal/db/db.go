package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func Initialize(dbPath string) (*sql.DB, error) {
    const createTableQuery = `CREATE TABLE IF NOT EXISTS requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        ip_address TEXT,
        user_agent TEXT,
        device_type TEXT,
        visited_at DATETIME
    )`

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(createTableQuery)
    if err != nil {
        return nil, err
    }

    return db, nil
}
