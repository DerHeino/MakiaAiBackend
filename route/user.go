package route

import (
	"crypto/sha1"
	"errors"
	"fmt"
	bg "health/background"
	c "health/clog"
	"health/model"
	"health/network"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

// Retrieves user from database and verifies correct password
// returns JWT-Token and nil error if successful
func VerifyLogin(userMap map[string]interface{}) (string, error) {
	var userCredential model.UserCredentials
	var retrieved model.UserCredentials

	if missing := validateParameters(model.CredentialParameters, userMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters %s\n", m)
		return "", fmt.Errorf("missing required parameter(s): %s", m)
	}

	if err := mapstructure.Decode(userMap, &userCredential); err != nil {
		c.ErrorLog.Printf(err.Error())
		return "", errors.New("failed: decoding error")
	}

	if err := network.GetUser(hashUsername(userCredential.Username), &retrieved); err != nil {
		return "", errors.New("failed to retrieve user credentials")
	}

	if pwderr := bcrypt.CompareHashAndPassword([]byte(retrieved.Password), []byte(userCredential.Password)); pwderr != nil {
		c.InfoLog.Println(pwderr.Error())
		return "", errors.New("invalid password")
	}

	return buildToken(retrieved.Username, os.Getenv("LOGIN_KEY"))
}

func buildToken(username string, secret string) (string, error) {
	var mySecret []byte = []byte(secret)

	type MyCustomClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	claim := MyCustomClaims{
		username,
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return "", errors.New("failed to build token")
	}

	return tokenString, nil
}

// Takes "Authorization" string from the HTTP header and
// evaluates its parsing ("Bearer" <Token>)
// Takes into account if admin rights are required for POST /project and /location
// calls ValidateToken from route/user.go
//
// returns a string and an error for ResponseWriter and error logging
// if the token is invalid in any way
// returns an empty string and error if token is valid

// needed for POST project and location
// the only reason this method is still here is for future admin/non-admin differentiation
func ValidateToken(auth string, admin bool, secret ...string) (string, error) {
	var mySecret []byte

	tokenString, err := validateTokenForm(auth)
	if err != nil {
		return "", errors.New("no Token found")
	}

	if len(secret) == 0 {
		mySecret = []byte(os.Getenv("LOGIN_KEY"))
	} else {
		mySecret = []byte(secret[0])
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySecret, nil
	})
	if err != nil {
		c.ErrorLog.Println(err.Error())
		return "", fmt.Errorf("token invalid: %s", tokenString)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := hashUsername(claims["username"].(string))
		if len(secret) > 0 {
			if !bg.GetRegMap().Exists(id) {
				return "", fmt.Errorf("token expired: %s", tokenString)
			}
		}
		return claims["username"].(string), network.UserValid(id, admin)
	}

	return "", fmt.Errorf("token invalid: %s", tokenString)
}

func validateTokenForm(auth string) (string, error) {
	parts := strings.Split(auth, "Bearer")
	if len(parts) != 2 {
		return "", errors.New("no Token found")
	}
	return strings.TrimSpace(parts[1]), nil
}

func hashUsername(username string) string {
	hash := fmt.Sprintf("%X", sha1.Sum([]byte(username)))
	return hash[3:27]
}

func BuildRegisterKey(username string) (string, error) {

	registerToken, err := buildToken(username, os.Getenv("REGISTER_KEY"))
	if err != nil {
		return "", err
	}

	regMap := bg.GetRegMap()
	if regMap.Exists(hashUsername(username)) {
		regMap.Update(hashUsername(username), registerToken)
	} else {
		regMap.Add(hashUsername(username), registerToken)
	}

	return registerToken, nil
}

func VerifyUser(inviter string, userMap map[string]interface{}) (string, error) {
	var userCredential model.UserCredentials

	if missing := validateParameters(model.UserParameters, userMap); len(missing) != 0 {
		m := strings.Join(missing, ", ")
		c.ErrorLog.Printf("missing parameters %s\n", m)
		return "", fmt.Errorf("missing required parameter(s): %s", m)
	}

	if err := mapstructure.Decode(userMap, &userCredential); err != nil {
		c.ErrorLog.Printf(err.Error())
		return "", errors.New("failed to register user: decoding error")
	}

	userCredential.Password = hashPassword(userCredential.Password)
	userCredential.Id = hashUsername(userCredential.Username)

	token, err := postUser(&userCredential, inviter)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %s", err.Error())
	}

	return token, nil
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		c.ErrorLog.Println(err)
		return password
	}
	return string(hash)
}

func postUser(user *model.UserCredentials, inviter string) (string, error) {

	err := network.SetUserFire(user)
	if err != nil {
		return "", fmt.Errorf("failed to upload user: %s", err)
	}
	regMap := bg.GetRegMap()
	regMap.Delete(inviter)
	return buildToken(user.Username, os.Getenv("LOGIN_KEY"))
}
