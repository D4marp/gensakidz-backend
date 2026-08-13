package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	os.MkdirAll(uploadsDir, 0o755)

	dbPath := envOr("DB_PATH", "gensakidz.db")
	db = openDB(dbPath)
	migrate(db)

	adminEmail := envOr("ADMIN_EMAIL", "admin@gensakidz.com")
	adminPassword := envOr("ADMIN_PASSWORD", "gensakidz2026")
	ensureAdminUser(adminEmail, adminPassword)
	seedIfEmpty()

	mux := http.NewServeMux()

	// static assets + uploaded images
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// auth
	mux.HandleFunc("GET /login", handleLoginPage)
	mux.HandleFunc("POST /login", handleLoginSubmit)
	mux.HandleFunc("GET /logout", handleLogout)

	// admin dashboard
	mux.HandleFunc("GET /admin", requireAuth(handleDashboard))
	mux.HandleFunc("GET /admin/settings", requireAuth(handleSettingsPage))
	mux.HandleFunc("POST /admin/settings", requireAuth(handleSettingsSave))

	registerCRUD(mux, "services", handleServicesList, handleServiceNew, handleServiceCreate, handleServiceEdit, handleServiceUpdate, handleServiceDelete)
	registerCRUD(mux, "jobs", handleJobsList, handleJobNew, handleJobCreate, handleJobEdit, handleJobUpdate, handleJobDelete)
	registerCRUD(mux, "articles", handleArticlesList, handleArticleNew, handleArticleCreate, handleArticleEdit, handleArticleUpdate, handleArticleDelete)
	registerCRUD(mux, "facilities", handleFacilitiesList, handleFacilityNew, handleFacilityCreate, handleFacilityEdit, handleFacilityUpdate, handleFacilityDelete)
	registerCRUD(mux, "gallery", handleGalleryList, handleGalleryNew, handleGalleryCreate, handleGalleryEdit, handleGalleryUpdate, handleGalleryDelete)
	registerCRUD(mux, "team", handleTeamList, handleTeamNew, handleTeamCreate, handleTeamEdit, handleTeamUpdate, handleTeamDelete)
	registerCRUD(mux, "branches", handleBranchesList, handleBranchNew, handleBranchCreate, handleBranchEdit, handleBranchUpdate, handleBranchDelete)
	registerCRUD(mux, "testimonials", handleTestimonialsList, handleTestimonialNew, handleTestimonialCreate, handleTestimonialEdit, handleTestimonialUpdate, handleTestimonialDelete)

	// public read-only JSON API for the frontend to consume
	mux.HandleFunc("GET /api/services", handleAPIServices)
	mux.HandleFunc("GET /api/jobs", handleAPIJobs)
	mux.HandleFunc("GET /api/articles", handleAPIArticles)
	mux.HandleFunc("GET /api/facilities", handleAPIFacilities)
	mux.HandleFunc("GET /api/gallery", handleAPIGallery)
	mux.HandleFunc("GET /api/team", handleAPITeam)
	mux.HandleFunc("GET /api/branches", handleAPIBranches)
	mux.HandleFunc("GET /api/testimonials", handleAPITestimonials)
	mux.HandleFunc("GET /api/settings", handleAPISettings)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})

	port := envOr("PORT", "8080")
	log.Printf("GenSA Kidz backend running on http://localhost:%s (admin login: %s)", port, adminEmail)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// registerCRUD wires the standard six routes (list/new/create/edit/update/delete)
// for an admin-managed content type, all behind requireAuth.
func registerCRUD(mux *http.ServeMux, path string,
	list, newForm, create, edit, update, del http.HandlerFunc) {
	base := "/admin/" + path
	mux.HandleFunc("GET "+base, requireAuth(list))
	mux.HandleFunc("GET "+base+"/new", requireAuth(newForm))
	mux.HandleFunc("POST "+base+"/new", requireAuth(create))
	mux.HandleFunc("GET "+base+"/{id}/edit", requireAuth(edit))
	mux.HandleFunc("POST "+base+"/{id}/edit", requireAuth(update))
	mux.HandleFunc("POST "+base+"/{id}/delete", requireAuth(del))
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
