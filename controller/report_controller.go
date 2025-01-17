package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/entity"
	"main.go/service"
	"net/http"
)

type ReportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

// GenerateReport untuk mendapatkan data laporan
func (rc *ReportController) GenerateReport(c *gin.Context) {
	var filters entity.ReportFilters
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filters"})
		return
	}

	// Ambil user_id dan role dari JWT
	userRole := c.GetString("role")
	isAdmin := userRole == "administrator"
	if !isAdmin {
		filters.UserID = c.GetUint("user_id")
	}

	summaries, err := rc.reportService.GenerateReport(filters, isAdmin, c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summaries})
}

// DownloadReport untuk mengunduh laporan
func (rc *ReportController) DownloadReport(c *gin.Context) {
	format := c.Query("format")
	filters := entity.ReportFilters{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query"})
		return
	}

	userRole := c.GetString("role")
	isAdmin := userRole == "administrator"
	if !isAdmin {
		filters.UserID = c.GetUint("user_id")
	}

	summaries, err := rc.reportService.GenerateReport(filters, isAdmin, c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	var filePath string
	if format == "csv" {
		filePath, err = rc.reportService.SaveReportToCSV(summaries)
	} else if format == "pdf" {
		filePath, err = rc.reportService.SaveReportToPDF(summaries)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}
	if format == "" {
		format = "csv" // Default ke CSV
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report"})
		return
	}

	c.File(filePath)
}
