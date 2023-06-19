package bcrypt

import "golang.org/x/crypto/bcrypt"

func PasswordHash(password string) (string, error) {
	// change password being random password by bcrypt with password salting
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(byteHash), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	//compare random password to the real password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
