package main

import (
	"github.com/ykkssyaa/Posts_Service/internal/config"
	"os"
)
import lg "github.com/ykkssyaa/Posts_Service/pkg/logger"

func main() {
	logger := lg.InitLogger()
	logger.Info.Print("Executing InitLogger.")

	envFile := ".env"
	if len(os.Args) >= 2 {
		envFile = os.Args[1]
	}
	logger.Info.Print("Executing InitConfig.")
	logger.Info.Printf("Reading %s \n", envFile)
	if err := config.InitConfig(envFile); err != nil {
		logger.Err.Fatalf(err.Error())
	}

}
