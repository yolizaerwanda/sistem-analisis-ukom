package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetData struct {
	Title   string
	Headers []string
	Rows    [][]interface{}
}

func GetSpreadsheetData() ([]*SpreadsheetData, error) {

	ctx := context.Background()

	credentialsFile := os.Getenv("GOOGLE_CREDENTIALS_FILE")
	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")

	if credentialsFile == "" {
		return nil, fmt.Errorf("GOOGLE_CREDENTIALS_FILE belum diatur")
	}

	if spreadsheetID == "" {
		return nil, fmt.Errorf("GOOGLE_SPREADSHEET_ID belum diatur")
	}

	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membaca credentials: %w",
			err,
		)
	}

	config, err := google.JWTConfigFromJSON(
		credentials,
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membaca konfigurasi Google: %w",
			err,
		)
	}

	client := config.Client(ctx)

	service, err := sheets.NewService(
		ctx,
		option.WithHTTPClient(client),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat Google Sheets service: %w",
			err,
		)
	}

	spreadsheet, err := service.Spreadsheets.
		Get(spreadsheetID).
		Do()

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil daftar sheet: %w",
			err,
		)
	}

	if len(spreadsheet.Sheets) == 0 {
		return nil, fmt.Errorf(
			"spreadsheet tidak memiliki sheet",
		)
	}

	var allData []*SpreadsheetData

	for _, sheet := range spreadsheet.Sheets {

		sheetName := sheet.Properties.Title

		fmt.Println()
		fmt.Println("====================================")
		fmt.Println("MEMPROSES SHEET:", sheetName)
		fmt.Println("====================================")

		titleRange := fmt.Sprintf(
			"%s!A1:Y3",
			sheetName,
		)

		titleResponse, err := service.Spreadsheets.Values.
			Get(spreadsheetID, titleRange).
			Do()

		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca judul spreadsheet pada sheet %q: %w",
				sheetName,
				err,
			)
		}

		title := extractReportTitle(titleResponse.Values)

		if title == "" {
			return nil, fmt.Errorf(
				"judul laporan tidak ditemukan pada sheet %q pada range A1:Y3",
				sheetName,
			)
		}

		readRange := fmt.Sprintf(
			"%s!A:ZZ",
			sheetName,
		)

		response, err := service.Spreadsheets.Values.
			Get(spreadsheetID, readRange).
			Do()

		if err != nil {
			return nil, fmt.Errorf(
				"gagal mengambil data spreadsheet pada sheet %q: %w",
				sheetName,
				err,
			)
		}

		if len(response.Values) == 0 {
			return nil, fmt.Errorf(
				"sheet %q tidak memiliki data",
				sheetName,
			)
		}

		headerRowIndex := findHeaderRow(response.Values)

		if headerRowIndex == -1 {
			return nil, fmt.Errorf(
				"header spreadsheet tidak ditemukan pada sheet %q",
				sheetName,
			)
		}

		headers := make(
			[]string,
			len(response.Values[headerRowIndex]),
		)

		for i, value := range response.Values[headerRowIndex] {

			headers[i] = strings.TrimSpace(
				fmt.Sprint(value),
			)
		}

		var rows [][]interface{}

		if headerRowIndex+1 < len(response.Values) {
			rows = response.Values[headerRowIndex+1:]
		}

		data := &SpreadsheetData{
			Title:   title,
			Headers: headers,
			Rows:    rows,
		}

		allData = append(allData, data)

		fmt.Println(
			"Sheet berhasil dibaca:",
			sheetName,
		)
	}

	if len(allData) == 0 {
		return nil, fmt.Errorf(
			"tidak ada data spreadsheet yang berhasil dibaca",
		)
	}

	fmt.Println()
	fmt.Println("====================================")
	fmt.Println(
		"TOTAL SHEET DIPROSES:",
		len(allData),
	)
	fmt.Println("====================================")

	return allData, nil
}

func extractReportTitle(rows [][]interface{}) string {

	for _, row := range rows {

		for _, value := range row {

			text := strings.TrimSpace(
				fmt.Sprint(value),
			)

			if text == "" {
				continue
			}

			upperText := strings.ToUpper(text)

			if strings.Contains(upperText, "PERIODE") &&
				strings.Contains(upperText, "KERJASAMA") {

				return text
			}
		}
	}

	return ""
}

func findHeaderRow(rows [][]interface{}) int {

	for i, row := range rows {

		if hasRequiredHeader(row) {
			return i
		}
	}

	return -1
}

func hasRequiredHeader(row []interface{}) bool {

	requiredHeaders := map[string]bool{
		"NAMA LENGKAP":   false,
		"ASAL INSTITUSI": false,
		"BATCH":          false,
		"KEHADIRAN (%)":  false,
		"RATA - RATA":    false,
	}

	for _, value := range row {

		header := strings.TrimSpace(
			strings.ToUpper(
				fmt.Sprint(value),
			),
		)

		if _, exists := requiredHeaders[header]; exists {
			requiredHeaders[header] = true
		}
	}

	for _, found := range requiredHeaders {

		if !found {
			return false
		}
	}

	return true
}