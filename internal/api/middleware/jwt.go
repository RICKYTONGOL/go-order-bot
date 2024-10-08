package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("dotfundfd3gadfbvnuujy1f000")
var adminHmacSampleSecret = []byte("adminDotfund0090SecretKey")
var ttl int64 = 3600 * 24 // 1 day

type Payload struct {
	UserId    string
	Email     string
	Phone     string
	ExpiredAt int64
}

type AdminPayload struct {
	UserId    string
	Email     string
	Role      string
	ExpiredAt int64
}

func GenerateJWT(userId string, email string, phone string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"email":     email,
		"phone":     phone,
		"expiredAt": time.Now().Unix() + ttl,
	})
	tokenString, e := token.SignedString(hmacSampleSecret)
	if e != nil {
		log.Print("error", e)
	}
	return tokenString
}

func ParseJWT(tokenString string) (payload Payload, err error) {
	if len(tokenString) < 7 {
		return payload, fmt.Errorf("token is not valid: %s", tokenString)
	}

	if tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload.UserId = claims["userId"].(string)
		payload.Email = claims["email"].(string)
		payload.Phone = claims["phone"].(string)
		payload.ExpiredAt = int64(claims["expiredAt"].(float64))
	} else {
		fmt.Println(err)
	}

	return
}

func GenerateAdminJWT(userId string, email string, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"email":     email,
		"role":      role,
		"expiredAt": time.Now().Unix() + ttl,
	})
	tokenString, e := token.SignedString(adminHmacSampleSecret)
	if e != nil {
		log.Print("error", e)
	}
	return tokenString
}

func ParseAdminJWT(tokenString string) (payload AdminPayload, err error) {
	if len(tokenString) < 7 {
		return payload, fmt.Errorf("token is not valid: %s", tokenString)
	}

	if tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return adminHmacSampleSecret, nil
	})
	if err != nil {
		return payload, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload.UserId = claims["userId"].(string)
		payload.Email = claims["email"].(string)
		payload.Role = claims["role"].(string)
		payload.ExpiredAt = int64(claims["expiredAt"].(float64))
	} else {
		fmt.Println(err)
	}
	return
}
