package config

import (
	"os"
	"strconv"
)

const defaultSessionMaxAge = 86400 * 30

func GetSessionMaxAge() int {
	if maxAge, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE")); err == nil {
		return maxAge
	}
	return defaultSessionMaxAge
}
