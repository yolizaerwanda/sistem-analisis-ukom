package main

import (
	"net/http"
	"database/sql"
	"os"
	"sistem-analisis-ukom/config"
	"sistem-analisis-ukom/middleware"
	"sistem-analisis-ukom/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
)

func showLogin(c *gin.Context) {
    errorMessage := c.Query("error")

    c.HTML(200, "login.html", gin.H{
        "error": errorMessage,
    })
}

func processLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	
	if username == "" {
		c.Redirect(302, "/login?error=Username+tidak+boleh+kosong")
		return
	}

	if password == "" {
		c.Redirect(302, "/login?error=Password+tidak+boleh+kosong")
		return
	}
	
	db := c.MustGet("db").(*sql.DB)

	userID, success := config.CheckLogin(db, username, password)

	if success {
		session := sessions.Default(c)

		session.Set("user_id", userID)
		session.Save()

		c.Redirect(302, "/analisis")
		return
	}

	c.Redirect(302, "/login?error=Username+atau+password+salah")
}

func databaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func logout(c *gin.Context) {
    session := sessions.Default(c)

    session.Clear()
    session.Save()

    c.Redirect(302, "/login")
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

	router.GET("/login", showLogin)
	router.POST("/login", processLogin)
	router.GET("/analisis", middleware.AuthRequired(), handlers.ShowAnalisis)
	router.POST("/analisis/process", middleware.AuthRequired(), handlers.ProcessAnalysis)
	router.GET("/api/filter/tahun", middleware.AuthRequired(), handlers.GetFilterTahun)
	router.GET("/api/filter/periode", middleware.AuthRequired(), handlers.GetFilterPeriode)
	router.GET("/api/filter/institusi", middleware.AuthRequired(), handlers.GetFilterInstitusi)
	router.GET("/api/filter/batch", middleware.AuthRequired(), handlers.GetFilterBatch)
	router.GET("/analisis/result", middleware.AuthRequired(), handlers.ShowResultAnalisis)
	router.GET("/logout", logout)

	router.Run(":8080")
}