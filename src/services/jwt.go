package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func JSONDumps(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshaling data to JSON: %v", err)
	}
	return string(jsonData), nil
}

func JSONLoads(jsonData string, data interface{}) error {
	err := json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data from JSON: %v", err)
	}
	return nil
}

func urlSafeBase64Encode(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.ReplaceAll(encoded, "=", "")

	return encoded
}

func urlSafeBase64Decode(s string) (string, error) {
	decoded := strings.ReplaceAll(s, "-", "+")
	decoded = strings.ReplaceAll(decoded, "_", "/")

	// Add back missing padding
	switch len(decoded) % 4 {
	case 2:
		decoded += "=="
	case 3:
		decoded += "="
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		return "", fmt.Errorf("error decoding string from base64: %v", err)
	}

	return string(decodedBytes), nil
}
