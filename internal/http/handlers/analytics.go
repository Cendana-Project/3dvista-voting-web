package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"voteweb/internal/domain"
)

type AnalyticsHandler struct {
	service domain.VoteService
	logger  *slog.Logger
}

func NewAnalyticsHandler(service domain.VoteService, logger *slog.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
		logger:  logger,
	}
}

// AnalyticsData represents the data structure for analytics page
type AnalyticsData struct {
	TotalInnovations int64
	TotalVotes       int64
	Innovations      []*InnovationStats
}

// InnovationStats represents statistics for a single innovation
type InnovationStats struct {
	*domain.Innovation
	VoteCount      int64
	VotePercentage float64
}

func (h *AnalyticsHandler) ShowAnalytics(c *gin.Context) {
	// Get all innovations with their vote counts
	innovations, err := h.service.ListInnovations(c.Request.Context())
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to list innovations", "error", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl.html", gin.H{
			"Title":   "Error",
			"Message": "An error occurred while loading analytics.",
		})
		return
	}

	// Fetch vote counts for each innovation
	var totalVotes int64
	var innovationsWithStats []*InnovationStats

	for _, innovation := range innovations {
		voteCount, err := h.service.GetVoteCount(c.Request.Context(), innovation.ID)
		if err != nil {
			h.logger.ErrorContext(c.Request.Context(), "failed to get vote count",
				"innovation_id", innovation.ID, "error", err)
			voteCount = 0
		}

		totalVotes += voteCount

		// Calculate percentage (we'll do this in template based on max votes)
		innovationsWithStats = append(innovationsWithStats, &InnovationStats{
			Innovation: innovation,
			VoteCount:  voteCount,
		})
	}

	// Find max vote count for percentage calculation
	var maxVotes int64
	for _, stats := range innovationsWithStats {
		if stats.VoteCount > maxVotes {
			maxVotes = stats.VoteCount
		}
	}

	// Calculate percentages
	for _, stats := range innovationsWithStats {
		if maxVotes > 0 {
			stats.VotePercentage = (float64(stats.VoteCount) / float64(maxVotes)) * 100
		} else {
			stats.VotePercentage = 0
		}
	}

	analytics := AnalyticsData{
		TotalInnovations: int64(len(innovations)),
		TotalVotes:       totalVotes,
		Innovations:      innovationsWithStats,
	}

	c.HTML(http.StatusOK, "analytics.tmpl.html", gin.H{
		"Title":     "Analytics Dashboard",
		"Analytics": analytics,
		"MaxVotes":  maxVotes,
	})
}

// GetAnalyticsData returns analytics data as JSON for client-side rendering
func (h *AnalyticsHandler) GetAnalyticsData(c *gin.Context) {
	innovations, err := h.service.ListInnovations(c.Request.Context())
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to list innovations", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load analytics data",
		})
		return
	}

	var totalVotes int64
	var innovationsWithStats []*InnovationStats

	for _, innovation := range innovations {
		voteCount, err := h.service.GetVoteCount(c.Request.Context(), innovation.ID)
		if err != nil {
			h.logger.ErrorContext(c.Request.Context(), "failed to get vote count",
				"innovation_id", innovation.ID, "error", err)
			voteCount = 0
		}

		totalVotes += voteCount
		innovationsWithStats = append(innovationsWithStats, &InnovationStats{
			Innovation: innovation,
			VoteCount:  voteCount,
		})
	}

	// Find max vote count for percentage calculation
	var maxVotes int64
	for _, stats := range innovationsWithStats {
		if stats.VoteCount > maxVotes {
			maxVotes = stats.VoteCount
		}
	}

	// Calculate percentages
	for _, stats := range innovationsWithStats {
		if maxVotes > 0 {
			stats.VotePercentage = (float64(stats.VoteCount) / float64(maxVotes)) * 100
		} else {
			stats.VotePercentage = 0
		}
	}

	// Get total unique voters
	totalVoters, err := h.service.GetTotalVoters(c.Request.Context())
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to get total voters", "error", err)
		totalVoters = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"total_innovations": len(innovations),
		"total_votes":       totalVotes,
		"total_voters":      totalVoters,
		"max_votes":         maxVotes,
		"innovations":       innovationsWithStats,
	})
}
