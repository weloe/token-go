package util

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func SendGetRequest(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get() failed: %v", err)
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("read response body failed: %v", err)
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll() failed: %v", err)
		return "", err
	}

	return string(body), nil
}

// SpliceUrl splice two url.
// Examples:
// u1 = "http://domain.com" u2 = "/sso/auth" return http://domain.com/sso/auth
func SpliceUrl(u1 string, u2 string) string {
	if u1 == "" {
		return u2
	}
	if u2 == "" {
		return u1
	}

	if strings.HasPrefix(u2, "http") {
		return u2
	}

	return u1 + u2
}

func HasUrl(urls []string, url string) bool {
	for _, s := range urls {
		if MatchUrl(s, url) {
			return true
		}
	}
	return false
}

func MatchUrl(pattern string, url string) bool {
	if pattern == "*" {
		return true
	}
	return pattern == url
}

func IsValidUrl(u1 string) bool {

	_, err := url.ParseRequestURI(u1)
	if err != nil {
		return false
	}

	u, err := url.Parse(u1)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	// check if the URL has a valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}

// AddQueryMap add map param for the path.
func AddQueryMap(path string, paramMap map[string]string) string {
	queryString := MapToQuery(paramMap)

	return AddQueryValue(path, queryString)
}

func AddQueryValue(path string, queryString string) string {
	index := strings.LastIndex(path, "?")
	// if the path is not included "?"
	if index == -1 {
		return path + "?" + queryString
	}
	// if the last is "?"
	if index == len(path)-1 {
		return path + queryString
	}

	// if "?" inside path, the last is not "&" and queryString's first string is not "&"
	if index < len(path)-1 {
		if strings.LastIndex(path, "&") != len(path)-1 && strings.Index(path, "&") != 0 {
			return path + "&" + queryString
		} else {
			return path + queryString
		}
	}

	return path
}

// AddQuery add query param for the path.
func AddQuery(path string, key string, value string) string {
	queryString := key + "=" + value
	return AddQueryValue(path, queryString)
}

// MapToQuery convert map to k=v array, and use "&" to join.
func MapToQuery(paramMap map[string]string) string {
	var queryString []string
	for k, v := range paramMap {
		queryString = append(queryString, k+"="+v)
	}
	query := strings.Join(queryString, "&")
	return query
}

func Encode(u string) string {
	return url.QueryEscape(u)
}
