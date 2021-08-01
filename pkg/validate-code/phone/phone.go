package phone

import "regexp"

func CheckPhone(phone string) bool {
	regular := "^1([358][0-9]|4[579]|66|7[0135678]|9[89])[0-9]{8}$"
	// CheckMobileNum 手机号码的验证
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}
