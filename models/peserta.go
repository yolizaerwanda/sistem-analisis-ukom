package models

type Peserta struct {
	IDPeserta             int     `db:"id_peserta"`
	NamaPeserta           string  `db:"nama_peserta"`
	Batch                 string  `db:"batch"`
	PersentaseKehadiran   *float64 `db:"persentase_kehadiran"`
	RataRataTO            *float64 `db:"rata_rata_to"`
}