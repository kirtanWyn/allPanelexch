package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// GenerateRandomString generates a random alphanumeric string of a given length
func GenerateRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// HashPassword hashes the password exactly like the PHP application: hash('sha256', $password . $saltKey . $salt)
func HashPassword(password, saltKey, salt string) string {
	hashStr := password + saltKey + salt
	hasher := sha256.New()
	hasher.Write([]byte(hashStr))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateTelegramOTP triggers the external Telegram OTP generation API
func GenerateTelegramOTP(userID int) error {
	authCodeURL := os.Getenv("AUTH_CODE_URL")
	if authCodeURL == "" {
		return errors.New("AUTH_CODE_URL not set in environment")
	}

	url := authCodeURL + "verification/telegram_generate_otp"
	
	payload := map[string]int{
		"user_id": userID,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	
	// Assuming the API returns {"status": true} on success
	if status, ok := result["status"].(bool); !ok || !status {
		msg := "Telegram OTP generation failed"
		if message, ok := result["message"].(string); ok {
			msg = message
		}
		return errors.New(msg)
	}
	
	return nil
}
