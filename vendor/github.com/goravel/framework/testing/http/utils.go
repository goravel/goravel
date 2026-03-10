package http

import "net/http"

func Cookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
	}
}

func Cookies(data map[string]string) []*http.Cookie {
	cookies := make([]*http.Cookie, 0, len(data))
	for name, value := range data {
		cookies = append(cookies, Cookie(name, value))
	}
	return cookies
}
