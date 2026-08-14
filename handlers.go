package main

import (
	"net/http"
	"strconv"
	"strings"
)

var iconOptions = []string{"speech", "hands", "run", "puzzle", "book", "heart", "growth", "clipboard"}

func idParam(r *http.Request) int64 {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id
}

// ---------- Auth ----------

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	render(w, "login", map[string]any{})
}

func handleLoginSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email, password := r.FormValue("email"), r.FormValue("password")
	if !checkLogin(email, password) {
		render(w, "login", map[string]any{"Error": "Email atau kata sandi salah."})
		return
	}
	createSession(w)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	destroySession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ---------- Dashboard ----------

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	count := func(table string) int {
		var n int
		db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n)
		return n
	}
	render(w, "dashboard", DashboardPage{Cards: []Card{
		{URL: "/admin/services", Label: "Layanan", Count: count("services")},
		{URL: "/admin/jobs", Label: "Lowongan Karir", Count: count("jobs")},
		{URL: "/admin/articles", Label: "Artikel", Count: count("articles")},
		{URL: "/admin/facilities", Label: "Fasilitas", Count: count("facilities")},
		{URL: "/admin/gallery", Label: "Foto Galeri", Count: count("gallery_photos")},
		{URL: "/admin/team", Label: "Anggota Tim", Count: count("team_members")},
		{URL: "/admin/branches", Label: "Cabang", Count: count("branches")},
		{URL: "/admin/testimonials", Label: "Testimoni", Count: count("testimonials")},
	}})
}

// ================= Services =================

