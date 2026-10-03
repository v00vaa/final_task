package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func checkToken(jwtToken, password string) bool {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(password), nil
	})
	if err != nil || token == nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hash := sha256.Sum256([]byte(password))
	expectedHash := fmt.Sprintf("%x", hash)
	return claims["hash"] == expectedHash
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}
			var valid bool
			valid = checkToken(jwtToken, pass)
			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func generateToken(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"hash": fmt.Sprintf("%x", hash),
		})

	return token.SignedString([]byte(password))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	systemPassword := os.Getenv("TODO_PASSWORD")
	if request.Password != systemPassword {
		writeJSON(w, map[string]string{
			"error": "Неверный пароль",
		})
		return
	}
	token, err := generateToken(systemPassword)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, map[string]string{
		"token": token,
	})
}
