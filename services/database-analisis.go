package services

import (
	"database/sql"
	"fmt"
)

type ResultAnalisis struct {
	IDAnalisis int
	Periode    string
	Tahun      string
	Institusi  string
	Batch      string

	TotalPeserta int

	Attendance []ResultAttendance
	Zona       ResultZona
	TO         []ResultTO
}

type ResultAttendance struct {
	No                  int
	Nama                string
	PersentaseKehadiran float64
}

type ResultZona struct {
	JumlahMerah  int
	JumlahKuning int
	JumlahHijau  int
	TotalPeserta int
}

type ResultTO struct {
	No                   int
	Nama                 string
	TotalPengerjaanTO    int
	TotalBelumPengerjaan int
	TotalTO              int
	Persentase            float64
}

func SaveAnalysis(db *sql.DB, data *SpreadsheetData) error {

	hasilPerBatch, peserta, err := AnalyzeAllBatches(data)
	if err != nil {
		return err
	}

	metadata, err := ParseMetadata(data.Title)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf(
			"gagal memulai transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	for batch, hasil := range hasilPerBatch {

		idAnalisis, err := getOrCreateAnalisisTx(
			tx,
			metadata,
			batch,
		)

		if err != nil {
			return fmt.Errorf(
				"gagal menyimpan analisis batch %s: %w",
				batch,
				err,
			)
		}

		fmt.Println(
			"Analisis batch:",
			batch,
			"id:",
			idAnalisis,
		)

		err = deleteAnalysisDetail(
			tx,
			idAnalisis,
		)

		if err != nil {
			return fmt.Errorf(
				"gagal menghapus data lama batch %s: %w",
				batch,
				err,
			)
		}

		pesertaBatch := FilterPesertaByBatch(
			peserta,
			batch,
		)

		pesertaIDs := make(map[string]int)

		for _, p := range pesertaBatch {

			idPeserta, err := createPeserta(
				tx,
				idAnalisis,
				p,
			)

			if err != nil {
				return fmt.Errorf(
					"gagal menyimpan peserta %s: %w",
					p.Nama,
					err,
				)
			}

			key := p.Nama + "|" + p.Batch

			pesertaIDs[key] = idPeserta

			fmt.Println(
				"Peserta:",
				p.Nama,
				"id:",
				idPeserta,
			)

			persentaseKehadiran := 0.0

			if p.PersentaseKehadiran != nil {
				persentaseKehadiran = *p.PersentaseKehadiran
			}

			if err := createAbsensi(
				tx,
				idAnalisis,
				idPeserta,
				persentaseKehadiran,
			); err != nil {
				return fmt.Errorf(
					"gagal menyimpan absensi %s: %w",
					p.Nama,
					err,
				)
			}
		}

		for _, hasilTO := range hasil.HasilTO {

			key := hasilTO.NamaPeserta +
				"|" +
				hasilTO.Batch

			idPeserta, exists := pesertaIDs[key]

			if !exists {
				return fmt.Errorf(
					"id peserta untuk %s tidak ditemukan",
					hasilTO.NamaPeserta,
				)
			}

			err = createTO(
				tx,
				idAnalisis,
				idPeserta,
				hasilTO,
			)

			if err != nil {
				return fmt.Errorf(
					"gagal menyimpan hasil TO %s: %w",
					hasilTO.NamaPeserta,
					err,
				)
			}
		}

		err = createZona(
			tx,
			idAnalisis,
			hasil,
		)

		if err != nil {
			return fmt.Errorf(
				"gagal menyimpan zona batch %s: %w",
				batch,
				err,
			)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"gagal commit transaction: %w",
			err,
		)
	}

	return nil
}

