package route

import (
	"errors"
	"fmt"
	"health/model"
	"log"
	"strings"

	"health/network"
	"time"

	"github.com/golang-jwt/jwt"
	mapstructure "github.com/mitchellh/mapstructure"

	"golang.org/x/crypto/bcrypt"
)

// It wouldn't be a bad idea to also store admin Id's in a local array
var userCredential model.UserCredentials

// Retrieves user from database and verifies correct password
// returns JWT-Token and nil error if successful
func VerifyLogin(userMap map[string]interface{}) (string, error) {
	defer clearModel(&userCredential)
	var retrieved model.UserCredentials

	if err := validateCredentials(userMap); err != nil {
		return "", err
	}

	if err := mapstructure.Decode(userMap, &userCredential); err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to convert user credentials")
	}

	if err := network.GetUser(userCredential.Username, &retrieved); err != nil {
		log.Println(err.Error())
		return "", errors.New("failed to retrieve user credentials")
	}

	if pwderr := bcrypt.CompareHashAndPassword([]byte(retrieved.Password), []byte(userCredential.Password)); pwderr != nil {
		log.Println(pwderr.Error())
		return "", errors.New("invalid password")
	}

	tokenString, tokenerr := BuildToken(retrieved.Username)
	return tokenString, tokenerr
}

func validateCredentials(userMap map[string]interface{}) error {

	missing := CountParameters(model.CredentialParameters, userMap)

	if len(missing) > 0 {
		return errors.New("missing required parameter(s): " + strings.Join(missing, ", "))
	}

	return nil
}

// needed for register
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		fmt.Println(err)
		return password
	}
	return string(hash)
}

func BuildToken(username string) (string, error) {
	var mySecret []byte = []byte("21062022") // for now

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
		log.Println(err.Error())
		return "", errors.New("failed to build new token")
	}

	return tokenString, nil
}

// needed for POST project and location
// the only reason this method is still here is for future admin/non-admin differentiation
func ValidateToken(tokenString string, admin bool) (bool, error) {
	var mySecret []byte = []byte("21062022") // for now

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if admin {
			return isAdmin(claims), err
		} else {
			return true, err
		}
	} else {
		fmt.Println(err.Error())
		return false, err
	}
}

func isAdmin(claims map[string]interface{}) bool {
	fmt.Println(claims)

	username := claims["username"].(string)

	for _, admin := range network.AdminList {
		if admin == username {
			return true
		}
	}

	return false
}
