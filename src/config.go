package main

import (
	c "health/clog"
	"log"
	"os"
)

// this file exists mainly to support both local and remote deployment
// since heroku does not support specific files
func setConfig() {

	if _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !ok {
		var key string

		if err := openFireKey(&key); err != nil {
			log.Println(err)
		} else {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", key)
		}
	}

	if _, ok := os.LookupEnv("PORT"); !ok {
		os.Setenv("PORT", "10000")
	}

	if _, ok := os.LookupEnv("LOGIN_KEY"); !ok {
		os.Setenv("LOGIN_KEY", "21062022")
	}

	if _, ok := os.LookupEnv("REGISTER_KEY"); !ok {
		os.Setenv("REGISTER_KEY", "secret_register_key")
	}
}

func openFireKey(key *string) error {

	bytes, err := os.ReadFile("key/healthcheck-key.json")
	if err != nil {
		return err
	}

	*key = string(bytes)

	return nil
}

func checkLocal() bool {

	if _, err := os.Stat("key/local-logs"); err == nil {
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
