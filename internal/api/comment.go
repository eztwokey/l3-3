package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/eztwokey/l3-3/internal/logic"
	"github.com/eztwokey/l3-3/internal/models"
	"github.com/eztwokey/l3-3/internal/storage"
)

func (a *Api) createComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Warn("comment: bind error", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	comment, err := a.logic.CreateComment(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (a *Api) listComments(c *gin.Context) {
	params := models.ListParams{
		Page:     1,
		PageSize: 20,
		SortBy:   "created_at",
		SortDir:  "desc",
	}

	if pid := c.Query("parent_id"); pid != "" {
		id, err := strconv.Atoi(pid)
		if err == nil {
			params.ParentID = &id
		}
	}

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			params.Page = v
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			params.PageSize = v
		}
	}

	if sort := c.Query("sort"); sort != "" {
		switch sort {
		case "created_at_asc":
			params.SortBy = "created_at"
			params.SortDir = "asc"
		case "created_at_desc":
			params.SortBy = "created_at"
			params.SortDir = "desc"
		}
	}

	params.Search = c.Query("search")

	list, err := a.logic.ListComments(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (a *Api) getTree(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tree, err := a.logic.GetTree(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, tree)
}

func (a *Api) deleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := a.logic.DeleteComment(c.Request.Context(), id); err != nil {
		if errors.Is(err, logic.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
