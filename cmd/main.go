package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/ykkssyaa/Posts_Service/graph"
	"github.com/ykkssyaa/Posts_Service/internal/config"
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	"github.com/ykkssyaa/Posts_Service/internal/db"
	"github.com/ykkssyaa/Posts_Service/internal/gateway"
	in_memory "github.com/ykkssyaa/Posts_Service/internal/gateway/in-memory"
	"github.com/ykkssyaa/Posts_Service/internal/gateway/postgres"
	"github.com/ykkssyaa/Posts_Service/internal/server/graphql"
	"github.com/ykkssyaa/Posts_Service/internal/service"
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

	logger.Info.Print("Connecting to Postgres.")

	options := db.PostgresOptions{
		Name:     os.Getenv("POSTGRES_DBNAME"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	logger.Info.Print(options)
	postgresDb, err := db.NewPostgresDB(options)

	if err != nil {
		logger.Err.Fatalf(err.Error())
	}

	var gateways *gateway.Gateways

	logger.Info.Print("Creating Gateways.")
	logger.Info.Print("USE_IN_MEMORY = ", os.Getenv("USE_IN_MEMORY"))

	if os.Getenv("USE_IN_MEMORY") == "true" {
		posts := in_memory.NewPostsInMemory(consts.PostsPullSize)
		comments := in_memory.NewCommentsInMemory(consts.CommentsPullSize)
		gateways = gateway.NewGateways(posts, comments)
	} else {
		posts := postgres.NewPostsPostgres(postgresDb)
		comments := postgres.NewCommentsPostgres(postgresDb)
		gateways = gateway.NewGateways(posts, comments)
	}

	logger.Info.Print("Creating Services.")
	services := service.NewServices(gateways, logger)

	logger.Info.Print("Creating graphql server.")
	port := os.Getenv("PORT")
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graphql.Resolver{
		PostsService:      services.Posts,
		CommentsService:   services.Comments,
		CommentsObservers: graphql.NewCommentsObserver(),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger.Info.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	logger.Err.Fatal(http.ListenAndServe(":"+port, nil))

}