func GetFilterData(db *sql.DB) (map[string][]string, error) {

	data := map[string][]string{
		"periode":   {},
		"tahun":     {},
		"institusi": {},
		"batch":     {},
	}

	rows, err := db.Query(`
		SELECT DISTINCT periode
		FROM analisis
		ORDER BY periode
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var value string

		if err := rows.Scan(&value); err != nil {
			return nil, err
		}

		data["periode"] = append(data["periode"], value)
	}

	rowsTahun, err := db.Query(`
		SELECT DISTINCT tahun
		FROM analisis
		ORDER BY tahun DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rowsTahun.Close()

	for rowsTahun.Next() {
		var value string

		if err := rowsTahun.Scan(&value); err != nil {
			return nil, err
		}

		data["tahun"] = append(data["tahun"], value)
	}

	rowsInstitusi, err := db.Query(`
		SELECT DISTINCT institusi
		FROM analisis
		ORDER BY institusi
	`)
	if err != nil {
		return nil, err
	}
	defer rowsInstitusi.Close()

	for rowsInstitusi.Next() {
		var value string

		if err := rowsInstitusi.Scan(&value); err != nil {
			return nil, err
		}

		data["institusi"] = append(data["institusi"], value)
	}

	rowsBatch, err := db.Query(`
		SELECT DISTINCT batch
		FROM analisis
		ORDER BY batch
	`)
	if err != nil {
		return nil, err
	}
	defer rowsBatch.Close()

	for rowsBatch.Next() {
		var value string

		if err := rowsBatch.Scan(&value); err != nil {
			return nil, err
		}

		data["batch"] = append(data["batch"], value)
	}

	return data, nil
}

func getOrCreateAnalisisTx(
	tx *sql.Tx,
	metadata *Metadata,
	batch string,
) (int, error) {

	var idAnalisis int

	err := tx.QueryRow(`
		SELECT id_analisis
		FROM analisis
		WHERE periode = ?
		  AND tahun = ?
		  AND institusi = ?
		  AND batch = ?
		LIMIT 1
	`,
		metadata.Periode,
		metadata.Tahun,
		metadata.Institusi,
		batch,
	).Scan(&idAnalisis)

	if err == nil {
		return idAnalisis, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := tx.Exec(`
		INSERT INTO analisis (
			periode,
			tahun,
			institusi,
			batch
		)
		VALUES (?, ?, ?, ?)
	`,
		metadata.Periode,
		metadata.Tahun,
		metadata.Institusi,
		batch,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func createPeserta(
	tx *sql.Tx,
	idAnalisis int,
	peserta PesertaData,
) (int, error) {

	result, err := tx.Exec(`
		INSERT INTO peserta (
			id_analisis,
			nama_peserta,
			batch,
			persentase_kehadiran,
			rata_rata_to
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		idAnalisis,
		peserta.Nama,
		peserta.Batch,
		peserta.PersentaseKehadiran,
		peserta.RataRataTO,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func createTO(
	tx *sql.Tx,
	idAnalisis int,
	idPeserta int,
	hasil HasilTO,
) error {

	_, err := tx.Exec(`
		INSERT INTO pengerjaan_to (
			id_analisis,
			id_peserta,
			total_pengerjaan_to,
			total_belum_pengerjaan_to,
			total_to
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		idAnalisis,
		idPeserta,
		hasil.TotalPengerjaanTO,
		hasil.TotalBelumPengerjaan,
		hasil.TotalTO,
	)

	return err
}

func createAbsensi(tx *sql.Tx, idAnalisis int, idPeserta int, persentase float64) error {
    kurangDari40 := persentase < 40

    _, err := tx.Exec(`
        INSERT INTO absensi (
            id_analisis,
            id_peserta,
            persentase_kehadiran,
            kehadiran_kurang_dari_40
        )
        VALUES (?, ?, ?, ?)
    `, idAnalisis, idPeserta, persentase, kurangDari40)

    return err
}

func createZona(
	tx *sql.Tx,
	idAnalisis int,
	hasil *HasilAnalisis,
) error {

	_, err := tx.Exec(`
		INSERT INTO zona (
			id_analisis,
			jumlah_merah,
			jumlah_kuning,
			jumlah_hijau,
			total_peserta
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		idAnalisis,
		hasil.JumlahMerah,
		hasil.JumlahKuning,
		hasil.JumlahHijau,
		hasil.TotalPeserta,
	)

	return err
}

func deleteAnalysisDetail(
	tx *sql.Tx,
	idAnalisis int,
) error {

	_, err := tx.Exec(`
		DELETE FROM pengerjaan_to
		WHERE id_analisis = ?
	`, idAnalisis)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM absensi
		WHERE id_analisis = ?
	`, idAnalisis)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM zona
		WHERE id_analisis = ?
	`, idAnalisis)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM peserta
		WHERE id_analisis = ?
	`, idAnalisis)

	if err != nil {
		return err
	}

	return nil
}

func GetFilterTahun(db *sql.DB) ([]string, error) {

	rows, err := db.Query(`
		SELECT DISTINCT tahun
		FROM analisis
		ORDER BY tahun DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil tahun: %w", err)
	}

	defer rows.Close()

	var result []string

	for rows.Next() {

		var tahun string

		if err := rows.Scan(&tahun); err != nil {
			return nil, fmt.Errorf("gagal membaca tahun: %w", err)
		}

		result = append(result, tahun)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca data tahun: %w", err)
	}

	return result, nil
}

func GetFilterPeriode(
	db *sql.DB,
	tahun string,
) ([]string, error) {

	rows, err := db.Query(`
		SELECT DISTINCT periode
		FROM analisis
		WHERE tahun = ?
		ORDER BY periode
	`, tahun)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil periode: %w", err)
	}

	defer rows.Close()

	var result []string

	for rows.Next() {

		var periode string

		if err := rows.Scan(&periode); err != nil {
			return nil, fmt.Errorf("gagal membaca periode: %w", err)
		}

		result = append(result, periode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca data periode: %w", err)
	}

	return result, nil
}

func GetFilterInstitusi(
	db *sql.DB,
	tahun string,
	periode string,
) ([]string, error) {

	rows, err := db.Query(`
		SELECT DISTINCT institusi
		FROM analisis
		WHERE tahun = ?
		  AND periode = ?
		ORDER BY institusi
	`, tahun, periode)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil institusi: %w", err)
	}

	defer rows.Close()

	var result []string

	for rows.Next() {

		var institusi string

		if err := rows.Scan(&institusi); err != nil {
			return nil, fmt.Errorf("gagal membaca institusi: %w", err)
		}

		result = append(result, institusi)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca data institusi: %w", err)
	}

	return result, nil
}

func GetFilterBatch(
	db *sql.DB,
	tahun string,
	periode string,
	institusi string,
) ([]string, error) {

	rows, err := db.Query(`
		SELECT DISTINCT batch
		FROM analisis
		WHERE tahun = ?
		  AND periode = ?
		  AND institusi = ?
		ORDER BY batch
	`, tahun, periode, institusi)

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil batch: %w", err)
	}

	defer rows.Close()

	var result []string

	for rows.Next() {

		var batch string

		if err := rows.Scan(&batch); err != nil {
			return nil, fmt.Errorf("gagal membaca batch: %w", err)
		}

		result = append(result, batch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal membaca data batch: %w", err)
	}

	return result, nil
}

func GetIDAnalisis(
	db *sql.DB,
	tahun string,
	periode string,
	institusi string,
	batch string,
) (int, error) {

	var idAnalisis int

	err := db.QueryRow(`
		SELECT id_analisis
		FROM analisis
		WHERE tahun = ?
		  AND periode = ?
		  AND institusi = ?
		  AND batch = ?
		LIMIT 1
	`,
		tahun,
		periode,
		institusi,
		batch,
	).Scan(&idAnalisis)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("data analisis tidak ditemukan")
		}

		return 0, fmt.Errorf(
			"gagal mencari analisis: %w",
			err,
		)
	}

	return idAnalisis, nil
}

func GetAttendanceResult(
	db *sql.DB,
	idAnalisis int,
) ([]ResultAttendance, error) {

	rows, err := db.Query(`
		SELECT
			p.nama_peserta,
			a.persentase_kehadiran
		FROM absensi a
		INNER JOIN peserta p
			ON p.id_peserta = a.id_peserta
		WHERE a.id_analisis = ?
		  AND a.kehadiran_kurang_dari_40 = 1
		ORDER BY p.nama_peserta
	`,
		idAnalisis,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil data kehadiran: %w",
			err,
		)
	}

	defer rows.Close()

	var result []ResultAttendance
	no := 1

	for rows.Next() {
		var data ResultAttendance

		if err := rows.Scan(
			&data.Nama,
			&data.PersentaseKehadiran,
		); err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data kehadiran: %w",
				err,
			)
		}

		data.No = no
		result = append(result, data)
		no++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca data kehadiran: %w",
			err,
		)
	}

	return result, nil
}

func GetZonaResult(
	db *sql.DB,
	idAnalisis int,
) (ResultZona, error) {

	var result ResultZona

	err := db.QueryRow(`
		SELECT
			jumlah_merah,
			jumlah_kuning,
			jumlah_hijau,
			total_peserta
		FROM zona
		WHERE id_analisis = ?
		LIMIT 1
	`,
		idAnalisis,
	).Scan(
		&result.JumlahMerah,
		&result.JumlahKuning,
		&result.JumlahHijau,
		&result.TotalPeserta,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return result, fmt.Errorf(
				"data zona tidak ditemukan",
			)
		}

		return result, fmt.Errorf(
			"gagal mengambil data zona: %w",
			err,
		)
	}

	return result, nil
}

func GetTOResult(
	db *sql.DB,
	idAnalisis int,
) ([]ResultTO, error) {

	rows, err := db.Query(`
		SELECT
			p.nama_peserta,
			pt.total_pengerjaan_to,
			pt.total_belum_pengerjaan_to,
			pt.total_to
		FROM pengerjaan_to pt
		INNER JOIN peserta p
			ON p.id_peserta = pt.id_peserta
		WHERE pt.id_analisis = ?
		ORDER BY p.nama_peserta
	`,
		idAnalisis,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil data pengerjaan TO: %w",
			err,
		)
	}

	defer rows.Close()

	var result []ResultTO

	no := 1

	for rows.Next() {

		var data ResultTO

		if err := rows.Scan(
			&data.Nama,
			&data.TotalPengerjaanTO,
			&data.TotalBelumPengerjaan,
			&data.TotalTO,
		); err != nil {

			return nil, fmt.Errorf(
				"gagal membaca data pengerjaan TO: %w",
				err,
			)
		}

		if data.TotalTO > 0 {
			data.Persentase =
				float64(data.TotalPengerjaanTO) /
					float64(data.TotalTO) *
					100
		}

		data.No = no

		result = append(result, data)

		no++
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca data pengerjaan TO: %w",
			err,
		)
	}

	return result, nil
}

func GetResultAnalisis(
	db *sql.DB,
	tahun string,
	periode string,
	institusi string,
	batch string,
) (*ResultAnalisis, error) {

	idAnalisis, err := GetIDAnalisis(
		db,
		tahun,
		periode,
		institusi,
		batch,
	)

	if err != nil {
		return nil, err
	}


	attendance, err := GetAttendanceResult(
		db,
		idAnalisis,
	)

	if err != nil {
		return nil, err
	}


	zona, err := GetZonaResult(
		db,
		idAnalisis,
	)

	if err != nil {
		return nil, err
	}


	to, err := GetTOResult(
		db,
		idAnalisis,
	)

	if err != nil {
		return nil, err
	}


	return &ResultAnalisis{

		IDAnalisis: idAnalisis,

		Periode:   periode,
		Tahun:     tahun,
		Institusi: institusi,
		Batch:     batch,

		TotalPeserta: zona.TotalPeserta,

		Attendance: attendance,

		Zona: zona,

		TO: to,

	}, nil
}