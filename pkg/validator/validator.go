package validator

import (
	"github.com/asaskevich/govalidator"
)

var shortUrlPattern string

func ValidatorInit(pattern string) {
	govalidator.SetFieldsRequiredByDefault(true)
	shortUrlPattern = pattern
}

func IsUrl(url string) bool {
	return govalidator.IsURL(url)
}

func IsShortUrl(shortUrl string) bool {
	return govalidator.Matches(shortUrl, shortUrlPattern)
}