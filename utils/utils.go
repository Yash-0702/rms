package utils

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	json "github.com/json-iterator/go"
	"golang.org/x/crypto/bcrypt"
)

func DecodeJSONBody(req io.Reader, data interface{}) error {
	return json.NewDecoder(req).Decode(data)
}

func EncodeJSONBody(res io.Writer, data interface{}) error {
	return json.NewEncoder(res).Encode(data)
}

func ResponseJSON(res http.ResponseWriter, statusCode int, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	if err := json.NewEncoder(res).Encode(data); err != nil {
		panic(err)
	}
}

func ResponseError(res http.ResponseWriter, statusCode int, err error, message string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	errMsg := "Something went wrong"
	if err != nil {
		errMsg = err.Error()
	}

	if err := json.NewEncoder(res).Encode(map[string]interface{}{
		"message": message,
		"error":   errMsg,
	}); err != nil {
		panic(err)
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func UnhashPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(email, user_id, sessionId, role string) (string, error) {
	claims := jwt.MapClaims{
		"email":     email,
		"user_id":   user_id,
		"sessionID": sessionId,
		"role":      role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
