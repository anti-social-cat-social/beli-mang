package hasher

import (
	"os"
	"strconv"
	// "errors"

	"golang.org/x/crypto/bcrypt"
)

// Get hashing cost using environtment variable
func getCost() (int, error) {
	envCost := os.Getenv("BCRYPT_COST")

	// Handle if the BCRYPT_COST is not a valid integer
	cost, err := strconv.Atoi(envCost)
	if err != nil {
		return 0, nil
	}

	return cost, nil
}

// Generate hashed password using password string
func HashPassword(password string) (string, error) {
	// Cost from environment
	cost, err := getCost()
	if err != nil {
		return "", err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// Check hashed password
func CheckPassword(hashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
