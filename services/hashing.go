package services

import(
	"golang.org/x/crypto/bcrypt"

)

func HashPassword(password string) (string, error) {
    // GenerateFromPassword returns a hash of the password using bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
    if err != nil {
        return "", err
    }
    // convert the hash to a string and return
    return string(hash), nil
}

func ComparePassword(verifyUserLoginPassword,userLoginPassword string)error{
    err := bcrypt.CompareHashAndPassword([]byte(verifyUserLoginPassword), []byte(userLoginPassword))
    return err
}