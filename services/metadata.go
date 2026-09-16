package services

import (
	"fmt"
	"regexp"
	"strings"
)

type Metadata struct {
	Periode   int
	Tahun     int
	Institusi string
}

func ParseMetadata(title string) (*Metadata, error) {

	title = strings.TrimSpace(title)

	re := regexp.MustCompile(
		`(?i)PERIODE\s+(\d+)\s+(\d{4})`,
	)

	matches := re.FindStringSubmatch(title)

	if len(matches) < 3 {
		return nil, fmt.Errorf(
			"periode dan tahun tidak ditemukan pada judul spreadsheet",
		)
	}

	var periode int
	var tahun int

	_, err := fmt.Sscanf(
		matches[1],
		"%d",
		&periode,
	)

	if err != nil {
		return nil, fmt.Errorf(
		"periode tidak valid: %w",
		err,
	)
	}

	_, err = fmt.Sscanf(
		matches[2],
		"%d",
		&tahun,
	)

	if err != nil {
		return nil, fmt.Errorf(
		"tahun tidak valid: %w",
		err,
	)

	}

	institusi := ParseInstitusi(title)

	if institusi == "" {
		return nil, fmt.Errorf(
			"institusi tidak ditemukan pada judul spreadsheet",
		)
	}

	return &Metadata{
		Periode:   periode,
		Tahun:     tahun,
		Institusi: institusi,
	}, nil
}

func ParseInstitusi(title string) string {
	title = strings.TrimSpace(title)

	re := regexp.MustCompile(`(?i)KERJASAMA\s+(.+?)\s+DAN\s+`)
	matches := re.FindStringSubmatch(title)

	if len(matches) < 2 {
		return ""
	}

	return strings.TrimSpace(matches[1])
}