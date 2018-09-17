package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jayden-chan/ctl-server/db"
)

// GenerateJWT generates a JSON Web Token for the specified customer
func GenerateJWT(userID string, access string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   userID,
		"access": access,
		"nbf":    (time.Now().Add(time.Second * 2)).Unix(),
		"iat":    time.Now().Unix(),
	})

	secretKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
	tokenString, err := token.SignedString(secretKey)

	return tokenString, err
}

// Authenticate checks to see whether the provided JWT is valid
// and that the associated customer actually exists
func Authenticate(req *http.Request) (success bool, user string, access string) {
	auth := req.Header.Get("Authorization")
	authWords := strings.Fields(auth)

	if len(authWords) != 2 || authWords[0] != "Bearer" {
		return
	}

	success, user, access = validateJWT(authWords[1])
	if !success {
		return
	}

	query := `SELECT id FROM users WHERE id = $1 AND access = $2`
	rows, err := db.Query(query, user, access)
	if err != nil {
		log.Println(err)
		return false, user, access
	}
	defer rows.Close()

	userExists := rows.Next()
	return userExists && success, user, access
}

func validateJWT(asString string) (bool, string, string) {
	token, err := jwt.Parse(asString, func(token *jwt.Token) (interface{}, error) {
		// Validate that the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Wrong signing method")
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secretKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
		return secretKey, nil
	})

	if err != nil {
		log.Println("JWT Validation error:", err)
		return false, "", ""
	}

	if token == nil {
		return false, "", ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user, claimOK := claims["user"].(string)
		if !claimOK {
			return false, "", ""
		}

		access, claimOK := claims["access"].(string)
		if !claimOK {
			return false, "", ""
		}
		return true, user, access
	}
	log.Println(err)
	return false, "", ""
}
