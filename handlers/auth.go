package handlers

import (
	"database/sql"
	"sistem-analisis-ukom/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	errorMessage := c.Query("error")

	c.HTML(200, "login.html", gin.H{
		"error": errorMessage,
	})
}

func ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" && password == "" {
		c.Redirect(302, "/login?error=Username+dan+password+tidak+boleh+kosong")
		return
	}

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

func Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	session.Save()

	c.Redirect(302, "/login")
}