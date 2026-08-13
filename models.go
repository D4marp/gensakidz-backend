package main

// FAQItem is a question/answer pair used inside a Service's extra FAQ list.
type FAQItem struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type Service struct {
	ID            int64
	Slug          string
	Title         string
	Icon          string
	Short         string
	Detail        []string
	ForWho        string
	Signs         []string
	Goal          string
	Process       []string
	Duration      string
	Professionals string
	WhatToBring   []string
	ExtraFAQ      []FAQItem
	ImagePath     string
	SortOrder     int
}

type Job struct {
	ID           int64
	Slug         string
	Title        string
	Branch       string
	Type         string // Full-time | Part-time
	Status       string // Dibuka | Ditutup
	Description  string
	Requirements []string
	SortOrder    int
}

type Article struct {
	ID        int64
	Slug      string
	Title     string
	Category  string
	Excerpt   string
	Content   []string
	ImagePath string
	SortOrder int
}

type Facility struct {
	ID          int64
	Title       string
	Description string
	Icon        string
	ImagePath   string
	SortOrder   int
}

type GalleryPhoto struct {
	ID        int64
	Category  string // aktivitas | fasilitas
	Caption   string
	ImagePath string
	SortOrder int
}

type TeamMember struct {
	ID        int64
	Name      string
	Role      string
	ImagePath string
	SortOrder int
}

type ScheduleItem struct {
	Days  string `json:"days"`
	Hours string `json:"hours"`
}

type Branch struct {
	ID        int64
	Slug      string
	Name      string
	Address   string
	WhatsApp  string
	Phone     string
	MapsQuery string
	MapsURL   string
	Schedules []ScheduleItem
}

type Testimonial struct {
	ID        int64
	Name      string
	Role      string
	Quote     string
	SortOrder int
}

type AdminUser struct {
	ID           int64
	Email        string
	PasswordHash string
}
