package main

import "database/sql"

// ---------- Services ----------

func listServices() []Service {
	rows, _ := db.Query(`SELECT id,slug,title,icon,short,detail,for_who,signs,goal,process,duration,professionals,what_to_bring,extra_faq,image_path,sort_order FROM services ORDER BY sort_order, id`)
	defer rows.Close()
	var out []Service
	for rows.Next() {
		var s Service
		var detail, signs, process, wtb, faq string
		rows.Scan(&s.ID, &s.Slug, &s.Title, &s.Icon, &s.Short, &detail, &s.ForWho, &signs, &s.Goal, &process, &s.Duration, &s.Professionals, &wtb, &faq, &s.ImagePath, &s.SortOrder)
		s.Detail, s.Signs, s.Process, s.WhatToBring = fromJSONStrings(detail), fromJSONStrings(signs), fromJSONStrings(process), fromJSONStrings(wtb)
		s.ExtraFAQ = fromJSONFAQ(faq)
		out = append(out, s)
	}
	return out
}

func getService(id int64) (Service, bool) {
	var s Service
	var detail, signs, process, wtb, faq string
	err := db.QueryRow(`SELECT id,slug,title,icon,short,detail,for_who,signs,goal,process,duration,professionals,what_to_bring,extra_faq,image_path,sort_order FROM services WHERE id=?`, id).
		Scan(&s.ID, &s.Slug, &s.Title, &s.Icon, &s.Short, &detail, &s.ForWho, &signs, &s.Goal, &process, &s.Duration, &s.Professionals, &wtb, &faq, &s.ImagePath, &s.SortOrder)
	if err != nil {
		return s, false
	}
	s.Detail, s.Signs, s.Process, s.WhatToBring = fromJSONStrings(detail), fromJSONStrings(signs), fromJSONStrings(process), fromJSONStrings(wtb)
	s.ExtraFAQ = fromJSONFAQ(faq)
	return s, true
}

