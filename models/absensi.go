package models

type Absensi struct {
	IDAbsensi               int      `db:"id_absensi"`
	IDAnalisis              int      `db:"id_analisis"`
	IDPeserta               int      `db:"id_peserta"`
	PersentaseKehadiran     *float64 `db:"persentase_kehadiran"`
	KehadiranKurangDari40   bool     `db:"kehadiran_kurang_dari_40"`
}