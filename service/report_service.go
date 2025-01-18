package service

import (
	"encoding/csv"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"main.go/entity"
	"main.go/repository"
	"os"
	"time"
)

type ReportService interface {
	GenerateReport(filters entity.ReportFilters, isAdmin bool, userID uint) ([]entity.TransactionSummary, error)
	SaveReportToCSV(summaries []entity.TransactionSummary) (string, error)
	SaveReportToPDF(summaries []entity.TransactionSummary) (string, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GenerateReport(filters entity.ReportFilters, isAdmin bool, userID uint) ([]entity.TransactionSummary, error) {
	if !isAdmin {
		filters.UserID = userID // Hanya tetapkan UserID jika bukan admin
	}
	return s.reportRepo.GetTransactionsSummary(filters, isAdmin, userID)
}

func (s *reportService) SaveReportToCSV(summaries []entity.TransactionSummary) (string, error) {
	filePath := fmt.Sprintf("Reports/report_%d.csv", time.Now().Unix())
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Transaction ID", "User ID", "User Name", "Product Name", "Category Name", "Quantity", "Total Price", "Transaction Date"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	for _, summary := range summaries {
		row := []string{
			fmt.Sprintf("%d", summary.TransactionID),
			fmt.Sprintf("%d", summary.UserID),
			summary.UserName,
			summary.ProductName,
			summary.CategoryName,
			fmt.Sprintf("%d", summary.Quantity),
			fmt.Sprintf("%.2f", summary.TotalPrice),
			summary.TransactionDate,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return filePath, nil
}

func (s *reportService) SaveReportToPDF(summaries []entity.TransactionSummary) (string, error) {
	filePath := fmt.Sprintf("Reports/report_%d.pdf", time.Now().Unix())
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 10, "Transaction Report", "0", 1, "C", false, 0, "")

	headers := []string{"Transaction ID", "User ID", "User Name", "Product Name", "Category Name", "Quantity", "Total Price", "Transaction Date"}
	for _, header := range headers {
		pdf.CellFormat(24, 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, summary := range summaries {
		pdf.CellFormat(24, 10, fmt.Sprintf("%d", summary.TransactionID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, fmt.Sprintf("%d", summary.UserID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, summary.UserName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, summary.ProductName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, summary.CategoryName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, fmt.Sprintf("%d", summary.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, fmt.Sprintf("%.2f", summary.TotalPrice), "1", 0, "C", false, 0, "")
		pdf.CellFormat(24, 10, summary.TransactionDate, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
