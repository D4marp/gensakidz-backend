package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func openDB(path string) *sql.DB {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("pragma: %v", err)
	}
	return conn
}

const schema = `
CREATE TABLE IF NOT EXISTS admin_users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS services (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT UNIQUE NOT NULL,
	title TEXT NOT NULL,
	icon TEXT NOT NULL,
	short TEXT NOT NULL,
	detail TEXT NOT NULL DEFAULT '[]',
	for_who TEXT NOT NULL DEFAULT '',
	signs TEXT NOT NULL DEFAULT '[]',
	goal TEXT NOT NULL DEFAULT '',
	process TEXT NOT NULL DEFAULT '[]',
	duration TEXT NOT NULL DEFAULT '',
	professionals TEXT NOT NULL DEFAULT '',
	what_to_bring TEXT NOT NULL DEFAULT '[]',
	extra_faq TEXT NOT NULL DEFAULT '[]',
	image_path TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS jobs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT UNIQUE NOT NULL,
	title TEXT NOT NULL,
	branch TEXT NOT NULL DEFAULT 'Lamongan',
	type TEXT NOT NULL DEFAULT 'Full-time',
	status TEXT NOT NULL DEFAULT 'Dibuka',
	description TEXT NOT NULL DEFAULT '',
	requirements TEXT NOT NULL DEFAULT '[]',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS articles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT UNIQUE NOT NULL,
	title TEXT NOT NULL,
	category TEXT NOT NULL DEFAULT '',
	excerpt TEXT NOT NULL DEFAULT '',
	content TEXT NOT NULL DEFAULT '[]',
	image_path TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS facilities (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	icon TEXT NOT NULL DEFAULT 'heart',
	image_path TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS gallery_photos (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	category TEXT NOT NULL DEFAULT 'aktivitas',
	caption TEXT NOT NULL DEFAULT '',
	image_path TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS team_members (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT '',
	image_path TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS branches (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT UNIQUE NOT NULL,
	name TEXT NOT NULL,
	address TEXT NOT NULL DEFAULT '',
	whatsapp TEXT NOT NULL DEFAULT '',
	phone TEXT NOT NULL DEFAULT '',
	maps_query TEXT NOT NULL DEFAULT '',
	maps_url TEXT NOT NULL DEFAULT '',
	schedules TEXT NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS testimonials (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT '',
	quote TEXT NOT NULL DEFAULT '',
	sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS site_settings (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL DEFAULT ''
);
`

func migrate(conn *sql.DB) {
	if _, err := conn.Exec(schema); err != nil {
		log.Fatalf("migrate: %v", err)
	}
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
