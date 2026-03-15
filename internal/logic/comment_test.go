package logic

import (
	"testing"

	"github.com/eztwokey/l3-3/internal/models"
)

func TestCreateComment_EmptyAuthor(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateComment(t.Context(), models.CreateCommentRequest{
		Author: "",
		Text:   "текст",
	})
	if err == nil {
		t.Fatal("expected error for empty author, got nil")
	}
}

func TestCreateComment_EmptyText(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateComment(t.Context(), models.CreateCommentRequest{
		Author: "Кирилл",
		Text:   "",
	})
	if err == nil {
		t.Fatal("expected error for empty text, got nil")
	}
}

func TestCreateComment_WhitespaceOnly(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateComment(t.Context(), models.CreateCommentRequest{
		Author: "   ",
		Text:   "   ",
	})
	if err == nil {
		t.Fatal("expected error for whitespace-only fields, got nil")
	}
}

func TestBuildTree_Single(t *testing.T) {
	flat := []models.Comment{
		{ID: 1, Author: "A", Text: "root"},
	}

	tree := buildTree(flat, 1)
	if tree.ID != 1 {
		t.Errorf("expected root id 1, got %d", tree.ID)
	}
	if len(tree.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(tree.Children))
	}
}

func TestBuildTree_WithChildren(t *testing.T) {
	pid1 := 1
	flat := []models.Comment{
		{ID: 1, Author: "A", Text: "root"},
		{ID: 2, ParentID: &pid1, Author: "B", Text: "child"},
	}

	tree := buildTree(flat, 1)
	if len(tree.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(tree.Children))
	}
	if tree.Children[0].ID != 2 {
		t.Errorf("expected child id 2, got %d", tree.Children[0].ID)
	}
}

func TestTotalPages(t *testing.T) {
	cases := []struct {
		total, pageSize, want int
	}{
		{0, 20, 0},
		{1, 20, 1},
		{20, 20, 1},
		{21, 20, 2},
		{100, 20, 5},
	}

	for _, tc := range cases {
		got := totalPages(tc.total, tc.pageSize)
		if got != tc.want {
			t.Errorf("totalPages(%d, %d) = %d, want %d", tc.total, tc.pageSize, got, tc.want)
		}
	}
}
