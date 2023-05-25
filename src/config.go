package main

import (
	"io/ioutil"
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
}

func openFireKey(key *string) error {

	bytes, err := ioutil.ReadFile("key/healthcheck-key.json")
	if err != nil {
		return err
	}

	*key = string(bytes)

	return nil
}
