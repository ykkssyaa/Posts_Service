package config

import (
	"github.com/joho/godotenv"
)

func InitConfig(file string) error {
	err := godotenv.Load(file)

	return err
}
