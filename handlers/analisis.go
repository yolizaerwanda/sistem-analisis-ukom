package handlers

import (
	"database/sql"
	"net/http"

	"sistem-analisis-ukom/services"

	"github.com/gin-gonic/gin"
)

func ProcessAnalysis(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	allData, err := services.GetSpreadsheetData()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data spreadsheet: " + err.Error(),
		})
		return
	}

	for _, data := range allData {

		err = services.SaveAnalysis(db, data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal memproses analisis: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analisis semua data berhasil diproses",
	})
}

func ShowAnalisis(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	filterData, err := services.GetFilterData(db)

	if err != nil {
		c.HTML(500, "analisis.html", gin.H{
			"error": "Gagal mengambil data filter.",
		})
		return
	}

	c.HTML(200, "analisis.html", gin.H{
		"PeriodeList":   filterData["periode"],
		"TahunList":     filterData["tahun"],
		"InstitusiList": filterData["institusi"],
		"BatchList":     filterData["batch"],
	})
}

func GetFilterTahun(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	tahun, err := services.GetFilterTahun(db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tahun,
	})
}

func GetFilterPeriode(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	tahun := c.Query("tahun")

	if tahun == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tahun wajib dipilih.",
		})
		return
	}

	periode, err := services.GetFilterPeriode(db, tahun)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": periode,
	})
}

func GetFilterInstitusi(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	tahun := c.Query("tahun")
	periode := c.Query("periode")

	if tahun == "" || periode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tahun dan periode wajib dipilih.",
		})
		return
	}

	institusi, err := services.GetFilterInstitusi(
		db,
		tahun,
		periode,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": institusi,
	})
}

func GetFilterBatch(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	tahun := c.Query("tahun")
	periode := c.Query("periode")
	institusi := c.Query("institusi")

	if tahun == "" || periode == "" || institusi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tahun, periode, dan institusi wajib dipilih.",
		})
		return
	}

	batch, err := services.GetFilterBatch(
		db,
		tahun,
		periode,
		institusi,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": batch,
	})
}

func ShowResultAnalisis(c *gin.Context) {

	db := c.MustGet("db").(*sql.DB)

	tahun := c.Query("tahun")
	periode := c.Query("periode")
	institusi := c.Query("institusi")
	batch := c.Query("batch")


	if tahun == "" ||
		periode == "" ||
		institusi == "" ||
		batch == "" {

		c.HTML(
			http.StatusBadRequest,
			"hasil-analisis.html",
			gin.H{
				"error": "Filter analisis belum lengkap.",
			},
		)

		return
	}


	result, err := services.GetResultAnalisis(
		db,
		tahun,
		periode,
		institusi,
		batch,
	)

	if err != nil {

		c.HTML(
			http.StatusInternalServerError,
			"hasil-analisis.html",
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}


	c.HTML(
		http.StatusOK,
		"hasil-analisis.html",
		gin.H{
			"Result": result,
		},
	)
}