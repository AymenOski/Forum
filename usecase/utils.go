package usecase

import "regexp"

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	res, _ := regexp.MatchString(emailRegex, email)
	return res
}

func isValidName(name string) bool {
	nameRegex := `^[a-zA-Z0-9]+$`
	res, _ := regexp.MatchString(nameRegex, name)
	return res
}
