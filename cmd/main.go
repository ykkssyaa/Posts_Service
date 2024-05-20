package main

import (
	"github.com/ykkssyaa/Posts_Service/internal/config"
)
import lg "github.com/ykkssyaa/Posts_Service/pkg/logger"

func main() {
	logger := lg.InitLogger()

	logger.Info.Print("Executing InitConfig.")
	if err := config.InitConfig(); err != nil {
		logger.Err.Fatalf(err.Error())
	}

}
