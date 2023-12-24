package utils

import (
	"os"
	"strconv"
)

func EnvInt(name string, defaultValue int) (value int) {
	value, _ = strconv.Atoi(os.Getenv(name))
	if value == 0 {
		value = defaultValue
	}

	return
}

func EnvStr(name, defaultValue string) (value string) {
	value = os.Getenv(name)
	if value == "" {
		value = defaultValue
	}

	return
}