func handleServicesList(w http.ResponseWriter, r *http.Request) {
	items := listServices()
	var rows []ListRow
	for _, s := range items {
		rows = append(rows, ListRow{
			ImageURL:  s.ImagePath,
			Cells:     []string{s.Title, s.Slug, s.Icon},
			EditURL:   "/admin/services/" + strconv.FormatInt(s.ID, 10) + "/edit",
			DeleteURL: "/admin/services/" + strconv.FormatInt(s.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Layanan", NewURL: "/admin/services/new", Columns: []string{"Judul", "Slug", "Ikon"}, Rows: rows, Count: len(items)})
}

func serviceForm(s Service, action, backURL, deleteURL string, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL, ImageURL: s.ImagePath,
		Fields: []FieldDef{
			{Name: "slug", Label: "Slug (URL, mis. terapi-wicara)", Type: "text", Value: s.Slug},
			{Name: "title", Label: "Judul Layanan", Type: "text", Value: s.Title},
			{Name: "icon", Label: "Ikon", Type: "select", Value: s.Icon, Options: iconOptions},
			{Name: "short", Label: "Deskripsi Singkat (untuk kartu layanan)", Type: "textarea", Value: s.Short},
			{Name: "detail", Label: "Poin Detail", Type: "list", Value: listToLines(s.Detail)},
			{Name: "for_who", Label: "Untuk Siapa?", Type: "textarea", Value: s.ForWho},
			{Name: "signs", Label: "Tanda-tanda Anak Membutuhkan", Type: "list", Value: listToLines(s.Signs)},
			{Name: "goal", Label: "Apa Tujuannya?", Type: "textarea", Value: s.Goal},
			{Name: "process", Label: "Gambaran Proses", Type: "list", Value: listToLines(s.Process)},
			{Name: "duration", Label: "Durasi Sesi", Type: "text", Value: s.Duration},
			{Name: "professionals", Label: "Tenaga Profesional", Type: "text", Value: s.Professionals},
			{Name: "what_to_bring", Label: "Yang Perlu Dibawa", Type: "list", Value: listToLines(s.WhatToBring)},
			{Name: "extra_faq", Label: "FAQ Tambahan (opsional)", Type: "faq", Value: faqToText(s.ExtraFAQ)},
			{Name: "image", Label: "Foto Layanan", Type: "file"},
		},
	}
}

func handleServiceNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", serviceForm(Service{}, "/admin/services/new", "/admin/services", "", "Tambah Layanan"))
}

func handleServiceCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	s := Service{
		Slug: r.FormValue("slug"), Title: r.FormValue("title"), Icon: r.FormValue("icon"),
		Short: r.FormValue("short"), Detail: linesToList(r.FormValue("detail")), ForWho: r.FormValue("for_who"),
		Signs: linesToList(r.FormValue("signs")), Goal: r.FormValue("goal"), Process: linesToList(r.FormValue("process")),
		Duration: r.FormValue("duration"), Professionals: r.FormValue("professionals"),
		WhatToBring: linesToList(r.FormValue("what_to_bring")), ExtraFAQ: textToFAQ(r.FormValue("extra_faq")),
		ImagePath: saveUpload(r, "image", ""), SortOrder: 999,
	}
	saveService(s)
	http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
}

func handleServiceEdit(w http.ResponseWriter, r *http.Request) {
	s, ok := getService(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(s.ID, 10)
	render(w, "form", serviceForm(s, "/admin/services/"+id+"/edit", "/admin/services", "/admin/services/"+id+"/delete", "Ubah Layanan"))
}

func handleServiceUpdate(w http.ResponseWriter, r *http.Request) {
	existing, ok := getService(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	existing.Slug, existing.Title, existing.Icon = r.FormValue("slug"), r.FormValue("title"), r.FormValue("icon")
	existing.Short, existing.Detail, existing.ForWho = r.FormValue("short"), linesToList(r.FormValue("detail")), r.FormValue("for_who")
	existing.Signs, existing.Goal, existing.Process = linesToList(r.FormValue("signs")), r.FormValue("goal"), linesToList(r.FormValue("process"))
	existing.Duration, existing.Professionals = r.FormValue("duration"), r.FormValue("professionals")
	existing.WhatToBring, existing.ExtraFAQ = linesToList(r.FormValue("what_to_bring")), textToFAQ(r.FormValue("extra_faq"))
	existing.ImagePath = saveUpload(r, "image", existing.ImagePath)
	saveService(existing)
	http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
}

func handleServiceDelete(w http.ResponseWriter, r *http.Request) {
	deleteService(idParam(r))
	http.Redirect(w, r, "/admin/services", http.StatusSeeOther)
}

// ================= Jobs =================

var jobTypeOptions = []string{"Full-time", "Part-time"}
var jobStatusOptions = []string{"Dibuka", "Ditutup"}

func handleJobsList(w http.ResponseWriter, r *http.Request) {
	items := listJobs()
	var rows []ListRow
	for _, j := range items {
		rows = append(rows, ListRow{
			Cells:     []string{j.Title, j.Type, j.Status},
			EditURL:   "/admin/jobs/" + strconv.FormatInt(j.ID, 10) + "/edit",
			DeleteURL: "/admin/jobs/" + strconv.FormatInt(j.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Lowongan Karir", NewURL: "/admin/jobs/new", Columns: []string{"Posisi", "Jenis", "Status"}, Rows: rows, Count: len(items)})
}

func jobForm(j Job, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL,
		Fields: []FieldDef{
			{Name: "slug", Label: "Slug", Type: "text", Value: j.Slug},
			{Name: "title", Label: "Nama Posisi", Type: "text", Value: j.Title},
			{Name: "branch", Label: "Cabang", Type: "text", Value: j.Branch},
			{Name: "type", Label: "Jenis Pekerjaan", Type: "select", Value: j.Type, Options: jobTypeOptions},
			{Name: "status", Label: "Status Lowongan", Type: "select", Value: j.Status, Options: jobStatusOptions},
			{Name: "description", Label: "Deskripsi Pekerjaan", Type: "textarea", Value: j.Description},
			{Name: "requirements", Label: "Persyaratan", Type: "list", Value: listToLines(j.Requirements)},
		},
	}
}

func handleJobNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", jobForm(Job{Branch: "Lamongan", Type: "Full-time", Status: "Dibuka"}, "/admin/jobs/new", "/admin/jobs", "", "Tambah Lowongan"))
}

func handleJobCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveJob(Job{
		Slug: r.FormValue("slug"), Title: r.FormValue("title"), Branch: r.FormValue("branch"),
		Type: r.FormValue("type"), Status: r.FormValue("status"), Description: r.FormValue("description"),
		Requirements: linesToList(r.FormValue("requirements")), SortOrder: 999,
	})
	http.Redirect(w, r, "/admin/jobs", http.StatusSeeOther)
}

func handleJobEdit(w http.ResponseWriter, r *http.Request) {
	j, ok := getJob(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(j.ID, 10)
	render(w, "form", jobForm(j, "/admin/jobs/"+id+"/edit", "/admin/jobs", "/admin/jobs/"+id+"/delete", "Ubah Lowongan"))
}

func handleJobUpdate(w http.ResponseWriter, r *http.Request) {
	j, ok := getJob(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	j.Slug, j.Title, j.Branch = r.FormValue("slug"), r.FormValue("title"), r.FormValue("branch")
	j.Type, j.Status, j.Description = r.FormValue("type"), r.FormValue("status"), r.FormValue("description")
	j.Requirements = linesToList(r.FormValue("requirements"))
	saveJob(j)
	http.Redirect(w, r, "/admin/jobs", http.StatusSeeOther)
}

func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	deleteJob(idParam(r))
	http.Redirect(w, r, "/admin/jobs", http.StatusSeeOther)
}

// ================= Articles =================

func handleArticlesList(w http.ResponseWriter, r *http.Request) {
	items := listArticles()
	var rows []ListRow
	for _, a := range items {
		rows = append(rows, ListRow{
			ImageURL:  a.ImagePath,
			Cells:     []string{a.Title, a.Category},
			EditURL:   "/admin/articles/" + strconv.FormatInt(a.ID, 10) + "/edit",
			DeleteURL: "/admin/articles/" + strconv.FormatInt(a.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Artikel", NewURL: "/admin/articles/new", Columns: []string{"Judul", "Kategori"}, Rows: rows, Count: len(items)})
}

func articleForm(a Article, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL, ImageURL: a.ImagePath,
		Fields: []FieldDef{
			{Name: "slug", Label: "Slug (URL)", Type: "text", Value: a.Slug},
			{Name: "title", Label: "Judul Artikel", Type: "text", Value: a.Title},
			{Name: "category", Label: "Kategori", Type: "text", Value: a.Category},
			{Name: "excerpt", Label: "Ringkasan (tampil di kartu artikel)", Type: "textarea", Value: a.Excerpt},
			{Name: "content", Label: "Isi Artikel", Type: "textarea-tall", Value: listToLines(a.Content)},
			{Name: "image", Label: "Gambar Sampul", Type: "file"},
		},
	}
}

func handleArticleNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", articleForm(Article{}, "/admin/articles/new", "/admin/articles", "", "Tambah Artikel"))
}

func handleArticleCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveArticle(Article{
		Slug: r.FormValue("slug"), Title: r.FormValue("title"), Category: r.FormValue("category"),
		Excerpt: r.FormValue("excerpt"), Content: linesToList(r.FormValue("content")),
		ImagePath: saveUpload(r, "image", ""), SortOrder: 999,
	})
	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func handleArticleEdit(w http.ResponseWriter, r *http.Request) {
	a, ok := getArticle(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(a.ID, 10)
	render(w, "form", articleForm(a, "/admin/articles/"+id+"/edit", "/admin/articles", "/admin/articles/"+id+"/delete", "Ubah Artikel"))
}

func handleArticleUpdate(w http.ResponseWriter, r *http.Request) {
	a, ok := getArticle(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	a.Slug, a.Title, a.Category = r.FormValue("slug"), r.FormValue("title"), r.FormValue("category")
	a.Excerpt, a.Content = r.FormValue("excerpt"), linesToList(r.FormValue("content"))
	a.ImagePath = saveUpload(r, "image", a.ImagePath)
	saveArticle(a)
	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

func handleArticleDelete(w http.ResponseWriter, r *http.Request) {
	deleteArticle(idParam(r))
	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}

// ================= Facilities =================

func handleFacilitiesList(w http.ResponseWriter, r *http.Request) {
	items := listFacilities()
	var rows []ListRow
	for _, f := range items {
		rows = append(rows, ListRow{
			ImageURL:  f.ImagePath,
			Cells:     []string{f.Title, f.Icon},
			EditURL:   "/admin/facilities/" + strconv.FormatInt(f.ID, 10) + "/edit",
			DeleteURL: "/admin/facilities/" + strconv.FormatInt(f.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Fasilitas", NewURL: "/admin/facilities/new", Columns: []string{"Nama Ruang", "Ikon"}, Rows: rows, Count: len(items)})
}

func facilityForm(f Facility, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL, ImageURL: f.ImagePath,
		Fields: []FieldDef{
			{Name: "title", Label: "Nama Ruang/Fasilitas", Type: "text", Value: f.Title},
			{Name: "description", Label: "Deskripsi", Type: "textarea", Value: f.Description},
			{Name: "icon", Label: "Ikon", Type: "select", Value: f.Icon, Options: iconOptions},
			{Name: "image", Label: "Foto", Type: "file"},
		},
	}
}

func handleFacilityNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", facilityForm(Facility{Icon: "heart"}, "/admin/facilities/new", "/admin/facilities", "", "Tambah Fasilitas"))
}

func handleFacilityCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveFacility(Facility{
		Title: r.FormValue("title"), Description: r.FormValue("description"), Icon: r.FormValue("icon"),
		ImagePath: saveUpload(r, "image", ""), SortOrder: 999,
	})
	http.Redirect(w, r, "/admin/facilities", http.StatusSeeOther)
}

func handleFacilityEdit(w http.ResponseWriter, r *http.Request) {
	f, ok := getFacility(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(f.ID, 10)
	render(w, "form", facilityForm(f, "/admin/facilities/"+id+"/edit", "/admin/facilities", "/admin/facilities/"+id+"/delete", "Ubah Fasilitas"))
}

func handleFacilityUpdate(w http.ResponseWriter, r *http.Request) {
	f, ok := getFacility(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	f.Title, f.Description, f.Icon = r.FormValue("title"), r.FormValue("description"), r.FormValue("icon")
	f.ImagePath = saveUpload(r, "image", f.ImagePath)
	saveFacility(f)
	http.Redirect(w, r, "/admin/facilities", http.StatusSeeOther)
}

func handleFacilityDelete(w http.ResponseWriter, r *http.Request) {
	deleteFacility(idParam(r))
	http.Redirect(w, r, "/admin/facilities", http.StatusSeeOther)
}

// ================= Gallery =================

var galleryCategoryOptions = []string{"aktivitas", "fasilitas"}

func handleGalleryList(w http.ResponseWriter, r *http.Request) {
	items := listGallery()
	var rows []ListRow
	for _, g := range items {
		rows = append(rows, ListRow{
			ImageURL:  g.ImagePath,
			Cells:     []string{g.Caption, g.Category},
			EditURL:   "/admin/gallery/" + strconv.FormatInt(g.ID, 10) + "/edit",
			DeleteURL: "/admin/gallery/" + strconv.FormatInt(g.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Foto Galeri", NewURL: "/admin/gallery/new", GroupHint: "kategori: aktivitas / fasilitas", Columns: []string{"Keterangan", "Kategori"}, Rows: rows, Count: len(items)})
}

func galleryForm(g GalleryPhoto, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL, ImageURL: g.ImagePath,
		Fields: []FieldDef{
			{Name: "category", Label: "Kategori", Type: "select", Value: g.Category, Options: galleryCategoryOptions},
			{Name: "caption", Label: "Keterangan Foto", Type: "text", Value: g.Caption},
			{Name: "image", Label: "Foto", Type: "file"},
		},
	}
}

func handleGalleryNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", galleryForm(GalleryPhoto{Category: "aktivitas"}, "/admin/gallery/new", "/admin/gallery", "", "Tambah Foto Galeri"))
}

func handleGalleryCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveGalleryPhoto(GalleryPhoto{
		Category: r.FormValue("category"), Caption: r.FormValue("caption"),
		ImagePath: saveUpload(r, "image", ""), SortOrder: 999,
	})
	http.Redirect(w, r, "/admin/gallery", http.StatusSeeOther)
}

func handleGalleryEdit(w http.ResponseWriter, r *http.Request) {
	g, ok := getGalleryPhoto(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(g.ID, 10)
	render(w, "form", galleryForm(g, "/admin/gallery/"+id+"/edit", "/admin/gallery", "/admin/gallery/"+id+"/delete", "Ubah Foto Galeri"))
}

func handleGalleryUpdate(w http.ResponseWriter, r *http.Request) {
	g, ok := getGalleryPhoto(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	g.Category, g.Caption = r.FormValue("category"), r.FormValue("caption")
	g.ImagePath = saveUpload(r, "image", g.ImagePath)
	saveGalleryPhoto(g)
	http.Redirect(w, r, "/admin/gallery", http.StatusSeeOther)
}

func handleGalleryDelete(w http.ResponseWriter, r *http.Request) {
	deleteGalleryPhoto(idParam(r))
	http.Redirect(w, r, "/admin/gallery", http.StatusSeeOther)
}

// ================= Team =================

func handleTeamList(w http.ResponseWriter, r *http.Request) {
	items := listTeam()
	var rows []ListRow
	for _, t := range items {
		rows = append(rows, ListRow{
			ImageURL:  t.ImagePath,
			Cells:     []string{t.Name, t.Role},
			EditURL:   "/admin/team/" + strconv.FormatInt(t.ID, 10) + "/edit",
			DeleteURL: "/admin/team/" + strconv.FormatInt(t.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Tim Kami", NewURL: "/admin/team/new", Columns: []string{"Nama", "Peran"}, Rows: rows, Count: len(items)})
}

func teamForm(t TeamMember, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL, ImageURL: t.ImagePath,
		Fields: []FieldDef{
			{Name: "name", Label: "Nama Lengkap & Gelar", Type: "text", Value: t.Name},
			{Name: "role", Label: "Peran (mis. Terapis Wicara)", Type: "text", Value: t.Role},
			{Name: "image", Label: "Foto Profil", Type: "file"},
		},
	}
}

func handleTeamNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", teamForm(TeamMember{}, "/admin/team/new", "/admin/team", "", "Tambah Anggota Tim"))
}

func handleTeamCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveTeamMember(TeamMember{
		Name: r.FormValue("name"), Role: r.FormValue("role"),
		ImagePath: saveUpload(r, "image", ""), SortOrder: 999,
	})
	http.Redirect(w, r, "/admin/team", http.StatusSeeOther)
}

func handleTeamEdit(w http.ResponseWriter, r *http.Request) {
	t, ok := getTeamMember(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(t.ID, 10)
	render(w, "form", teamForm(t, "/admin/team/"+id+"/edit", "/admin/team", "/admin/team/"+id+"/delete", "Ubah Anggota Tim"))
}

func handleTeamUpdate(w http.ResponseWriter, r *http.Request) {
	t, ok := getTeamMember(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	t.Name, t.Role = r.FormValue("name"), r.FormValue("role")
	t.ImagePath = saveUpload(r, "image", t.ImagePath)
	saveTeamMember(t)
	http.Redirect(w, r, "/admin/team", http.StatusSeeOther)
}

func handleTeamDelete(w http.ResponseWriter, r *http.Request) {
	deleteTeamMember(idParam(r))
	http.Redirect(w, r, "/admin/team", http.StatusSeeOther)
}

// ================= Branches =================

func handleBranchesList(w http.ResponseWriter, r *http.Request) {
	items := listBranches()
	var rows []ListRow
	for _, b := range items {
		rows = append(rows, ListRow{
			Cells:     []string{b.Name, b.Address},
			EditURL:   "/admin/branches/" + strconv.FormatInt(b.ID, 10) + "/edit",
			DeleteURL: "/admin/branches/" + strconv.FormatInt(b.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Cabang", NewURL: "/admin/branches/new", Columns: []string{"Nama", "Alamat"}, Rows: rows, Count: len(items)})
}

func schedulesToText(items []ScheduleItem) string {
	var lines []string
	for _, it := range items {
		lines = append(lines, it.Days+" | "+it.Hours)
	}
	return listToLines(lines)
}

func textToSchedules(s string) []ScheduleItem {
	var out []ScheduleItem
	for _, line := range linesToList(s) {
		days, hours, ok := strings.Cut(line, "|")
		if ok {
			out = append(out, ScheduleItem{Days: strings.TrimSpace(days), Hours: strings.TrimSpace(hours)})
		}
	}
	return out
}

func branchForm(b Branch, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL,
		Fields: []FieldDef{
			{Name: "slug", Label: "Slug", Type: "text", Value: b.Slug},
			{Name: "name", Label: "Nama Cabang", Type: "text", Value: b.Name},
			{Name: "address", Label: "Alamat", Type: "textarea", Value: b.Address},
			{Name: "whatsapp", Label: "Nomor WhatsApp (format 62xxx)", Type: "text", Value: b.WhatsApp},
			{Name: "phone", Label: "Telepon (opsional)", Type: "text", Value: b.Phone},
			{Name: "maps_query", Label: "Google Maps Query (untuk embed peta)", Type: "text", Value: b.MapsQuery},
			{Name: "maps_url", Label: "Google Maps URL (link \"Buka di Maps\")", Type: "text", Value: b.MapsURL},
			{Name: "schedules", Label: "Jadwal Layanan", Type: "list", Value: schedulesToText(b.Schedules)},
		},
	}
}

func handleBranchNew(w http.ResponseWriter, r *http.Request) {
	f := branchForm(Branch{}, "/admin/branches/new", "/admin/branches", "", "Tambah Cabang")
	f.Fields[len(f.Fields)-1].Value = ""
	render(w, "form", f)
}

func handleBranchCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveBranch(Branch{
		Slug: r.FormValue("slug"), Name: r.FormValue("name"), Address: r.FormValue("address"),
		WhatsApp: r.FormValue("whatsapp"), Phone: r.FormValue("phone"),
		MapsQuery: r.FormValue("maps_query"), MapsURL: r.FormValue("maps_url"),
		Schedules: textToSchedules(r.FormValue("schedules")),
	})
	http.Redirect(w, r, "/admin/branches", http.StatusSeeOther)
}

func handleBranchEdit(w http.ResponseWriter, r *http.Request) {
	b, ok := getBranch(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(b.ID, 10)
	f := branchForm(b, "/admin/branches/"+id+"/edit", "/admin/branches", "/admin/branches/"+id+"/delete", "Ubah Cabang")
	f.Fields[len(f.Fields)-1].Label = "Jadwal Layanan (format: Hari | Jam, satu baris per jadwal)"
	render(w, "form", f)
}

func handleBranchUpdate(w http.ResponseWriter, r *http.Request) {
	b, ok := getBranch(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	b.Slug, b.Name, b.Address = r.FormValue("slug"), r.FormValue("name"), r.FormValue("address")
	b.WhatsApp, b.Phone = r.FormValue("whatsapp"), r.FormValue("phone")
	b.MapsQuery, b.MapsURL = r.FormValue("maps_query"), r.FormValue("maps_url")
	b.Schedules = textToSchedules(r.FormValue("schedules"))
	saveBranch(b)
	http.Redirect(w, r, "/admin/branches", http.StatusSeeOther)
}

func handleBranchDelete(w http.ResponseWriter, r *http.Request) {
	deleteBranch(idParam(r))
	http.Redirect(w, r, "/admin/branches", http.StatusSeeOther)
}

// ================= Testimonials =================

func handleTestimonialsList(w http.ResponseWriter, r *http.Request) {
	items := listTestimonials()
	var rows []ListRow
	for _, t := range items {
		rows = append(rows, ListRow{
			Cells:     []string{t.Name, t.Role},
			EditURL:   "/admin/testimonials/" + strconv.FormatInt(t.ID, 10) + "/edit",
			DeleteURL: "/admin/testimonials/" + strconv.FormatInt(t.ID, 10) + "/delete",
		})
	}
	render(w, "list", ListPage{Title: "Testimoni", NewURL: "/admin/testimonials/new", Columns: []string{"Nama", "Peran"}, Rows: rows, Count: len(items)})
}

func testimonialForm(t Testimonial, action, backURL, deleteURL, title string) FormPage {
	return FormPage{
		Title: title, Action: action, BackURL: backURL, DeleteURL: deleteURL,
		Fields: []FieldDef{
			{Name: "name", Label: "Nama Orang Tua", Type: "text", Value: t.Name},
			{Name: "role", Label: "Keterangan (mis. Orang tua pasien terapi wicara)", Type: "text", Value: t.Role},
			{Name: "quote", Label: "Isi Testimoni", Type: "textarea", Value: t.Quote},
		},
	}
}

func handleTestimonialNew(w http.ResponseWriter, r *http.Request) {
	render(w, "form", testimonialForm(Testimonial{}, "/admin/testimonials/new", "/admin/testimonials", "", "Tambah Testimoni"))
}

func handleTestimonialCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	saveTestimonial(Testimonial{Name: r.FormValue("name"), Role: r.FormValue("role"), Quote: r.FormValue("quote"), SortOrder: 999})
	http.Redirect(w, r, "/admin/testimonials", http.StatusSeeOther)
}

func handleTestimonialEdit(w http.ResponseWriter, r *http.Request) {
	t, ok := getTestimonial(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	id := strconv.FormatInt(t.ID, 10)
	render(w, "form", testimonialForm(t, "/admin/testimonials/"+id+"/edit", "/admin/testimonials", "/admin/testimonials/"+id+"/delete", "Ubah Testimoni"))
}

func handleTestimonialUpdate(w http.ResponseWriter, r *http.Request) {
	t, ok := getTestimonial(idParam(r))
	if !ok {
		http.NotFound(w, r)
		return
	}
	r.ParseMultipartForm(20 << 20)
	t.Name, t.Role, t.Quote = r.FormValue("name"), r.FormValue("role"), r.FormValue("quote")
	saveTestimonial(t)
	http.Redirect(w, r, "/admin/testimonials", http.StatusSeeOther)
}

func handleTestimonialDelete(w http.ResponseWriter, r *http.Request) {
	deleteTestimonial(idParam(r))
	http.Redirect(w, r, "/admin/testimonials", http.StatusSeeOther)
}

// ================= Settings =================

func handleSettingsPage(w http.ResponseWriter, r *http.Request) {
	var fields []FieldDef
	for _, k := range settingsKeys {
		typ := "text"
		if k.Key == "misi" || k.Key == "hero_subtitle" || k.Key == "about_text" || k.Key == "visi" || k.Key == "tim_intro" {
			typ = "textarea"
		}
		fields = append(fields, FieldDef{Name: k.Key, Label: k.Label, Type: typ, Value: getSetting(k.Key, "")})
	}
	render(w, "settings", SettingsPage{Fields: fields})
}

func handleSettingsSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for _, k := range settingsKeys {
		setSetting(k.Key, r.FormValue(k.Key))
	}
	var fields []FieldDef
	for _, k := range settingsKeys {
		typ := "text"
		if k.Key == "misi" || k.Key == "hero_subtitle" || k.Key == "about_text" || k.Key == "visi" || k.Key == "tim_intro" {
			typ = "textarea"
		}
		fields = append(fields, FieldDef{Name: k.Key, Label: k.Label, Type: typ, Value: getSetting(k.Key, "")})
	}
	render(w, "settings", SettingsPage{Flash: "Pengaturan berhasil disimpan.", Fields: fields})
}
