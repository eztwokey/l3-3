package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	ParentID  *int      `json:"parent_id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Children  []Comment `json:"children"`
}

type CreateCommentRequest struct {
	ParentID *int   `json:"parent_id"`
	Author   string `json:"author"`
	Text     string `json:"text"`
}

type ListParams struct {
	ParentID *int
	Page     int
	PageSize int
	SortBy   string
	SortDir  string
	Search   string
}

type CommentList struct {
	Comments   []Comment `json:"comments"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}
