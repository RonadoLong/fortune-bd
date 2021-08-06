package bcrypt2

import "golang.org/x/crypto/bcrypt"

func CryptPassword(pass string) string {
	passBcrypt, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(passBcrypt)
}
