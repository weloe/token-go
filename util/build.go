package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	// set version number (4)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// set variant bits (2)
	uuid[8] = (uuid[8] & 0xbf) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GenerateSimpleUUID() (string, error) {
	uuid, err := GenerateUUID()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(uuid, "-", ""), nil
}

func GenerateRandomString32() (string, error) {
	data := make([]byte, 24)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data)[:32], nil
}

func GenerateRandomString64() (string, error) {
	data := make([]byte, 48)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data)[:64], nil
}

func GenerateRandomString128() (string, error) {
	data := make([]byte, 96)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data)[:128], nil
}
