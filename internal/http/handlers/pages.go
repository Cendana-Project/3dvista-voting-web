package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"voteweb/internal/domain"
	"voteweb/internal/http/middleware"
)

type PageHandler struct {
	service  domain.VoteService
	logger   *slog.Logger
	autoVote bool
}

func NewPageHandler(service domain.VoteService, logger *slog.Logger, autoVote bool) *PageHandler {
	return &PageHandler{
		service:  service,
		logger:   logger,
		autoVote: autoVote,
	}
}

func (h *PageHandler) ShowInnovation(c *gin.Context) {
	groupSlug := c.Param("group")
	slug := c.Param("slug")

	// Get innovation
	innovation, err := h.service.GetInnovation(c.Request.Context(), groupSlug, slug)
	if err != nil {
		if err == domain.ErrInnovationNotFound {
			c.HTML(http.StatusNotFound, "error.tmpl.html", gin.H{
				"Title":   "Innovation Not Found",
				"Message": "The innovation you're looking for does not exist.",
			})
			return
		}
		h.logger.ErrorContext(c.Request.Context(), "failed to get innovation",
			"group_slug", groupSlug,
			"slug", slug,
			"error", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl.html", gin.H{
			"Title":   "Error",
			"Message": "An error occurred while loading the page.",
		})
		return
	}

	// Get current vote count
	voteCount, err := h.service.GetVoteCount(c.Request.Context(), innovation.ID)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to get vote count",
			"innovation_id", innovation.ID,
			"error", err)
		voteCount = 0
	}

	// Check if user has already voted
	clientIP := c.GetString("client_ip")
	hasVoted := false
	if clientIP != "" {
		voted, err := h.service.CheckHasVoted(c.Request.Context(), innovation.ID, clientIP)
		if err != nil {
			h.logger.ErrorContext(c.Request.Context(), "failed to check vote status",
				"innovation_id", innovation.ID,
				"error", err)
		} else {
			hasVoted = voted
		}
	}

	// Auto-vote on first view if enabled, user hasn't voted, and likely not a bot
	if h.autoVote && !hasVoted && clientIP != "" && !isBot(c.GetHeader("User-Agent")) && c.Query("preview") != "1" {
		if _, err := h.service.SubmitVote(c.Request.Context(), domain.VoteRequest{
			GroupSlug: groupSlug,
			Slug:      slug,
			ClientIP:  clientIP,
			UserAgent: c.GetHeader("User-Agent"),
		}); err != nil {
			h.logger.ErrorContext(c.Request.Context(), "auto-vote failed",
				"group_slug", groupSlug,
				"slug", slug,
				"error", err)
		} else {
			// refresh state
			if cnt, err := h.service.GetVoteCount(c.Request.Context(), innovation.ID); err == nil {
				voteCount = cnt
			}
			hasVoted = true
		}
	}

	// Get CSRF token for the page
	csrfToken := middleware.GetCSRFToken(c)

	c.HTML(http.StatusOK, "innovation.tmpl.html", gin.H{
		"Innovation": innovation,
		"VoteCount":  voteCount,
		"CSRFToken":  csrfToken,
		"HasVoted":   hasVoted,
	})
}

// basic bot detection by user-agent substrings
func isBot(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	if ua == "" {
		return false
	}
	bots := []string{"bot", "crawler", "spider", "bingpreview", "slurp", "curl", "wget", "headless", "python-requests"}
	for _, b := range bots {
		if strings.Contains(ua, b) {
			return true
		}
	}
	return false
}
