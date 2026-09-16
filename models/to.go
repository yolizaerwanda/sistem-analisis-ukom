package models

type TO struct {
	IDTO                    int `db:"id_to"`
	IDAnalisis              int `db:"id_analisis"`
	IDPeserta               int `db:"id_peserta"`
	TotalPengerjaanTO       int `db:"total_pengerjaan_to"`
	TotalBelumPengerjaanTO  int `db:"total_belum_pengerjaan_to"`
	TotalTO                 int `db:"total_to"`
}