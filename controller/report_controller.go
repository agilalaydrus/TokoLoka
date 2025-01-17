package controller

import (
	"github.com/gin-gonic/gin"
	"main.go/entity"
	"main.go/service"
	"net/http"
	"strconv"
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

	// Validasi wajib adanya start_date dan end_date
	if filters.StartDate == "" || filters.EndDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	// Ambil user_id dan role dari JWT
	userRole := c.GetString("role")
	isAdmin := userRole == "administrator"
	if !isAdmin {
		filters.UserID = c.GetUint("user_id")
	}

	// Ambil page dan limit dari query parameter
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	// Default nilai jika tidak diisi
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default 10 data per halaman
	}

	filters.Page = page
	filters.Limit = limit

	// Panggil service untuk menghasilkan laporan
	summaries, err := rc.reportService.GenerateReport(filters, isAdmin, c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	// Gunakan variabel response untuk respons JSON
	response := gin.H{
		"data":  summaries,
		"page":  page,
		"limit": limit,
		"total": len(summaries), // Total data pada halaman ini
	}
	c.JSON(http.StatusOK, response)
}

// DownloadReport untuk mengunduh laporan
func (rc *ReportController) DownloadReport(c *gin.Context) {
	// Validasi format laporan (CSV atau PDF)
	format := c.Query("format")
	if format != "csv" && format != "pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Use 'csv' or 'pdf'"})
		return
	}

	// Ambil filter dari query parameter
	filters := entity.ReportFilters{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	// Validasi wajib adanya start_date dan end_date
	if filters.StartDate == "" || filters.EndDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	// Ambil user_id dan role dari JWT
	userRole := c.GetString("role")
	isAdmin := userRole == "administrator"
	if !isAdmin {
		// Jika bukan admin, batasi hanya transaksi milik user
		filters.UserID = c.GetUint("user_id")
	}

	// Pagination (opsional untuk report download)
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	filters.Page = page
	filters.Limit = limit

	// Panggil service untuk mendapatkan data laporan
	summaries, err := rc.reportService.GenerateReport(filters, isAdmin, c.GetUint("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	// Simpan laporan sesuai format
	var filePath string
	if format == "csv" {
		filePath, err = rc.reportService.SaveReportToCSV(summaries)
	} else if format == "pdf" {
		filePath, err = rc.reportService.SaveReportToPDF(summaries)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report"})
		return
	}

	// Kirim file ke klien
	c.File(filePath)
}
