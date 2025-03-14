package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rk-the-dev/golib-core/pkg/logger" // Import our common logger
	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt creates a random salt of the specified length
func GenerateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		logger.Error("Failed to generate salt", map[string]interface{}{"error": err})
		return "", err
	}
	logger.Debug("Successfully generated salt", nil)
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashPassword hashes a password using bcrypt with a given salt
func HashPassword(password, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", map[string]interface{}{"error": err})
		return "", err
	}
	logger.Debug("Password hashed successfully", nil)
	return string(hashedPassword), nil
}

// ComparePassword checks if a given password matches the stored hash
func ComparePassword(hashedPassword, password, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	if err != nil {
		logger.Warn("Password comparison failed", map[string]interface{}{"error": err})
		return false
	}
	logger.Debug("Password matched successfully", nil)
	return true
}

// HashSHA256 hashes a string using SHA-256
func HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	logger.Debug("SHA-256 hash generated successfully", nil)
	return base64.URLEncoding.EncodeToString(hash[:])
}

// EncryptData encrypts plaintext using AES-GCM encryption
func EncryptData(plaintext, encryptionKey string) (string, error) {
	keyHash := sha256.Sum256([]byte(encryptionKey))
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		logger.Error("Failed to create AES cipher", map[string]interface{}{"error": err})
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("Failed to create AES-GCM mode", map[string]interface{}{"error": err})
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error("Failed to generate nonce", map[string]interface{}{"error": err})
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	logger.Debug("Data encrypted successfully", nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptData decrypts AES-GCM encrypted data
func DecryptData(encryptedData, encryptionKey string) (string, error) {
	keyHash := sha256.Sum256([]byte(encryptionKey))
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		logger.Error("Failed to create AES cipher", map[string]interface{}{"error": err})
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("Failed to create AES-GCM mode", map[string]interface{}{"error": err})
		return "", err
	}

	data, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		logger.Error("Failed to decode base64 encrypted data", map[string]interface{}{"error": err})
		return "", err
	}

	if len(data) < aesGCM.NonceSize() {
		logger.Error("Invalid ciphertext length", nil)
		return "", errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:aesGCM.NonceSize()], data[aesGCM.NonceSize():]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Error("Failed to decrypt data", map[string]interface{}{"error": err})
		return "", err
	}

	logger.Debug("Data decrypted successfully", nil)
	return string(plaintext), nil
}

// GenerateJWT generates a new JWT token with claims and expiry
func GenerateJWT(secretKey string, claims jwt.MapClaims, expiry time.Duration) (string, error) {
	claims["exp"] = time.Now().Add(expiry).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Error("Failed to generate JWT", map[string]interface{}{"error": err})
		return "", err
	}
	logger.Debug("JWT generated successfully", nil)
	return signedToken, nil
}

// VerifyJWT verifies a JWT token using the provided secret key
func VerifyJWT(secretKey, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warn("Invalid JWT signing method", nil)
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Error("JWT verification failed", map[string]interface{}{"error": err})
		return nil, err
	}
	logger.Debug("JWT verified successfully", nil)
	return token, nil
}
