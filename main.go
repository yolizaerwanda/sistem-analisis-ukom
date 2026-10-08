package main

import (
	"database/sql"
	"net/http"
	"os"

	"sistem-analisis-ukom/config"
	"sistem-analisis-ukom/handlers"
	"sistem-analisis-ukom/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func databaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "create-user" {
		createUser()
		return
	}

	err := godotenv.Load()

	if err != nil {
		panic("Gagal membaca file .env")
	}

	db, err := config.ConnectDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	router := gin.Default()

	sessionSecret := os.Getenv("SESSION_SECRET")
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("ukom_session", store))

	router.Use(databaseMiddleware(db))

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")

	router.GET("/login", handlers.ShowLogin)
	router.POST("/login", handlers.ProcessLogin)
	router.GET("/logout", handlers.Logout)

	router.GET("/analisis", middleware.AuthRequired(), handlers.ShowAnalisis)
	router.POST("/analisis/process", middleware.AuthRequired(), handlers.ProcessAnalysis)
	router.GET("/api/filter/tahun", middleware.AuthRequired(), handlers.GetFilterTahun)
	router.GET("/api/filter/periode", middleware.AuthRequired(), handlers.GetFilterPeriode)
	router.GET("/api/filter/institusi", middleware.AuthRequired(), handlers.GetFilterInstitusi)
	router.GET("/api/filter/batch", middleware.AuthRequired(), handlers.GetFilterBatch)
	router.GET("/analisis/result", middleware.AuthRequired(), handlers.ShowResultAnalisis)

	router.Run(":8080")
}