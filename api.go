package main

import (
	"encoding/json"
	"net/http"
)

// The public API is read-only JSON, meant for the Next.js frontend (or
// anything else) to consume once it's wired to fetch from here instead of
// its static TypeScript data files. CORS is wide open since this only ever
// serves public marketing content, never anything user-specific.

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(v)
}

func handleAPIServices(w http.ResponseWriter, r *http.Request)     { writeJSON(w, listServices()) }
func handleAPIJobs(w http.ResponseWriter, r *http.Request)         { writeJSON(w, listJobs()) }
func handleAPIArticles(w http.ResponseWriter, r *http.Request)     { writeJSON(w, listArticles()) }
func handleAPIFacilities(w http.ResponseWriter, r *http.Request)   { writeJSON(w, listFacilities()) }
func handleAPIGallery(w http.ResponseWriter, r *http.Request)      { writeJSON(w, listGallery()) }
func handleAPITeam(w http.ResponseWriter, r *http.Request)         { writeJSON(w, listTeam()) }
func handleAPIBranches(w http.ResponseWriter, r *http.Request)     { writeJSON(w, listBranches()) }
func handleAPITestimonials(w http.ResponseWriter, r *http.Request) { writeJSON(w, listTestimonials()) }

func handleAPISettings(w http.ResponseWriter, r *http.Request) {
	out := map[string]string{}
	for _, k := range settingsKeys {
		out[k.Key] = getSetting(k.Key, "")
	}
	writeJSON(w, out)
}
