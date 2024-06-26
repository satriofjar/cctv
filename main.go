package main

import (
	"cctv/models"
	"fmt"
	"os"

	"golang.org/x/term"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Periksa jumlah argumen baris perintah
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./example <input>")
		// Mengembalikan kesalahan jika jumlah argumen tidak tepat
		fmt.Println("Error: Expected exactly one argument")
		return
	}

	input := os.Args[1]
	if input != "createuser" {
		fmt.Println("Unknown argument")
		return
	}

	// Connect to database
	models.ConnectDatabase()

	createUser()

}

func createUser() {
	// Take input username, password, and confirim password from user
	var (
		userName, pwd, cPwd string
	)

	fmt.Print("Username : ")
	fmt.Scan(&userName)

	for true {
		fmt.Print("Password: ")

		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Failed to read password:", err)
			return
		}

		fmt.Print("\nConfirm Password: ")

		confirmPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Failed to read password:", err)
			return
		}

		fmt.Println()

		pwd = string(password)
		cPwd = string(confirmPassword)
		// Check password and confirmPassword is match or not
		if pwd != cPwd {
			fmt.Println("\nPassword and Confirm password do not match")
			continue
		}
		break
	}

	// Check is user already exist
	var existingUser models.User

	if err := models.DB.Where("user_name = ?", userName).First(&existingUser).Error; err == nil {
		fmt.Println("\nUsername already exists")
		return
	}

	// create hash from password
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		return
	}

	// Create the user in the database
	user := models.User{
		UserName:     userName,
		HashPassword: string(hash),
	}

	if err := models.DB.Create(&user).Error; err != nil {
		fmt.Println("Failed to create user")
		return
	}

	fmt.Println("User successfully created!")
	return
}
