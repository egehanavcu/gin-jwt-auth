package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gin-jwt-auth/src/config"
	"strings"
	"time"
)

func GenerateJWT(data interface{}) (string, error) {
	header, header_err := GenerateJWTPart(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})

	if header_err != nil {
		return "", header_err
	}

	payload, payload_err := GenerateJWTPart(data)
	if payload_err != nil {
		return "", payload_err
	}
	signature := SignJWT(header, payload)
	return fmt.Sprintf("%s.%s.%s", header, payload, signature), nil
}

func GenerateJWTPart(data interface{}) (string, error) {
	jsonData, err := JSONDumps(data)
	if err != nil {
		return "", fmt.Errorf("error encoding data to JSON: %v", err)
	}
	return URLSafeBase64Encode(jsonData), nil
}

func SignJWT(header, payload string) string {
	secretKey := []byte(config.Get("JWTSecret").(string))
	message := []byte(header + "." + payload)
	h := hmac.New(sha256.New, secretKey)
	h.Write(message)
	signature := h.Sum(nil)
	return URLSafeBase64Encode(string(signature))
}

func HandleJWT(token string, allowedRoles ...string) (interface{}, error) {
	data, err := ParseJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	exp, ok := data["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if float64(time.Now().Unix()) > exp {
		return data, fmt.Errorf("token has expired")
	}

	if allowedRoles != nil {
		if !isRoleAllowed(data["role"].(string), allowedRoles) {
			return nil, fmt.Errorf("unauthorized")
		}
	}

	return data, nil
}

func ParseJWT(jwt string) (map[string]any, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}
	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	// Check signature
	if SignJWT(header, payload) != signature {
		return nil, fmt.Errorf("invalid JWT signature")
	}

	// Decode payload
	var data map[string]any
	base64DecodedPayload, err := URLSafeBase64Decode(payload)
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %v", err)
	}

	err = JSONLoads(base64DecodedPayload, &data)
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %v", err)
	}

	return data, nil
}

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

func URLSafeBase64Encode(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.ReplaceAll(encoded, "=", "")

	return encoded
}

func URLSafeBase64Decode(s string) (string, error) {
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

func isRoleAllowed(userRole string, allowedRoles []string) bool {
	if len(allowedRoles) == 0 {
		return true
	}

	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}

	return false
}
