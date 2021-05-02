package main

import (
	"log"
	"net/http"

	"github.com/bccfilkom/drophere-go/routes"
	"github.com/bccfilkom/drophere-go/utils/env_driver"
	"github.com/go-chi/chi"
)

func main() {
	appEnv, err := env_driver.NewAppEnvironmentDriver()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	app := routes.Router{
		Router: router,
	}

	app.NewChiRoutes()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", appEnv.Port)
	err = http.ListenAndServe(":"+appEnv.Port, router)
	if err != nil {
		log.Fatal(err)
	}
}
