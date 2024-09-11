package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cchalop1.com/deploy/internal/adapter"
	"cchalop1.com/deploy/internal/api/graph"
	"cchalop1.com/deploy/internal/api/service"
	"cchalop1.com/deploy/internal/application"
	"cchalop1.com/deploy/internal/web"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

var flags struct {
	noBrowser bool
	help      bool
	redeploy  struct {
		deployId string
	}
}

func main() {
	// app := api.NewApplication()

	databaseAdapter := adapter.NewDatabaseAdapter()
	filesystemAdapter := adapter.NewFilesystemAdapter()
	dockerAdapter := adapter.NewDockerAdapter()

	databaseAdapter.Init()

	deployService := service.DeployService{
		DatabaseAdapter:   databaseAdapter,
		DockerAdapter:     dockerAdapter,
		FilesystemAdapter: filesystemAdapter,
		EventAdapter:      adapter.NewAdapterEvent(),
	}

	projectID, err := application.CreateProjectForCurrentFolder(&deployService)

	if err != nil {
		panic("Error creating project: " + err.Error())
	}

	getArgsOptions()

	if flags.help {
		showHelp()
	}

	port := "8080"
	mux := http.NewServeMux()
	// Add any other routes if needed, for example:
	// mux.HandleFunc("/some-path", someHandler)
	web.CreateMiddlewareWebFiles(mux)
	createGraphServer(&deployService, port, mux)

	// Wrap the mux with the CORS middleware
	handlerWithCORS := web.EnableCORS(mux)

	if !flags.noBrowser {
		adapter.OpenBrowser("http://localhost:" + port + "/project/" + projectID)
	}

	log.Println("Server started on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))

	// if flags.redeploy.deployId != "" {
	// 	application.ReDeployApplication(&deployService, flags.redeploy.deployId)
	// 	os.Exit(0)
	// } else {
	// 	api.InitValidator(app)
	// 	// api.CreateRoutes(app, &deployService)
	// 	app.StartServer()
	// }
}

func createGraphServer(deployService *service.DeployService, port string, mux *http.ServeMux) {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(deployService)}))

	mux.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	mux.Handle("/graphql", srv)
}

func getArgsOptions() {
	flag.BoolVar(&flags.noBrowser, "no-browser", false, "Do not open the browser")
	flag.StringVar(&flags.redeploy.deployId, "redeploy", "", "Redeploy application by deploy id")
	flag.BoolVar(&flags.help, "help", false, "Show help")
	flag.Parse()
}

func showHelp() {
	fmt.Println("Usage: main [options]")
	fmt.Println("  -no-browser    Do not open the browser")
	fmt.Println("  -redeploy <id> Redeploy application by deploy id")
	os.Exit(0)

}
