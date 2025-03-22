package utils

import "os"

func Response(statusCode int, message any, data any) map[string]any {
	var status string
	switch {
	case statusCode >= 200 && statusCode <= 299:
		status = "success"
	case statusCode >= 300 && statusCode <= 399:
		status = "redirect"
	case statusCode == 400:
		status = "error"
	case statusCode == 404:
		status = "not found"
	case statusCode >= 405 && statusCode <= 499:
		status = "error"
	case statusCode == 401 || statusCode == 403:
		status = "unauthorized"
	case statusCode >= 500:
		status = "error"
		message = "This is from us!, please contact admin"
	default:
		status = "error"
		message = "This is from us!, please contact admin"
	}
	res := map[string]any{
		"status":     status,
		"message":    message,
		"data":       data,
		"statusCode": statusCode,
	}

	return res

}

func GetEnv(key, fallback string) string {
	// get value from env
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
