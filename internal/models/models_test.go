package models

import (
	"encoding/json"
	"testing"
)

func TestCreateCommentRequest_JSON(t *testing.T) {
	input := `{"parent_id":1,"author":"Кирилл","text":"Привет"}`

	var req CreateCommentRequest
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if req.Author != "Кирилл" {
		t.Errorf("expected author 'Кирилл', got %q", req.Author)
	}
	if req.ParentID == nil || *req.ParentID != 1 {
		t.Errorf("expected parent_id 1, got %v", req.ParentID)
	}
}

func TestCreateCommentRequest_NullParent(t *testing.T) {
	input := `{"author":"Кирилл","text":"Корневой"}`

	var req CreateCommentRequest
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if req.ParentID != nil {
		t.Errorf("expected nil parent_id, got %v", req.ParentID)
	}
}

func TestCommentList_JSON(t *testing.T) {
	list := CommentList{
		Comments:   []Comment{},
		Total:      0,
		Page:       1,
		PageSize:   20,
		TotalPages: 0,
	}

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if result["total"].(float64) != 0 {
		t.Errorf("expected total 0, got %v", result["total"])
	}
}
