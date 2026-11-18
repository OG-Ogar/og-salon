package models

import "time"

// SocialPost is a shared customer/barber feed entry.
type SocialPost struct {
	ID              int             `json:"id"`
	AuthorRole      string          `json:"author_role"`
	AuthorID        int             `json:"author_id"`
	AuthorName      string          `json:"author_name"`
	AuthorImage     string          `json:"author_image"`
	AuthorType      string          `json:"author_type"`
	Text            string          `json:"text"`
	BackgroundColor string          `json:"background_color"`
	TextStyle       string          `json:"text_style"`
	Image           string          `json:"image"`
	Video           string          `json:"video"`
	Music           string          `json:"music"`
	CreatedAt       time.Time       `json:"created_at"`
	LikeCount       int             `json:"like_count"`
	LikedByMe       bool            `json:"liked_by_me"`
	Comments        []SocialComment `json:"comments"`
}

// AuthorLabel is the role shown under a post's author name.
func (p SocialPost) AuthorLabel() string {
	if p.AuthorRole != SenderBarber {
		return "Client"
	}
	if NormalizeProviderType(p.AuthorType) == ProviderHairstylist {
		return "Stylist"
	}
	return "Barber"
}

// SocialComment is a short, plain-text reply on a feed post.
type SocialComment struct {
	ID          int       `json:"id"`
	PostID      int       `json:"post_id"`
	AuthorRole  string    `json:"author_role"`
	AuthorID    int       `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	AuthorImage string    `json:"author_image"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"created_at"`
}

// HairstyleStatus is a barber's available hairstyle shown in the status row.
type HairstyleStatus struct {
	HairstyleID int
	BarberID    int
	BarberName  string
	BarberImage string
	ShopName    string
	Name        string
	Image       string
	Video       string
	Price       float64
}
