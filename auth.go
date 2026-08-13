package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const sessionCookie = "gensa_admin_session"

func ensureAdminUser(email, password string) {
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM admin_users`).Scan(&count)
	if count > 0 {
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.Exec(`INSERT INTO admin_users (email, password_hash) VALUES (?,?)`, email, string(hash))
}

func checkLogin(email, password string) bool {
	var hash string
	err := db.QueryRow(`SELECT password_hash FROM admin_users WHERE email=?`, email).Scan(&hash)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func newSessionToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func createSession(w http.ResponseWriter) {
	token := newSessionToken()
	var userID int64
	db.QueryRow(`SELECT id FROM admin_users LIMIT 1`).Scan(&userID)
	db.Exec(`INSERT INTO sessions (token, user_id) VALUES (?,?)`, token, userID)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
}

func destroySession(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		db.Exec(`DELETE FROM sessions WHERE token=?`, c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/", MaxAge: -1})
}

func isLoggedIn(r *http.Request) bool {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return false
	}
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM sessions WHERE token=?`, c.Value).Scan(&count)
	return count > 0
}

// requireAuth wraps an admin handler, redirecting to /login when the
// request has no valid session cookie.
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
