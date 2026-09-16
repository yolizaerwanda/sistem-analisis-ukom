package config

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func CheckLogin(db *sql.DB, username string, password string) (int, bool) {
	var userID int
	var passwordHash string

	err := db.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		username,
	).Scan(&userID, &passwordHash)

	if err != nil {
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	)

	if err != nil {
		return 0, false
	}

	return userID, true
}