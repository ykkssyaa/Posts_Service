package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ykkssyaa/Posts_Service/graph"
	"github.com/ykkssyaa/Posts_Service/internal/config"
	"github.com/ykkssyaa/Posts_Service/internal/server/graphql"
	"net/http"
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

	port := os.Getenv("PORT")

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graphql.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger.Info.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	logger.Err.Fatal(http.ListenAndServe(":"+port, nil))

}
