package config

import (
	"log"
	"os"
	"strconv"
)

func GetServerPort() int{
	v := getEnvironmentValue("PORT")
	intV, _ := strconv.Atoi(v)
	return intV
}

func getEnvironmentValue(key string) string{
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing enviroment variable: %s \n", "PORT")
	}

	return value
}