package crypto

import "golang.org/x/crypto/bcrypt"

var cost = 15

// CreatePassword function creates a password to store in the database
func CreatePassword(plainTextPassword string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), cost)

	if err != nil {
		return "", err
	}

	return string(password), nil
}

// ValidatePassword function validates a password against encrypted version from database
func ValidatePassword(plainTextPassword string, encryptedPassword string) bool {
	bPlainTextPassword := []byte(plainTextPassword)
	bEncryptedPassword := []byte(encryptedPassword)

	err := bcrypt.CompareHashAndPassword(bEncryptedPassword, bPlainTextPassword)

	return err == nil
}
