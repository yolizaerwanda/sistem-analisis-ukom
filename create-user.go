package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sistem-analisis-ukom/config"
)

func createUser() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		fmt.Println("Username dan password tidak boleh kosong")
		return
	}

	hash, err := config.HashPassword(password)

	if err != nil {
		fmt.Println("Gagal membuat password hash:", err)
		return
	}

	db, err := config.ConnectDatabase()

	if err != nil {
		fmt.Println("Gagal terhubung ke database:", err)
		return
	}

	defer db.Close()

	_, err = db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		username,
		string(hash),
	)

	if err != nil {
		fmt.Println("Gagal membuat user:", err)
		return
	}

	fmt.Println("User berhasil dibuat.")
}