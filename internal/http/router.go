package http

import (
	"html/template"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"voteweb/internal/config"
	"voteweb/internal/domain"
	"voteweb/internal/http/handlers"
	"voteweb/internal/http/middleware"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(cfg *config.Config, pool *pgxpool.Pool, service domain.VoteService, logger *slog.Logger) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	router := gin.New()

	// Load templates
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int64) int64 { return a + b },
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},
	})
	router.LoadHTMLGlob("web/templates/*")

	// Serve static files
	router.Static("/static", "./web/static")

	// Global middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Recover(logger))
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.ProxiedIP(cfg.TrustProxy, cfg.AllowedProxyCIDRs))
	router.Use(middleware.CSRF())

	// Logging middleware
	router.Use(func(c *gin.Context) {
		logger.InfoContext(c.Request.Context(),
			"request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.GetString("client_ip"),
			"request_id", c.GetString("request_id"))
		c.Next()
	})

	// Health check
	healthHandler := handlers.NewHealthHandler(pool)
	router.GET("/healthz", healthHandler.HealthCheck)

	// List handler
	listHandler := handlers.NewListHandler(service, logger)
	router.GET("/", listHandler.ShowList)

	// Admin login page (public, no auth required)
	adminHandler := handlers.NewAdminHandler()
	router.GET("/admin/login", adminHandler.ShowLogin)

	// Admin dashboard viewer (client-side, no server-side auth required)
	router.GET("/admin/dashboard", adminHandler.ShowDashboardViewer)

	// Admin routes - protected with X-ADMIN-CODE header (must be before /:group/:slug)
	if cfg.AdminCode != "" {
		adminGroup := router.Group("/admin")

		analyticsHandler := handlers.NewAnalyticsHandler(service, logger)

		// API endpoint for JSON data
		adminGroup.Use(middleware.AdminAuth(cfg.AdminCode))
		adminGroup.GET("/analytics", analyticsHandler.ShowAnalytics)
		adminGroup.GET("/api/data", analyticsHandler.GetAnalyticsData)

		logger.Info("Admin routes enabled",
			"login_path", "/admin/login",
			"dashboard_path", "/admin/dashboard",
			"analytics_path", "/admin/analytics",
			"api_path", "/admin/api/data")
	}

	// API handlers
	voteHandler := handlers.NewVoteHandler(service, logger)
	router.POST("/api/vote/:group/:slug", voteHandler.SubmitVote)

	// Page handler (catch-all, must be last)
	pageHandler := handlers.NewPageHandler(service, logger)
	router.GET("/:group/:slug", pageHandler.ShowInnovation)

	return router
}
