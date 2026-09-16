package models

type Zona struct {
	IDZona        int `db:"id_zona"`
	IDAnalisis    int `db:"id_analisis"`
	JumlahMerah   int `db:"jumlah_merah"`
	JumlahKuning  int `db:"jumlah_kuning"`
	JumlahHijau   int `db:"jumlah_hijau"`
	TotalPeserta  int `db:"total_peserta"`
}