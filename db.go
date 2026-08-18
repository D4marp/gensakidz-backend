package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// mysqlDSN builds a go-sql-driver/mysql DSN from discrete connection
// parameters — easier to wire up via docker-compose / PaaS env vars than a
// single connection-string env var.
func mysqlDSN(host, port, user, password, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, password, host, port, name)
}

// openDB connects to MySQL, retrying for a bit since in docker-compose the
// backend container often starts before MySQL is ready to accept connections.
func openDB(dsn string) *sql.DB {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	var lastErr error
	for i := 0; i < 20; i++ {
		if lastErr = conn.Ping(); lastErr == nil {
			return conn
		}
		log.Printf("waiting for MySQL... (%v)", lastErr)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("could not connect to MySQL: %v", lastErr)
	return nil
}

const schema = `
CREATE TABLE IF NOT EXISTS admin_users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS sessions (
	token VARCHAR(64) PRIMARY KEY,
	user_id INT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS services (
	id INT AUTO_INCREMENT PRIMARY KEY,
	slug VARCHAR(255) UNIQUE NOT NULL,
	title VARCHAR(255) NOT NULL,
	icon VARCHAR(50) NOT NULL,
	short TEXT NOT NULL,
	detail TEXT NOT NULL,
	for_who TEXT NOT NULL,
	signs TEXT NOT NULL,
	goal TEXT NOT NULL,
	process TEXT NOT NULL,
	duration VARCHAR(255) NOT NULL DEFAULT '',
	professionals VARCHAR(255) NOT NULL DEFAULT '',
	what_to_bring TEXT NOT NULL,
	extra_faq TEXT NOT NULL,
	image_path VARCHAR(500) NOT NULL DEFAULT '',
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS jobs (
	id INT AUTO_INCREMENT PRIMARY KEY,
	slug VARCHAR(255) UNIQUE NOT NULL,
	title VARCHAR(255) NOT NULL,
	branch VARCHAR(100) NOT NULL DEFAULT 'Lamongan',
	type VARCHAR(50) NOT NULL DEFAULT 'Full-time',
	status VARCHAR(50) NOT NULL DEFAULT 'Dibuka',
	description TEXT NOT NULL,
	requirements TEXT NOT NULL,
	image_path VARCHAR(500) NOT NULL DEFAULT '',
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS articles (
	id INT AUTO_INCREMENT PRIMARY KEY,
	slug VARCHAR(255) UNIQUE NOT NULL,
	title VARCHAR(255) NOT NULL,
	category VARCHAR(100) NOT NULL DEFAULT '',
	excerpt TEXT NOT NULL,
	content TEXT NOT NULL,
	image_path VARCHAR(500) NOT NULL DEFAULT '',
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS facilities (
	id INT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	icon VARCHAR(50) NOT NULL DEFAULT 'heart',
	image_path VARCHAR(500) NOT NULL DEFAULT '',
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gallery_photos (
	id INT AUTO_INCREMENT PRIMARY KEY,
	category VARCHAR(50) NOT NULL DEFAULT 'aktivitas',
	caption VARCHAR(255) NOT NULL DEFAULT '',
	image_path VARCHAR(500) NOT NULL DEFAULT '',
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS team_members (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	role VARCHAR(255) NOT NULL DEFAULT '',
	image_path VARCHAR(500) NOT NULL DEFAULT '',
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS branches (
	id INT AUTO_INCREMENT PRIMARY KEY,
	slug VARCHAR(255) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	address TEXT NOT NULL,
	whatsapp VARCHAR(50) NOT NULL DEFAULT '',
	phone VARCHAR(50) NOT NULL DEFAULT '',
	maps_query VARCHAR(500) NOT NULL DEFAULT '',
	maps_url VARCHAR(500) NOT NULL DEFAULT '',
	schedules TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS testimonials (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	role VARCHAR(255) NOT NULL DEFAULT '',
	quote TEXT NOT NULL,
	sort_order INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS site_settings (
	setting_key VARCHAR(100) PRIMARY KEY,
	value TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

func migrate(conn *sql.DB) {
	for _, stmt := range splitStatements(schema) {
		if _, err := conn.Exec(stmt); err != nil {
			log.Fatalf("migrate: %v\nstatement: %s", err, stmt)
		}
	}
	addColumnIfMissing(conn, "jobs", "image_path", "VARCHAR(500) NOT NULL DEFAULT ''")
}

// addColumnIfMissing runs an ALTER TABLE ADD COLUMN only if the column
// doesn't already exist — MySQL (unlike MariaDB) has no "ADD COLUMN IF NOT
// EXISTS" syntax, so this checks information_schema first.
func addColumnIfMissing(conn *sql.DB, table, column, definition string) {
	var count int
	err := conn.QueryRow(
		`SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?`,
		table, column,
	).Scan(&count)
	if err != nil {
		log.Fatalf("check column %s.%s: %v", table, column, err)
	}
	if count > 0 {
		return
	}
	if _, err := conn.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition)); err != nil {
		log.Fatalf("add column %s.%s: %v", table, column, err)
	}
}

// splitStatements splits the schema block on ";" — MySQL's driver doesn't
// support multi-statement Exec by default, unlike SQLite.
func splitStatements(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			stmt := trimSchema(s[start:i])
			if stmt != "" {
				out = append(out, stmt)
			}
			start = i + 1
		}
	}
	return out
}

func trimSchema(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == '\n' || s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == '\n' || s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// --- small JSON array helpers ---

func toJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func fromJSONStrings(s string) []string {
	var out []string
	if s == "" {
		return out
	}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}

func fromJSONFAQ(s string) []FAQItem {
	var out []FAQItem
	if s == "" {
		return out
	}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}

func fromJSONSchedules(s string) []ScheduleItem {
	var out []ScheduleItem
	if s == "" {
		return out
	}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}
