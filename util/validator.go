package util

import "github.com/go-playground/validator/v10"

func ProductNameValidator(fl validator.FieldLevel) bool {
	productName := fl.Field().String()
	allowedProducts := []string{"WhiteGods", "motor", "mobil"}
	for _, product := range allowedProducts {
		if productName == product {
			return true
		}
	}
	return false
}
