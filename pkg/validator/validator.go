package validator

import (
	"github.com/asaskevich/govalidator"
)

var short_url_pattern string

func ValidatorInit(pattern string) {
	govalidator.SetFieldsRequiredByDefault(true)
	short_url_pattern = pattern
}

func IsUrl(url string) bool {
	return govalidator.IsURL(url)
}

func IsShortUrl(shortUrl string) bool {
	return govalidator.Matches(shortUrl, short_url_pattern)
}