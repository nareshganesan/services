package globals

import (
	"fmt"
	"time"
	// "strconv"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// GenerateJWT returns JWT token for the string
// stores 1. issuer (doamin), uid, expiration and iat (time of issue) in the token
func GenerateJWT(uid string) string {
	l := Gbl.Log
	if uid == "" {
		l.WithFields(logrus.Fields{
			"uid": uid,
		}).Error("uid should not be empty!")
		return ""
	}
	jwtSecret := Config.Tokens.Auth.Secret
	issuer := Config.Owner.Domain
	appSigningKey := []byte(jwtSecret)

	// Ref: https://github.com/dgrijalva/jwt-go/blob/master/claims.go
	// Create the Claims
	claims := make(jwt.MapClaims)
	claims["issuer"] = issuer
	claims["uid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(Config.Tokens.Auth.Maxage*24)).Unix()
	claims["iat"] = time.Now().Unix()
	// jwt.GetSigningMethod("HS512")
	token := jwt.NewWithClaims(jwt.GetSigningMethod(Config.Tokens.Auth.Algorithm), claims)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwt, err := token.SignedString(appSigningKey)
	fmt.Printf("%v %v", jwt, err)
	if err != nil {
		l.Error("Error generating token")
		return ""
	}
	l.WithFields(logrus.Fields{
		"jwt": jwt,
	}).Info("Token generated")
	return jwt
}

// ParseJWT checks if JWT token is valid
func ParseJWT(tokenString string) *jwt.MapClaims {
	l := Gbl.Log
	claims, status := isValidToken(tokenString)
	if !status {
		l.WithFields(logrus.Fields{
			"jwt": tokenString,
		}).Error("Error parsing token")
		return nil
	}
	l.WithFields(logrus.Fields{
		"claims": claims,
	}).Info("Token parsed successfully")
	return claims
}

// RefreshToken returns a new JWT token given a valid JWT token
func RefreshToken(tokenString string) string {
	l := Gbl.Log
	claims, status := isValidToken(tokenString)
	if !status {
		l.WithFields(logrus.Fields{
			"existing_jwt": tokenString,
		}).Error("Could not refresh token")
		return ""
	}
	uid := (*claims)["uid"].(string)
	token := GenerateJWT(uid)
	if token != "" {
		l.WithFields(logrus.Fields{
			"jwt": token,
		}).Info("Refresh token generated")
		return token
	}
	l.WithFields(logrus.Fields{
		"existing_jwt": tokenString,
	}).Error("Error generating refresh token")
	return ""
}

// Helper for checking if a token is valid
func isValidToken(tokenString string) (*jwt.MapClaims, bool) {
	l := Gbl.Log
	jwtSecret := Config.Tokens.Auth.Secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		l.WithFields(logrus.Fields{
			"jwt": tokenString,
		}).Error("Cannot parse token")
		return nil, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		l.WithFields(logrus.Fields{
			"jwt": tokenString,
		}).Info("token is valid")
		return &claims, true
	}
	l.WithFields(logrus.Fields{
		"jwt": tokenString,
	}).Error("Error parsing token claim")
	return nil, false
}
