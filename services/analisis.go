package services

import (
	"fmt"
	"strconv"
	"strings"
)

type PesertaData struct {
	Nomor               string
	Nama                string
	Institusi           string
	Batch               string
	PersentaseKehadiran *float64
	RataRataTO          *float64

	TO map[string]*float64
}

type HasilAnalisis struct {
	TotalPeserta int

	JumlahMerah  int
	JumlahKuning int
	JumlahHijau  int

	TotalTOAktif int

	HasilTO      []HasilTO
	HasilAbsensi []HasilAbsensi
}

type HasilTO struct {
	NamaPeserta          string
	Batch                string
	TotalPengerjaanTO    int
	TotalBelumPengerjaan int
	TotalTO              int
	Persentase           float64
}

type HasilAbsensi struct {
	NamaPeserta         string
	Batch               string
	PersentaseKehadiran float64
}

func parseFloat(value interface{}) (*float64, error) {
	if value == nil {
		return nil, nil
	}

	text := strings.TrimSpace(fmt.Sprint(value))

	if text == "" {
		return nil, nil
	}

	switch strings.ToUpper(text) {
	case "#DIV/0!", "#N/A", "#VALUE!", "#REF!", "#NAME?", "#NUM!", "#NULL!":
		return nil, nil
	}

	text = strings.ReplaceAll(text, ",", ".")
	text = strings.TrimSuffix(text, "%")
	text = strings.TrimSpace(text)

	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return nil, err
	}

	return &number, nil
}

func createHeaderMap(headers []string) map[string]int {
	result := make(map[string]int)

	for i, header := range headers {
		key := strings.TrimSpace(
			strings.ToUpper(header),
		)

		result[key] = i
	}

	return result
}

func ParsePeserta(data *SpreadsheetData) ([]PesertaData, error) {
	headerMap := createHeaderMap(data.Headers)

	required := []string{
		"NAMA LENGKAP",
		"ASAL INSTITUSI",
		"BATCH",
		"KEHADIRAN (%)",
		"RATA - RATA",
	}

	for _, header := range required {
		if _, exists := headerMap[header]; !exists {
			return nil, fmt.Errorf(
				"kolom %q tidak ditemukan",
				header,
			)
		}
	}

	var peserta []PesertaData

	for _, row := range data.Rows {
		nama := getCell(row, headerMap["NAMA LENGKAP"])

		if strings.TrimSpace(nama) == "" {
			continue
		}

		kehadiran, err := parseFloat(
			getCell(row, headerMap["KEHADIRAN (%)"]),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"nilai kehadiran %s tidak valid: %w",
				nama,
				err,
			)
		}

		rataRata, err := parseFloat(
			getCell(row, headerMap["RATA - RATA"]),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"nilai rata-rata %s tidak valid: %w",
				nama,
				err,
			)
		}

		p := PesertaData{
			Nomor:               getCell(row, headerMap["NO"]),
			Nama:                nama,
			Institusi:           getCell(row, headerMap["ASAL INSTITUSI"]),
			Batch:               getCell(row, headerMap["BATCH"]),
			PersentaseKehadiran: kehadiran,
			RataRataTO:          rataRata,
			TO:                  make(map[string]*float64),
		}

		for i, header := range data.Headers {
			if !IsTOColumn(header) {
				continue
			}

			value, err := parseFloat(getCell(row, i))
			if err != nil {
				return nil, fmt.Errorf(
					"nilai TO %q untuk %s tidak valid: %w",
					header,
					nama,
					err,
				)
			}

			p.TO[header] = value
		}

		peserta = append(peserta, p)
	}

	return peserta, nil
}

func getCell(row []interface{}, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}

	return strings.TrimSpace(
		fmt.Sprint(row[index]),
	)
}

