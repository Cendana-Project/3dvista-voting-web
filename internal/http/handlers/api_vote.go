package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"voteweb/internal/domain"
)

type VoteHandler struct {
	service domain.VoteService
	logger  *slog.Logger
}

func NewVoteHandler(service domain.VoteService, logger *slog.Logger) *VoteHandler {
	return &VoteHandler{
		service: service,
		logger:  logger,
	}
}

func (h *VoteHandler) SubmitVote(c *gin.Context) {
	// Voting system is closed
	c.JSON(http.StatusForbidden, gin.H{
		"error":   "voting_closed",
		"message": "Sistem voting telah ditutup. Terima kasih atas partisipasi Anda.",
	})
	return

	// Code below is disabled while voting is closed
	// Uncomment to re-enable voting functionality
	/*
	groupSlug := c.Param("group")
	slug := c.Param("slug")

	// Get client IP from context (set by ProxiedIP middleware)
	clientIP, exists := c.Get("client_ip")
	if !exists {
		clientIP = c.ClientIP()
	}

	req := domain.VoteRequest{
		GroupSlug: groupSlug,
		Slug:      slug,
		ClientIP:  clientIP.(string),
		UserAgent: c.GetHeader("User-Agent"),
	}

	result, err := h.service.SubmitVote(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInnovationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Innovation not found",
			})
			return
		}

		h.logger.ErrorContext(c.Request.Context(), "failed to submit vote",
			"group_slug", groupSlug,
			"slug", slug,
			"error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit vote",
		})
		return
	}

	// If already voted, return 409 Conflict
	if result.AlreadyVoted {
		c.JSON(http.StatusConflict, gin.H{
			"error":      "already_voted",
			"message":    result.Message,
			"vote_count": result.VoteCount,
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    result.Message,
		"vote_count": result.VoteCount,
	})
	*/
}



