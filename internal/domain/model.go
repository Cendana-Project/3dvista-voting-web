package domain

import "time"

// Innovation represents an innovation entry
type Innovation struct {
	ID                string    `json:"id"`
	GroupSlug         string    `json:"group_slug"`
	Slug              string    `json:"slug"`
	Name              string    `json:"name"`
	Division          *string   `json:"division,omitempty"`
	EntityName        *string   `json:"entity_name,omitempty"`
	PIC               *string   `json:"pic,omitempty"`
	Description       *string   `json:"description,omitempty"`
	LogoInnovationURL *string   `json:"logo_innovation_url,omitempty"`
	LogoEntityURL     *string   `json:"logo_entity_url,omitempty"`
	VideoURL          *string   `json:"video_url,omitempty"`
	SlideURL          *string   `json:"slide_url,omitempty"`
	IgURL             *string   `json:"ig_url,omitempty"`
	YtURL             *string   `json:"yt_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Vote represents a vote record
type Vote struct {
	ID           int64     `json:"id"`
	InnovationID string    `json:"innovation_id"`
	VoterIPHash  []byte    `json:"-"`
	UserAgent    string    `json:"user_agent,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// VoteRequest represents a vote submission request
type VoteRequest struct {
	GroupSlug string
	Slug      string
	ClientIP  string
	UserAgent string
}

// VoteResponse represents the result of a vote operation
type VoteResponse struct {
	Success      bool   `json:"success"`
	AlreadyVoted bool   `json:"already_voted,omitempty"`
	VoteCount    int64  `json:"vote_count"`
	Message      string `json:"message,omitempty"`
}
