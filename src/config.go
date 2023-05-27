package main

import (
	c "health/clog"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var dirs = [3]string{"_key/", "etc/secrets/", ""}

// sets environment variables
// tries to read from different locations first
// if it fails, it will set its own environment variables
// note that the program will exit if no GOOGLE_APPLICATION_CREDENTIALS is set
func setConfig() {

	if _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !ok {
		var key string

		if err := openFireKey(&key); err == nil {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", key)
		} else {
			log.Println(err)
			return
		}
	}

	loadEnv("key.env")

	if _, ok := os.LookupEnv("PORT"); !ok {
		os.Setenv("PORT", "10000")
	}

	if _, ok := os.LookupEnv("LOGIN_KEY"); !ok {
		os.Setenv("LOGIN_KEY", "default_log_key")
	}

	if _, ok := os.LookupEnv("REGISTER_KEY"); !ok {
		os.Setenv("REGISTER_KEY", "default_reg_key")
	}
}

func loadEnv(file string) {
	for _, dir := range dirs {
		filepath := dir + file
		err := godotenv.Load(filepath)
		if err == nil {
			break
		}
	}
}

func openFireKey(key *string) error {
	var err error
	var bytes []byte

	for _, dir := range dirs {
		filepath := dir + "healthcheck-key.json"
		bytes, err = os.ReadFile(filepath)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	*key = string(bytes)

	return nil
}

func checkLocal() bool {

	if value := os.Getenv("LOCAL_LOGS"); value == "true" {
		return true
	}
	return false
}

// Custom logger
// Creates or opens a logs.txt file for local testing
// and uses Stderr in case of deployment.
func initLogs() {
	var err error
	var out *os.File

	if checkLocal() {
		out, err = os.OpenFile("clog/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			out = os.Stderr
		}
	} else {
		out = os.Stderr
	}

	c.InfoLog = log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	c.ErrorLog = log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	c.WarningLog = log.New(out, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	c.DebugLog = log.New(out, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	if err != nil {
		c.ErrorLog.Println(err.Error())
		c.WarningLog.Println("os.Stderr used for output")
	}
}
