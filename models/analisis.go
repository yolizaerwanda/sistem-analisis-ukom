package models

type Analisis struct {
	IDAnalisis int    `db:"id_analisis"`
	Periode    int    `db:"periode"`
	Tahun      int    `db:"tahun"`
	Institusi  string `db:"institusi"`
	Batch      string `db:"batch"`
}