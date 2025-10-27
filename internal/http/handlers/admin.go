package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.tmpl.html", gin.H{
		"Title": "Admin Login",
	})
}

func (h *AdminHandler) ShowDashboardViewer(c *gin.Context) {
	c.HTML(http.StatusOK, "analytics_viewer.tmpl.html", gin.H{
		"Title": "Analytics Dashboard",
	})
}
