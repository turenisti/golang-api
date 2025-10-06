package utils

import (
	"encoding/json"
	"strings"
)

// ParseBodyToJSON parses JSON body and returns as map for structured logging
func ParseBodyToJSON(body []byte, contentType string) interface{} {
	if len(body) == 0 {
		return nil
	}

	// Only parse if content type is JSON
	if !strings.Contains(strings.ToLower(contentType), "application/json") {
		return string(body) // Return raw string for non-JSON
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// If parsing fails, return raw string
		return string(body)
	}

	return result
}

// Sensitive fields for masking in scheduling-report-api
var sensitiveFields = map[string]bool{
	"password":        true,
	"token":           true,
	"authorization":   true,
	"secret":          true,
	"key":             true,
	"connection_url":  true, // Database connection strings
	"connection_config": true, // Database credentials
	"delivery_config": true,   // Delivery credentials (SMTP, SFTP, S3)
	"smtp_password":   true,
	"sftp_password":   true,
	"api_key":         true,
	"access_key":      true,
	"secret_key":      true,
}

// MaskSensitiveData masks sensitive fields in the parsed JSON
func MaskSensitiveData(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			isSensitive := sensitiveFields[strings.ToLower(key)]

			if isSensitive {
				if str, ok := value.(string); ok && len(str) > 0 {
					// Mask string values - show first 2 and last 2 chars
					if len(str) <= 4 {
						result[key] = "****"
					} else {
						result[key] = str[:2] + strings.Repeat("*", len(str)-4) + str[len(str)-2:]
					}
				} else {
					result[key] = "****"
				}
			} else {
				// Recursively process nested objects
				result[key] = MaskSensitiveData(value)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = MaskSensitiveData(item)
		}
		return result
	default:
		return v
	}
}

// FormatHeaders formats HTTP headers for logging, masking sensitive ones
func FormatHeaders(headers map[string][]string) map[string]string {
	result := make(map[string]string)
	sensitiveHeaders := []string{"authorization", "cookie", "x-api-key", "token"}

	for key, values := range headers {
		isSensitive := false
		for _, sensitive := range sensitiveHeaders {
			if strings.Contains(strings.ToLower(key), sensitive) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			result[key] = "****"
		} else {
			result[key] = strings.Join(values, ", ")
		}
	}

	return result
}

// ParseQueryParams extracts and formats query parameters for logging
func ParseQueryParams(queryString string) map[string]string {
	if queryString == "" {
		return nil
	}

	params := make(map[string]string)
	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		if parts := strings.SplitN(pair, "=", 2); len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}

	return params
}
