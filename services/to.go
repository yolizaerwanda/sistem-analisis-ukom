package services 
import ( "fmt" 
		 "strings" 
		)

func IsTOColumn(header string) bool {

	header = strings.TrimSpace(
		strings.ToUpper(header),
	)

	if header == "SCORE PRE-TEST" {
		return true
	}

	if strings.HasPrefix(header, "POST TEST ") {
		return true
	}

	if strings.HasPrefix(header, "MT ") {
		return true
	}

	if strings.HasPrefix(header, "SCORE TO") {
		return true
	}

	if strings.HasPrefix(header, "TO INTENSIF") {
		return true
	}

	if header == "FINAL PREPARATION" {
		return true
	}

	return false
}

func GetActiveTOColumnsByBatch(
	data *SpreadsheetData,
	batch string,
) []int {
	var activeColumns []int

	nameColumnIndex := -1

	batchColumnIndex := -1

	for i, header := range data.Headers {
		headerUpper := strings.TrimSpace(strings.ToUpper(header))

		if headerUpper == "NAMA LENGKAP" {
			nameColumnIndex = i
		}

		if headerUpper == "BATCH" {
			batchColumnIndex = i
		}
	}

	if nameColumnIndex == -1 {
		fmt.Println("Kolom NAMA LENGKAP tidak ditemukan")
		return activeColumns
	}

	if batchColumnIndex == -1 {
		fmt.Println("Kolom BATCH tidak ditemukan")
		return activeColumns
	}

	fmt.Println()
	fmt.Println("========== CEK KOLOM TO PER BATCH ==========")
	fmt.Println("BATCH:", batch)

	for columnIndex, header := range data.Headers {

		if !IsTOColumn(header) {
			continue
		}

		hasValue := hasValueByBatch(
			data.Rows,
			columnIndex,
			nameColumnIndex,
			batchColumnIndex,
			batch,
		)

		fmt.Printf(
			"INDEX=%d | HEADER=%q | AKTIF=%v\n",
			columnIndex,
			header,
			hasValue,
		)

		if hasValue {
			activeColumns = append(activeColumns, columnIndex)
		}
	}

	fmt.Println("----------------------------------")
	fmt.Println("BATCH:", batch)
	fmt.Println("TOTAL TO AKTIF:", len(activeColumns))
	fmt.Println("===========================================")
	fmt.Println()

	return activeColumns
}

func hasValueByBatch(
	rows [][]interface{},
	columnIndex int,
	nameColumnIndex int,
	batchColumnIndex int,
	targetBatch string,
) bool {

	for _, row := range rows {

		if nameColumnIndex >= len(row) ||
			batchColumnIndex >= len(row) ||
			columnIndex >= len(row) {
			continue
		}

		nama := strings.TrimSpace(
			fmt.Sprint(row[nameColumnIndex]),
		)

		if nama == "" || nama == "<nil>" {
			continue
		}

		batch := strings.TrimSpace(
			fmt.Sprint(row[batchColumnIndex]),
		)

		if batch == "" || batch == "<nil>" {
			continue
		}

		if !strings.EqualFold(batch, targetBatch) {
			continue
		}

		if row[columnIndex] == nil {
			continue
		}

		value := strings.TrimSpace(
			fmt.Sprint(row[columnIndex]),
		)

		if value == "" || value == "<nil>" {
			continue
		}

		return true
	}

	return false
}