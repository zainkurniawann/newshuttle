package main

import (
	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"shuttle/databases"
)

type Role string

const (
	SuperAdmin Role = "superadmin"
)

type User struct {
	ID       int    `json:"id" db:"user_id"`
	UUID     string `json:"uuid" db:"user_uuid"`
	Username string `json:"username" db:"user_username"`
	Email    string `json:"email" db:"user_email"`
	Password string `json:"password" db:"user_password"`
	Role     Role   `json:"role" db:"user_role"`
	RoleCode string `json:"role_code" db:"user_role_code"`
}

type SuperAdminDetails struct {
	UUID      string `json:"uuid" db:"user_id"`
	Picture   string `json:"picture" db:"user_picture"`
	FirstName string `json:"first_name" db:"user_first_name"`
	LastName  string `json:"last_name" db:"user_last_name"`
	Gender    string `json:"user_gender" db:"user_gender"`
	Phone     string `json:"phone" db:"user_phone"`
	Address   string `json:"address" db:"user_address"`
}

func main() {
	db, err := databases.PostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE user_role = $1", SuperAdmin)
	if err != nil {
		log.Fatal("Failed to count users:", err)
		return
	}

	if count > 0 {
		color.Yellow("Superadmin already exists, no need to seed.")
		return
	}

	hashedPassword, err := hashPassword("12345678")
	if err != nil {
		log.Fatal("Error hashing password:", err)
		return
	}

	user := User{
		ID:       0,
		UUID:     "00000000-0000-0000-0000-000000000000",
		Username: "",
		Email:    "faker@gmail.com",
		Password: hashedPassword,
		Role:     SuperAdmin,
		RoleCode: "SA",
	}

	details := SuperAdminDetails{
		UUID:      user.UUID,
		Picture:   "",
		FirstName: "",
		LastName:  "",
		Gender:    "",
		Phone:     "",
		Address:   "",
	}

	_, err = db.Exec(`
		INSERT INTO users (user_id, user_uuid, user_username, user_email, user_password, user_role, user_role_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.UUID, user.Username, user.Email, user.Password, user.Role, user.RoleCode)
	if err != nil {
		log.Fatal("Failed to insert superadmin user:", err)
		return
	}

	_, err = db.Exec(`
		INSERT INTO super_admin_details (user_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, details.UUID, details.Picture, details.FirstName, details.LastName, details.Gender, details.Phone, details.Address)
	if err != nil {
		log.Fatal("Failed to insert superadmin details:", err)
		return
	}

	color.Green("Users seeded successfully!")
	color.Yellow("Please login and create a new user with superadmin role immediately and do hard delete this user.")
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
