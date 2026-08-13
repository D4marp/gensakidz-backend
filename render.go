package main

import (
	"embed"
	"html/template"
	"net/http"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

type Card struct {
	URL, Label string
	Count      int
}

type DashboardPage struct {
	Cards []Card
}

type ListRow struct {
	ImageURL           string
	Cells              []string
	EditURL, DeleteURL string
}

type ListPage struct {
	Title, NewURL, GroupHint, Flash string
	Columns                         []string
	Rows                            []ListRow
	Count                           int
}

type FieldDef struct {
	Name, Label, Type, Value string
	Options                  []string
}

type FormPage struct {
	Title, Subtitle, Action, BackURL, ImageURL, DeleteURL string
	Fields                                                []FieldDef
}

type SettingsPage struct {
	Flash  string
	Fields []FieldDef
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- small text <-> list helpers used by handlers ---

func linesToList(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func listToLines(list []string) string {
	return strings.Join(list, "\n")
}

func faqToText(items []FAQItem) string {
	var b strings.Builder
	for i, it := range items {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString("Q: " + it.Q + "\nA: " + it.A)
	}
	return b.String()
}

func textToFAQ(s string) []FAQItem {
	var out []FAQItem
	blocks := strings.Split(strings.TrimSpace(s), "\n\n")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.SplitN(block, "\n", 2)
		if len(lines) < 2 {
			continue
		}
		q := strings.TrimSpace(strings.TrimPrefix(lines[0], "Q:"))
		a := strings.TrimSpace(strings.TrimPrefix(lines[1], "A:"))
		out = append(out, FAQItem{Q: q, A: a})
	}
	return out
}
