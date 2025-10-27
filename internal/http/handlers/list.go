package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"voteweb/internal/domain"
)

type ListHandler struct {
	service domain.VoteService
	logger  *slog.Logger
}

func NewListHandler(service domain.VoteService, logger *slog.Logger) *ListHandler {
	return &ListHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ListHandler) ShowList(c *gin.Context) {
	innovations, err := h.service.ListInnovations(c.Request.Context())
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to list innovations", "error", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl.html", gin.H{
			"Title":   "Error",
			"Message": "An error occurred while loading innovations.",
		})
		return
	}

	// Group innovations by group_slug
	grouped := make(map[string][]*domain.Innovation)
	for _, innovation := range innovations {
		if grouped[innovation.GroupSlug] == nil {
			grouped[innovation.GroupSlug] = []*domain.Innovation{}
		}
		grouped[innovation.GroupSlug] = append(grouped[innovation.GroupSlug], innovation)
	}

	c.HTML(http.StatusOK, "list.tmpl.html", gin.H{
		"Title":       "Innovation Voting System",
		"Innovations": innovations,
		"Grouped":     grouped,
	})
}