func CalculateTO(
	data *SpreadsheetData,
	peserta []PesertaData,
) []HasilTO {

	var hasil []HasilTO

	for _, p := range peserta {

		activeColumns := GetActiveTOColumnsByBatch(
			data,
			p.Batch,
		)

		totalTO := len(activeColumns)

		if totalTO == 0 {
			continue
		}

		dikerjakan := 0

		for _, columnIndex := range activeColumns {

			header := data.Headers[columnIndex]

			value, exists := p.TO[header]

			if !exists || value == nil {
				continue
			}

			dikerjakan++
		}

		belum := totalTO - dikerjakan

		persentase :=
			float64(dikerjakan) /
				float64(totalTO) *
				100

		if dikerjakan*2 < totalTO {
			hasil = append(hasil, HasilTO{
				NamaPeserta:          p.Nama,
				Batch:                p.Batch,
				TotalPengerjaanTO:    dikerjakan,
				TotalBelumPengerjaan: belum,
				TotalTO:              totalTO,
				Persentase:           persentase,
			})
		}
	}

	return hasil
}

func CalculateZona(
	peserta []PesertaData,
) (merah, kuning, hijau int) {

	for _, p := range peserta {

		if p.RataRataTO == nil {
			continue
		}

		score := *p.RataRataTO

		if score <= 48 {
			merah++
		} else if score < 55 {
			kuning++
		} else {
			hijau++
		}
	}

	return
}

func CalculateAbsensi(
	peserta []PesertaData,
) []HasilAbsensi {

	hasil := []HasilAbsensi{}

	for _, p := range peserta {

		persentase := 0.0

		if p.PersentaseKehadiran != nil {
			persentase = *p.PersentaseKehadiran
		}

		if persentase < 40 {
			hasil = append(hasil, HasilAbsensi{
				NamaPeserta:         p.Nama,
				Batch:               p.Batch,
				PersentaseKehadiran: persentase,
			})
		}
	}

	return hasil
}

func FilterPesertaByBatch(
	peserta []PesertaData,
	batch string,
) []PesertaData {

	var hasil []PesertaData

	batch = strings.TrimSpace(batch)

	for _, p := range peserta {

		if strings.TrimSpace(p.Batch) == batch {
			hasil = append(hasil, p)
		}
	}

	return hasil
}

func GetUniqueBatches(
	peserta []PesertaData,
) []string {

	batchMap := make(map[string]bool)

	for _, p := range peserta {

		batch := strings.TrimSpace(p.Batch)

		if batch == "" {
			continue
		}

		batchMap[batch] = true
	}

	var batches []string

	for batch := range batchMap {
		batches = append(batches, batch)
	}

	return batches
}

func AnalyzePeserta(
	data *SpreadsheetData,
	peserta []PesertaData,
	batch string,
) *HasilAnalisis {

	merah, kuning, hijau :=
		CalculateZona(peserta)

	hasilTO :=
		CalculateTO(data, peserta)

	hasilAbsensi :=
		CalculateAbsensi(peserta)

	activeTO :=
		GetActiveTOColumnsByBatch(data, batch)

	return &HasilAnalisis{
		TotalPeserta: len(peserta),

		JumlahMerah:  merah,
		JumlahKuning: kuning,
		JumlahHijau:  hijau,

		TotalTOAktif: len(activeTO),

		HasilTO:      hasilTO,
		HasilAbsensi: hasilAbsensi,
	}
}

func AnalyzeAllBatches(
	data *SpreadsheetData,
) (map[string]*HasilAnalisis, []PesertaData, error) {

	peserta, err := ParsePeserta(data)
	if err != nil {
		return nil, nil, err
	}

	batches := GetUniqueBatches(peserta)

	hasil := make(map[string]*HasilAnalisis)

	for _, batch := range batches {

		pesertaBatch :=
			FilterPesertaByBatch(
				peserta,
				batch,
			)

		hasil[batch] =
			AnalyzePeserta(
				data,
				pesertaBatch,
				batch,
			)
	}

	return hasil, peserta, nil
}