func saveService(s Service) {
	if s.ID == 0 {
		db.Exec(`INSERT INTO services (slug,title,icon,short,detail,for_who,signs,goal,process,duration,professionals,what_to_bring,extra_faq,image_path,sort_order) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			s.Slug, s.Title, s.Icon, s.Short, toJSON(s.Detail), s.ForWho, toJSON(s.Signs), s.Goal, toJSON(s.Process), s.Duration, s.Professionals, toJSON(s.WhatToBring), toJSON(s.ExtraFAQ), s.ImagePath, s.SortOrder)
		return
	}
	db.Exec(`UPDATE services SET slug=?,title=?,icon=?,short=?,detail=?,for_who=?,signs=?,goal=?,process=?,duration=?,professionals=?,what_to_bring=?,extra_faq=?,image_path=?,sort_order=? WHERE id=?`,
		s.Slug, s.Title, s.Icon, s.Short, toJSON(s.Detail), s.ForWho, toJSON(s.Signs), s.Goal, toJSON(s.Process), s.Duration, s.Professionals, toJSON(s.WhatToBring), toJSON(s.ExtraFAQ), s.ImagePath, s.SortOrder, s.ID)
}

func deleteService(id int64) { db.Exec(`DELETE FROM services WHERE id=?`, id) }

// ---------- Jobs ----------

func listJobs() []Job {
	rows, _ := db.Query(`SELECT id,slug,title,branch,type,status,description,requirements,sort_order FROM jobs ORDER BY sort_order, id`)
	defer rows.Close()
	var out []Job
	for rows.Next() {
		var j Job
		var req string
		rows.Scan(&j.ID, &j.Slug, &j.Title, &j.Branch, &j.Type, &j.Status, &j.Description, &req, &j.SortOrder)
		j.Requirements = fromJSONStrings(req)
		out = append(out, j)
	}
	return out
}

func getJob(id int64) (Job, bool) {
	var j Job
	var req string
	err := db.QueryRow(`SELECT id,slug,title,branch,type,status,description,requirements,sort_order FROM jobs WHERE id=?`, id).
		Scan(&j.ID, &j.Slug, &j.Title, &j.Branch, &j.Type, &j.Status, &j.Description, &req, &j.SortOrder)
	if err != nil {
		return j, false
	}
	j.Requirements = fromJSONStrings(req)
	return j, true
}

func saveJob(j Job) {
	if j.ID == 0 {
		db.Exec(`INSERT INTO jobs (slug,title,branch,type,status,description,requirements,sort_order) VALUES (?,?,?,?,?,?,?,?)`,
			j.Slug, j.Title, j.Branch, j.Type, j.Status, j.Description, toJSON(j.Requirements), j.SortOrder)
		return
	}
	db.Exec(`UPDATE jobs SET slug=?,title=?,branch=?,type=?,status=?,description=?,requirements=?,sort_order=? WHERE id=?`,
		j.Slug, j.Title, j.Branch, j.Type, j.Status, j.Description, toJSON(j.Requirements), j.SortOrder, j.ID)
}

func deleteJob(id int64) { db.Exec(`DELETE FROM jobs WHERE id=?`, id) }

// ---------- Articles ----------

func listArticles() []Article {
	rows, _ := db.Query(`SELECT id,slug,title,category,excerpt,content,image_path,sort_order FROM articles ORDER BY sort_order, id`)
	defer rows.Close()
	var out []Article
	for rows.Next() {
		var a Article
		var content string
		rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Category, &a.Excerpt, &content, &a.ImagePath, &a.SortOrder)
		a.Content = fromJSONStrings(content)
		out = append(out, a)
	}
	return out
}

func getArticle(id int64) (Article, bool) {
	var a Article
	var content string
	err := db.QueryRow(`SELECT id,slug,title,category,excerpt,content,image_path,sort_order FROM articles WHERE id=?`, id).
		Scan(&a.ID, &a.Slug, &a.Title, &a.Category, &a.Excerpt, &content, &a.ImagePath, &a.SortOrder)
	if err != nil {
		return a, false
	}
	a.Content = fromJSONStrings(content)
	return a, true
}

func saveArticle(a Article) {
	if a.ID == 0 {
		db.Exec(`INSERT INTO articles (slug,title,category,excerpt,content,image_path,sort_order) VALUES (?,?,?,?,?,?,?)`,
			a.Slug, a.Title, a.Category, a.Excerpt, toJSON(a.Content), a.ImagePath, a.SortOrder)
		return
	}
	db.Exec(`UPDATE articles SET slug=?,title=?,category=?,excerpt=?,content=?,image_path=?,sort_order=? WHERE id=?`,
		a.Slug, a.Title, a.Category, a.Excerpt, toJSON(a.Content), a.ImagePath, a.SortOrder, a.ID)
}

func deleteArticle(id int64) { db.Exec(`DELETE FROM articles WHERE id=?`, id) }

// ---------- Facilities ----------

func listFacilities() []Facility {
	rows, _ := db.Query(`SELECT id,title,description,icon,image_path,sort_order FROM facilities ORDER BY sort_order, id`)
	defer rows.Close()
	var out []Facility
	for rows.Next() {
		var f Facility
		rows.Scan(&f.ID, &f.Title, &f.Description, &f.Icon, &f.ImagePath, &f.SortOrder)
		out = append(out, f)
	}
	return out
}

func getFacility(id int64) (Facility, bool) {
	var f Facility
	err := db.QueryRow(`SELECT id,title,description,icon,image_path,sort_order FROM facilities WHERE id=?`, id).
		Scan(&f.ID, &f.Title, &f.Description, &f.Icon, &f.ImagePath, &f.SortOrder)
	return f, err == nil
}

func saveFacility(f Facility) {
	if f.ID == 0 {
		db.Exec(`INSERT INTO facilities (title,description,icon,image_path,sort_order) VALUES (?,?,?,?,?)`,
			f.Title, f.Description, f.Icon, f.ImagePath, f.SortOrder)
		return
	}
	db.Exec(`UPDATE facilities SET title=?,description=?,icon=?,image_path=?,sort_order=? WHERE id=?`,
		f.Title, f.Description, f.Icon, f.ImagePath, f.SortOrder, f.ID)
}

func deleteFacility(id int64) { db.Exec(`DELETE FROM facilities WHERE id=?`, id) }

// ---------- Gallery ----------

func listGallery() []GalleryPhoto {
	rows, _ := db.Query(`SELECT id,category,caption,image_path,sort_order FROM gallery_photos ORDER BY category, sort_order, id`)
	defer rows.Close()
	var out []GalleryPhoto
	for rows.Next() {
		var g GalleryPhoto
		rows.Scan(&g.ID, &g.Category, &g.Caption, &g.ImagePath, &g.SortOrder)
		out = append(out, g)
	}
	return out
}

func getGalleryPhoto(id int64) (GalleryPhoto, bool) {
	var g GalleryPhoto
	err := db.QueryRow(`SELECT id,category,caption,image_path,sort_order FROM gallery_photos WHERE id=?`, id).
		Scan(&g.ID, &g.Category, &g.Caption, &g.ImagePath, &g.SortOrder)
	return g, err == nil
}

func saveGalleryPhoto(g GalleryPhoto) {
	if g.ID == 0 {
		db.Exec(`INSERT INTO gallery_photos (category,caption,image_path,sort_order) VALUES (?,?,?,?)`,
			g.Category, g.Caption, g.ImagePath, g.SortOrder)
		return
	}
	db.Exec(`UPDATE gallery_photos SET category=?,caption=?,image_path=?,sort_order=? WHERE id=?`,
		g.Category, g.Caption, g.ImagePath, g.SortOrder, g.ID)
}

func deleteGalleryPhoto(id int64) { db.Exec(`DELETE FROM gallery_photos WHERE id=?`, id) }

// ---------- Team ----------

func listTeam() []TeamMember {
	rows, _ := db.Query(`SELECT id,name,role,image_path,sort_order FROM team_members ORDER BY sort_order, id`)
	defer rows.Close()
	var out []TeamMember
	for rows.Next() {
		var t TeamMember
		rows.Scan(&t.ID, &t.Name, &t.Role, &t.ImagePath, &t.SortOrder)
		out = append(out, t)
	}
	return out
}

func getTeamMember(id int64) (TeamMember, bool) {
	var t TeamMember
	err := db.QueryRow(`SELECT id,name,role,image_path,sort_order FROM team_members WHERE id=?`, id).
		Scan(&t.ID, &t.Name, &t.Role, &t.ImagePath, &t.SortOrder)
	return t, err == nil
}

func saveTeamMember(t TeamMember) {
	if t.ID == 0 {
		db.Exec(`INSERT INTO team_members (name,role,image_path,sort_order) VALUES (?,?,?,?)`,
			t.Name, t.Role, t.ImagePath, t.SortOrder)
		return
	}
	db.Exec(`UPDATE team_members SET name=?,role=?,image_path=?,sort_order=? WHERE id=?`,
		t.Name, t.Role, t.ImagePath, t.SortOrder, t.ID)
}

func deleteTeamMember(id int64) { db.Exec(`DELETE FROM team_members WHERE id=?`, id) }

// ---------- Branches ----------

func listBranches() []Branch {
	rows, _ := db.Query(`SELECT id,slug,name,address,whatsapp,phone,maps_query,maps_url,schedules FROM branches ORDER BY id`)
	defer rows.Close()
	var out []Branch
	for rows.Next() {
		var b Branch
		var sched string
		rows.Scan(&b.ID, &b.Slug, &b.Name, &b.Address, &b.WhatsApp, &b.Phone, &b.MapsQuery, &b.MapsURL, &sched)
		b.Schedules = fromJSONSchedules(sched)
		out = append(out, b)
	}
	return out
}

func getBranch(id int64) (Branch, bool) {
	var b Branch
	var sched string
	err := db.QueryRow(`SELECT id,slug,name,address,whatsapp,phone,maps_query,maps_url,schedules FROM branches WHERE id=?`, id).
		Scan(&b.ID, &b.Slug, &b.Name, &b.Address, &b.WhatsApp, &b.Phone, &b.MapsQuery, &b.MapsURL, &sched)
	if err != nil {
		return b, false
	}
	b.Schedules = fromJSONSchedules(sched)
	return b, true
}

func saveBranch(b Branch) {
	if b.ID == 0 {
		db.Exec(`INSERT INTO branches (slug,name,address,whatsapp,phone,maps_query,maps_url,schedules) VALUES (?,?,?,?,?,?,?,?)`,
			b.Slug, b.Name, b.Address, b.WhatsApp, b.Phone, b.MapsQuery, b.MapsURL, toJSON(b.Schedules))
		return
	}
	db.Exec(`UPDATE branches SET slug=?,name=?,address=?,whatsapp=?,phone=?,maps_query=?,maps_url=?,schedules=? WHERE id=?`,
		b.Slug, b.Name, b.Address, b.WhatsApp, b.Phone, b.MapsQuery, b.MapsURL, toJSON(b.Schedules), b.ID)
}

func deleteBranch(id int64) { db.Exec(`DELETE FROM branches WHERE id=?`, id) }

// ---------- Testimonials ----------

func listTestimonials() []Testimonial {
	rows, _ := db.Query(`SELECT id,name,role,quote,sort_order FROM testimonials ORDER BY sort_order, id`)
	defer rows.Close()
	var out []Testimonial
	for rows.Next() {
		var t Testimonial
		rows.Scan(&t.ID, &t.Name, &t.Role, &t.Quote, &t.SortOrder)
		out = append(out, t)
	}
	return out
}

func getTestimonial(id int64) (Testimonial, bool) {
	var t Testimonial
	err := db.QueryRow(`SELECT id,name,role,quote,sort_order FROM testimonials WHERE id=?`, id).
		Scan(&t.ID, &t.Name, &t.Role, &t.Quote, &t.SortOrder)
	return t, err == nil
}

func saveTestimonial(t Testimonial) {
	if t.ID == 0 {
		db.Exec(`INSERT INTO testimonials (name,role,quote,sort_order) VALUES (?,?,?,?)`,
			t.Name, t.Role, t.Quote, t.SortOrder)
		return
	}
	db.Exec(`UPDATE testimonials SET name=?,role=?,quote=?,sort_order=? WHERE id=?`,
		t.Name, t.Role, t.Quote, t.SortOrder, t.ID)
}

func deleteTestimonial(id int64) { db.Exec(`DELETE FROM testimonials WHERE id=?`, id) }

// ---------- Site settings ----------

func getSetting(key, def string) string {
	var v string
	err := db.QueryRow(`SELECT value FROM site_settings WHERE setting_key=?`, key).Scan(&v)
	if err == sql.ErrNoRows || err != nil {
		return def
	}
	return v
}

func setSetting(key, value string) {
	db.Exec(`INSERT INTO site_settings (setting_key,value) VALUES (?,?) ON DUPLICATE KEY UPDATE value=VALUES(value)`, key, value)
}

var settingsKeys = []struct{ Key, Label string }{
	{"hero_badge", "Hero — Label kecil (mis. Lamongan · Sejak 2020)"},
	{"hero_title", "Hero — Judul utama"},
	{"hero_subtitle", "Hero — Paragraf deskripsi"},
	{"about_text", "Tentang — Paragraf ringkas (badge home)"},
	{"visi", "Visi"},
	{"misi", "Misi (satu poin per baris)"},
	{"tim_intro", "Tim — Paragraf pengantar"},
	{"ig_url", "Link Instagram"},
	{"fb_url", "Link Facebook"},
	{"tiktok_url", "Link TikTok"},
	{"linkedin_url", "Link LinkedIn"},
	{"apply_email", "Email lamaran karir"},
}
