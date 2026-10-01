package validator

import (
	"net/url"
	"regexp"
)

var shortURLPattern *regexp.Regexp

func ValidatorInit(pattern string) {
	shortURLPattern = regexp.MustCompile(pattern)
}

func IsUrl(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func IsShortUrl(shortURL string) bool {
	return shortURLPattern.MatchString(shortURL)
}
